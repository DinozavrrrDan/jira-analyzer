package analytic

import (
	"Jira-analyzer/common/logger"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func (analyticServer *AnalyticServer) StartServer() error {
	analyticServer.logger.Log(logger.INFO, "Server start server...")

	router := mux.NewRouter()

	analyticServer.handlers(router)

	err := http.ListenAndServe(analyticServer.configReader.GetAnalyticHost()+":"+analyticServer.configReader.GetAnalyticHost(), router)
	if err != nil {
		analyticServer.logger.Log(logger.ERROR, fmt.Sprintf("StartServer: %w", err))
		return fmt.Errorf("StartServer: %w", err)
	}

	return nil
}
