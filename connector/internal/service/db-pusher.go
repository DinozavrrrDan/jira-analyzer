package service

/*
   Подключение к базе данных должно осуществляться через драйвер
   database/sql
   При загрузке данных из JIRA в БД должна поддерживаться
   атомарность, то есть, если при скачивании части данных произошла
   ошибка, то никакие данные не будут записаны в БД (все или ничего)
*/

import (
	"connector/config"
	"connector/internal/models"
	"connector/pkg/logger"
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

type DatabasePusherService struct {
	cfg      *config.Config
	log      *logger.Logger
	database *sql.DB
}

func NewDatabasePusher(log *logger.Logger, cfg *config.Config) *DatabasePusherService {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cfg.HostDB,
		cfg.PortDB,
		cfg.UserDB,
		cfg.PasswordDB,
		cfg.NameDB)
	newDatabase, err := sql.Open("postgres", psqlInfo)

	if err != nil {
		log.Log(logger.ERROR, "Error while open db:"+err.Error())
	}

	return &DatabasePusherService{
		cfg:      cfg,
		log:      log,
		database: newDatabase,
	}
}

func (databasePusher *DatabasePusherService) PushIssue(issues []models.TransformedIssue) {

	for _, issue := range issues {
		projectId, err := databasePusher.getProjectId(issue.Project)

		if err != nil {
			databasePusher.log.Log(logger.ERROR, err.Error())

			return
		}

		authorId, err := databasePusher.getAuthorId(issue.Author)

		if err != nil {
			databasePusher.log.Log(logger.ERROR, err.Error())

			return
		}

		assigneeId, err := databasePusher.getAssigneeId(issue.Assignee)

		if err != nil {
			databasePusher.log.Log(logger.ERROR, err.Error())

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
				databasePusher.log.Log(logger.ERROR, err.Error())
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
				databasePusher.log.Log(logger.ERROR, fmt.Sprintf("ERROR: %v", err.Error()))

				return
			}
		}
	}
}
