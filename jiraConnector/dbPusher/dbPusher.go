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
	_ "github.com/lib/pq"
	"io"
	"net/http"
	"strings"
	"time"
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

		issueId, err := databasePusher.getIssueId(issue.Key)

		if err != nil {
			databasePusher.logger.Log(logger.ERROR, fmt.Sprintf("ERROR: %v", err.Error()))
			return
		}

		requestString := databasePusher.configReader.GetJiraUrl() + "/rest/api/2/issue/" + issue.Key + "?expand=changelog"

		response, err := httpClient.Get(requestString)
		if err != nil {
			databasePusher.logger.Log(logger.ERROR, fmt.Sprintf("ERROR: %v", err.Error()))
			return
		}

		body, err := io.ReadAll(response.Body)

		if err != nil {
			databasePusher.logger.Log(logger.ERROR, fmt.Sprintf("ERROR: %v", err.Error()))
			return
		}

		var issueHistories models.IssueHistories
		err = json.Unmarshal(body, &issueHistories)

		if err != nil {
			databasePusher.logger.Log(logger.ERROR, fmt.Sprintf("ERROR: %v", err.Error()))
			return
		}

		for _, history := range issueHistories.Changelog.Histories {
			for _, statusChange := range history.StatusChanges {
				if strings.Compare(statusChange.Field, "status") == 0 {

					createdTime, _ := time.Parse("2006-01-02T15:04:05.999-0700", history.ChangeTime)

					if databasePusher.skipStatusChange(issueId, createdTime) {
						break
					}

					newAuthorId, _ := databasePusher.getAuthorId(history.Author.Name)

					err := databasePusher.insertInfoIntoStatusChanges(issueId, newAuthorId, createdTime, statusChange.FromStatus, statusChange.ToStatus)
					if err != nil {
						databasePusher.logger.Log(logger.ERROR, err.Error())
						return
					}
				}
			}
		}
	}
}

func (databasePusher *DatabasePusher) skipStatusChange(issueId int, createdTime time.Time) bool {
	var count int
	_ = databasePusher.database.QueryRow("SELECT COUNT(*) FROM statuschange WHERE issueid=$1 AND changetime=$2", issueId, createdTime).Scan(&count)
	return count != 0
}
