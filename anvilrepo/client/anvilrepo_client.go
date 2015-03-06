package client

import (
	"bytes"
	"fmt"
	"github.com/benschw/anvil-mgr/anvilrepo/api"
	"github.com/benschw/opin-go/rest"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
)

var _ = log.Print

type AnvilRepoClient struct {
	Host string
}

func (c *AnvilRepoClient) AddArtifact(user string, module string, path string) (api.Artifact, error) {
	var artifact api.Artifact

	url := fmt.Sprintf("%s/api/repo/%s/%s", c.Host, user, module)

	req, err := newfileUploadRequest(url, "artifact", path)
	if err != nil {
		return artifact, err
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return artifact, err
	}

	err = rest.ProcessResponseEntity(resp, &artifact, http.StatusCreated)
	return artifact, err
}

func (c *AnvilRepoClient) FindAllArtifacts() ([]api.Artifact, error) {
	var artifacts []api.Artifact

	url := fmt.Sprintf("%s/api/repo", c.Host)
	r, err := rest.MakeRequest("GET", url, nil)
	if err != nil {
		return artifacts, err
	}
	err = rest.ProcessResponseEntity(r, &artifacts, http.StatusOK)
	return artifacts, nil
}
func (c *AnvilRepoClient) FindAllArtifactVersions(user string, module string) ([]api.Artifact, error) {
	var artifacts []api.Artifact

	url := fmt.Sprintf("%s/api/repo/%s/%s", c.Host, user, module)
	r, err := rest.MakeRequest("GET", url, nil)
	if err != nil {
		return artifacts, err
	}
	err = rest.ProcessResponseEntity(r, &artifacts, http.StatusOK)
	return artifacts, nil
}

func newfileUploadRequest(uri string, paramName, path string) (*http.Request, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile(paramName, filepath.Base(path))
	if err != nil {
		return nil, err
	}
	_, err = io.Copy(part, file)

	//	_ = writer.WriteField("f", "b")

	if err = writer.Close(); err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", uri, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	return req, nil
}
