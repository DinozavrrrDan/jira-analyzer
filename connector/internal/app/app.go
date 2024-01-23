package app

import (
	"database/sql"
	"fmt"
	"github.com/DinozvrrDan/jira-analyzer/connector/config"
	"github.com/DinozvrrDan/jira-analyzer/connector/internal/handler"
	"github.com/DinozvrrDan/jira-analyzer/connector/internal/repository"
	"github.com/DinozvrrDan/jira-analyzer/connector/internal/service"
	"github.com/DinozvrrDan/jira-analyzer/connector/pkg/logger"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"net/http"
)

type App struct {
	log *logger.Logger
	cfg *config.Config

	db     *sql.DB
	server *http.Server
}

func NewApp(cfg *config.Config, log *logger.Logger) (*App, error) {
	deps := service.ServicesDependencies{
		JiraRepositoryUrl: cfg.Connector.JiraUrl,
	}

	var db *sql.DB
	sqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.DB.HostDB,
		cfg.DB.PortDB,
		cfg.DB.UserDB,
		cfg.DB.PasswordDB,
		cfg.DB.NameDB)
	db, err := sql.Open("postgres", sqlInfo)
	if err != nil {
		return nil, fmt.Errorf(err.Error())
	}

	repositories := repository.NewRepositories(db)

	services := service.NewServices(deps, log, cfg)

	handlers := handler.NewHandler(services, repositories, log, cfg)

	router := mux.NewRouter()
	handlers.GetRouter(router)

	server := &http.Server{
		Addr:    cfg.Server.ConnectorHTTP.ConnectorHost + ":" + cfg.Server.ConnectorHTTP.ConnectorPort,
		Handler: router,
	}

	return &App{
		log:    log,
		cfg:    cfg,
		db:     db,
		server: server,
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
	if err := app.server.Close(); err != nil {
		return fmt.Errorf(err.Error())
	}
	return app.db.Close()
}
