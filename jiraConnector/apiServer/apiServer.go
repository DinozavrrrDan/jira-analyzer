package apiServer

import (
	"encoding/json"
	"net/http"
	"strconv"

	"Jira-analyzer/jiraConnector/configReader"
	"Jira-analyzer/jiraConnector/connector"
	"Jira-analyzer/jiraConnector/dbPusher"
	"Jira-analyzer/jiraConnector/logger"
	"Jira-analyzer/jiraConnector/transformer"
)

//Как я понял он создает просто все части

type ApiServer struct {
	configReader   *configReader.ConfigRaeder
	logger         *logger.JiraLogger
	jiraConnector  *connector.Connector
	databasePusher *dbPusher.DatabasePusher
	transformer    *transformer.Transformer
	//	dbPusher     *dbPusher.DatabasePusher
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

func (server *ApiServer) UpdateProject(responseWriter http.ResponseWriter, request *http.Request) {

}

func (server *ApiServer) Project(responseWriter http.ResponseWriter, request *http.Request) {
	if request.Method != "GET" {
		server.logger.Log(logger.ERROR, "Incorrect")
		return
	}

	limit, page, search := getProjectParametersFromRequest(request)
	projets /*, err*/ := server.jiraConnector.GetProjects(limit, page, search)
	//if err != nil {}
	response, _ := json.Marshal(projets)
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
