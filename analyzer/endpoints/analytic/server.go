package analytic

import (
	"Jira-analyzer/common/logger"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func (analyticServer *AnalyticServer) StartServer() {
	analyticServer.logger.Log(logger.INFO, "Server start server...")

	router := mux.NewRouter()

	analyticServer.handlers(router)

	err := http.ListenAndServe(analyticServer.configReader.GetAnalyticHost()+":"+analyticServer.configReader.GetAnalyticHost(), router)
	if err != nil {
		analyticServer.logger.Log(logger.ERROR, fmt.Sprintf("StartServer: %v", err))
		panic(fmt.Errorf("StartServer: %v", err))
	}
}
