package repository

import "github.com/DinozvrrDan/jira-analyzer/connector/internal/models"

type IConnectorRepository interface {
	PushIssue(issues []models.TransformedIssue) error
}
