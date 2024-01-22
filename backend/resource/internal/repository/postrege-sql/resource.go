package postrege_sql

import (
	"database/sql"
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

	stmt, _ := resourceRepository.db.
		Prepare("SELECT id, projectId, authorId, assigneeId, key, summary, description, type, priority, createdTime, closedTime, updatedTime, timeSpent FROM issues where id = ?")

	err := stmt.QueryRow(id).Scan(
		&issueInfo.ID,
		&issueInfo.Project.ID,
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
		&issueInfo.TimeSpent)

	if err != nil {
		return issueInfo, err
	}

	err = resourceRepository.db.QueryRow("SELECT name FROM author where id = ?", authorID).Scan(&issueInfo.Author)
	if err != nil {
		return issueInfo, err
	}

	err = resourceRepository.db.QueryRow("SELECT id FROM author where name = ?", assigneeID).Scan(&issueInfo.Assignee)
	if err != nil {
		return issueInfo, err
	}

	return issueInfo, nil
}

func (resourceRepository *ResourceRepository) GetHistoryInfo(id int) ([]models.HistoryInfo, error) {
	var historyInfos []models.HistoryInfo

	rows, err := resourceRepository.db.Query(
		"SELECT "+
			"authorId,"+
			"changeTime"+
			"fromStatus,"+
			"toStatus "+
			"FROM statusChanges "+
			"WHERE issueId = ?", id,
	)

	if err != nil {
		return historyInfos, err
	}

	for rows.Next() {
		historyInfo := models.HistoryInfo{}
		err := rows.Scan(&historyInfo.AuthorID, &historyInfo.ChangeTime, &historyInfo.FromStatus, &historyInfo.ToStatus)

		if err != nil {
			return historyInfos, err
		}

		historyInfos = append(historyInfos, historyInfo)
	}

	return historyInfos, nil
}

func (resourceRepository *ResourceRepository) GetProjectInfo(id int) (models.ProjectInfo, error) {
	projectInfo := models.ProjectInfo{}

	stmt, _ := resourceRepository.db.Prepare("SELECT id, title FROM project WHERE id = ?")
	err := stmt.QueryRow(id).Scan(&projectInfo.ID, &projectInfo.Title)

	if err != nil {
		return projectInfo, err
	}

	return projectInfo, nil
}

func (resourceRepository *ResourceRepository) InsertProject(projectInfo models.ProjectInfo) (int, error) {
	var projectId int

	if err := resourceRepository.db.QueryRow("INSERT INTO project (title) VALUES(?) RETURNING id",
		projectInfo.Title).Scan(&projectId); err != nil {
		return projectId, err
	}

	return projectId, nil
}

func (resourceRepository *ResourceRepository) InsertIssue(issueInfo models.IssueInfo) (int, error) {
	var issueId, authorId, assigneeId int

	err := resourceRepository.db.QueryRow("SELECT id FROM author WHERE name = ?", issueInfo.Author).Scan(&authorId)
	if err != nil {
		return issueId, err
	}

	err = resourceRepository.db.QueryRow("SELECT id FROM author where name = ?", issueInfo.Assignee).Scan(&assigneeId)
	if err != nil {
		return issueId, err
	}

	stmt, _ := resourceRepository.db.Prepare("INSERT INTO issues (projectId, authorId, assigneeId, key, summary, description, type, priority, status, createdTime, closedTime, updatedTime, timeSpent) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")

	err = stmt.QueryRow(issueInfo.Project.ID,
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
		issueInfo.TimeSpent).Scan(&issueId)

	if err != nil {
		return issueId, err
	}

	return issueId, nil
}

func (resourceRepository *ResourceRepository) InsertHistory(historyInfo models.HistoryInfo) (int, error) {
	var historyID int

	stmt, _ := resourceRepository.db.
		Prepare("INSERT INTO StatusChanges (issueId,authorId,changeTime,fromStatus,toStatus) VALUES (?, ?, now(), ?, ?)")

	err := stmt.QueryRow(historyInfo.IssueID, historyInfo.AuthorID, historyInfo.FromStatus, historyInfo.ToStatus).Err()

	if err != nil {
		return historyID, err
	}

	stmt, _ = resourceRepository.db.
		Prepare("UPDATE Issue SET status = ?, updatedTime = now(), timespent = now()-createdTime WHERE id = ?")

	err = stmt.QueryRow(historyInfo.ToStatus, historyInfo.IssueID).Err()

	if err != nil {
		return historyID, err
	}

	return historyID, nil
}
