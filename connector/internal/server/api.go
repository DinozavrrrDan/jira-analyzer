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

	"github.com/gorilla/mux"
)

type ApiServer struct {
	connectorSvc   service.Connector
	transformerSvc service.Transformer
	dbPusherSvc    service.DatabasePusher
	cfg            *config.Config
	log            *logger.Logger
}

func NewApiServer(services *service.Services,
	log *logger.Logger, cfg *config.Config) *ApiServer {
	return &ApiServer{
		cfg:            cfg,
		log:            log,
		connectorSvc:   services.Connector,
		transformerSvc: services.Transformer,
		dbPusherSvc:    services.DatabasePusher,
	}
}

func (server *ApiServer) updateProject(writer http.ResponseWriter, request *http.Request) {

	projectName := request.URL.Query().Get("project")

	if len(projectName) == 0 {
		errorWriter(writer, server, "error: no projects in request.", http.StatusBadRequest)
		return
	}

	issues, err := server.connectorSvc.GetProjectIssues(projectName)
	response, err := json.MarshalIndent(issues, "", "\t")

	writer.Write(response)

	if err != nil {
		errorWriter(writer, server, err.Error(), http.StatusBadRequest)
		return
	}

	transformedIssues := server.transformerSvc.TransformData(issues)
	server.dbPusherSvc.PushIssue(transformedIssues)
}

func (server *ApiServer) project(writer http.ResponseWriter, request *http.Request) {

	limit, page, search, err := getProjectParametersFromRequest(request)
	if err != nil {
		server.log.Log(logger.WARNING, "error: Incorrect project parameter. Set default value.")
	}

	writer.Header().Set("Content-Type", "application/json")

	projects, pages, err := server.connectorSvc.GetProjects(limit, page, search)

	if err != nil {
		errorWriter(writer, server, err.Error(), http.StatusBadRequest)
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
	_, err = writer.Write(response)

	if err != nil {
		errorWriter(writer, server, err.Error(), http.StatusBadRequest)
		return
	}
}

func getProjectParametersFromRequest(request *http.Request) (int, int, string, error) {
	defaultLimit := 20
	defaultPage := 1
	defaultSearch := ""

	var err error
	urlQuery := request.URL.Query()

	if limit, ok := urlQuery["limit"]; ok {
		defaultLimit, err = strconv.Atoi(limit[0])
		if err != nil {
			return defaultLimit, defaultPage, defaultSearch, err
		}
	}

	if page, ok := urlQuery["page"]; ok {
		defaultPage, err = strconv.Atoi(page[0])
		if err != nil {
			return defaultLimit, defaultPage, defaultSearch, err
		}
	}

	if search, ok := urlQuery["limit"]; ok {
		defaultSearch = search[0]
	}

	return defaultLimit, defaultPage, defaultSearch, nil
}

func errorWriter(w http.ResponseWriter, server *ApiServer, message string, status int) {
	server.log.Log(logger.ERROR, message)
	w.WriteHeader(status)
}

func (server *ApiServer) StartServer() {
	server.log.Log(logger.INFO, "Server start server...")
	router := mux.NewRouter()

	server.handlers(router)
	err := http.ListenAndServe(server.cfg.ConnectorHost+":"+server.cfg.ConnectorPort, router)

	if err != nil {
		server.log.Log(logger.ERROR, "error while start a server")
	}
}

func (server *ApiServer) handlers(router *mux.Router) {
	router.HandleFunc(server.cfg.ApiPrefix+
		server.cfg.ConnectorPrefix+
		"/updateProject",
		server.updateProject).Methods("POST")
	router.HandleFunc(server.cfg.ApiPrefix+
		server.cfg.ConnectorPrefix+
		"/projects",
		server.project).Methods("GET")
}
