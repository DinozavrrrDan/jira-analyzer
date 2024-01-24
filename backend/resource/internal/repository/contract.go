package repository

import "github.com/DinozvrrDan/jira-analyzer/backend/resource/internal/models"

type IResourceRepository interface {
	GetIssueInfo(id int) (models.IssueInfo, error)
	GetProjectInfo(title string) (models.ProjectInfo, error)
	InsertProject(projectInfo models.ProjectInfo) (int64, error)
	InsertIssue(issueInfo models.IssueInfo) (int64, error)
	DeleteProject(title string) error
<<<<<<< HEAD
=======
	DeleteAuthor(name string) error
	DeleteAssignee(name string) error
	DeleteIssue(issue models.IssueInfo) error
>>>>>>> 6a15cb1650a9c1e304607e8f9b48d77b20ebf674
	GetProjects() ([]models.Project, error)
}
