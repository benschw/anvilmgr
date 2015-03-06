package anvilrepo

import (
	"fmt"
	"github.com/benschw/opin-go/rest"
	"log"
	"net/http"
)

var _ = log.Print

type AnvilArtifactResource struct {
	Repo *RepoClient
}

func (r *AnvilArtifactResource) AddArtifact(res http.ResponseWriter, req *http.Request) {
	user, module, err := getPathArgs(req)
	if err != nil {
		log.Println(err)
		rest.SetBadRequestResponse(res)
		return
	}

	// capture input
	file, header, err := req.FormFile("artifact")
	if err != nil {
		log.Println(err)
		rest.SetBadRequestResponse(res)
		return
	}
	defer file.Close()

	// check if artifact already exists
	id := fmt.Sprintf("%s/%s/%s", user, module, header.Filename)
	if _, err = r.Repo.getArtifact(id); err == nil {
		rest.SetConflictResponse(res)
		return
	}
	// add artifact
	artifact, err := r.Repo.addArtifact(user, module, header.Filename, file)
	if err != nil {
		rest.SetInternalServerErrorResponse(res, err)
		return
	}

	location := fmt.Sprintf("api/repo/%s/%s/%s", artifact.Id)
	if err := rest.SetCreatedResponse(res, artifact, location); err != nil {
		rest.SetInternalServerErrorResponse(res, err)
		return
	}
}

func (r *AnvilArtifactResource) GetModuleArtifacts(res http.ResponseWriter, req *http.Request) {
	user, module, err := getPathArgs(req)
	if err != nil {
		log.Println(err)
		rest.SetBadRequestResponse(res)
		return
	}
	artifacts, err := r.Repo.getModuleArtifacts(user, module)
	if err != nil {
		rest.SetInternalServerErrorResponse(res, err)
		return
	}
	if err := rest.SetOKResponse(res, artifacts); err != nil {
		rest.SetInternalServerErrorResponse(res, err)
		return
	}

}

func (r *AnvilArtifactResource) FindAll(res http.ResponseWriter, req *http.Request) {
	artifacts, err := r.Repo.getAllArtifacts()
	if err != nil {
		rest.SetInternalServerErrorResponse(res, err)
		return
	}
	if err := rest.SetOKResponse(res, artifacts); err != nil {
		rest.SetInternalServerErrorResponse(res, err)
		return
	}
}

func getPathArgs(req *http.Request) (string, string, error) {
	user, err := rest.PathString(req, "user")
	if err != nil {
		return "", "", err
	}
	module, err := rest.PathString(req, "module")
	if err != nil {
		return "", "", err
	}
	return user, module, nil
}
