package app

import (
	"database/sql"
	"fmt"
	"github.com/DinozvrrDan/jira-analyzer/backend/resource/config"
	resourceHandler "github.com/DinozvrrDan/jira-analyzer/backend/resource/internal/handler"
	"github.com/DinozvrrDan/jira-analyzer/backend/resource/internal/repository"
	"github.com/DinozvrrDan/jira-analyzer/backend/resource/pkg/logger"
	"github.com/gorilla/mux"
	"net/http"
)

type App struct {
	log *logger.Logger
	cfg *config.Config

	db     *sql.DB
	server *http.Server
}

func NewApp(cfg *config.Config, log *logger.Logger) (*App, error) {

	var db *sql.DB
	sqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.DB.HostDB,
		cfg.DB.PortDB,
		cfg.DB.UserDB,
		cfg.DB.PasswordDB,
		cfg.DB.NameDB)
	db, _ = sql.Open("postgres", sqlInfo)
	/*	if err != nil {
			return nil, fmt.Errorf(err.Error())
		}
	*/
	repositories := repository.NewRepositories(db)

	resourceHandlers := resourceHandler.NewHandler(repositories, log, cfg)

	resourceRouter := mux.NewRouter()

	resourceHandlers.GetRouter(resourceRouter)

	resourceServer := &http.Server{
		Addr:    cfg.Server.ResourceHTTP.ResourceHost + ":" + cfg.Server.ResourceHTTP.ResourcePort,
		Handler: resourceRouter,
	}

	return &App{
		log:    log,
		cfg:    cfg,
		db:     db,
		server: resourceServer,
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
