package endpoints

import (
	"Jira-analyzer/analyzer/models"
	"Jira-analyzer/common/configReader"
	"Jira-analyzer/common/logger"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type ResourceHandler struct {
	configReader *configReader.ConfigReader
	logger       *logger.Logger
	database     *sql.DB
}

func CreateNewResourceHandler() *ResourceHandler {
	newReader := configReader.CreateNewConfigReader()
	sqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		newReader.GetHostDB(),
		newReader.GetPortDB(),
		newReader.GetUserDb(),
		newReader.GetPasswordDB(),
		newReader.GetDatabaseName())
	newDatabase, err := sql.Open("postgres", sqlInfo)
	newLogger := logger.CreateNewLogger()

	if err != nil {
		newLogger.Log(logger.ERROR, err.Error())

		return &ResourceHandler{}
	}

	return &ResourceHandler{
		configReader: newReader,
		logger:       newLogger,
		database:     newDatabase,
	}
}

func (resourceHandler *ResourceHandler) handleGetIssue(responseWriter http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	id, err := strconv.Atoi(vars["id"])

	if err != nil {
		resourceHandler.logger.Log(logger.ERROR, err.Error())
		responseWriter.WriteHeader(http.StatusBadRequest)

		return
	}

	issue, err := resourceHandler.getIssue(id)
	if err != nil {
		resourceHandler.logger.Log(logger.ERROR, err.Error())
		responseWriter.WriteHeader(http.StatusBadRequest)

		return
	}

	project, err := resourceHandler.getProject(issue.Project.ID)
	if err != nil {
		resourceHandler.logger.Log(logger.ERROR, err.Error())
		responseWriter.WriteHeader(http.StatusBadRequest)

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
		resourceHandler.logger.Log(logger.ERROR, err.Error())
		responseWriter.WriteHeader(http.StatusBadRequest)

		return
	}

	responseWriter.WriteHeader(http.StatusOK)
	resourceHandler.logger.Log(logger.INFO, "HandleGetIssue successfully")

	_, err = responseWriter.Write(response)

	if err != nil {
		responseWriter.WriteHeader(http.StatusBadRequest)

		return
	}

	responseWriter.WriteHeader(http.StatusOK)
}

