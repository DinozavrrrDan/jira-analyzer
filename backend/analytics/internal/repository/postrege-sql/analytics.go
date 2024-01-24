package postrege_sql

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/DinozvrrDan/jira-analyzer/backend/analytics/internal/models"
)

type AnalyticsRepository struct {
	db *sql.DB
}

func NewAnalyticsRepository(db *sql.DB) *AnalyticsRepository {
	return &AnalyticsRepository{
		db: db,
	}
}

func (analyticsRepository *AnalyticsRepository) GetGraphsOneData(title string) ([]models.GraphOneData, error) {
	var graph []models.GraphOneData
	id, _ := analyticsRepository.getProjectId(title)
	rows, err := analyticsRepository.db.Query("SELECT id, FLOOR((EXTRACT(EPOCH FROM (closedTime)) - EXTRACT(EPOCH FROM (createdTime))))::bigint AS openTime  FROM issues WHERE projectid = $1 and status IN ('Closed', 'Resolved') ORDER BY openTime", id)
	if err != nil {
		return nil, fmt.Errorf("GetGraphsFourData: %w", err)
	}
	for rows.Next() {
		var tmp models.GraphOneData
		var time int64
		err := rows.Scan(&tmp.Id, &time)
		if err != nil {
			return nil, fmt.Errorf("GetGraphsFourData: %w", err)
		}

		if time < 3600 {

		} else if time < 10800 {

		} else if time < 21600 {

		} else if time < 43200 {

		} else if time < 86400 {

		} else if time < 604800 {

		} else if time < 1209600 {

		} else if time < 2592000 {

		} else {

		}

		graph = append(graph, tmp)
	}
	return graph, nil
}

func (analyticsRepository *AnalyticsRepository) GetGraphsOneCategories(title string) ([]string, error) {
	var categories []string
	id, _ := analyticsRepository.getProjectId(title)
	rows, err := analyticsRepository.db.Query("SELECT id FROM issues WHERE projectid = $1", id)
	if err != nil {
		return nil, fmt.Errorf("GetGraphsOneData: %w", err)
	}
	for rows.Next() {
		var tmp string
		err := rows.Scan(&tmp)
		if err != nil {
			return nil, fmt.Errorf("GetGraphsOneCategories: %w", err)
		}
		categories = append(categories, tmp)
	}
	return categories, nil
}

func (analyticsRepository *AnalyticsRepository) GetGraphsFourData(title string) ([]models.GraphFourData, error) {
	var graph []models.GraphFourData
	id, _ := analyticsRepository.getProjectId(title)
	rows, err := analyticsRepository.db.Query("SELECT type, COUNT(*) AS count FROM issues WHERE projectid = $1 GROUP BY TYPE ORDER BY TYPE", id)
	if err != nil {
		return nil, fmt.Errorf("GetGraphsFourData: %w", err)
	}
	for rows.Next() {
		var tmp models.GraphFourData
		err := rows.Scan(&tmp.Type, &tmp.Count)
		if err != nil {
			return nil, fmt.Errorf("GetGraphsFourData: %w", err)
		}
		graph = append(graph, tmp)
	}
	return graph, nil
}

func (analyticsRepository *AnalyticsRepository) GetGraphsFourCategories(title string) ([]string, error) {
	var categories []string
	id, _ := analyticsRepository.getProjectId(title)
	rows, err := analyticsRepository.db.Query("SELECT type FROM issues WHERE projectid = $1 GROUP BY TYPE", id)
	if err != nil {
		return nil, fmt.Errorf("GetGraphsFourData: %w", err)
	}
	for rows.Next() {
		var tmp string
		err := rows.Scan(&tmp)
		if err != nil {
			return nil, fmt.Errorf("GetGraphsFourCategories: %w", err)
		}
		categories = append(categories, tmp)
	}
	return categories, nil
}

func (analyticsRepository *AnalyticsRepository) getProjectId(projectTitle string) (int64, error) {
	var projectId int64

	row := analyticsRepository.db.QueryRow("SELECT id FROM project where title = $1", projectTitle)

	if err := row.Scan(&projectId); err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return projectId, fmt.Errorf("getProjectId error: %w", err)
		}
	}

	return projectId, nil
}
