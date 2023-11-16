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

func CheckIssueExists(db *sql.DB, table, column string, value string) bool {
	row := db.QueryRow("SELECT * FROM $1 where $2 = $3", table, column, value)
	err := row.Scan(&value)
	if err != nil {
		return false
	} else {
		return true
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
	for i, issue := range issues {
		if CheckIssueExists(databasePusher.database, "Issue", "key", issue.Key) {
			stmt, err :=
				databasePusher.database.Prepare("UPDATE Issue set summary = ?, description = ?, type = ?, priority = ?, status = ?, closedtime = ?, updatedtime = ?, timespent = ? where key = ?")
			if err != nil {
				panic(err)
			}
			stmt.Exec(issue.Summary, issue.Description, issue.Type, issue.Priority, issue.Status, issue.ClosedTime, issue.UpdatedTime, issue.Timespent, issue.Key)

			stmt, err =
				databasePusher.database.Prepare("UPDATE Project set title = ? where key = ?")
			if err != nil {
				panic(err)
			}
			stmt.Exec(issue.Project, i)
		}

		authorId, err := CountRows(databasePusher.database, "Author")
		if err != nil {
			panic(err)
		}
		projectId, err := CountRows(databasePusher.database, "Project")
		if err != nil {
			panic(err)
		}
		assigneeid, err := CountRows(databasePusher.database, "Issue")
		if err != nil {
			panic(err)
		}

		query, _ :=
			databasePusher.database.Prepare("INSERT INTO Issue (projectid, authorid, assigneeid, key, summary, description, type, priority, status, createdtime, closedtime, updatedtime, timespent) values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")
		query.Exec(projectId, authorId, assigneeid, issue.Assignee, issue.Key, issue.Summary, issue.Description, issue.Type, issue.Priority, issue.Status, issue.CreatedTime, issue.ClosedTime, issue.UpdatedTime, issue.Timespent)

		query, _ =
			databasePusher.database.Prepare("INSERT INTO Author (id, name) values (?, ?)")
		query.Exec(i, issue.Author)

		query, _ =
			databasePusher.database.Prepare("INSERT INTO Project (id, title) values (?, ?)")
		query.Exec(i, issue.Project)
	}
}
