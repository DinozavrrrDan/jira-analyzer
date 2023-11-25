package endpoints

import (
	"Jira-analyzer/common/logger"
	"Jira-analyzer/jiraConnector/logger"

	_ "github.com/lib/pq"
)

func (resourseHandler *ResourseHandler) GetIssueInfo(id int) IssueInfo {
	issueInfo := IssueInfo{}
	var authorId, assigneeId int

	if err := resourseHandler.database.QueryRow("id, projectId, authorId, assigneeId, key, summary, description, type, priority, createdTime, closedTime, updatedTime, timeSpent FROM issues where id = ?", id).Scan(&issueInfo.Id, &issueInfo.Project.Id, &authorId, &assigneeId, &issueInfo.Key, &issueInfo.Summary, &issueInfo.Description, &issueInfo.Type, &issueInfo.Priority, &issueInfo.Status, &issueInfo.CreatedTime, &issueInfo.ClosedTime, &issueInfo.UpdatedTime, &issueInfo.TimeSpent); err != nil {
		resourseHandler.logger.Log(logger.ERROR, err.Error())
		return issueInfo
	}

	//еще достать автора и assignee

	resourseHandler.logger.Log(logger.INFO, "GetIssueInfo successfully")
	return issueInfo
}

func (resourseHandler *ResourseHandler) GetHistoryInfo(id int) HistoryInfo {
	historyInfo := HistoryInfo{}

	resourseHandler.logger.Log(logger.INFO, "GetHistoryInfo successfully")
	return historyInfo
}

func (resourseHandler *ResourseHandler) GetProjectInfo(id int) ProjectInfo {
	projectInfo := ProjectInfo{}

	if err := resourseHandler.database.QueryRow("id, title FROM project where id = ?",
		id).Scan(&projectInfo.Id, &projectInfo.Title); err != nil {
		resourseHandler.logger.Log(logger.ERROR, err.Error())
		return projectInfo
	}

	resourseHandler.logger.Log(logger.INFO, "GetProjectInfo successfully")
	return projectInfo
}

func (resourseHandler *ResourseHandler) InsertProject(projectInfo ProjectInfo) int {
	var projectId int

	if err := resourseHandler.database.QueryRow("INSERT INTO project (title) VALUES(?) RETURNING id",
		projectInfo.Title).Scan(&projectId); err != nil {
		resourseHandler.logger.Log(logger.ERROR, err.Error())
		return projectId
	}

	resourseHandler.logger.Log(logger.INFO, "InsertProject successfully")
	return projectId
}

func (resourseHandler *ResourseHandler) InsertIssue(issueInfo IssueInfo) int {
	var issueId, authorId, assigneeId int

	//достать authorId и assigneeId из таблицы author

	if err := resourseHandler.database.QueryRow("INSERT INTO issues (projectId, authorId, assigneeId, key, summary, description, type, priority, status, createdTime, closedTime, updatedTime, timeSpent) values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)", issueInfo.Project.Id, authorId, assigneeId, issueInfo.Key, issueInfo.Summary, issueInfo.Description, issueInfo.Type, issueInfo.Priority, issueInfo.Status, issueInfo.CreatedTime, issueInfo.ClosedTime, issueInfo.UpdatedTime, issueInfo.TimeSpent).Scan(&issueId); err != nil {
		resourseHandler.logger.Log(logger.ERROR, err.Error())
		return issueId
	}

	resourseHandler.logger.Log(logger.INFO, "InsertIssue successfully")
	return issueId
}
