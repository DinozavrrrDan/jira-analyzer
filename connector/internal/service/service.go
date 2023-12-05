package service

import (
	"connector/config"
	"connector/internal/entities"
	"connector/pkg/logger"
)

type Connector interface {
	GetProjectIssues(projectName string) ([]entities.Issue, error)
	GetProjects(limit int, page int, search string) ([]entities.Project, entities.Page, error)
}

type Transformer interface {
	TransformData(issues []entities.Issue) []entities.TransformedIssue
}

type DatabasePusher interface {
	PushIssue(issues []entities.TransformedIssue)
}

type Services struct {
	Connector      Connector
	Transformer    Transformer
	DatabasePusher DatabasePusher
}

type ServicesDependencies struct {
	JiraRepositoryUrl string
}

func NewServices(deps ServicesDependencies, log *logger.Logger, cfg *config.Reader) *Services {
	return &Services{
		Connector:      NewConnectorService(deps.JiraRepositoryUrl, log, cfg),
		Transformer:    NewTransformerService(),
		DatabasePusher: NewDatabasePusher(log, cfg),
	}
}
