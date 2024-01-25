package postrege_sql

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/DinozvrrDan/jira-analyzer/backend/resource/internal/models"
	_ "github.com/lib/pq"
	"time"
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

const openStatus = "Open"
const closedStatus = "Closed"
const resolvedStatus = "Resolved"
const reopenedStatus = "Reopened"
const progressStatus = "In Progress"

func (resourceRepository *ResourceRepository) GetProjectInfo(title string) (models.ProjectInfo, error) {
	id, _ := resourceRepository.getProjectId(title)
	projectInfo := models.ProjectInfo{}

	err := resourceRepository.db.QueryRow("SELECT id, title,"+
		" (SELECT COUNT(*) from issues where projectid = $1) as issuesCount,"+
		" (SELECT COUNT(*) from issues where projectid = $1 and status = $2) as openedIssuesCount,"+
		" (SELECT COUNT(*) from issues where projectid = $1 and status = $3) as closedIssuesCount,"+
		" (SELECT COUNT(*) from issues where projectid = $1 and status = $4) as resolvedIssuesCount,"+
		" (SELECT COUNT(*) from issues where projectid = $1 and status = $5) as reopenedIssueCount,"+
		" (SELECT COUNT(*) from issues where projectid = $1 and status = $6) as progressIssuesCount"+
		" FROM project WHERE title = $7",
		id, openStatus, closedStatus, resolvedStatus, reopenedStatus, progressStatus, title).
		Scan(
			&projectInfo.Id,
			&projectInfo.Title,
			&projectInfo.IssueCount,
			&projectInfo.OpenedIssuesCount,
			&projectInfo.ClosedIssuesCount,
			&projectInfo.ResolvedIssuesCount,
			&projectInfo.ReopenedIssueCount,
			&projectInfo.ProgressIssuesCount)

	if err != nil {
		return models.ProjectInfo{}, fmt.Errorf("GetProjectInfo: %v", err)
	}

	avrTime, err := resourceRepository.GetAvrgActiveTime(id)

	if err != nil {
		return models.ProjectInfo{}, fmt.Errorf("GetProjectInfo: %v", err)
	}
	projectInfo.AverageTime = int(avrTime)

	avrIssue, err := resourceRepository.GetNewProjectsAvrgAmount(id)

	if err != nil {
		return models.ProjectInfo{}, fmt.Errorf("GetProjectInfo: %v", err)
	}
	projectInfo.AverageIssue = int(avrIssue)

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
	id, _ := resourceRepository.getProjectId(title)

	ids, _ := resourceRepository.getAuthorIds(id)
	if err := resourceRepository.DeleteAuthors(ids); err != nil {
		fmt.Println(err.Error())
		return fmt.Errorf("DeleteProject error: %w", err)
	}

	if err := resourceRepository.DeleteIssues(id); err != nil {
		return fmt.Errorf("DeleteProject error: %w", err)
	}

	if _, err := resourceRepository.db.Exec("DELETE FROM project WHERE id = $1", id); err != nil {
		return fmt.Errorf("DeleteProject error: %w", err)
	}
	return nil
}

func (resourceRepository *ResourceRepository) getProjectId(projectTitle string) (int64, error) {
	var projectId int64

	err := resourceRepository.db.QueryRow("SELECT id FROM project where title = $1", projectTitle).Scan(&projectId)
	if err != nil {
		return 0, fmt.Errorf("getProjectId error: %w", err)
	}

	return projectId, nil
}

func (resourceRepository *ResourceRepository) DeleteAuthors(ids []int64) error {
	for _, id := range ids {
		if _, err := resourceRepository.db.Exec("DELETE FROM author WHERE id = $1", id); err != nil {
			if !errors.Is(err, sql.ErrNoRows) {
				return fmt.Errorf("DeleteAuthor error: %w", err)
			}
		}
	}

	return nil
}

func (resourceRepository *ResourceRepository) DeleteIssues(projectId int64) error {

	if _, err := resourceRepository.db.Exec("DELETE FROM issues WHERE projectid = $1", projectId); err != nil {
		return fmt.Errorf("DeleteIssues error: %w", err)
	}
	return nil
}

func (resourceRepository *ResourceRepository) GetProjects() ([]models.Project, error) {
	rows, err := resourceRepository.db.Query("SELECT id, title FROM project")
	if err != nil {
		return nil, fmt.Errorf("GetProjects error: %w", err)
	}
	defer rows.Close()

	var projects []models.Project

	for rows.Next() {
		var tmp models.Project
		err := rows.Scan(&tmp.Id, &tmp.Name)
		if err != nil {
			return nil, fmt.Errorf("GetProjects error: %w", err)
		}
		tmp.Existence = true

		projects = append(projects, tmp)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("GetProjects error: %w", err)
	}

	return projects, nil

}

func (resourceRepository *ResourceRepository) GetAvrgActiveTime(projectId int64) (float64, error) {
	var result float64
	var resolvedProjectsAmount int64

	rows, err := resourceRepository.db.Query("SELECT timespent FROM issues WHERE projectid = $1 AND status = $2", projectId, "Resolved")
	if err != nil {
		return result, fmt.Errorf("GetAvrgActiveTime error: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var projectTimeSpent float64
		err := rows.Scan(&projectTimeSpent)
		if err != nil {
			return result, fmt.Errorf("GetAvrgActiveTime error: %w", err)
		}
		resolvedProjectsAmount++
		result += projectTimeSpent
	}

	if err = rows.Err(); err != nil {
		return result, fmt.Errorf("GetAvrgActiveTime error: %w", err)
	}

	result /= float64(resolvedProjectsAmount)

	const minute = 60.0
	const hour = 60.0
	return result / (minute * hour), nil
}

func (resourceRepository *ResourceRepository) GetNewProjectsAvrgAmount(projectId int64) (float64, error) {
	var result float64
	actualTime := time.Now()

	rows, err := resourceRepository.db.Query("SELECT createdtime FROM issues WHERE projectid = $1", projectId)
	if err != nil {
		return result, fmt.Errorf("GetAvrgActiveTime error: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var projectCreatedTime time.Time
		err := rows.Scan(&projectCreatedTime)
		if err != nil {
			return result, fmt.Errorf("GetAvrgActiveTime error: %w", err)
		}

		if actualTime.Sub(projectCreatedTime) <= 60*60*24*7*1e9 {
			result += 1.0
		}
	}

	if err = rows.Err(); err != nil {
		return result, fmt.Errorf("GetAvrgActiveTime error: %w", err)
	}

	return result / 7.0, nil
}

func (resourceRepository *ResourceRepository) getAuthorIds(projectId int64) ([]int64, error) {
	var authorId, assigneeId int64
	var result []int64

	rows, err := resourceRepository.db.Query("SELECT authorid, assigneeid FROM issues where projectid = $1", projectId)
	if err != nil {
		return result, fmt.Errorf("getAuthorIds error: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&authorId, &assigneeId)
		if err != nil {
			return result, fmt.Errorf("getAuthorIds error: %w", err)
		}
		result = append(result, authorId)
		result = append(result, assigneeId)
	}

	return result, nil
}
