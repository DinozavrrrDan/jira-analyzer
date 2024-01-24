package app

import (
	"fmt"
	"github.com/DinozvrrDan/jira-analyzer/backend/gateway/config"
	"github.com/DinozvrrDan/jira-analyzer/backend/gateway/internal/handler"
	"github.com/DinozvrrDan/jira-analyzer/backend/gateway/pkg/logger"
	"net/http"
)

type App struct {
	log    *logger.Logger
	cfg    *config.Config
	server *http.Server
}

func NewApp(cfg *config.Config, log *logger.Logger) (*App, error) {

	gatewayMux := http.NewServeMux()
	handlers := handler.NewHandler(log, cfg)

	handlers.GetProxy(gatewayMux)

	gatewayServer := &http.Server{
		Addr:    cfg.Gateway.Host + ":" + cfg.Gateway.Port,
		Handler: gatewayMux,
	}

	return &App{
		log:    log,
		cfg:    cfg,
		server: gatewayServer,
	}, nil
}

func (app *App) Run() error {
	err := app.server.ListenAndServe()
	if err != nil {
		return fmt.Errorf(err.Error())
	}

	return nil
}

func (app *App) Close() error {
	return app.server.Close()
}
