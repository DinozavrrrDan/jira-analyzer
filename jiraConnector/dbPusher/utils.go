package dbPusher

import (
	"Jira-analyzer/common/logger"
	"time"
)

func (databasePusher *DatabasePusher) insertInfoIntoIssues(projectId, authorId, assigneeId int, key, summary, description, Type, priority, status string, createdTime, closedTime, updatedTime time.Time, timeSpent int) {
	stmt, _ :=
		databasePusher.database.Prepare("INSERT INTO issues (projectId, authorId, assigneeId, key, summary, description, type, priority, status, createdTime, closedTime, updatedTime, timeSpent) values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")
	_, err := stmt.Exec(
		projectId,
		authorId,
		assigneeId,
		key,
		summary,
		description,
		Type,
		priority,
		status,
		createdTime,
		closedTime,
		updatedTime,
		timeSpent)

	if err != nil {
		databasePusher.logger.Log(logger.ERROR, err.Error())
	}
}

func (databasePusher *DatabasePusher) insertInfoIntoStatusChanges(issueIdId, authorId int, changeTime time.Time, fromStatus, toStatus string) {
	stmt, _ :=
		databasePusher.database.Prepare("INSERT INTO statuschanges (issueId, authorId, changeTime, fromStatus, toStatus) VALUES (?, ?, ?, ?, ?)")
	_, err := stmt.Exec(
		issueIdId,
		authorId,
		changeTime,
		fromStatus,
		toStatus)

	if err != nil {
		databasePusher.logger.Log(logger.ERROR, err.Error())
	}
}

// updateIssue обвновляет данные issue заданного key в таблицк issues
func (databasePusher *DatabasePusher) updateIssue(projectId, authorId, assigneeId int, key, summary, description, Type, priority, status string, createdTime, closedTime, updatedTime time.Time, timespent int) {
	stmt, _ :=
		databasePusher.database.Prepare("UPDATE issues set projectId = ?, authorId = ?, assigneeId = ? summary = ?, description = ?, type = ?, priority = ?, status = ?, createdTime = ?, closedtime = ?, updatedtime = ?, timespent = ? where key = ?")

	_, err := stmt.Exec(
		projectId,
		authorId,
		assigneeId,
		summary,
		description,
		Type,
		priority,
		status,
		createdTime,
		closedTime,
		updatedTime,
		timespent,
		key)

	if err != nil {
		databasePusher.logger.Log(logger.ERROR, err.Error())
	}
}

// getIssueId получает id по ключу задачи из таблицы issues
func (databasePusher *DatabasePusher) getIssueId(issueKey string) int {
	var issueID int
	err := databasePusher.database.QueryRow("SELECT id FROM issues where key = ?", issueKey).Scan(&issueID)

	if err != nil {
		databasePusher.logger.Log(logger.ERROR, err.Error())
	}

	return issueID
}

// getProjectId получает id по названию проекта из таблицы project
func (databasePusher *DatabasePusher) getProjectId(projectTitle string) int {
	var projectId int
	err := databasePusher.database.QueryRow("SELECT id FROM project where title = ?", projectTitle).Scan(&projectId)

	if err != nil {
		databasePusher.logger.Log(logger.ERROR, err.Error())
	}

	if projectId == 0 {
		err = databasePusher.database.QueryRow("INSERT INTO project (title) VALUES(?) RETURNING id", projectTitle).
			Scan(&projectId)
		if err != nil {
			databasePusher.logger.Log(logger.ERROR, err.Error())
		}
	}

	return projectId
}

// getAuthorId получает id по имени автора из таблицы author
func (databasePusher *DatabasePusher) getAuthorId(authorName string) int {
	var authorId int
	err := databasePusher.database.QueryRow("SELECT id FROM author where name = ?", authorName).
		Scan(&authorId)

	if err != nil {
		databasePusher.logger.Log(logger.ERROR, err.Error())
	}

	if authorId == 0 {
		err = databasePusher.database.QueryRow("INSERT INTO author (name) VALUES(?) RETURNING id", authorName).
			Scan(&authorId)

		if err != nil {
			databasePusher.logger.Log(logger.ERROR, err.Error())
		}
	}

	return authorId
}

// getAssigneeId получает id по имени assignee из таблицы author
func (databasePusher *DatabasePusher) getAssigneeId(assignee string) int {
	var assigneeId int
	err := databasePusher.database.QueryRow("SELECT id FROM author where name = ?", assignee).
		Scan(&assigneeId)

	if err != nil {
		databasePusher.logger.Log(logger.ERROR, err.Error())
	}

	if assigneeId == 0 {
		err = databasePusher.database.QueryRow("INSERT INTO author (name) VALUES(?) RETURNING id", assignee).Scan(&assigneeId)
		if err != nil {
			databasePusher.logger.Log(logger.ERROR, err.Error())
		}
	}

	return assigneeId
}

// checkIssueExists проверяет наличие issue заданного issueKey
func (databasePusher *DatabasePusher) checkIssueExists(issueKey string) bool {
	var issueId int
	err := databasePusher.database.QueryRow("SELECT id FROM issues where key = ?", issueKey).
		Scan(&issueId)

	if err != nil {
		databasePusher.logger.Log(logger.ERROR, err.Error())
	}

	return !(issueId == 0)
}
