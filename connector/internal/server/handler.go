package server

import (
	"encoding/json"
	"fmt"
	"github.com/DinozvrrDan/jira-analyzer/connector/api"
	"github.com/DinozvrrDan/jira-analyzer/connector/internal/service"
	"github.com/DinozvrrDan/jira-analyzer/connector/pkg/logger"
	"net/http"
	"strconv"
)

type Handler struct {
	connectorSvc   service.Connector
	transformerSvc service.Transformer
	dbPusherSvc    service.DatabasePusher
	log            *logger.Logger
}

func NewHandler(services *service.Services, log *logger.Logger) *Handler {
	return &Handler{
		log:            log,
		connectorSvc:   services.Connector,
		transformerSvc: services.Transformer,
		dbPusherSvc:    services.DatabasePusher,
	}
}

func (handler *Handler) UpdateProject(writer http.ResponseWriter, request *http.Request) {

	projectName := request.URL.Query().Get("getProjects")

	if len(projectName) == 0 {
		errorWriter(writer, handler.log, "error: no projects in request.", http.StatusBadRequest)
		return
	}

	issues, err := handler.connectorSvc.GetProjectIssues(projectName)
	response, err := json.MarshalIndent(issues, "", "\t")

	writer.Write(response)

	if err != nil {
		errorWriter(writer, handler.log, err.Error(), http.StatusBadRequest)
		return
	}

	transformedIssues := handler.transformerSvc.TransformData(issues)
	handler.dbPusherSvc.PushIssue(transformedIssues)
}

func (handler *Handler) GetProjects(writer http.ResponseWriter, request *http.Request) {

	limit, page, search, err := getProjectParametersFromRequest(request)
	if err != nil {
		handler.log.Log(logger.WARNING, "error: Incorrect getProjects parameter. Set default value.")
	}

	writer.Header().Set("Content-Type", "application/json")

	projects, pages, err := handler.connectorSvc.GetProjects(limit, page, search)

	if err != nil {
		errorWriter(writer, handler.log, err.Error(), http.StatusBadRequest)
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
		errorWriter(writer, handler.log, err.Error(), http.StatusBadRequest)
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
