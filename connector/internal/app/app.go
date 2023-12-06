package app

import (
	"github.com/DinozvrrDan/jira-analyzer/connector/config"
	"github.com/DinozvrrDan/jira-analyzer/connector/internal/server"
	"github.com/DinozvrrDan/jira-analyzer/connector/internal/service"
	"github.com/DinozvrrDan/jira-analyzer/connector/pkg/logger"
)

type App struct {
	services *service.Services
	server   *server.ApiServer
	log      *logger.Logger
	cfg      *config.Config
}

func NewApp(cfg *config.Config, log *logger.Logger) *App {
	deps := service.ServicesDependencies{
		JiraRepositoryUrl: cfg.Connector.JiraUrl,
	}

	services := service.NewServices(deps, log, cfg)

	return &App{services: services,
		server: server.NewApiServer(services, log, cfg),
		log:    log,
		cfg:    cfg}
}

func (app *App) Run() {
	app.server.StartServer()
}

func (app *App) Close() {

}
