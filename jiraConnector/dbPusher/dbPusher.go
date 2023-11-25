package dbPusher

/*
Подключение к базе данных должно осуществляться через драйвер
database/sql
При загрузке данных из JIRA в БД должна поддерживаться
атомарность, то есть, если при скачивании части данных произошла
ошибка, то никакие данные не будут записаны в БД (все или ничего)
*/

import (
	"Jira-analyzer/common/configReader"
	"Jira-analyzer/common/logger"
	"Jira-analyzer/jiraConnector/models"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq"
)

type DatabasePusher struct {
	configReader *configReader.ConfigReader
	logger       *logger.JiraLogger
	database     *sql.DB
}

func CreateNewDatabasePusher() *DatabasePusher {
	newReader := configReader.CreateNewConfigReader()
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		newReader.GetHostDB(),
		newReader.GetPortDB(),
		newReader.GetUserDb(),
		newReader.GetPasswordDB(),
		newReader.GetDatabaseName())
	newDatabase, err := sql.Open("postgres", psqlInfo)

	if err != nil {
		panic(err)
	}
	return &DatabasePusher{
		configReader: newReader,
		logger:       logger.CreateNewLogger(),
		database:     newDatabase,
	}
}

func (databasePusher *DatabasePusher) PushIssue(issues []models.TransformedIssue) {
	for _, issue := range issues {
		projectId := databasePusher.getProjectId(issue.Project)
		authorId := databasePusher.getAuthorId(issue.Author)
		assigneeId := databasePusher.getAssigneeId(issue.Assignee)

		exists := databasePusher.checkIssueExists(issue.Key)
		if exists {
			databasePusher.updateIssue(projectId, authorId, assigneeId, issue.Key, issue.Summary, issue.Description, issue.Type, issue.Priority, issue.Status, issue.CreatedTime, issue.ClosedTime, issue.UpdatedTime, issue.Timespent)
		} else {
			databasePusher.insertInfoIntoIssues(projectId, authorId, assigneeId, issue.Key, issue.Summary, issue.Description, issue.Type, issue.Priority, issue.Status, issue.CreatedTime, issue.ClosedTime, issue.UpdatedTime, issue.Timespent)
		}
	}
}

func (databasePusher *DatabasePusher) insertInfoIntoIssues(projectId, authorId, assigneeId int, key, summary, description, Type, priority, status string, createdTime, closedTime, updatedTime time.Time, timeSpent int) {
	stmt, _ :=
		databasePusher.database.Prepare("INSERT INTO issues (projectId, authorId, assigneeId, key, summary, description, type, priority, status, createdTime, closedTime, updatedTime, timeSpent) values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")
	_, err := stmt.Exec(projectId, authorId, assigneeId, key, summary, description, Type, priority, status, createdTime, closedTime, updatedTime, timeSpent)
	if err != nil {
		databasePusher.logger.Log(logger.ERROR, err.Error())
	}
}

// updateIssue обвновляет данные issue заданного key в таблицк issues
func (databasePusher *DatabasePusher) updateIssue(projectId, authorId, assigneeId int, key, summary, description, Type, priority, status string, createdTime, closedTime, updatedTime time.Time, timespent int) {
	stmt, _ :=
		databasePusher.database.Prepare("UPDATE issues set projectId = ?, authorId = ?, assigneeId = ? summary = ?, description = ?, type = ?, priority = ?, status = ?, createdTime = ?, closedtime = ?, updatedtime = ?, timespent = ? where key = ?")
	_, err := stmt.Exec(projectId, authorId, assigneeId, summary, description, Type, priority, status, createdTime, closedTime, updatedTime, timespent, key)
	if err != nil {
		databasePusher.logger.Log(logger.ERROR, err.Error())
	}
}

// getProjectId получает id по названию проекта из таблицы project
func (databasePusher *DatabasePusher) getProjectId(projectTitle string) int {
	var projectId int
	err := databasePusher.database.QueryRow("SELECT id FROM project where title = ?", projectTitle).Scan(&projectId)
	if err != nil {
		databasePusher.logger.Log(logger.ERROR, err.Error())
	}

	if projectId == 0 {
		err = databasePusher.database.QueryRow("INSERT INTO project (title) VALUES(?) RETURNING id", projectTitle).Scan(&projectId)
		if err != nil {
			databasePusher.logger.Log(logger.ERROR, err.Error())
		}
	}
	return projectId
}

// getAuthorId получает id по имени автора из таблицы author
func (databasePusher *DatabasePusher) getAuthorId(authorName string) int {
	var authorId int
	err := databasePusher.database.QueryRow("SELECT id FROM author where name = ?", authorName).Scan(&authorId)
	if err != nil {
		databasePusher.logger.Log(logger.ERROR, err.Error())
	}

	if authorId == 0 {
		err = databasePusher.database.QueryRow("INSERT INTO author (name) VALUES(?) RETURNING id", authorName).Scan(&authorId)
		if err != nil {
			databasePusher.logger.Log(logger.ERROR, err.Error())
		}
	}
	return authorId
}

// getAssigneeId получает id по имени assignee из таблицы author
func (databasePusher *DatabasePusher) getAssigneeId(assignee string) int {
	var assigneeId int
	err := databasePusher.database.QueryRow("SELECT id FROM author where name = ?", assignee).Scan(&assigneeId)
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
	err := databasePusher.database.QueryRow("SELECT id FROM issues where key = ?", issueKey).Scan(&issueId)
	if err != nil {
		databasePusher.logger.Log(logger.ERROR, err.Error())
	}

	return !(issueId == 0)
}
