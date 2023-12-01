package endpoints

import (
	"Jira-analyzer/analyzer/models"
	"Jira-analyzer/common/logger"

	_ "github.com/lib/pq"
)

func (resourceHandler *ResourceHandler) GetIssueInfo(id int) (models.IssueInfo, error) {
	issueInfo := models.IssueInfo{}

	var authorID, assigneeID int

	stmt, _ := resourceHandler.database.
		Prepare("SELECT id, projectId, authorId, assigneeId, key, summary, description, type, priority, createdTime, closedTime, updatedTime, timeSpent FROM issues where id = ?")

	err := stmt.QueryRow(id).Scan(
		&issueInfo.ID,
		&issueInfo.Project.ID,
		&authorID,
		&assigneeID,
		&issueInfo.Key,
		&issueInfo.Summary,
		&issueInfo.Description,
		&issueInfo.Type,
		&issueInfo.Priority,
		&issueInfo.Status,
		&issueInfo.CreatedTime,
		&issueInfo.ClosedTime,
		&issueInfo.UpdatedTime,
		&issueInfo.TimeSpent)

	if err != nil {
		return issueInfo, err
	}

	err = resourceHandler.database.QueryRow("SELECT name FROM author where id = ?", authorID).Scan(&issueInfo.Author)
	if err != nil {
		return issueInfo, err
	}

	err = resourceHandler.database.QueryRow("SELECT id FROM author where name = ?", assigneeID).Scan(&issueInfo.Assignee)
	if err != nil {
		return issueInfo, err
	}

	resourceHandler.logger.Log(logger.INFO, "GetIssueInfo successfully")

	return issueInfo, nil
}

func (resourceHandler *ResourceHandler) GetHistoryInfo(id int) ([]models.HistoryInfo, error) {
	var historyInfos []models.HistoryInfo

	rows, err := resourceHandler.database.Query(
		"SELECT "+
			"authorId,"+
			"changeTime"+
			"fromStatus,"+
			"toStatus "+
			"FROM statusChanges "+
			"WHERE issueId = ?", id,
	)

	if err != nil {
		resourceHandler.logger.Log(logger.ERROR, err.Error())

		return historyInfos, err
	}

	for rows.Next() {
		historyInfo := models.HistoryInfo{}
		err := rows.Scan(&historyInfo.AuthorID, &historyInfo.ChangeTime, &historyInfo.FromStatus, &historyInfo.ToStatus)

		if err != nil {
			resourceHandler.logger.Log(logger.ERROR, err.Error())

			return historyInfos, err
		}

		historyInfos = append(historyInfos, historyInfo)
	}

	resourceHandler.logger.Log(logger.INFO, "GetHistoryInfo successfully")

	return historyInfos, nil
}

func (resourceHandler *ResourceHandler) GetProjectInfo(id int) (models.ProjectInfo, error) {
	projectInfo := models.ProjectInfo{}

	stmt, _ := resourceHandler.database.Prepare("SELECT id, title FROM project WHERE id = ?")
	err := stmt.QueryRow(id).Scan(&projectInfo.ID, &projectInfo.Title)

	if err != nil {
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

	err := resourceHandler.database.QueryRow("SELECT id FROM author WHERE name = ?", issueInfo.Author).Scan(&authorId)
	if err != nil {
		return issueId, err
	}

	err = resourceHandler.database.QueryRow("SELECT id FROM author where name = ?", issueInfo.Assignee).Scan(&assigneeId)
	if err != nil {
		return issueId, err
	}

	stmt, _ := resourceHandler.database.Prepare("INSERT INTO issues (projectId, authorId, assigneeId, key, summary, description, type, priority, status, createdTime, closedTime, updatedTime, timeSpent) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")

	err = stmt.QueryRow(issueInfo.Project.ID,
		authorId,
		assigneeId,
		issueInfo.Key,
		issueInfo.Summary,
		issueInfo.Description,
		issueInfo.Type,
		issueInfo.Priority,
		issueInfo.Status,
		issueInfo.CreatedTime,
		issueInfo.ClosedTime,
		issueInfo.UpdatedTime,
		issueInfo.TimeSpent).Scan(&issueId)

	if err != nil {
		return issueId, err
	}

	resourceHandler.logger.Log(logger.INFO, "InsertIssue successfully")

	return issueId, nil
}

func (resourceHandler *ResourceHandler) InsertHistory(historyInfo models.HistoryInfo) (int, error) {
	var historyID int

	stmt, _ := resourceHandler.database.
		Prepare("INSERT INTO StatusChanges (issueId,authorId,changeTime,fromStatus,toStatus) VALUES (?, ?, now(), ?, ?)")

	err := stmt.QueryRow(historyInfo.IssueID, historyInfo.AuthorID, historyInfo.FromStatus, historyInfo.ToStatus).Err()

	if err != nil {
		resourceHandler.logger.Log(logger.ERROR, err.Error())

		return historyID, err
	}

	stmt, _ = resourceHandler.database.
		Prepare("UPDATE Issue SET status = ?, updatedTime = now(), timespent = now()-createdTime WHERE id = ?")

	err = stmt.QueryRow(historyInfo.ToStatus, historyInfo.IssueID).Err()

	if err != nil {
		resourceHandler.logger.Log(logger.ERROR, err.Error())

		return historyID, err
	}

	resourceHandler.logger.Log(logger.INFO, "InsertHistoryInfo successfully")

	return historyID, nil
}
