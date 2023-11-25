package endpoints

import (
	"Jira-analyzer/analyzer/models"
	"Jira-analyzer/common/logger"

	_ "github.com/lib/pq"
)

func (resourceHandler *ResourceHandler) GetIssueInfo(id int) (models.IssueInfo, error) {
	issueInfo := models.IssueInfo{}
	var authorId, assigneeId int

	if err := resourceHandler.database.QueryRow("id, projectId, authorId, assigneeId, key, summary, description, type, priority, createdTime, closedTime, updatedTime, timeSpent FROM issues where id = ?", id).Scan(&issueInfo.Id, &issueInfo.Project.Id, &authorId, &assigneeId, &issueInfo.Key, &issueInfo.Summary, &issueInfo.Description, &issueInfo.Type, &issueInfo.Priority, &issueInfo.Status, &issueInfo.CreatedTime, &issueInfo.ClosedTime, &issueInfo.UpdatedTime, &issueInfo.TimeSpent); err != nil {
		return issueInfo, err
	}

	if err := resourceHandler.database.QueryRow("SELECT name FROM author where id = ?", authorId).Scan(&issueInfo.Author); err != nil {
		return issueInfo, err
	}
	if err := resourceHandler.database.QueryRow("SELECT id FROM author where name = ?", assigneeId).Scan(&issueInfo.Assignee); err != nil {
		return issueInfo, err
	}

	resourceHandler.logger.Log(logger.INFO, "GetIssueInfo successfully")
	return issueInfo, nil
}

func (resourceHandler *ResourceHandler) GetHistoryInfo(id int) (models.HistoryInfo, error) {
	historyInfo := models.HistoryInfo{}

	resourceHandler.logger.Log(logger.INFO, "GetHistoryInfo successfully")
	return historyInfo, nil
}

func (resourceHandler *ResourceHandler) GetProjectInfo(id int) (models.ProjectInfo, error) {
	projectInfo := models.ProjectInfo{}

	if err := resourceHandler.database.QueryRow("id, title FROM project where id = ?",
		id).Scan(&projectInfo.Id, &projectInfo.Title); err != nil {
		return projectInfo, err
	}

	resourceHandler.logger.Log(logger.INFO, "GetProjectInfo successfully")
	return projectInfo, nil
}

func (resourceHandler *ResourceHandler) InsertProject(projectInfo models.ProjectInfo) (int, error) {
	var projectId int

	if err := resourceHandler.database.QueryRow("INSERT INTO project (title) VALUES(?) RETURNING id",
		projectInfo.Title).Scan(&projectId); err != nil {
		return projectId, err
	}

	resourceHandler.logger.Log(logger.INFO, "InsertProject successfully")
	return projectId, nil
}

func (resourceHandler *ResourceHandler) InsertIssue(issueInfo models.IssueInfo) (int, error) {
	var issueId, authorId, assigneeId int

	if err := resourceHandler.database.QueryRow("SELECT id FROM author where name = ?", issueInfo.Author).Scan(&authorId); err != nil {
		return issueId, err
	}
	if err := resourceHandler.database.QueryRow("SELECT id FROM author where name = ?", issueInfo.Assignee).Scan(&assigneeId); err != nil {
		return issueId, err
	}

	if err := resourceHandler.database.QueryRow("INSERT INTO issues (projectId, authorId, assigneeId, key, summary, description, type, priority, status, createdTime, closedTime, updatedTime, timeSpent) values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)", issueInfo.Project.Id, authorId, assigneeId, issueInfo.Key, issueInfo.Summary, issueInfo.Description, issueInfo.Type, issueInfo.Priority, issueInfo.Status, issueInfo.CreatedTime, issueInfo.ClosedTime, issueInfo.UpdatedTime, issueInfo.TimeSpent).Scan(&issueId); err != nil {
		return issueId, err
	}

	resourceHandler.logger.Log(logger.INFO, "InsertIssue successfully")
	return issueId, nil
}
