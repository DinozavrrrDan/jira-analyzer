package analytic

import (
	"Jira-analyzer/common/configReader"
	"Jira-analyzer/common/logger"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

type AnalyticServer struct {
	configReader *configReader.ConfigReader
	logger       *logger.JiraLogger
	database     *sql.DB
}

func CreateNewResourceServer() *AnalyticServer {
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
	return &AnalyticServer{
		configReader: newReader,
		logger:       logger.CreateNewLogger(),
		database:     newDatabase,
	}
}

func (analyticServer *AnalyticServer) getGraph(responseWriter http.ResponseWriter, request *http.Request) {
	projectId, err := strconv.Atoi(request.URL.Query().Get("project"))
	if err != nil {
		analyticServer.logger.Log(logger.ERROR, err.Error())
		responseWriter.WriteHeader(http.StatusBadRequest)
		return
	}

	data, err := json.MarshalIndent(analyticServer.GraphFive(projectId), "", "\t")

	_, err = responseWriter.Write(data)
	responseWriter.WriteHeader(http.StatusOK)
}

func (analyticServer *AnalyticServer) handlers(router *mux.Router) {
	router.HandleFunc(analyticServer.configReader.GetAnalyticPref()+"{group:[1-2]}", analyticServer.getGraph).
		Queries("project", "{projectName}").
		Methods("GET")
}
