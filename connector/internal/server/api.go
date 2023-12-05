package server

import (
	"connector/api"
	"connector/config"
	"connector/internal/service"
	"connector/pkg/logger"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

type ApiServer struct {
	connectorSvc   service.Connector
	transformerSvc service.Transformer
	dbPusherSvc    service.DatabasePusher
	cfg            *config.Reader
	log            *logger.Logger
}

func NewApiServer(services *service.Services,
	log *logger.Logger, cfg *config.Reader) *ApiServer {
	return &ApiServer{
		cfg:            cfg,
		log:            log,
		connectorSvc:   services.Connector,
		transformerSvc: services.Transformer,
		dbPusherSvc:    services.DatabasePusher,
	}
}

func (server *ApiServer) updateProject(responseWriter http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodPost {
		server.log.Log(logger.ERROR, "Incorrect")
		responseWriter.WriteHeader(http.StatusBadRequest)

		return
	}

	projectName := request.URL.Query().Get("project")

	if len(projectName) == 0 {
		server.log.Log(logger.ERROR, "Incorrect")
		responseWriter.WriteHeader(http.StatusBadRequest)

		return
	}

	issues, err := server.connectorSvc.GetProjectIssues(projectName)
	fmt.Println(issues)
	response, err := json.MarshalIndent(issues, "", "\t")
	responseWriter.Write(response)

	if err != nil {
		server.log.Log(logger.ERROR, err.Error())
		responseWriter.WriteHeader(http.StatusBadRequest)

		return
	}

	transformedIssues := server.transformerSvc.TransformData(issues)
	server.dbPusherSvc.PushIssue(transformedIssues)
}

func (server *ApiServer) project(responseWriter http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodGet {
		server.log.Log(logger.ERROR, "Incorrect")

		return
	}

	limit, page, search := getProjectParametersFromRequest(request)

	responseWriter.Header().Set("Content-Type", "application/json")

	projects, pages, err := server.connectorSvc.GetProjects(limit, page, search)
	if err != nil {
		server.log.Log(logger.ERROR, err.Error())
		responseWriter.WriteHeader(http.StatusBadRequest)

		return
	}

	var issueResponse = api.ResponseStruct{
		Links: api.ListOfReferences{
			Issues:    api.Link{Href: "/api/v1/issues"},
			Projects:  api.Link{Href: "/api/v1/projects"},
			Histories: api.Link{Href: "/api/v1/histories"},
			Self:      api.Link{Href: fmt.Sprintf("/api/v1/issues/%d", 1)},
		},
		Info:    projects,
		Message: "Hello from connector",
		Name:    "",
		PageInfo: api.Page{
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

func (server *ApiServer) StartServer() {
	fmt.Println(server.cfg.GetConnectorHost() + ":" + server.cfg.GetConnectorPort())
	server.log.Log(logger.INFO, "Server start server...")
	server.handlers()
	err := http.ListenAndServe(server.cfg.GetConnectorHost()+":"+server.cfg.GetConnectorPort(), nil)

	if err != nil {
		server.log.Log(logger.ERROR, "Error while start a server")
	}
}

func (server *ApiServer) handlers() {
	http.HandleFunc(server.cfg.GetApiPrefix()+
		server.cfg.GetConnectorPref()+
		"/updateProject",
		server.updateProject)

	http.HandleFunc(server.cfg.GetApiPrefix()+
		server.cfg.GetConnectorPref()+
		"/projects",
		server.project)
}
