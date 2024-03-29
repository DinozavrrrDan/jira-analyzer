package repository

import "github.com/DinozvrrDan/jira-analyzer/backend/resource/internal/models"

type IResourceRepository interface {
	GetIssueInfo(id int) (models.IssueInfo, error)
	GetProjectInfo(title string) (models.ProjectInfo, error)
	InsertProject(projectInfo models.ProjectInfo) (int64, error)
	InsertIssue(issueInfo models.IssueInfo) (int64, error)
	DeleteProject(title string) error
	GetProjects() ([]models.Project, error)
}
