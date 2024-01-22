package repository

import (
	"github.com/DinozvrrDan/jira-analyzer/backend/resource/internal/models"
)

type IResourceRepository interface {
	GetIssueInfo(id int) (models.IssueInfo, error)
	GetHistoryInfo(id int) ([]models.HistoryInfo, error)
	GetProjectInfo(id int) (models.ProjectInfo, error)
	InsertProject(projectInfo models.ProjectInfo) (int, error)
	InsertIssue(issueInfo models.IssueInfo) (int, error)
	InsertHistory(historyInfo models.HistoryInfo) (int, error)
}
