package app

import (
	"connector/config"
	"connector/internal/server"
	"connector/internal/service"
	"connector/pkg/logger"
)

type App struct {
	services *service.Services
	log      *logger.Logger
	cfg      *config.Config
}

func NewApp(cfg *config.Config, log *logger.Logger) *App {
	deps := service.ServicesDependencies{
		JiraRepositoryUrl: cfg.JiraUrl,
	}

	services := service.NewServices(deps, log, cfg)

	return &App{services: services, log: log, cfg: cfg}
}

func (app *App) Run() {
	apiServer := server.NewApiServer(app.services, app.log, app.cfg)
	apiServer.StartServer()
}

func (app *App) Close() {

}
