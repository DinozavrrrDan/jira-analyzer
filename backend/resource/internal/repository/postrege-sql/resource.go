package postrege_sql

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/DinozvrrDan/jira-analyzer/backend/resource/internal/models"
	_ "github.com/lib/pq"
)

type ResourceRepository struct {
	db *sql.DB
}

func NewResourceRepository(db *sql.DB) *ResourceRepository {
	return &ResourceRepository{
		db: db,
	}
}

func (resourceRepository *ResourceRepository) GetIssueInfo(id int) (models.IssueInfo, error) {
	issueInfo := models.IssueInfo{}

	var authorID, assigneeID int

	row := resourceRepository.db.QueryRow(
		"SELECT "+
			"id, projectId, authorId, assigneeId, "+
			"key, summary, description, type, priority, "+
			"createdTime, closedTime, updatedTime, timeSpent "+
			"FROM issues where id = ?", id)

	if err := row.Scan(
		&issueInfo.Id,
		&issueInfo.Project.Id,
		&authorID,
		&assigneeID,
		&issueInfo.Key,
		&issueInfo.Summary,
		&issueInfo.Description,
		&issueInfo.Type,
		&issueInfo.Priority,
		&issueInfo.Status,
		&issueInfo.CreatedTime,
		&issueInfo.ClosedTime,
		&issueInfo.UpdatedTime,
		&issueInfo.TimeSpent); err != nil {

		if errors.Is(err, sql.ErrNoRows) {
			return issueInfo, fmt.Errorf("GetIssueInfo: no such issue: %v", id)
		}
		return issueInfo, fmt.Errorf("GetIssueInfo %d: %v", id, err)
	}

	row = resourceRepository.db.QueryRow("SELECT name FROM author where id = ?", authorID)
	if err := row.Scan(issueInfo.Author); err != nil {

		if errors.Is(err, sql.ErrNoRows) {
			return issueInfo, fmt.Errorf("GetIssueInfo: no such author %d", authorID)
		}
		return issueInfo, fmt.Errorf("GetIssueInfo: %v", err)
	}

	row = resourceRepository.db.QueryRow("SELECT name FROM author where id = ?", assigneeID)
	if err := row.Scan(issueInfo.Assignee); err != nil {

		if errors.Is(err, sql.ErrNoRows) {
			return issueInfo, fmt.Errorf("GetIssueInfo: no such author %d", assigneeID)
		}
		return issueInfo, fmt.Errorf("GetIssueInfo: %v", err)
	}

	return issueInfo, nil
}

func (resourceRepository *ResourceRepository) GetProjectInfo(id int) (models.ProjectInfo, error) {
	projectInfo := models.ProjectInfo{}

	row := resourceRepository.db.QueryRow("SELECT id, title FROM project WHERE id = ?", id)
	if err := row.Scan(&projectInfo.Id, &projectInfo.Title); err != nil {

		if errors.Is(err, sql.ErrNoRows) {
			return projectInfo, fmt.Errorf("GetProjectInfo: no such project %d", id)
		}
		return projectInfo, fmt.Errorf("GetProjectInfo: %v", err)
	}

	return projectInfo, nil
}

func (resourceRepository *ResourceRepository) InsertProject(projectInfo models.ProjectInfo) (int64, error) {
	var projectId int64

	result, err := resourceRepository.db.Exec("INSERT INTO project (title) VALUES($1)", projectInfo.Title)
	if err != nil {
		return 0, fmt.Errorf("InsertProject: %v", err)
	}
	projectId, err = result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("InsertProject: %v", err)
	}

	return projectId, nil
}

func (resourceRepository *ResourceRepository) InsertIssue(issueInfo models.IssueInfo) (int64, error) {
	var issueId, authorId, assigneeId int64

	row := resourceRepository.db.QueryRow("SELECT id FROM author WHERE name = ?", issueInfo.Author)
	if err := row.Scan(&authorId); err != nil {

		if errors.Is(err, sql.ErrNoRows) {
			return issueId, fmt.Errorf("InsertIssue: no such author %s", issueInfo.Author)
		}
		return issueId, fmt.Errorf("InsertIssue: %v", err)
	}

	row = resourceRepository.db.QueryRow("SELECT id FROM author WHERE name = ?", issueInfo.Assignee)
	if err := row.Scan(&assigneeId); err != nil {

		if err == sql.ErrNoRows {
			return issueId, fmt.Errorf("InsertIssue: no such assignee %s", issueInfo.Assignee)
		}
		return issueId, fmt.Errorf("InsertIssue: %w", err)
	}

	result, err := resourceRepository.db.Exec(
		"INSERT INTO issues (projectId, authorId, assigneeId, key, summary, description, type, priority, status, createdTime, closedTime, updatedTime, timeSpent) values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)",
		issueInfo.Project.Id,
		authorId,
		assigneeId,
		issueInfo.Key,
		issueInfo.Summary,
		issueInfo.Description,
		issueInfo.Type,
		issueInfo.Priority,
		issueInfo.Status,
		issueInfo.CreatedTime,
		issueInfo.ClosedTime,
		issueInfo.UpdatedTime,
		issueInfo.TimeSpent)

	if err != nil {
		return issueId, fmt.Errorf("insertIssue: %w", err)
	}
	issueId, err = result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("insertIssue: %w", err)
	}

	return issueId, nil
}

func (resourceRepository *ResourceRepository) DeleteProject(title string) error {
	if err := resourceRepository.db.QueryRow("DELETE FROM project WHERE title = $1", title); err != nil {
		return fmt.Errorf("DeleteProject error: %w", err)
	}

	return nil
}

func (resourceRepository *ResourceRepository) DeleteAuthor(name string) error {
	if err := resourceRepository.db.QueryRow("DELETE FROM author WHERE name = $1", name); err != nil {
		return fmt.Errorf("DeleteAuthor error: %w", err)
	}

	return nil
}

func (resourceRepository *ResourceRepository) DeleteAssignee(name string) error {
	if err := resourceRepository.db.QueryRow("DELETE FROM author WHERE name = $1", name); err != nil {
		return fmt.Errorf("DeleteAssignee error: %w", err)
	}

	return nil
}

func (resourceRepository *ResourceRepository) DeleteIssue(issue models.IssueInfo) error {
	if err := resourceRepository.DeleteProject(issue.Project.Title); err != nil {
		return fmt.Errorf("DeleteProject error: %w", err)
	}

	if err := resourceRepository.DeleteAuthor(issue.Author); err != nil {
		return fmt.Errorf("DeleteProject error: %w", err)
	}

	if err := resourceRepository.DeleteAssignee(issue.Assignee); err != nil {
		return fmt.Errorf("DeleteProject error: %w", err)
	}

	if err := resourceRepository.db.QueryRow("DELETE FROM issues WHERE key = $1", issue.Key); err != nil {
		return fmt.Errorf("DeleteProject error: %w", err)
	}

	return nil
}
