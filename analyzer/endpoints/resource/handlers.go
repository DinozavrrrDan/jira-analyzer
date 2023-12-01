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
	logger       *logger.JiraLogger
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

	if err != nil {
		panic(err)
	}

	return &ResourceHandler{
		configReader: newReader,
		logger:       logger.CreateNewLogger(),
		database:     newDatabase,
	}
}

func (resourceHandler *ResourceHandler) getIssue(responseWriter http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	id, err := strconv.Atoi(vars["id"])

	if err != nil {
		resourceHandler.logger.Log(logger.ERROR, err.Error())
		responseWriter.WriteHeader(http.StatusBadRequest)

		return
	}

	//issue, err := resourceHandler.GetIssueInfo(id)
	if err != nil {
		resourceHandler.logger.Log(logger.ERROR, err.Error())
		responseWriter.WriteHeader(http.StatusBadRequest)

		return
	}

	//project, err := resourceHandler.GetProjectInfo(issue.Project.Id)
	if err != nil {
		resourceHandler.logger.Log(logger.ERROR, err.Error())
		responseWriter.WriteHeader(http.StatusBadRequest)

		return
	}

	var issueResponce = models.ResponseStrucrt{
		Links: models.ListOfReferens{
			Issues:    models.Link{Href: "/api/v1/issues"},
			Projects:  models.Link{Href: "/api/v1/projects"},
			Histories: models.Link{Href: "/api/v1/histories"},
			Self:      models.Link{Href: fmt.Sprintf("/api/v1/issues/%d", id)},
		},
		Message: "",
		Name:    "",
		Status:  true,
	}

	response, err := json.MarshalIndent(issueResponce, "", "\t")

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

func (resourceHandler *ResourceHandler) getHistory(responseWriter http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	id, err := strconv.Atoi(vars["id"])

	if err != nil {
		resourceHandler.logger.Log(logger.ERROR, err.Error())
		responseWriter.WriteHeader(http.StatusBadRequest)

		return
	}

	//history, err := resourceHandler.GetHistoryInfo(id)
	if err != nil {
		resourceHandler.logger.Log(logger.ERROR, err.Error())
		responseWriter.WriteHeader(http.StatusBadRequest)

		return
	}

	var historyResponce = models.ResponseStrucrt{
		Links: models.ListOfReferens{
			Issues:    models.Link{Href: "/api/v1/issues"},
			Projects:  models.Link{Href: "/api/v1/projects"},
			Histories: models.Link{Href: "/api/v1/histories"},
			Self:      models.Link{Href: fmt.Sprintf("/api/v1/issues/%d", id)},
		},
		Message: "",
		Name:    "",
		Status:  true,
	}

	response, err := json.MarshalIndent(historyResponce, "", "\t")

	if err != nil {
		resourceHandler.logger.Log(logger.ERROR, err.Error())
		responseWriter.WriteHeader(http.StatusBadRequest)

		return
	}

	responseWriter.WriteHeader(http.StatusOK)
	_, err = responseWriter.Write(response)

	if err != nil {
		resourceHandler.logger.Log(logger.ERROR, err.Error())
		responseWriter.WriteHeader(http.StatusBadRequest)

		return
	}

	responseWriter.WriteHeader(http.StatusOK)
}

func (resourceHandler *ResourceHandler) getProject(responseWriter http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	id, err := strconv.Atoi(vars["id"])

	if err != nil {
		resourceHandler.logger.Log(logger.ERROR, err.Error())
		responseWriter.WriteHeader(http.StatusBadRequest)

		return
	}

	//project, err := resourceHandler.GetProjectInfo(id)
	if err != nil {
		resourceHandler.logger.Log(logger.ERROR, err.Error())
		responseWriter.WriteHeader(http.StatusBadRequest)

		return
	}

	var projectResponce = models.ResponseStrucrt{
		Links: models.ListOfReferens{
			Issues:    models.Link{Href: "/api/v1/issues"},
			Projects:  models.Link{Href: "/api/v1/projects"},
			Histories: models.Link{Href: "/api/v1/histories"},
			Self:      models.Link{Href: fmt.Sprintf("/api/v1/issues/%d", id)},
		},
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

	if err != nil {
		resourceHandler.logger.Log(logger.ERROR, err.Error())
		responseWriter.WriteHeader(http.StatusBadRequest)

		return
	}

	responseWriter.WriteHeader(http.StatusOK)
}

func (resourceHandler *ResourceHandler) postIssue(responseWriter http.ResponseWriter, request *http.Request) {
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

	id, err := resourceHandler.InsertIssue(issueInfo)
	if err != nil {
		resourceHandler.logger.Log(logger.ERROR, err.Error())
		responseWriter.WriteHeader(http.StatusBadRequest)
		//как-то напишем об ошибке
		//statusCode = http.Status - подобрать верный статус
	} else {
		responseWriter.WriteHeader(http.StatusOK)
		//statusCode = http.Status - подобрать верный статус
	}

	var issuesResponce = models.ResponseStrucrt{
		Links: models.ListOfReferens{
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

func (resourceHandler *ResourceHandler) postHistory(responseWriter http.ResponseWriter, request *http.Request) {
	body, err := io.ReadAll(request.Body)

	if err != nil {
		resourceHandler.logger.Log(logger.ERROR, err.Error())
		responseWriter.WriteHeader(http.StatusBadRequest)

		return
	}

	var historyInfo models.HistoryInfo
	err = json.Unmarshal(body, &historyInfo)

	if err != nil {
		resourceHandler.logger.Log(logger.ERROR, err.Error())
		responseWriter.WriteHeader(http.StatusBadRequest)

		return
	}

	id, err := resourceHandler.InsertHistory(historyInfo)
	if err != nil {
		resourceHandler.logger.Log(logger.ERROR, err.Error())
		responseWriter.WriteHeader(http.StatusBadRequest)
		//как-то напишем об ошибке
		//statusCode = http.Status - подобрать верный статус
	} else {
		responseWriter.WriteHeader(http.StatusOK)
		//statusCode = http.Status - подобрать верный статус
	}

	var historyResponce = models.ResponseStrucrt{
		Links: models.ListOfReferens{
			Issues:    models.Link{Href: "/api/v1/issues"},
			Projects:  models.Link{Href: "/api/v1/projects"},
			Histories: models.Link{Href: "/api/v1/histories"},
			Self:      models.Link{Href: fmt.Sprintf("/api/v1/issues/%d", id)},
		},
		Message: "",
		Name:    "",
		Status:  true,
	}

	response, err := json.MarshalIndent(historyResponce, "", "\t")
	if err != nil {
		resourceHandler.logger.Log(logger.ERROR, err.Error())
		responseWriter.WriteHeader(http.StatusBadRequest)

		return
	}

	_, err = responseWriter.Write(response)

	if err != nil {
		resourceHandler.logger.Log(logger.ERROR, err.Error())
		responseWriter.WriteHeader(http.StatusBadRequest)

		return
	}

	responseWriter.WriteHeader(http.StatusCreated)
}

func (resourceHandler *ResourceHandler) postProject(responseWriter http.ResponseWriter, request *http.Request) {
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

	id, err := resourceHandler.InsertProject(projectInfo)
	if err != nil {
		resourceHandler.logger.Log(logger.ERROR, err.Error())
		responseWriter.WriteHeader(http.StatusBadRequest)
		//как-то напишем об ошибке
		//statusCode = http.Status - подобрать верный статус
	} else {
		resourceHandler.logger.Log(logger.INFO, err.Error())
		responseWriter.WriteHeader(http.StatusOK)
		//statusCode = http.Status - подобрать верный статус
	}

	var projectResponce = models.ResponseStrucrt{
		Links: models.ListOfReferens{
			Issues:    models.Link{Href: "/api/v1/issues"},
			Projects:  models.Link{Href: "/api/v1/projects"},
			Histories: models.Link{Href: "/api/v1/histories"},
			Self:      models.Link{Href: fmt.Sprintf("/api/v1/issues/%d", id)},
		},
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

	_, err = responseWriter.Write(response)
	if err != nil {
		responseWriter.WriteHeader(http.StatusBadRequest)

		return
	}

	responseWriter.WriteHeader(http.StatusCreated)
}

func (server *ResourceHandler) handlers(router *mux.Router) {
	router.HandleFunc(server.configReader.GetApiPrefix()+server.configReader.GetResourcePrefix()+
		"issues/{id:[0-9]+}", server.getIssue).Methods("GET")
	router.HandleFunc(server.configReader.GetApiPrefix()+server.configReader.GetResourcePrefix()+
		"projects/{id:[0-9]+}", server.getProject).Methods("GET")
	router.HandleFunc(server.configReader.GetApiPrefix()+server.configReader.GetResourcePrefix()+
		"histories/{id:[0-9]+}", server.getHistory).Methods("GET")

	router.HandleFunc(server.configReader.GetApiPrefix()+server.configReader.GetResourcePrefix()+
		"issues/", server.postIssue).Methods("POST")
	router.HandleFunc(server.configReader.GetApiPrefix()+server.configReader.GetResourcePrefix()+
		"projects/", server.postProject).Methods("POST")
	router.HandleFunc(server.configReader.GetApiPrefix()+server.configReader.GetResourcePrefix()+
		"histories/", server.postHistory).Methods("POST")

}