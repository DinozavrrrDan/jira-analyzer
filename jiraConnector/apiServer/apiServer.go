package apiServer

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"Jira-analyzer/analyzer/models"
	"Jira-analyzer/common/configReader"
	"Jira-analyzer/common/logger"
	"Jira-analyzer/jiraConnector/connector"
	"Jira-analyzer/jiraConnector/dbPusher"
	"Jira-analyzer/jiraConnector/transformer"
)

type ApiServer struct {
	configReader   *configReader.ConfigReader
	logger         *logger.Logger
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
	if request.Method != http.MethodPost {
		server.logger.Log(logger.ERROR, "Incorrect")
		responseWriter.WriteHeader(http.StatusBadRequest)

		return
	}

	projectName := request.URL.Query().Get("project")

	if len(projectName) == 0 {
		server.logger.Log(logger.ERROR, "Incorrect")
		responseWriter.WriteHeader(http.StatusBadRequest)

		return
	}

	issues, err := server.jiraConnector.GetProjectIssues(projectName)
	fmt.Println(issues)
	response, err := json.MarshalIndent(issues, "", "\t")
	responseWriter.Write(response)

	if err != nil {
		server.logger.Log(logger.ERROR, err.Error())
		responseWriter.WriteHeader(http.StatusBadRequest)

		return
	}

	transformedIssues := server.transformer.TransformData(issues)
	server.databasePusher.PushIssues(transformedIssues)
}

func (server *ApiServer) project(responseWriter http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodGet {
		server.logger.Log(logger.ERROR, "Incorrect")

		return
	}

	limit, page, search := getProjectParametersFromRequest(request)

	responseWriter.Header().Set("Content-Type", "application/json")

	projects, pages, err := server.jiraConnector.GetProjects(limit, page, search)
	if err != nil {
		server.logger.Log(logger.ERROR, err.Error())
		responseWriter.WriteHeader(http.StatusBadRequest)

		return
	}

	var issueResponse = models.ResponseStruct{
		Links: models.ListOfReferences{
			Issues:    models.Link{Href: "/api/v1/issues"},
			Projects:  models.Link{Href: "/api/v1/projects"},
			Histories: models.Link{Href: "/api/v1/histories"},
			Self:      models.Link{Href: fmt.Sprintf("/api/v1/issues/%d", 1)},
		},
		Info:    projects,
		Message: "Hello from connector",
		Name:    "",
		PageInfo: models.Page{
			TotalPageCount:     pages.TotalPageCount,
			CurrentPageNumber:  pages.CurrentPageNumber,
			TotalProjectsCount: pages.TotalProjectsCount,
		},
		Status: true,
	}

	response, _ := json.MarshalIndent(issueResponse, "", "\t")
	_, err = responseWriter.Write(response)

	if err != nil {
		responseWriter.WriteHeader(http.StatusBadRequest)

		return
	}
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
		defaultPage, _ = strconv.Atoi(page)
	}

	search := request.URL.Query().Get("search")
	if len(search) != 0 {
		defaultSearch = search
	}

	return defaultLimit, defaultPage, defaultSearch
}

func (server *ApiServer) StartServer() error {
	server.logger.Log(logger.INFO, "Server start server...")
	server.handlers()
	err := http.ListenAndServe(server.configReader.GetConnectorHost()+":"+server.configReader.GetConnectorPort(), nil)

	if err != nil {
		server.logger.Log(logger.ERROR, "Error while start a server")
		return fmt.Errorf("StartServer: %w", err)
	}

	return nil
}

func (server *ApiServer) handlers() {
	http.HandleFunc(server.configReader.GetApiPrefix()+
		server.configReader.GetConnectorPref()+
		"/updateProject",
		server.updateProject)

	http.HandleFunc(server.configReader.GetApiPrefix()+
		server.configReader.GetConnectorPref()+
		"/projects",
		server.project)
}
