package handler

import (
	"github.com/DinozvrrDan/jira-analyzer/connector/config"
	v1 "github.com/DinozvrrDan/jira-analyzer/connector/internal/handler/http/v1"
	"github.com/DinozvrrDan/jira-analyzer/connector/internal/repository"
	"github.com/DinozvrrDan/jira-analyzer/connector/internal/service"
	"github.com/DinozvrrDan/jira-analyzer/connector/pkg/logger"
	"github.com/gorilla/mux"
)

type Handler struct {
	v1  *v1.Handler
	cfg *config.Config
}

func NewHandler(services *service.Services, repositories *repository.Repositories, log *logger.Logger, cfg *config.Config) *Handler {
	return &Handler{
		v1:  v1.NewHandler(services, repositories, log, cfg),
		cfg: cfg,
	}
}

func (handler *Handler) GetRouter(router *mux.Router) {
	v1Router := router.PathPrefix(handler.cfg.Server.ApiServer.ApiPrefix).Subrouter()
	handler.v1.GetHandler(v1Router)
}
