package endpoints

import (
	"Jira-analyzer/common/logger"
	"net/http"

	"github.com/gorilla/mux"
)

func (server *ResourceHandler) StartServer() {
	server.logger.Log(logger.INFO, "Server start server...")

	router := mux.NewRouter()

	server.handlers(router)
	//fmt.Print("Host" + server.configReader.GetResourceHost())
	//fmt.Print("Pr" + server.configReader.GetResourcePrefix())
	err := http.ListenAndServe(server.configReader.GetResourceHost()+":"+server.configReader.GetResourcePort(), router)
	if err != nil {
		server.logger.Log(logger.ERROR, "Error while start a server")
	}
}
