package v1

import (
	"encoding/json"
	"fmt"
	"github.com/DinozvrrDan/jira-analyzer/backend/analytics/config"
	"github.com/DinozvrrDan/jira-analyzer/backend/analytics/internal/models"
	repository "github.com/DinozvrrDan/jira-analyzer/backend/analytics/internal/repository"
	"github.com/gorilla/mux"
	"github.com/magellon17/logger"
	"net/http"
	"strconv"
	"strings"
)

type AnalyticsHandler struct {
	analyticsRep repository.IAnalyticsRepository
	log          *logger.Logger
	cfg          *config.Config
}

func NewAnalyticsHandler(repositories *repository.Repositories, log *logger.Logger, cfg *config.Config) *AnalyticsHandler {
	return &AnalyticsHandler{
		log:          log,
		analyticsRep: repositories.AnalyticsRepository,
		cfg:          cfg,
	}
}

func (handler *AnalyticsHandler) GetAnalyticsHandler(router *mux.Router) {
	router.HandleFunc("/get/{group:[1-6]}", handler.getGraph).Methods(http.MethodGet)
	router.HandleFunc("/compare/{group:[1-6]}", handler.compareGraph).Methods(http.MethodGet)
}

func (handler *AnalyticsHandler) getGraph(writer http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	group, err := strconv.Atoi(vars["group"])
	if err != nil {
		errorWriter(writer, handler.log, "error: invalid group request.", http.StatusBadRequest)
		return
	}

	if group == 1 {
		handler.getGraph1(writer, request)
	} else if group == 4 {
		handler.getGraph4(writer, request)
	} else if group == 5 {
		handler.getGraph5(writer, request)
	} else if group == 6 {
		handler.getGraph6(writer, request)
	} else {
		errorWriter(writer, handler.log, "error: invalid group number.", http.StatusBadRequest)
		return
	}
}

