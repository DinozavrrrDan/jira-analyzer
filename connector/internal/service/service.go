package service

import (
	"github.com/DinozvrrDan/jira-analyzer/connector/config"
	"github.com/DinozvrrDan/jira-analyzer/connector/internal/models"
	"github.com/DinozvrrDan/jira-analyzer/connector/internal/repository"
	"github.com/magellon17/logger"
)

type Connector interface {
	GetProjectIssues(projectName string) ([]models.Issue, error)
	GetProjects(limit int, page int, search string) ([]models.Project, models.Page, error)
}

type Transformer interface {
	TransformData(issues []models.Issue) []models.TransformedIssue
}

type Services struct {
	Connector   Connector
	Transformer Transformer
}

type ServicesDependencies struct {
	JiraRepositoryUrl string
}

func NewServices(repositories *repository.Repositories, deps ServicesDependencies, log *logger.Logger, cfg *config.Config) *Services {
	return &Services{
		Connector:   NewConnectorService(repositories.ConnectorRepository, deps.JiraRepositoryUrl, log, cfg),
		Transformer: NewTransformerService(),
	}
}
