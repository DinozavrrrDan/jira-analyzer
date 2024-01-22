package repository

import "github.com/DinozvrrDan/jira-analyzer/connector/internal/models"

type IConnectorRepository interface {
	PushIssues(issues []models.TransformedIssue) error
}
