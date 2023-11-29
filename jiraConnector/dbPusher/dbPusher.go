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
	"encoding/json"
	"fmt"
	"io"
	"net/http"

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
	httpClient := &http.Client{}

	for _, issue := range issues {
		projectId := databasePusher.getProjectId(issue.Project)
		authorId := databasePusher.getAuthorId(issue.Author)
		assigneeId := databasePusher.getAssigneeId(issue.Assignee)
		issueId := databasePusher.getIssueId(issue.Key)

		exists := databasePusher.checkIssueExists(issue.Key)
		if exists {
			databasePusher.updateIssue(
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
		} else {
			databasePusher.insertInfoIntoIssues(
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
		}

		requestString := databasePusher.configReader.GetJiraUrl() + "/rest/api/2/issue/" + issue.Key + "?expand=changelog"
		response, err := httpClient.Get(requestString)
		if err != nil {
			databasePusher.logger.Log(logger.ERROR, err.Error())
			return
		}

		body, err := io.ReadAll(response.Body)

		if err != nil {
			databasePusher.logger.Log(logger.ERROR, err.Error())
			return
		}

		var issueHistories models.IssueHistories
		err = json.Unmarshal(body, &issueHistories)

		if err != nil {
			databasePusher.logger.Log(logger.ERROR, err.Error())
			return
		}

		for _, history := range issueHistories.Changelog.Histories {
			for _, statusChange := range history.StatusChanges {
				changeTime := history.ChangeTime
				newAuthorId := databasePusher.getAuthorId(history.Author.Name)

				databasePusher.insertInfoIntoStatusChanges(issueId, newAuthorId, changeTime, statusChange.FromStatus, statusChange.ToStatus)
			}
		}
	}
}
