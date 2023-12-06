package endpoints

import (
	"Jira-analyzer/analyzer/models"
	"Jira-analyzer/common/logger"
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

func (resourceHandler *ResourceHandler) getIssue(id int) (models.IssueInfo, error) {
	issueInfo := models.IssueInfo{}

	var authorID, assigneeID int

	row := resourceHandler.database.QueryRow(
		"SELECT "+
			"id, projectId, authorId, assigneeId, "+
			"key, summary, description, type, priority, "+
			"createdTime, closedTime, updatedTime, timeSpent "+
			"FROM issues where id = ?", id)

	if err := row.Scan(
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
		&issueInfo.TimeSpent); err != nil {

		if err == sql.ErrNoRows {
			return issueInfo, fmt.Errorf("GetIssueInfo: no such issue: %v", id)
		}
		return issueInfo, fmt.Errorf("GetIssueInfo %d: %v", id, err)
	}

	row = resourceHandler.database.QueryRow("SELECT name FROM author where id = ?", authorID)
	if err := row.Scan(issueInfo.Author); err != nil {

		if err == sql.ErrNoRows {
			return issueInfo, fmt.Errorf("GetIssueInfo: no such author %d", authorID)
		}
		return issueInfo, fmt.Errorf("GetIssueInfo: %v", err)
	}

	row = resourceHandler.database.QueryRow("SELECT name FROM author where id = ?", assigneeID)
	if err := row.Scan(issueInfo.Assignee); err != nil {

		if err == sql.ErrNoRows {
			return issueInfo, fmt.Errorf("GetIssueInfo: no such author %d", assigneeID)
		}
		return issueInfo, fmt.Errorf("GetIssueInfo: %v", err)
	}

	resourceHandler.logger.Log(logger.INFO, "GetIssueInfo successfully")

	return issueInfo, nil
}

func (resourceHandler *ResourceHandler) getProject(id int) (models.ProjectInfo, error) {
	projectInfo := models.ProjectInfo{}

	row := resourceHandler.database.QueryRow("SELECT id, title FROM project WHERE id = ?", id)
	if err := row.Scan(&projectInfo.ID, &projectInfo.Title); err != nil {

		if err == sql.ErrNoRows {
			return projectInfo, fmt.Errorf("GetProjectInfo: no such project %d", id)
		}
		return projectInfo, fmt.Errorf("GetProjectInfo: %v", err)
	}

	resourceHandler.logger.Log(logger.INFO, "GetProjectInfo successfully")

	return projectInfo, nil
}

func (resourceHandler *ResourceHandler) insertProject(projectInfo models.ProjectInfo) (int64, error) {
	var projectId int64

	result, err := resourceHandler.database.Exec("INSERT INTO project (title) VALUES(?)", projectInfo.Title)
	if err != nil {
		return 0, fmt.Errorf("InsertProject: %v", err)
	}
	projectId, err = result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("InsertProject: %v", err)
	}

	resourceHandler.logger.Log(logger.INFO, "InsertProject successfully")

	return projectId, nil
}

func (resourceHandler *ResourceHandler) deleteProject(projectInfo models.ProjectInfo) (int64, error) {
	var projectId int64

	result, err := resourceHandler.database.Exec("DELETE FROM project WHERE title=?", projectInfo.Title)
	if err != nil {
		return 0, fmt.Errorf("deleteProject: %v", err)
	}
	projectId, err = result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("deleteProject: %v", err)
	}

	resourceHandler.logger.Log(logger.INFO, "deleteProject successfully")

	return projectId, nil
}

func (resourceHandler *ResourceHandler) insertIssue(issueInfo models.IssueInfo) (int64, error) {
	var issueId, authorId, assigneeId int64

	row := resourceHandler.database.QueryRow("SELECT id FROM author WHERE name = ?", issueInfo.Author)
	if err := row.Scan(&authorId); err != nil {

		if err == sql.ErrNoRows {
			return issueId, fmt.Errorf("InsertIssue: no such author %s", issueInfo.Author)
		}
		return issueId, fmt.Errorf("InsertIssue: %v", err)
	}

	row = resourceHandler.database.QueryRow("SELECT id FROM author WHERE name = ?", issueInfo.Assignee)
	if err := row.Scan(&assigneeId); err != nil {

		if err == sql.ErrNoRows {
			return issueId, fmt.Errorf("InsertIssue: no such assignee %s", issueInfo.Assignee)
		}
		return issueId, fmt.Errorf("InsertIssue: %v", err)
	}

	result, err := resourceHandler.database.Exec(
		"INSERT INTO issues "+
			"(projectId, authorId, assigneeId,"+
			" key, summary, description, type, priority, status,"+
			" createdTime, closedTime, updatedTime, timeSpent)"+
			" values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)",
		issueInfo.Project.ID,
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
		return issueId, fmt.Errorf("insertIssue: %v", err)
	}
	issueId, err = result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("insertIssue: %v", err)
	}

	resourceHandler.logger.Log(logger.INFO, "insertIssue successfully")

	return issueId, nil
}

func (resourceHandler *ResourceHandler) getAllProjects(limit int) ([]models.ProjectInfo, error) {
	var projects []models.ProjectInfo

	rows, err := resourceHandler.database.Query(
		`SELECT 
    	projects.id, 
    	projects.title, 
		FROM 
    	projects
		GROUP BY 
    	projects.id, 
    	projects.title
		ORDER BY 
    	projects.id
		LIMIT 
    	$1`,
		limit,
	)
	if err != nil {
		return nil, fmt.Errorf("getAllProjects: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var project = models.ProjectInfo{}
		if err := rows.Scan(&project.ID, &project.Title); err != nil {
			return projects, fmt.Errorf("getAllProjects: %v", err)
		}
		projects = append(projects, project)
	}
	if err := rows.Err(); err != nil {
		return projects, fmt.Errorf("albumsByArtist: %v", err)
	}

	return projects, nil
}
