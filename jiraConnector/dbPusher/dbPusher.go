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
	"time"

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
	for _, issue := range issues {
		exists, assigneeId := databasePusher.CheckIssueExists("issues", "key", issue.Key)
		if exists {
			databasePusher.updateIssue(issue.Summary, issue.Description, issue.Type, issue.Priority, issue.Status, issue.Key, issue.ClosedTime, issue.UpdatedTime, issue.Timespent)

			projectId := databasePusher.getProjectId(assigneeId)
			databasePusher.updateProject(issue.Project, projectId)

			authorId := databasePusher.getAuthorId(assigneeId)
			databasePusher.updateAuthor(issue.Author, authorId)

			databasePusher.updateStatusChanges(777, "FromStatus", "ToStatus", 2)
		}

		newProjectId, err := databasePusher.CountRows("project")
		if err != nil {
			databasePusher.logger.Log(logger.ERROR, err.Error())
		}
		newAuthorId, err := databasePusher.CountRows("author")
		if err != nil {
			databasePusher.logger.Log(logger.ERROR, err.Error())
		}
		newAssigneeid, err := databasePusher.CountRows("issues")
		if err != nil {
			databasePusher.logger.Log(logger.ERROR, err.Error())
		}

		stmt, _ :=
			databasePusher.database.Prepare("INSERT INTO issues (projectId, authorId, assigneeId, key, summary, description, type, priority, status, createdTime, closedTime, updatedTime, timeSpent) values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")
		_, err = stmt.Exec(newProjectId, newAuthorId, newAssigneeid, issue.Assignee, issue.Key, issue.Summary, issue.Description, issue.Type, issue.Priority, issue.Status, issue.CreatedTime, issue.ClosedTime, issue.UpdatedTime, issue.Timespent)
		if err != nil {
			databasePusher.logger.Log(logger.ERROR, err.Error())
		}

		stmt, _ =
			databasePusher.database.Prepare("INSERT INTO project (id, title) values (?, ?)")
		_, err = stmt.Exec(newProjectId, issue.Project)
		if err != nil {
			databasePusher.logger.Log(logger.ERROR, err.Error())
		}

		stmt, _ =
			databasePusher.database.Prepare("INSERT INTO author (id, name) values (?, ?)")
		_, err = stmt.Exec(newAuthorId, issue.Author)
		if err != nil {
			databasePusher.logger.Log(logger.ERROR, err.Error())
		}

		stmt, _ =
			databasePusher.database.Prepare("INSERT INTO statusChange (issueId, authorId, changeTime, fromStatus, toStatus) values (?, ?, ?, ?, ?)")
		_, err = stmt.Exec(newAssigneeid, newAuthorId, 777, "idk", "idk")
		if err != nil {
			databasePusher.logger.Log(logger.ERROR, err.Error())
		}
	}
}

// updateAuthor обновляет имя автора заданного id
func (databasePusher *DatabasePusher) updateAuthor(Author string, authorId int) {
	stmt, _ := databasePusher.database.Prepare("UPDATE author set name = ? where id = ?")
	_, err := stmt.Exec(Author, authorId)
	if err != nil {
		databasePusher.logger.Log(logger.ERROR, err.Error())
	}
}

// updateProject обновляет название проекта заданного id
func (databasePusher *DatabasePusher) updateProject(Project string, projectId int) {
	stmt, _ :=
		databasePusher.database.Prepare("UPDATE project set title = ? where id = ?")
	_, err := stmt.Exec(Project, projectId)
	if err != nil {
		databasePusher.logger.Log(logger.ERROR, err.Error())
	}
}

// updateStatusChanges обновляет ChangeTime, FromStatus, ToStatus таблицы StatusChanges заданного AuthorId
func (databasePusher *DatabasePusher) updateStatusChanges(ChangeTime int, FromStatus, ToStatus string, AuthorId int) {
	stmt, _ := databasePusher.database.Prepare("UPDATE statusChanges set changeTime = ?, fromStatus = ?, toStatus = ? where authorId = ?")
	_, err := stmt.Exec(ChangeTime, FromStatus, ToStatus, AuthorId)
	if err != nil {
		databasePusher.logger.Log(logger.ERROR, err.Error())
	}
}

// updateIssue получает id проекта из таблицы issues по assigneeId
func (databasePusher *DatabasePusher) updateIssue(Summary, Description, Type, Priority, Status, Key string, ClosedTime, UpdatedTime time.Time, Timespent int) {
	stmt, _ :=
		databasePusher.database.Prepare("UPDATE issues set summary = ?, description = ?, type = ?, priority = ?, status = ?, closedtime = ?, updatedtime = ?, timespent = ? where key = ?")
	_, err := stmt.Exec(Summary, Description, Type, Priority, Status, ClosedTime, UpdatedTime, Timespent, Key)
	if err != nil {
		databasePusher.logger.Log(logger.ERROR, err.Error())
	}
}

// getProjectId получает id проекта из таблицы issues по assigneeId
func (databasePusher *DatabasePusher) getProjectId(assigneeId int) int {
	var projectId int
	databasePusher.database.QueryRow("SELECT projectId FROM $1 where $2 = $3", "issues", "assigneeId", assigneeId).Scan(&projectId)
	return projectId
}

// getAuthorId получает id автора из таблицы issues по assigneeId
func (databasePusher *DatabasePusher) getAuthorId(assigneeId int) int {
	var authorId int
	err := databasePusher.database.QueryRow("SELECT authorId FROM $1 where $2 = $3", "issues", "assigneeId", assigneeId).Scan(&authorId)

	if err != nil {
		databasePusher.logger.Log(logger.ERROR, err.Error())
	}
	return authorId
}

// CheckIssueExists проверяет есть ли задача с заданным issueKey в заданной таблице и возвращает ее assigneeId
func (databasePusher *DatabasePusher) CheckIssueExists(table, column string, issueKey string) (bool, int) {
	var assigneeId int
	err := databasePusher.database.QueryRow("SELECT assigneeId FROM $1 where $2 = $3", table, column, issueKey).Scan(&assigneeId)

	if err != nil {
		return false, assigneeId
	} else {
		return true, assigneeId
	}
}

// CountRows считает количество данных в заданной таблице
func (databasePusher *DatabasePusher) CountRows(table string) (int, error) {
	stmt, err := databasePusher.database.Prepare("SELECT COUNT(*) FROM ?")
	if err != nil {
		return 0, err
	}

	var result int
	err = stmt.QueryRow(table).Scan(&result)
	return result, err
}
