package endpoints

import (
	"Jira-analyzer/jiraConnector/logger"

	_ "github.com/lib/pq"
)

func (resourseHandler *ResourseHandler) GetIssueInfo(id int) (IssueInfo, error) {
	issueInfo := IssueInfo{}
	var authorId, assigneeId int

	if err := resourseHandler.database.QueryRow("id, projectId, authorId, assigneeId, key, summary, description, type, priority, createdTime, closedTime, updatedTime, timeSpent FROM issues where id = ?", id).Scan(&issueInfo.Id, &issueInfo.Project.Id, &authorId, &assigneeId, &issueInfo.Key, &issueInfo.Summary, &issueInfo.Description, &issueInfo.Type, &issueInfo.Priority, &issueInfo.Status, &issueInfo.CreatedTime, &issueInfo.ClosedTime, &issueInfo.UpdatedTime, &issueInfo.TimeSpent); err != nil {
		return issueInfo, err
	}

	//уведомить об успехе еще можно
	return issueInfo, nil
}

func (resourseHandler *ResourseHandler) GetHistoryInfo(id int) (HistoryInfo, error) {
	historyInfo := HistoryInfo{}

	//уведомить об успехе еще можно
	return historyInfo, nil
}

func (resourseHandler *ResourseHandler) GetProjectInfo(id int) (ProjectInfo, error) {
	projectInfo := ProjectInfo{}

	if err := resourseHandler.database.QueryRow("id, title FROM project where id = ?",
		id).Scan(&projectInfo.Id, &projectInfo.Title); err != nil {
		return projectInfo, err
	}

	resourseHandler.logger.Log(logger.INFO, "GetProjectInfoByID successfully")
	return projectInfo, nil
}