func (handler *AnalyticsHandler) compareGraph(writer http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	group, err := strconv.Atoi(vars["group"])

	if err != nil {
		errorWriter(writer, handler.log, "error: invalid group request.", http.StatusBadRequest)
		return
	}
	if group != 1 {
		errorWriter(writer, handler.log, "error: invalid group.", http.StatusBadRequest)
		return
	}
	projectName := request.URL.Query()["project"]
	if len(projectName) == 0 {
		errorWriter(writer, handler.log, "error: no projects in request.", http.StatusBadRequest)
		return
	}

	names := strings.Split(projectName[0], ",")

	var compareGraphs []models.GraphOne
	categories := []string{"1 hour", "3 hours", "6 hours", "12 hour", "24 hour",
		"3 days", "7 days", "14 days", "more 30 days"}

	for i := 0; i < len(names); i++ {
		graphData, err := handler.analyticsRep.GetGraphsOneData(names[i])

		if err != nil {
			errorWriter(writer, handler.log, err.Error(), http.StatusBadRequest)
			return
		}

		compareGraphs = append(compareGraphs,
			models.GraphOne{
				GraphOneData: graphData,
				Categories:   categories,
			})
	}

	var projectResponse = models.ResponseStruct{
		Links: models.ListOfReferences{
			Issues:   models.Link{Href: "/api/v1/resource/issues"},
			Projects: models.Link{Href: "/api/v1/resource/project"},
			Graphs:   models.Link{Href: "/api/v1/graph"},
			Self:     models.Link{Href: fmt.Sprintf("/api/v1/graph/compare/%d", group)},
		},
		Info: models.CompareGraphsOne{
			GraphsOne:  compareGraphs,
			Counter:    len(names),
			Categories: categories,
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
}

func (handler *AnalyticsHandler) getGraph4(writer http.ResponseWriter, request *http.Request) {
	projectName := request.URL.Query()["project"]
	if len(projectName) == 0 {
		errorWriter(writer, handler.log, "error: no projects in request.", http.StatusBadRequest)
		return
	}

	graphsData, err := handler.analyticsRep.GetGraphsFourData(projectName[0])
	categories, err := handler.analyticsRep.GetGraphsFourCategories(projectName[0])

	if err != nil {
		errorWriter(writer, handler.log, err.Error(), http.StatusBadRequest)
		return
	}

	var projectResponse = models.ResponseStruct{
		Links: models.ListOfReferences{
			Issues:   models.Link{Href: "/api/v1/resource/issues"},
			Projects: models.Link{Href: "/api/v1/resource/project"},
			Graphs:   models.Link{Href: "/api/v1/graph"},
			Self:     models.Link{Href: fmt.Sprintf("/api/v1/graph/%d", 4)},
		},
		Info: models.GraphFour{
			GraphFourData: graphsData,
			Categories:    categories,
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
}

func (handler *AnalyticsHandler) getGraph5(writer http.ResponseWriter, request *http.Request) {
	projectName := request.URL.Query()["project"]
	if len(projectName) == 0 {
		errorWriter(writer, handler.log, "error: no projects in request.", http.StatusBadRequest)
		return
	}

	graphsData, err := handler.analyticsRep.GetGraphsFiveData(projectName[0])

	if err != nil {
		errorWriter(writer, handler.log, err.Error(), http.StatusBadRequest)
		return
	}
	fmt.Println("graphData")
	categories, err := handler.analyticsRep.GetGraphsFiveCategories(projectName[0])

	if err != nil {
		errorWriter(writer, handler.log, err.Error(), http.StatusBadRequest)
		return
	}
	fmt.Println("categories")

	var projectResponse = models.ResponseStruct{
		Links: models.ListOfReferences{
			Issues:   models.Link{Href: "/api/v1/resource/issues"},
			Projects: models.Link{Href: "/api/v1/resource/project"},
			Graphs:   models.Link{Href: "/api/v1/graph"},
			Self:     models.Link{Href: fmt.Sprintf("/api/v1/graph/%d", 5)},
		},
		Info: models.GraphFive{
			GraphFiveData: graphsData,
			Categories:    categories,
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
}

func (handler *AnalyticsHandler) getGraph6(writer http.ResponseWriter, request *http.Request) {
	projectName := request.URL.Query()["project"]
	if len(projectName) == 0 {
		errorWriter(writer, handler.log, "error: no projects in request.", http.StatusBadRequest)
		return
	}

	graphsData, err := handler.analyticsRep.GetGraphsSixData(projectName[0])

	if err != nil {
		errorWriter(writer, handler.log, err.Error(), http.StatusBadRequest)
		return
	}
	fmt.Println("graphData")
	categories, err := handler.analyticsRep.GetGraphsSixCategories(projectName[0])

	if err != nil {
		errorWriter(writer, handler.log, err.Error(), http.StatusBadRequest)
		return
	}
	fmt.Println("categories")

	var projectResponse = models.ResponseStruct{
		Links: models.ListOfReferences{
			Issues:   models.Link{Href: "/api/v1/resource/issues"},
			Projects: models.Link{Href: "/api/v1/resource/project"},
			Graphs:   models.Link{Href: "/api/v1/graph"},
			Self:     models.Link{Href: fmt.Sprintf("/api/v1/graph/%d", 6)},
		},
		Info: models.GraphSix{
			GraphSixData: graphsData,
			Categories:   categories,
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
}

func (handler *AnalyticsHandler) getGraph1(writer http.ResponseWriter, request *http.Request) {

	projectName := request.URL.Query()["project"]
	if len(projectName) == 0 {
		errorWriter(writer, handler.log, "error: no projects in request.", http.StatusBadRequest)
		return
	}

	graphData, err := handler.analyticsRep.GetGraphsOneData(projectName[0])

	if err != nil {
		errorWriter(writer, handler.log, err.Error(), http.StatusBadRequest)
		return
	}

	categories := []string{"1 hour", "3 hours", "6 hours", "12 hour", "24 hour", "3 days", "7 days", "14 days", "more 30 days"}

	var projectResponse = models.ResponseStruct{
		Links: models.ListOfReferences{
			Issues:   models.Link{Href: "/api/v1/resource/issues"},
			Projects: models.Link{Href: "/api/v1/resource/project"},
			Graphs:   models.Link{Href: "/api/v1/graph"},
			Self:     models.Link{Href: fmt.Sprintf("/api/v1/graph/%d", 1)},
		},
		Info: models.GraphOne{
			GraphOneData: graphData,
			Categories:   categories,
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
}
