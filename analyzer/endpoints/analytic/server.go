package analytic

import (
	"Jira-analyzer/common/configReader"
	"Jira-analyzer/common/logger"
	"database/sql"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type AnalyticServer struct {
	configReader *configReader.ConfigReader
	logger       *logger.Logger
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
	vars := mux.Vars(request)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		analyticServer.logger.Log(logger.ERROR, err.Error())
		responseWriter.WriteHeader(http.StatusBadRequest)
		return
	}
	projectId, err := strconv.Atoi(request.URL.Query().Get("project"))
	if err != nil {
		analyticServer.logger.Log(logger.ERROR, err.Error())
		responseWriter.WriteHeader(http.StatusBadRequest)
		return
	}

	if id == 1 {
		//Функции с графиками
	} else if id == 2 {
		//Функции с графиками
	} else {
		analyticServer.logger.Log(logger.ERROR, err.Error())
		responseWriter.WriteHeader(http.StatusForbidden)
		return
	}

	_, err = responseWriter.Write( /*Какая-то дата*/ )
	responseWriter.WriteHeader(http.StatusOK)
}

func (server *AnalyticServer) StartServer() {
	server.logger.Log(logger.INFO, "Server start server...")

	router := mux.NewRouter()

	server.handlers(router)

	err := http.ListenAndServe(server.configReader.GetAnalyticHost()+":"+server.configReader.GetAnalyticHost(), router)
	if err != nil {
		server.logger.Log(logger.ERROR, "Error while start a server")
	}
}

func (server *AnalyticServer) handlers(router *mux.Router) {
	router.HandleFunc(server.configReader.GetAnalyticPref()+"{group:[1-2]}", server.getGraph).
		Queries("project", "{projectName}").
		Methods("GET")
}
