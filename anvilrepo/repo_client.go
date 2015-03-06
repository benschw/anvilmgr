package anvilrepo

import (
	"errors"
	"fmt"
	"github.com/benschw/anvil-mgr/anvilrepo/api"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"os"
	"strings"
)

var _ = log.Print

type RepoClient struct {
	Path string
}

func (c *RepoClient) addArtifact(user string, module string, fileName string, file multipart.File) (*api.Artifact, error) {
	version, err := parseFileName(fileName)
	if err != nil {
		return nil, err
	}

	// ensure path is created
	modulePath := fmt.Sprintf("%s/%s/%s", c.Path, user, module)
	if _, err := os.Stat(modulePath); err != nil {
		os.MkdirAll(modulePath, 0755)
	}

	// land artifact in repo
	out, err := os.Create(modulePath + "/" + fileName)
	if err != nil {
		return nil, err
	}
	defer out.Close()

	if _, err = io.Copy(out, file); err != nil {
		return nil, err
	}

	return api.NewArtifact(user, module, version, fileName)
}

func (c *RepoClient) getArtifact(id string) (*api.Artifact, error) {
	artifact, err := api.ArtifactFromId(id)
	if err != nil {
		return nil, err
	}

	filePath := fmt.Sprintf("%s/%s", c.Path, id)
	if _, err := os.Stat(filePath); err != nil {
		return nil, err
	}
	// log.Printf("%+v, %s", artifact, filePath)
	return artifact, nil
}

func (c *RepoClient) getModuleArtifacts(user string, module string) ([]*api.Artifact, error) {
	return getArtifacts(c.Path, user, module)

}

func (c *RepoClient) getAllArtifacts() ([]*api.Artifact, error) {
	return getModulePaths(c.Path)
}

func getModulePaths(path string) ([]*api.Artifact, error) {
	var artifacts []*api.Artifact

	users, err := ioutil.ReadDir(path)
	if err != nil {
		return artifacts, err
	}

	for _, u := range users {
		ms, err := ioutil.ReadDir(path + "/" + u.Name())
		if err != nil {
			return artifacts, err
		}
		for _, m := range ms {
			moduleArtifacts, err := getArtifacts(path, u.Name(), m.Name())
			if err != nil {
				return artifacts, err
			}
			artifacts = append(artifacts, moduleArtifacts...)
		}
	}
	return artifacts, nil
}

func getArtifacts(path string, user string, module string) ([]*api.Artifact, error) {
	var artifacts []*api.Artifact

	modulePath := fmt.Sprintf("%s/%s/%s", path, user, module)
	files, err := ioutil.ReadDir(modulePath)
	if err != nil {
		return artifacts, err
	}

	for _, f := range files {
		v, err := parseFileName(f.Name())
		if err != nil {
			return artifacts, err
		}
		artifact, err := api.NewArtifact(user, module, v, f.Name())
		if err != nil {
			return artifacts, err
		}

		artifacts = append(artifacts, artifact)
	}
	return artifacts, nil
}

func parseFileName(fileName string) (string, error) {
	iStart := strings.LastIndex(fileName, "-")
	iEnd := strings.LastIndex(fileName, ".tar.gz")

	if iStart == -1 || iEnd == -1 {
		return "", errors.New(fmt.Sprintf("Artifact %s doesn't have a parsable version: %d:%d", fileName, iStart, iEnd))
	}
	version := fileName[iStart+1 : iEnd]
	return version, nil
}
