package service

import (
	"github.com/DinozvrrDan/jira-analyzer/connector/config"
	"github.com/DinozvrrDan/jira-analyzer/connector/internal/models"
	"github.com/DinozvrrDan/jira-analyzer/connector/pkg/logger"
)

type Connector interface {
	GetProjectIssues(projectName string) ([]models.Issue, error)
	GetProjects(limit int, page int, search string) ([]models.Project, models.Page, error)
}

type Transformer interface {
	TransformData(issues []models.Issue) []models.TransformedIssue
}

type DatabasePusher interface {
	PushIssue(issues []models.TransformedIssue)
}

type Services struct {
	Connector      Connector
	Transformer    Transformer
	DatabasePusher DatabasePusher
}

type ServicesDependencies struct {
	JiraRepositoryUrl string
}

func NewServices(deps ServicesDependencies, log *logger.Logger, cfg *config.Config) *Services {
	return &Services{
		Connector:      NewConnectorService(deps.JiraRepositoryUrl, log, cfg),
		Transformer:    NewTransformerService(),
		DatabasePusher: NewDatabasePusher(log, cfg),
	}
}
