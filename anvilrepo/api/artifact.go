package api

import (
	"errors"
	"fmt"
	"strings"
)

func NewArtifact(user string, module string, version string, fileName string) (*Artifact, error) {
	return &Artifact{
		Id:       fmt.Sprintf("%s/%s/%s", user, module, fileName),
		User:     user,
		Module:   module,
		Version:  version,
		FileName: fileName,
	}, nil
}

type Artifact struct {
	Id       string `json:"id"`
	User     string `json:"user"`
	Module   string `json:"module"`
	Version  string `json:"version"`
	FileName string `json:"filename"`
}

func ArtifactFromId(id string) (*Artifact, error) {
	s := strings.Split(id, "/")
	user, module, fileName := s[0], s[1], s[2]

	version, err := getVersion(id)
	if err != nil {
		return nil, err
	}
	return &Artifact{
		User:     user,
		Module:   module,
		Version:  version,
		FileName: fileName,
	}, nil
}

func getVersion(fileName string) (string, error) {
	iStart := strings.LastIndex(fileName, "-")
	iEnd := strings.LastIndex(fileName, ".tar.gz")

	if iStart == -1 || iEnd == -1 {
		return "", errors.New(fmt.Sprintf("Artifact %s doesn't have a parsable version: %d:%d", fileName, iStart, iEnd))
	}
	version := fileName[iStart+1 : iEnd]
	return version, nil
}
