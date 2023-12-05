package analytic

import (
	"Jira-analyzer/common/logger"
	"github.com/gorilla/mux"
	"net/http"
)

func (analyticServer *AnalyticServer) StartServer() {
	analyticServer.logger.Log(logger.INFO, "Server start server...")

	router := mux.NewRouter()

	analyticServer.handlers(router)

	err := http.ListenAndServe(analyticServer.configReader.GetAnalyticHost()+":"+analyticServer.configReader.GetAnalyticHost(), router)
	if err != nil {
		analyticServer.logger.Log(logger.ERROR, "Error while start a server")
	}
}
