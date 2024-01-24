package v1

import (
	"github.com/DinozvrrDan/jira-analyzer/backend/analytics/config"
	repository "github.com/DinozvrrDan/jira-analyzer/backend/analytics/internal/repository"
	"github.com/DinozvrrDan/jira-analyzer/backend/analytics/pkg/logger"
	"github.com/gorilla/mux"
)

type AnalyticsHandler struct {
	analyticsRep repository.IAnalyticsRepository
	log          *logger.Logger
	cfg          *config.Config
}

func NewAnalyticsHandler(repositories *repository.Repositories, log *logger.Logger, cfg *config.Config) *AnalyticsHandler {
	return &AnalyticsHandler{
		log:          log,
		analyticsRep: repositories.AnalyticsRepository,
		cfg:          cfg,
	}
}

func (handler *AnalyticsHandler) GetAnalyticsHandler(router *mux.Router) {
	//	router.HandleFunc().Methods(http.MethodPost)
}
