package v1

import (
	"github.com/DinozvrrDan/jira-analyzer/backend/resource/config"
	"github.com/DinozvrrDan/jira-analyzer/backend/resource/internal/repository"
	"github.com/DinozvrrDan/jira-analyzer/backend/resource/pkg/logger"
	"github.com/gorilla/mux"
)

type Handler struct {
	resourceHandler *ResourceHandler
	cfg             *config.Config
}

func NewHandler(repositories *repository.Repositories, log *logger.Logger, cfg *config.Config) *Handler {
	return &Handler{
		resourceHandler: NewResourceHandler(repositories, log, cfg),
		cfg:             cfg,
	}
}

func (handler *Handler) GetHandler(router *mux.Router) {
	resourceRouter := router.PathPrefix(handler.cfg.Server.ResourceHTTP.ResourcePrefix).Subrouter()
	handler.resourceHandler.GetResourceHandler(resourceRouter)
}
