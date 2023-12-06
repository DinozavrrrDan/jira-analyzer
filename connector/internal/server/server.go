package server

import (
	"github.com/DinozvrrDan/jira-analyzer/connector/config"
	"github.com/DinozvrrDan/jira-analyzer/connector/internal/service"
	"github.com/DinozvrrDan/jira-analyzer/connector/pkg/logger"
	"github.com/gorilla/mux"
	"net/http"
)

type ApiServer struct {
	cfg     *config.Config
	log     *logger.Logger
	handler *Handler
}

func NewApiServer(services *service.Services,
	log *logger.Logger, cfg *config.Config) *ApiServer {
	return &ApiServer{
		cfg:     cfg,
		log:     log,
		handler: NewHandler(services, log),
	}
}

func (server *ApiServer) StartServer() {
	server.log.Log(logger.INFO, "Server start server...")
	router := mux.NewRouter()

	server.handlers(router)
	err := http.ListenAndServe(server.cfg.Server.ConnectorHTTP.ConnectorHost+":"+server.cfg.Server.ConnectorHTTP.ConnectorPort, router)

	if err != nil {
		server.log.Log(logger.ERROR, "error while start a server")
	}
}

func (server *ApiServer) StopServer() {

}

func (server *ApiServer) handlers(router *mux.Router) {
	router.HandleFunc(server.cfg.Server.ApiServer.ApiPrefix+
		server.cfg.Server.ConnectorHTTP.ConnectorPrefix+
		"/updateProject",
		server.handler.UpdateProject).Methods("POST")
	router.HandleFunc(server.cfg.Server.ApiServer.ApiPrefix+
		server.cfg.Server.ConnectorHTTP.ConnectorPrefix+
		"/projects",
		server.handler.GetProjects).Methods("GET")
}
