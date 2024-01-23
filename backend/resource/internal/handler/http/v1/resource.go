package v1

import (
	"encoding/json"
	"fmt"
	"github.com/DinozvrrDan/jira-analyzer/backend/resource/config"
	"github.com/DinozvrrDan/jira-analyzer/backend/resource/internal/models"
	repository2 "github.com/DinozvrrDan/jira-analyzer/backend/resource/internal/repository"
	"github.com/DinozvrrDan/jira-analyzer/backend/resource/pkg/logger"
	"github.com/gorilla/mux"
	"io"
	"net/http"
	"strconv"
)

type ResourceHandler struct {
	resourceRep repository2.IResourceRepository
	log         *logger.Logger
	cfg         *config.Config
}

func NewResourceHandler(repositories *repository2.Repositories, log *logger.Logger, cfg *config.Config) *ResourceHandler {
	return &ResourceHandler{
		log:         log,
		resourceRep: repositories.ResourceRepository,
		cfg:         cfg,
	}
}

func (handler *ResourceHandler) GetResourceHandler(router *mux.Router) {
	router.HandleFunc("/issues/{id:[0-9]+}",
		handler.getIssue).Methods("GET")
	router.HandleFunc("/projects/{id:[0-9]+}",
		handler.getProject).Methods("GET")

	router.HandleFunc("/issues/",
		handler.postIssue).Methods("POST")
	router.HandleFunc("/projects/",
		handler.postProject).Methods("POST")

	router.HandleFunc("/projects/",
		handler.deleteProject).Methods("DELETE")

}

func (handler *ResourceHandler) getIssue(writer http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	id, err := strconv.Atoi(vars["id"])

	if err != nil {
		errorWriter(writer, handler.log, err.Error(), http.StatusBadRequest)
		return
	}

	issue, err := handler.resourceRep.GetIssueInfo(id)

	if err != nil {
		errorWriter(writer, handler.log, err.Error(), http.StatusBadRequest)
		return
	}

	project, err := handler.resourceRep.GetProjectInfo(issue.Project.Id)

	if err != nil {
		errorWriter(writer, handler.log, err.Error(), http.StatusBadRequest)
		return
	}

	var issueResponse = models.ResponseStruct{
		Links: models.ListOfReferences{
			Issues:    models.Link{Href: "/api/v1/issues"},
			Projects:  models.Link{Href: "/api/v1/projects"},
			Histories: models.Link{Href: "/api/v1/histories"},
			Self:      models.Link{Href: fmt.Sprintf("/api/v1/issues/%d", id)},
		},
		Info:    project,
		Message: "",
		Name:    "",
		Status:  true,
	}

	response, err := json.MarshalIndent(issueResponse, "", "\t")

	if err != nil {
		errorWriter(writer, handler.log, err.Error(), http.StatusBadRequest)
		return
	}

	writer.WriteHeader(http.StatusOK)

	_, err = writer.Write(response)

	if err != nil {
		errorWriter(writer, handler.log, err.Error(), http.StatusBadRequest)
		return
	}

	writer.WriteHeader(http.StatusOK)
}

func (handler *ResourceHandler) getProject(writer http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	id, err := strconv.Atoi(vars["id"])

	if err != nil {
		errorWriter(writer, handler.log, err.Error(), http.StatusBadRequest)
		return
	}

	project, err := handler.resourceRep.GetProjectInfo(id)

	if err != nil {
		errorWriter(writer, handler.log, err.Error(), http.StatusBadRequest)
		return
	}

	var projectResponse = models.ResponseStruct{
		Links: models.ListOfReferences{
			Issues:    models.Link{Href: "/api/v1/issues"},
			Projects:  models.Link{Href: "/api/v1/projects"},
			Histories: models.Link{Href: "/api/v1/histories"},
			Self:      models.Link{Href: fmt.Sprintf("/api/v1/issues/%d", id)},
		},
		Info:    project,
		Message: "",
		Name:    "",
		Status:  true,
	}

	response, err := json.MarshalIndent(projectResponse, "", "\t")

	if err != nil {
		errorWriter(writer, handler.log, err.Error(), http.StatusBadRequest)
		return
	}

	writer.WriteHeader(http.StatusOK)
	_, err = writer.Write(response)

	if err != nil {
		errorWriter(writer, handler.log, err.Error(), http.StatusBadRequest)
		return
	}
	writer.WriteHeader(http.StatusOK)
}

