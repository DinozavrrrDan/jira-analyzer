package handler

import (
	"github.com/DinozvrrDan/jira-analyzer/backend/resource/config"
	"github.com/DinozvrrDan/jira-analyzer/backend/resource/internal/handler/http/v1"
	"github.com/DinozvrrDan/jira-analyzer/backend/resource/internal/repository"
	"github.com/gorilla/mux"
	"github.com/magellon17/logger"
)

type Handler struct {
	v1  *v1.Handler
	cfg *config.Config
}

func NewHandler(repositories *repository.Repositories, log *logger.Logger, cfg *config.Config) *Handler {
	return &Handler{
		v1:  v1.NewHandler(repositories, log, cfg),
		cfg: cfg,
	}
}

func (handler *Handler) GetRouter(router *mux.Router) {
	v1Router := router.PathPrefix(handler.cfg.Server.ApiServer.ApiPrefix).Subrouter()
	handler.v1.GetHandler(v1Router)
}
