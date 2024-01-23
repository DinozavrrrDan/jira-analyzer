package repository

import "github.com/DinozvrrDan/jira-analyzer/backend/resource/internal/models"

type IResourceRepository interface {
	GetIssueInfo(id int) (models.IssueInfo, error)
	GetProjectInfo(id int) (models.ProjectInfo, error)
	InsertProject(projectInfo models.ProjectInfo) (int64, error)
	InsertIssue(issueInfo models.IssueInfo) (int64, error)
	DeleteProject(title string) error
	DeleteAuthor(name string) error
	DeleteAssignee(name string) error
	DeleteIssue(issue models.IssueInfo) error
}
