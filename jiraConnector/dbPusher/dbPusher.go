package dbPusher

/*
Подключение к базе данных должно осуществляться через драйвер
database/sql
При загрузке данных из JIRA в БД должна поддерживаться
атомарность, то есть, если при скачивании части данных произошла
ошибка, то никакие данные не будут записаны в БД (все или ничего)
*/

import (
	"Jira-analyzer/jiraConnector/configReader"
	"Jira-analyzer/jiraConnector/logger"
	"Jira-analyzer/jiraConnector/models"
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

type DatabasePusher struct {
	configReader *configReader.ConfigRaeder
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
	for i, issue := range issues {
		_, err := databasePusher.database.Exec("INSERT INTO Issue (projectid, authorid, assigneeid, key, summary, description, type, priority, status, createdtime, closedtime, updatedtime, timespent) values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)", i, i, i, issue.Assignee, issue.Key, issue.Summary, issue.Description, issue.Type, issue.Priority, issue.Status, issue.CreatedTime, issue.ClosedTime, issue.Timespent) //хз пока насчет id, пусть будет пока такой
		if err != nil {
			panic(err)
		}

		_, err = databasePusher.database.Exec("INSERT INTO Author (id, name) values ($1, $2)", i, issue.Author)
		if err != nil {
			panic(err)
		}

		_, err = databasePusher.database.Exec("INSERT INTO Projects (id, title) values ($1, $2)", i, issue.Project)
		if err != nil {
			panic(err)
		}
	}
}
