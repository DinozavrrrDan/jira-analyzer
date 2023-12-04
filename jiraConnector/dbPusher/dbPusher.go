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
		projectId, err := databasePusher.getProjectId(issue.Project)

		if err != nil {
			databasePusher.logger.Log(logger.ERROR, err.Error())

			return
		}

		authorId, err := databasePusher.getAuthorId(issue.Author)

		if err != nil {
			databasePusher.logger.Log(logger.ERROR, err.Error())

			return
		}

		assigneeId, err := databasePusher.getAssigneeId(issue.Assignee)

		if err != nil {
			databasePusher.logger.Log(logger.ERROR, err.Error())

			return
		}

		exists := databasePusher.checkIssueExists(issue.Key)
		if exists {
			err := databasePusher.updateIssue(
				projectId,
				authorId,
				assigneeId,
				issue.Key,
				issue.Summary,
				issue.Description,
				issue.Type,
				issue.Priority,
				issue.Status,
				issue.CreatedTime,
				issue.ClosedTime,
				issue.UpdatedTime,
				issue.Timespent)
			if err != nil {
				databasePusher.logger.Log(logger.ERROR, err.Error())
				return
			}
		} else {
			err := databasePusher.insertInfoIntoIssues(
				projectId,
				authorId,
				assigneeId,
				issue.Key,
				issue.Summary,
				issue.Description,
				issue.Type,
				issue.Priority,
				issue.Status,
				issue.CreatedTime,
				issue.ClosedTime,
				issue.UpdatedTime,
				issue.Timespent)
			if err != nil {
				databasePusher.logger.Log(logger.ERROR, fmt.Sprintf("ERROR: %v", err.Error()))

				return
			}
		}
	}
}
