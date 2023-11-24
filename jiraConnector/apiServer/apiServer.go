package apiServer

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"Jira-analyzer/jiraConnector/configReader"
	"Jira-analyzer/jiraConnector/connector"
	"Jira-analyzer/jiraConnector/dbPusher"
	"Jira-analyzer/jiraConnector/logger"
	"Jira-analyzer/jiraConnector/transformer"
)

type ApiServer struct {
	configReader   *configReader.ConfigReader
	logger         *logger.JiraLogger
	jiraConnector  *connector.Connector
	databasePusher *dbPusher.DatabasePusher
	transformer    *transformer.Transformer
}

func CreateNewApiServer() *ApiServer {
	newReader := configReader.CreateNewConfigReader()
	return &ApiServer{
		configReader:   newReader,
		logger:         logger.CreateNewLogger(),
		jiraConnector:  connector.CreateNewJiraConnector(),
		databasePusher: dbPusher.CreateNewDatabasePusher(),
		transformer:    transformer.CreateNewTransformer(),
	}
}

func (server *ApiServer) updateProject(responseWriter http.ResponseWriter, request *http.Request) {
	if request.Method != "POST" {
		server.logger.Log(logger.ERROR, "Incorrect")
		responseWriter.WriteHeader(400)
		return
	}

	projectName := request.URL.Query().Get("project")

	if len(projectName) == 0 {
		server.logger.Log(logger.ERROR, "Incorrect")
		responseWriter.WriteHeader(400)
		return
	}

	issues, err := server.jiraConnector.GetProjectIssues(projectName)
	if err != nil {
		server.logger.Log(logger.ERROR, err.Error())
		responseWriter.WriteHeader(400)
		return
	}

	transformewIssues := server.transformer.TrasformData(issues)
	server.databasePusher.PushIssue(transformewIssues)
}

func (server *ApiServer) project(responseWriter http.ResponseWriter, request *http.Request) {
	fmt.Print("project work")
	if request.Method != "GET" {
		server.logger.Log(logger.ERROR, "Incorrect")
		return
	}

	limit, page, search := getProjectParametersFromRequest(request)

	responseWriter.Header().Set("Content-Type", "application/json")

	server.logger.Log(logger.INFO, "RETURN PROJECTS")

	projets, err := server.jiraConnector.GetProjects(limit, page, search)
	if err != nil {
		server.logger.Log(logger.ERROR, err.Error())
		responseWriter.WriteHeader(400)
		return
	}
	response, _ := json.Marshal(projets)
	fmt.Printf(string(response))
	responseWriter.Write(response)
}

func getProjectParametersFromRequest(request *http.Request) (int, int, string) {
	defaultLimit := 20
	defaultPage := 1
	defaultSearch := ""

	limit := request.URL.Query().Get("limit")
	if len(limit) != 0 {
		defaultLimit, _ = strconv.Atoi(limit) //нужно ли обрабатывать ошибки
	}

	page := request.URL.Query().Get("page")
	if len(page) != 0 {
		defaultLimit, _ = strconv.Atoi(page)
	}

	search := request.URL.Query().Get("search")
	if len(search) != 0 {
		defaultSearch = search
	}

	return defaultLimit, defaultPage, defaultSearch
}

func (server *ApiServer) StartServer() {
	server.logger.Log(logger.INFO, "Server start server...")

	http.HandleFunc("/api/v1/connector/updateProject", server.updateProject)
	http.HandleFunc("/api/v1/connector/projects", server.project)

	err := http.ListenAndServe("localhost:8003", nil)

	if err != nil {
		server.logger.Log(logger.ERROR, "Error while start a server")
	}
}
