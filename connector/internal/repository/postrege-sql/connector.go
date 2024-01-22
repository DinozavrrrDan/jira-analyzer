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

	_, err := connectorRepository.db.Exec(
		"INSERT INTO issues "+
			"(projectId, authorId, assigneeId,"+
			" key, summary, description, type, priority, status,"+
			" createdTime, closedTime, updatedTime, timeSpent)"+
			" values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)",
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
func (connectorRepository *ConnectorRepository) updateIssue(projectID, authorId, assigneeId int64, key, summary, description, Type, priority, status string, createdTime, closedTime, updatedTime time.Time, timespent int64) error {

	_, err := connectorRepository.db.Exec("UPDATE issues set"+
		" projectid = ?, authorid = ?, assigneeid = ?,"+
		" summary = ?, description = ?, type = ?, priority = ?, status = ?,"+
		" createdtime = ?, closedtime = ?, updatedtime = ?, timespent = ? where key = ?",
		projectID,
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
			return issueId, fmt.Errorf("getIssueId error: no issue", issueId)
		}
		return issueId, fmt.Errorf("getIssueId error: %w", issueId, err)
	}

	return issueId, nil
}

// getProjectId получает id по названию проекта из таблицы project
func (connectorRepository *ConnectorRepository) getProjectId(projectTitle string) (int64, error) {
	var projectId int64
	row := connectorRepository.db.QueryRow("SELECT id FROM project where title = ?", projectTitle)

	if err := row.Scan(&projectId); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return projectId, fmt.Errorf("getProjectId error: no project", projectId)
		}
		return projectId, fmt.Errorf("getProjectId error: %w", projectId, err)
	}

	if projectId == 0 {
		result, err := connectorRepository.db.Exec("INSERT INTO project (title) VALUES(?)", projectTitle)
		if err != nil {
			return projectId, fmt.Errorf("getProjectId error: %w", err.Error())
		}

		projectId, err := result.LastInsertId()
		if err != nil {
			return projectId, fmt.Errorf("getProjectId error: %w", err.Error())
		}
	}

	return projectId, nil
}

// getAuthorId получает id по имени автора из таблицы author
func (connectorRepository *ConnectorRepository) getAuthorId(authorName string) (int64, error) {
	var authorId int64
	row := connectorRepository.db.QueryRow("SELECT id FROM author where name = ?", authorName)

	if err := row.Scan(&authorId); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return authorId, fmt.Errorf("getAuthorId error: no author", authorId)
		}
		return authorId, fmt.Errorf("getAuthorId error: %w", authorId, err)
	}

	if authorId == 0 {
		result, err := connectorRepository.db.Exec("INSERT INTO author (name) VALUES(?)", authorName)
		if err != nil {
			return authorId, fmt.Errorf("getAuthorId error: %w", err.Error())
		}

		authorId, err := result.LastInsertId()
		if err != nil {
			return authorId, fmt.Errorf("getAuthorId error: %w", err.Error())
		}
	}

	return authorId, nil
}

// getAssigneeId получает id по имени assignee из таблицы author
func (connectorRepository *ConnectorRepository) getAssigneeId(assigneeName string) (int64, error) {
	var assigneeId int64
	row := connectorRepository.db.QueryRow("SELECT id FROM author where name = ?", assigneeName)

	if err := row.Scan(&assigneeId); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return assigneeId, fmt.Errorf("getAssigneeId error: no assignee", assigneeId)
		}
		return assigneeId, fmt.Errorf("getAssigneeId error: %w", err)
	}

	if assigneeId == 0 {
		result, err := connectorRepository.db.Exec("INSERT INTO author (name) VALUES(?)", assigneeName)
		if err != nil {
			return assigneeId, fmt.Errorf("getAssigneeId error: %w", err.Error())
		}

		authorId, err := result.LastInsertId()
		if err != nil {
			return authorId, fmt.Errorf("getAssigneeId error: LastInsertId: %w", err.Error())
		}
	}

	return assigneeId, nil
}

// checkIssueExists проверяет наличие issue заданного issueKey
func (connectorRepository *ConnectorRepository) checkIssueExists(issueKey string) (bool, error) {
	var issueId int64
	row := connectorRepository.db.QueryRow("SELECT id FROM issues where key = ?", issueKey)

	if err := row.Scan(&issueId); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, fmt.Errorf("checkIssueExists error: no issue", issueId)
		}
		return false, fmt.Errorf("checkIssueExists error: %w", issueId, err)
	}

	return issueId != 0, nil
}
