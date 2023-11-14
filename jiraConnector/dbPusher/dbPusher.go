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

	}
	return &DatabasePusher{
		configReader: newReader,
		logger:       logger.CreateNewLogger(),
		database:     newDatabase,
	}
}

func (databasePusher *DatabasePusher) PushIssue(issues []models.TransformedIssue) {
}