func (resourceHandler *ResourceHandler) handleGetProject(responseWriter http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	id, err := strconv.Atoi(vars["id"])

	if err != nil {
		resourceHandler.logger.Log(logger.ERROR, err.Error())
		responseWriter.WriteHeader(http.StatusBadRequest)

		return
	}

	project, err := resourceHandler.getProject(id)
	if err != nil {
		resourceHandler.logger.Log(logger.ERROR, err.Error())
		responseWriter.WriteHeader(http.StatusBadRequest)

		return
	}

	var projectResponce = models.ResponseStruct{
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

	response, err := json.MarshalIndent(projectResponce, "", "\t")

	if err != nil {
		resourceHandler.logger.Log(logger.ERROR, err.Error())
		responseWriter.WriteHeader(http.StatusBadRequest)

		return
	}

	responseWriter.WriteHeader(http.StatusOK)
	_, err = responseWriter.Write(response)
}

func (resourceHandler *ResourceHandler) handleGetAllProjects(responseWriter http.ResponseWriter, request *http.Request) {

	limit, err := strconv.Atoi(request.URL.Query().Get("limit"))
	if err != nil {
		limit = 20
	}

	projects, err := resourceHandler.getAllProjects(limit)

	if err != nil {
		resourceHandler.logger.Log(logger.ERROR, err.Error())
		responseWriter.WriteHeader(http.StatusBadRequest)

		return
	}

	var projectResponce = models.ResponseStruct{
		Links: models.ListOfReferences{
			Issues:    models.Link{Href: "/api/v1/issues"},
			Projects:  models.Link{Href: "/api/v1/projects"},
			Histories: models.Link{Href: "/api/v1/histories"},
			Self:      models.Link{Href: ""},
		},
		Info:    projects,
		Message: "",
		Name:    "",
		Status:  true,
	}

	response, err := json.MarshalIndent(projectResponce, "", "\t")

	if err != nil {
		resourceHandler.logger.Log(logger.ERROR, fmt.Sprintf("Error: %s", err.Error()))
		responseWriter.WriteHeader(http.StatusBadRequest)

		return
	}

	responseWriter.WriteHeader(http.StatusOK)
	_, err = responseWriter.Write(response)
}

func (resourceHandler *ResourceHandler) handlePostIssue(responseWriter http.ResponseWriter, request *http.Request) {
	body, err := io.ReadAll(request.Body)

	if err != nil {
		resourceHandler.logger.Log(logger.ERROR, err.Error())
		responseWriter.WriteHeader(http.StatusBadRequest)

		return
	}

	var issueInfo models.IssueInfo
	err = json.Unmarshal(body, &issueInfo)

	if err != nil {
		resourceHandler.logger.Log(logger.ERROR, err.Error())
		responseWriter.WriteHeader(http.StatusBadRequest)

		return
	}

	id, err := resourceHandler.insertIssue(issueInfo)
	if err != nil {
		resourceHandler.logger.Log(logger.ERROR, err.Error())
		responseWriter.WriteHeader(http.StatusBadRequest)
		//как-то напишем об ошибке
		//statusCode = http.Status - подобрать верный статус
	} else {
		responseWriter.WriteHeader(http.StatusOK)
		//statusCode = http.Status - подобрать верный статус
	}

	var issuesResponce = models.ResponseStruct{
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

	response, err := json.MarshalIndent(issuesResponce, "", "\t")
	if err != nil {
		resourceHandler.logger.Log(logger.ERROR, err.Error())
		responseWriter.WriteHeader(http.StatusBadRequest)

		return
	}

	_, err = responseWriter.Write(response)
	if err != nil {
		responseWriter.WriteHeader(http.StatusBadRequest)

		return
	}

	responseWriter.WriteHeader(http.StatusCreated)
}

func (resourceHandler *ResourceHandler) handlePostProject(responseWriter http.ResponseWriter, request *http.Request) {
	body, err := io.ReadAll(request.Body)
	if err != nil {
		resourceHandler.logger.Log(logger.ERROR, err.Error())
		responseWriter.WriteHeader(http.StatusBadRequest)

		return
	}

	var projectInfo models.ProjectInfo
	err = json.Unmarshal(body, &projectInfo)

	if err != nil {
		resourceHandler.logger.Log(logger.ERROR, err.Error())
		responseWriter.WriteHeader(http.StatusBadRequest)

		return
	}

	id, err := resourceHandler.insertProject(projectInfo)
	if err != nil {
		resourceHandler.logger.Log(logger.ERROR, err.Error())
		responseWriter.WriteHeader(http.StatusBadRequest)
	} else {
		resourceHandler.logger.Log(logger.INFO, err.Error())
		responseWriter.WriteHeader(http.StatusOK)
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
		resourceHandler.logger.Log(logger.ERROR, err.Error())
		responseWriter.WriteHeader(http.StatusBadRequest)

		return
	}

	_, err = responseWriter.Write(response)
	if err != nil {
		responseWriter.WriteHeader(http.StatusBadRequest)

		return
	}

	responseWriter.WriteHeader(http.StatusCreated)
}

func (resourceHandler *ResourceHandler) handleDeleteProject(responseWriter http.ResponseWriter, request *http.Request) {
	body, err := io.ReadAll(request.Body)
	if err != nil {
		resourceHandler.logger.Log(logger.ERROR, err.Error())
		responseWriter.WriteHeader(http.StatusBadRequest)

		return
	}

	var projectInfo models.ProjectInfo
	err = json.Unmarshal(body, &projectInfo)

	if err != nil {
		resourceHandler.logger.Log(logger.ERROR, err.Error())
		responseWriter.WriteHeader(http.StatusBadRequest)

		return
	}

	_, err = resourceHandler.deleteProject(projectInfo)
	if err != nil {
		resourceHandler.logger.Log(logger.ERROR, err.Error())
		responseWriter.WriteHeader(http.StatusBadRequest)
	}

	responseWriter.WriteHeader(http.StatusOK)
}

func (server *ResourceHandler) handlers(router *mux.Router) {
	router.HandleFunc(server.configReader.GetApiPrefix()+server.configReader.GetResourcePrefix()+
		"/issues/{id:[0-9]+}", server.handleGetIssue).Methods("GET")
	router.HandleFunc(server.configReader.GetApiPrefix()+server.configReader.GetResourcePrefix()+
		"/projects/{id:[0-9]+}", server.handleGetProject).Methods("GET")
	router.HandleFunc(server.configReader.GetApiPrefix()+server.configReader.GetResourcePrefix()+
		"/projects", server.handleGetAllProjects).Methods("GET")

	router.HandleFunc(server.configReader.GetApiPrefix()+server.configReader.GetResourcePrefix()+
		"/issues/", server.handlePostIssue).Methods("POST")
	router.HandleFunc(server.configReader.GetApiPrefix()+server.configReader.GetResourcePrefix()+
		"/projects/", server.handlePostProject).Methods("POST")

	//router.HandleFunc(server.configReader.GetApiPrefix()+server.configReader.GetResourcePrefix()+
	//"/projects/", server.handleDeleteProject).Methods("DELETE")
}
