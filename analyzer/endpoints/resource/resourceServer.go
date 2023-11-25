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

type ResourceServer struct {
	configReader *configReader.ConfigReader
	logger       *logger.JiraLogger
	database     *sql.DB
}

func CreateNewResourceServer() *ResourceServer {
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
	return &ResourceServer{
		configReader: newReader,
		logger:       logger.CreateNewLogger(),
		database:     newDatabase,
	}
}

func (resourceServer *ResourceServer) getIssue(responseWriter http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		resourceServer.logger.Log(logger.ERROR, err.Error())
		responseWriter.WriteHeader(http.StatusBadRequest)
		return
	}

	issue, err := GetIssueInfoByID(id)
	if err != nil {
		resourceServer.logger.Log(logger.ERROR, err.Error())
		responseWriter.WriteHeader(http.StatusBadRequest)
		return
	}

	project, err := GetProjectInfoByID(issue.ProjectID)
	if err != nil {
		resourceServer.logger.Log(logger.ERROR, err.Error())
		responseWriter.WriteHeader(400)
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
		resourceServer.logger.Log(logger.ERROR, err.Error())
		responseWriter.WriteHeader(400)
		return
	}
	responseWriter.WriteHeader(200)
	resourceServer.logger.Log(logger.INFO, "HandleGetIssue successfully")
	_, err = responseWriter.Write(response)
	if err != nil {
		responseWriter.WriteHeader(http.StatusBadRequest)
		return
	}
	responseWriter.WriteHeader(http.StatusOK)
}

func (resourceServer *ResourceServer) getHistory(responseWriter http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		resourceServer.logger.Log(logger.ERROR, err.Error())
		http.Error(responseWriter, err.Error(), http.StatusBadRequest)
		return
	}

	history, err := GetAllHistoryInfoByIssueID(id)
	if err != nil {
		resourceServer.logger.Log(logger.ERROR, err.Error())
		responseWriter.WriteHeader(400)
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
		resourceServer.logger.Log(logger.ERROR, err.Error())
		responseWriter.WriteHeader(400)
		return
	}
	responseWriter.WriteHeader(200)
	_, err = responseWriter.Write(response)
	if err != nil {
		responseWriter.WriteHeader(http.StatusBadRequest)
		return
	}
	responseWriter.WriteHeader(http.StatusOK)
}

func (resourceServer *ResourceServer) getProject(responseWriter http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		resourceServer.logger.Log(logger.ERROR, err.Error())
		responseWriter.WriteHeader(http.StatusBadRequest)
		return
	}

	project, err := GetProjectInfoByID(id)
	if err != nil {
		resourceServer.logger.Log(logger.ERROR, err.Error())
		responseWriter.WriteHeader(400)
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
		resourceServer.logger.Log(logger.ERROR, err.Error())
		responseWriter.WriteHeader(400)
		return
	}
	responseWriter.WriteHeader(200)
	_, err = responseWriter.Write(response)
	if err != nil {
		responseWriter.WriteHeader(http.StatusBadRequest)
		return
	}
	responseWriter.WriteHeader(http.StatusOK)
}

func (resourceServer *ResourceServer) postIssue(responseWriter http.ResponseWriter, request *http.Request) {
	body, err := io.ReadAll(request.Body)

	if err != nil {
		resourceServer.logger.Log(logger.ERROR, err.Error())
		responseWriter.WriteHeader(http.StatusBadRequest)
		return
	}

	var requestDataIssue models.IssueInfo
	err = json.Unmarshal(body, &requestDataIssue)

	if err != nil {
		resourceServer.logger.Log(logger.ERROR, err.Error())
		responseWriter.WriteHeader(http.StatusBadRequest)
		return
	}
	var statusCode int
	//id, err := PutIssueIntoDB функция которая будет помещать узел в БД БОРЯ
	if err != nil {
		//как-то напишем об ошибке
		//statusCode = http.Status - подобрать верный статус
	} else {
		//statusCode = http.Status - подобрать верный статус
	}
	statusCode = statusCode + 1 // заглушка

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
		resourceServer.logger.Log(logger.ERROR, err.Error())
		responseWriter.WriteHeader(http.StatusBadRequest)
		return
	}
	_, err = responseWriter.Write(response)
	if err != nil {
		responseWriter.WriteHeader(http.StatusBadRequest)
		return
	}
	responseWriter.WriteHeader(http.StatusOK)
}

func (resourceServer *ResourceServer) postHistory(responseWriter http.ResponseWriter, request *http.Request) {
	body, err := io.ReadAll(request.Body)

	if err != nil {
		resourceServer.logger.Log(logger.ERROR, err.Error())
		responseWriter.WriteHeader(http.StatusBadRequest)
		return
	}

	var requestDataIssue models.HistoryInfo
	err = json.Unmarshal(body, &requestDataIssue)

	if err != nil {
		resourceServer.logger.Log(logger.ERROR, err.Error())
		responseWriter.WriteHeader(http.StatusBadRequest)
		return
	}
	var statusCode int
	//id, err := PutIssueIntoDB функция которая будет помещать узел в БД БОРЯ
	if err != nil {
		//как-то напишем об ошибке
		//statusCode = http.Status - подобрать верный статус
	} else {
		//statusCode = http.Status - подобрать верный статус
	}
	statusCode = statusCode + 1 // заглушка
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
		resourceServer.logger.Log(logger.ERROR, err.Error())
		responseWriter.WriteHeader(http.StatusBadRequest)
		return
	}
	_, err = responseWriter.Write(response)
	if err != nil {
		responseWriter.WriteHeader(http.StatusBadRequest)
		return
	}
	responseWriter.WriteHeader(http.StatusOK)
}

func (resourceServer *ResourceServer) postProject(responseWriter http.ResponseWriter, request *http.Request) {
	body, err := io.ReadAll(request.Body)
	if err != nil {
		resourceServer.logger.Log(logger.ERROR, err.Error())
		responseWriter.WriteHeader(http.StatusBadRequest)
		return
	}

	var requestDataIssue models.ProjectInfo
	err = json.Unmarshal(body, &requestDataIssue)

	if err != nil {
		resourceServer.logger.Log(logger.ERROR, err.Error())
		responseWriter.WriteHeader(http.StatusBadRequest)
		return
	}
	var statusCode int
	//id, err := PutIssueIntoDB функция которая будет помещать узел в БД БОРЯ
	if err != nil {

	} else {
		//statusCode = http.Status - подобрать верный статус
	}
	statusCode = statusCode + 1 // заглушка
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
		resourceServer.logger.Log(logger.ERROR, err.Error())
		responseWriter.WriteHeader(http.StatusBadRequest)
		return
	}
	_, err = responseWriter.Write(response)
	if err != nil {
		responseWriter.WriteHeader(http.StatusBadRequest)
		return
	}
	responseWriter.WriteHeader(http.StatusOK)
}

func (server *ResourceServer) StartServer() {
	server.logger.Log(logger.INFO, "Server start server...")

	router := mux.NewRouter()

	server.handlers(router)

	err := http.ListenAndServe(server.configReader.GetResourceHost()+":"+server.configReader.GetResourceHost(), router)
	if err != nil {
		server.logger.Log(logger.ERROR, "Error while start a server")
	}

}

func (server *ResourceServer) handlers(router *mux.Router) {
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
