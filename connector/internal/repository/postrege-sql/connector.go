package postrege_sql

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/DinozvrrDan/jira-analyzer/connector/internal/models"
	"time"
)

type ConnectorRepository struct {
	db *sql.DB
}

func NewConnectorRepository(db *sql.DB) *ConnectorRepository {
	return &ConnectorRepository{
		db: db,
	}
}

func (connectorRepository *ConnectorRepository) PushIssues(issues []models.TransformedIssue) error {

	for _, issue := range issues {
		projectId, err := connectorRepository.getProjectId(issue.Project)

		if err != nil {
			return fmt.Errorf("error while getting projectId: %w", err)
		}

		authorId, err := connectorRepository.getAuthorId(issue.Author)

		if err != nil {
			return fmt.Errorf("error while getting authorId: %w", err)
		}

		assigneeId, err := connectorRepository.getAssigneeId(issue.Assignee)

		if err != nil {
			return fmt.Errorf("error while getting assigneeId: %w", err)
		}

		exists, err := connectorRepository.checkIssueExists(issue.Key)
		if err != nil {
			return fmt.Errorf("error while checking issue exists: %w", err)
		}
		if exists {
			err := connectorRepository.updateIssue(
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
				issue.TimeSpent)
			if err != nil {
				return fmt.Errorf("error while updating issue: %w", err)
			}
		} else {
			err := connectorRepository.insertIssue(
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
				issue.TimeSpent)
			if err != nil {
				return fmt.Errorf("error while inserting issue: %w", err)
			}
		}
	}
	return nil
}

func (connectorRepository *ConnectorRepository) insertIssue(projectId, authorId, assigneeId int64, key, summary, description, Type, priority, status string, createdTime, closedTime, updatedTime time.Time, timeSpent int64) error {
	_, err := connectorRepository.db.Exec("INSERT INTO issues (projectId, authorId, assigneeId, key, summary, description, type, priority, status, createdTime, closedTime, updatedTime, timeSpent) values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)",
		projectId,
		authorId,
		assigneeId,
		key,
		summary,
		description,
		Type,
		priority,
		status,
		createdTime,
		closedTime,
		updatedTime,
		timeSpent)

	if err != nil {
		return fmt.Errorf("insertIssue error: %w", err)
	}

	return nil
}

// updateIssue обвновляет данные issue заданного key в таблицк issues
func (connectorRepository *ConnectorRepository) updateIssue(projectId, authorId, assigneeId int64, key, summary, description, Type, priority, status string, createdTime, closedTime, updatedTime time.Time, timespent int64) error {

	_, err := connectorRepository.db.Exec("UPDATE issues set projectid = $1, authorid = $2, assigneeid = $3, summary = $4, description = $5, type = $6, priority = $7, status = $8, createdtime = $9, closedtime = $10, updatedtime = $11, timespent = $12 where key = $13",
		projectId,
		authorId,
		assigneeId,
		summary,
		description,
		Type,
		priority,
		status,
		createdTime,
		closedTime,
		updatedTime,
		timespent,
		key)

	if err != nil {
		return fmt.Errorf("updateIssue error: %w", err)
	}

	return nil
}

// getIssueId получает id по ключу задачи из таблицы issues
func (connectorRepository *ConnectorRepository) getIssueId(issueKey string) (int64, error) {
	var issueId int64
	row := connectorRepository.db.QueryRow("SELECT id FROM issues where key = ?", issueKey)

	if err := row.Scan(&issueId); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return issueId, fmt.Errorf("getIssueId error: no issue")
		}
		return issueId, fmt.Errorf("getIssueId error: %w", err)
	}

	return issueId, nil
}

// getProjectId получает id по названию проекта из таблицы project
func (connectorRepository *ConnectorRepository) getProjectId(projectTitle string) (int64, error) {
	var projectId int64

	row := connectorRepository.db.QueryRow("SELECT id FROM project where title = $1", projectTitle)

	if err := row.Scan(&projectId); err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return projectId, fmt.Errorf("getProjectId error: %w", err)
		}
	}

	if projectId == 0 {
		err := connectorRepository.db.QueryRow("INSERT INTO project (title) VALUES($1)", projectTitle).Scan(projectId)
		if !errors.Is(err, sql.ErrNoRows) {
			return projectId, fmt.Errorf("getProjectId error: %w", err)
		}
	}

	return projectId, nil
}

// getAuthorId получает id по имени автора из таблицы author
func (connectorRepository *ConnectorRepository) getAuthorId(authorName string) (int64, error) {
	var authorId int64
	row := connectorRepository.db.QueryRow("SELECT id FROM author where name = $1", authorName)

	if err := row.Scan(&authorId); err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return authorId, fmt.Errorf("getAuthorId error: %w", err)
		}
	}

	if authorId == 0 {
		err := connectorRepository.db.QueryRow("INSERT INTO author (name) VALUES($1)", authorName).Scan(authorId)
		if !errors.Is(err, sql.ErrNoRows) {
			return authorId, fmt.Errorf("getAuthorId error: %w", err)
		}
	}

	return authorId, nil
}

// getAssigneeId получает id по имени assignee из таблицы author
func (connectorRepository *ConnectorRepository) getAssigneeId(assigneeName string) (int64, error) {
	var assigneeId int64
	row := connectorRepository.db.QueryRow("SELECT id FROM author where name = $1", assigneeName)

	if err := row.Scan(&assigneeId); err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return assigneeId, fmt.Errorf("getAssigneeId error: %w", err)
		}
	}

	if assigneeId == 0 {
		err := connectorRepository.db.QueryRow("INSERT INTO author (name) VALUES($1)", assigneeName).Scan(assigneeId)
		if !errors.Is(err, sql.ErrNoRows) {
			return assigneeId, fmt.Errorf("getAuthorId error: %w", err)
		}
	}

	return assigneeId, nil
}

// checkIssueExists проверяет наличие issue заданного issueKey
func (connectorRepository *ConnectorRepository) checkIssueExists(issueKey string) (bool, error) {
	var issueId int64
	row := connectorRepository.db.QueryRow("SELECT id FROM issues where key = $1", issueKey)

	if err := row.Scan(&issueId); err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return false, fmt.Errorf("checkIssueExists error: %w", err)
		}
	}

	return issueId != 0, nil
}

// CheckProjectExists проверяет наличие issue заданного title
func (connectorRepository *ConnectorRepository) CheckProjectExists(title string) (bool, error) {
	var projectId int64
	row := connectorRepository.db.QueryRow("SELECT id FROM project where title = $1", title)

	if err := row.Scan(&projectId); err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return false, fmt.Errorf("checkIssueExists error: %w", err)
		}
	}

	return projectId != 0, nil
}
