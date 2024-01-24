package v1

import (
	"encoding/json"
	"fmt"
	"github.com/DinozvrrDan/jira-analyzer/backend/analytics/config"
	"github.com/DinozvrrDan/jira-analyzer/backend/analytics/internal/models"
	repository "github.com/DinozvrrDan/jira-analyzer/backend/analytics/internal/repository"
	"github.com/DinozvrrDan/jira-analyzer/backend/analytics/pkg/logger"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
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
	router.HandleFunc("/make/{group:[1-6]}", handler.makeGraph).Methods(http.MethodGet)
	router.HandleFunc("/delete", handler.deleteGraphs).Methods(http.MethodGet)
	router.HandleFunc("/compare/{group:[1-6]}", handler.getGraph1).Methods(http.MethodGet)
}

func (handler *AnalyticsHandler) getGraph(writer http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	group, err := strconv.Atoi(vars["group"])
	if err != nil {
		errorWriter(writer, handler.log, "error: invalid group request.", http.StatusBadRequest)
		return
	}
	projectName := request.URL.Query()["project"]
	if len(projectName) == 0 {
		errorWriter(writer, handler.log, "error: no projects in request.", http.StatusBadRequest)
		return
	}

	if group == 1 {
		fmt.Println("GROUP 1")
		handler.getGraph1(writer, request)
	} else if group == 2 {
		fmt.Println("GROUP 2")
	} else if group == 4 {
		handler.getGraph4(writer, request)
	} else {
		errorWriter(writer, handler.log, "error: invalid group.", http.StatusBadRequest)
		return
	}
}
func (handler *AnalyticsHandler) makeGraph(writer http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	group, err := strconv.Atoi(vars["group"])
	if err != nil {
		errorWriter(writer, handler.log, "error: invalid group request.", http.StatusBadRequest)
		return
	}
	projectName := request.URL.Query()["project"]
	if len(projectName) == 0 {
		errorWriter(writer, handler.log, "error: no projects in request.", http.StatusBadRequest)
		return
	}

	if group == 1 {
		fmt.Println("GROUP 1")
	} else if group == 2 {
		fmt.Println("GROUP 2")
	} else if group == 4 {
		handler.getGraph4(writer, request)
	} else {
		errorWriter(writer, handler.log, "error: invalid group.", http.StatusBadRequest)
		return
	}
}
func (handler *AnalyticsHandler) deleteGraphs(writer http.ResponseWriter, request *http.Request) {
	projectName := request.URL.Query()["project"]
	if len(projectName) == 0 {
		errorWriter(writer, handler.log, "error: no projects in request.", http.StatusBadRequest)
		return
	}
	handler.getGraph4(writer, request)
}
func (handler *AnalyticsHandler) compareGraph(writer http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	group, err := strconv.Atoi(vars["group"])
	if err != nil {
		errorWriter(writer, handler.log, "error: invalid group request.", http.StatusBadRequest)
		return
	}
	projectName := request.URL.Query()["project"]
	if len(projectName) == 0 {
		errorWriter(writer, handler.log, "error: no projects in request.", http.StatusBadRequest)
		return
	}

	if group == 1 {
		fmt.Println("GROUP 1")
	} else if group == 2 {
		fmt.Println("GROUP 2")
	} else if group == 4 {
		handler.getGraph4(writer, request)
	} else {
		errorWriter(writer, handler.log, "error: invalid group.", http.StatusBadRequest)
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
			Issues:    models.Link{Href: "/api/v1/issues"},
			Projects:  models.Link{Href: "/api/v1/projects"},
			Histories: models.Link{Href: "/api/v1/histories"},
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
	categories, err := handler.analyticsRep.GetGraphsOneCategories(projectName[0])
	if err != nil {
		errorWriter(writer, handler.log, err.Error(), http.StatusBadRequest)
		return
	}

	var projectResponse = models.ResponseStruct{
		Links: models.ListOfReferences{
			Issues:    models.Link{Href: "/api/v1/issues"},
			Projects:  models.Link{Href: "/api/v1/projects"},
			Histories: models.Link{Href: "/api/v1/histories"},
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

func (handler *AnalyticsHandler) getGraphOneCompare(writer http.ResponseWriter, request *http.Request) {
	fmt.Println("GET GRAPH")
	projectName := request.URL.Query()["project"]
	if len(projectName) == 0 {
		fmt.Println("lenth 0")
		errorWriter(writer, handler.log, "error: no projects in request.", http.StatusBadRequest)
		return
	}

	/*var graphs models.CompareGraphsOne
	for i := 0; i < len(projectName); i++ {
		graphData, err := handler.analyticsRep.GetGraphsOneData(projectName[i])
		graph := models.GraphOne{GraphOneData: graphData}
		if err != nil {
			errorWriter(writer, handler.log, err.Error(), http.StatusBadRequest)
			return
		}
		graphsOne := append(graphs, graph)
	}
	graphs, err := handler.analyticsRep.GetGraphsOneData(projectName[0])*/
	//
	//if err != nil {
	//	errorWriter(writer, handler.log, err.Error(), http.StatusBadRequest)
	//	return
	//}

	var projectResponse = models.ResponseStruct{
		Links: models.ListOfReferences{
			Issues:    models.Link{Href: "/api/v1/issues"},
			Projects:  models.Link{Href: "/api/v1/projects"},
			Histories: models.Link{Href: "/api/v1/histories"},
		},
		//Info:    graphs,
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
