package endpoints

import (
	"Jira-analyzer/common/logger"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func (server *ResourceHandler) StartServer() error {
	server.logger.Log(logger.INFO, "Server start server...")

	router := mux.NewRouter()

	server.handlers(router)
	err := http.ListenAndServe(server.configReader.GetResourceHost()+":"+server.configReader.GetResourcePort(), router)
	if err != nil {
		server.logger.Log(logger.ERROR, fmt.Sprintf("StartServer: %w", err))
		return fmt.Errorf("StartServer: %w", err)
	}

	return nil
}
