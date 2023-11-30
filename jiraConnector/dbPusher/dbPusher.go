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
<<<<<<< HEAD
=======
	"io"
	"net/http"
>>>>>>> f088cd021ec45f2588741546a422d8f32f0e2981

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

func CheckIssueExists(db *sql.DB, table, column string, value string) (bool, string) {
	row := db.QueryRow("SELECT assigneeId FROM $1 where $2 = $3", table, column, value)
	err := row.Scan(&value)
	if err != nil {
		return false, value
	} else {
		return true, value
	}
}

func CountRows(db *sql.DB, table string) (int, error) {
	stmt, err := db.Prepare("SELECT COUNT(*) FROM ?")
	if err != nil {
		return 0, err
	}

	var count int
	err = stmt.QueryRow(table).Scan(&count)
	return count, err
}

func (databasePusher *DatabasePusher) PushIssue(issues []models.TransformedIssue) {
	httpClient := &http.Client{}

	for _, issue := range issues {
<<<<<<< HEAD
		exists, id := CheckIssueExists(databasePusher.database, "issues", "key", issue.Key)
		if exists {
			stmt, err :=
				databasePusher.database.Prepare("UPDATE issues set summary = ?, description = ?, type = ?, priority = ?, status = ?, closedtime = ?, updatedtime = ?, timespent = ? where key = ?")
			if err != nil {
				panic(err)
			}
			stmt.Exec(issue.Summary, issue.Description, issue.Type, issue.Priority, issue.Status, issue.ClosedTime, issue.UpdatedTime, issue.Timespent, issue.Key)

			stmt, err =
				databasePusher.database.Prepare("UPDATE project set title = ? where id = ?")
			if err != nil {
				panic(err)
			}
			projectId := databasePusher.database.QueryRow("SELECT projectId FROM $1 where $2 = $3", "issues", "assigneeId", id)
			stmt.Exec(issue.Project, projectId)

			stmt, err =
				databasePusher.database.Prepare("UPDATE author set name = ? where id = ?")
			if err != nil {
				panic(err)
			}
			authorId := databasePusher.database.QueryRow("SELECT authorId FROM $1 where $2 = $3", "issues", "assigneeId", id)
			stmt.Exec(issue.Author, authorId)

			stmt, err =
				databasePusher.database.Prepare("UPDATE statusChanges set changeTime = ?, fromStatus = ?, toStatus = ? where id = ?")
			if err != nil {
				panic(err)
			}
			stmt.Exec(777, "idk", "idk", authorId)
		}

		newProjectId, err := CountRows(databasePusher.database, "project")
		if err != nil {
			panic(err)
		}
		newAuthorId, err := CountRows(databasePusher.database, "author")
		if err != nil {
			panic(err)
		}
		newAssigneeid, err := CountRows(databasePusher.database, "issues")
		if err != nil {
			panic(err)
		}

		stmt, err :=
			databasePusher.database.Prepare("INSERT INTO issues (projectId, authorId, assigneeId, key, summary, description, type, priority, status, createdTime, closedTime, updatedTime, timeSpent) values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")
		if err != nil {
			panic(err)
		}
		stmt.Exec(newProjectId, newAuthorId, newAssigneeid, issue.Assignee, issue.Key, issue.Summary, issue.Description, issue.Type, issue.Priority, issue.Status, issue.CreatedTime, issue.ClosedTime, issue.UpdatedTime, issue.Timespent)

		stmt, err =
			databasePusher.database.Prepare("INSERT INTO project (id, title) values (?, ?)")
		if err != nil {
			panic(err)
		}
		stmt.Exec(newProjectId, issue.Project)

		stmt, err =
			databasePusher.database.Prepare("INSERT INTO author (id, name) values (?, ?)")
		if err != nil {
			panic(err)
		}
		stmt.Exec(newAuthorId, issue.Author)

		stmt, err =
			databasePusher.database.Prepare("INSERT INTO statusChange (issueId, authorId, changeTime, fromStatus, toStatus) values (?, ?, ?, ?, ?)")
		if err != nil {
			panic(err)
		}
		stmt.Exec(newAssigneeid, newAuthorId, 777, "idk", "idk")
=======
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
>>>>>>> f088cd021ec45f2588741546a422d8f32f0e2981
	}
}