func (handler *ResourceHandler) postIssue(writer http.ResponseWriter, request *http.Request) {
	body, err := io.ReadAll(request.Body)

	if err != nil {
		errorWriter(writer, handler.log, err.Error(), http.StatusBadRequest)
		return
	}

	var issueInfo models.IssueInfo
	err = json.Unmarshal(body, &issueInfo)

	if err != nil {
		errorWriter(writer, handler.log, err.Error(), http.StatusBadRequest)
		return
	}

	id, err := handler.resourceRep.InsertIssue(issueInfo)
	if err != nil {
		errorWriter(writer, handler.log, err.Error(), http.StatusBadRequest)
		return
	} else {
		writer.WriteHeader(http.StatusOK)
	}

	var issuesResponse = models.ResponseStruct{
		Links: models.ListOfReferences{
			Issues:    models.Link{Href: "/api/v1/issues"},
			Projects:  models.Link{Href: "/api/v1/projects"},
			Histories: models.Link{Href: "/api/v1/histories"},
			Self:      models.Link{Href: fmt.Sprintf("/api/v1/issues/%d", id)},
		},
		Message: "",
		Name:    "",
		Status:  true,
	}

	response, err := json.MarshalIndent(issuesResponse, "", "\t")

	if err != nil {
		errorWriter(writer, handler.log, err.Error(), http.StatusBadRequest)
		return
	}

	_, err = writer.Write(response)
	if err != nil {
		errorWriter(writer, handler.log, err.Error(), http.StatusBadRequest)
		return
	}

	writer.WriteHeader(http.StatusCreated)
}

func (handler *ResourceHandler) postProject(writer http.ResponseWriter, request *http.Request) {
	body, err := io.ReadAll(request.Body)

	if err != nil {
		errorWriter(writer, handler.log, err.Error(), http.StatusBadRequest)
		return
	}

	var projectInfo models.ProjectInfo
	err = json.Unmarshal(body, &projectInfo)

	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	id, err := handler.resourceRep.InsertProject(projectInfo)
	if err != nil {
		errorWriter(writer, handler.log, err.Error(), http.StatusBadRequest)
		return
	} else {
		writer.WriteHeader(http.StatusOK)
	}

	var projectResponse = models.ResponseStruct{
		Links: models.ListOfReferences{
			Issues:    models.Link{Href: "/api/v1/issues"},
			Projects:  models.Link{Href: "/api/v1/projects"},
			Histories: models.Link{Href: "/api/v1/histories"},
			Self:      models.Link{Href: fmt.Sprintf("/api/v1/issues/%d", id)},
		},
		Message: "",
		Name:    "",
		Status:  true,
	}

	response, err := json.MarshalIndent(projectResponse, "", "\t")

	if err != nil {
		errorWriter(writer, handler.log, err.Error(), http.StatusBadRequest)
		return
	}

	_, err = writer.Write(response)

	if err != nil {
		errorWriter(writer, handler.log, err.Error(), http.StatusBadRequest)
		return
	}
	writer.WriteHeader(http.StatusCreated)
}

func (handler *ResourceHandler) deleteProject(writer http.ResponseWriter, request *http.Request) {
	body, err := io.ReadAll(request.Body)

	if err != nil {
		errorWriter(writer, handler.log, err.Error(), http.StatusBadRequest)
		return
	}

	var projectInfo models.ProjectInfo
	err = json.Unmarshal(body, &projectInfo)

	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	err = handler.resourceRep.DeleteProject(projectInfo.Title)
	if err != nil {
		errorWriter(writer, handler.log, err.Error(), http.StatusBadRequest)
		return
	} else {
		writer.WriteHeader(http.StatusOK)
	}

	var projectResponse = models.ResponseStruct{
		Links: models.ListOfReferences{
			Issues:    models.Link{Href: "/api/v1/issues"},
			Projects:  models.Link{Href: "/api/v1/projects"},
			Histories: models.Link{Href: "/api/v1/histories"},
			Self:      models.Link{Href: fmt.Sprintf("/api/v1/issues/%d", 0)},
		},
		Message: "",
		Name:    "",
		Status:  true,
	}

	response, err := json.MarshalIndent(projectResponse, "", "\t")

	if err != nil {
		errorWriter(writer, handler.log, err.Error(), http.StatusBadRequest)
		return
	}

	_, err = writer.Write(response)

	if err != nil {
		errorWriter(writer, handler.log, err.Error(), http.StatusBadRequest)
		return
	}
	writer.WriteHeader(http.StatusCreated)
}
