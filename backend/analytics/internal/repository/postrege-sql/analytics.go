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
	var count int

	times := [9]int{0, 3600, 10800, 21600, 43200, 86400, 604800, 1209600, 2592000}
	timesString := [8]string{"1 hour", "3 hours", "6 hours", "12 hour", "24 hour", "3 days", "7 days", "14 days"}

	for i := 0; i < 8; i++ {
		currentTimeStart := times[i]
		currentTimeEnd := times[i+1]

		err := analyticsRepository.db.QueryRow("SELECT COUNT(*) AS count FROM issues "+
			" WHERE projectid = $1 and status IN ('Closed', 'Resolved') and timespent > $2 and timespent < $3",
			id, currentTimeStart, currentTimeEnd).
			Scan(&count)

		if err != nil {
			fmt.Println("rows err")
			return nil, fmt.Errorf("GetGraphsFourData: %w", err)
		}

		graph = append(graph, models.GraphOneData{
			SpentTime: timesString[i],
			Count:     count,
		})

	}

	err := analyticsRepository.db.QueryRow("SELECT COUNT(*) AS spentTime FROM issues "+
		" WHERE projectid = $1 and status IN ('Closed', 'Resolved') and timespent > $2", id, times[8]).Scan(&count)
	if err != nil {
		return nil, fmt.Errorf("GetGraphsFourData: %w", err)
	}

	graph = append(graph, models.GraphOneData{
		SpentTime: "more 30 days",
		Count:     count,
	})

	return graph, nil
}

func (analyticsRepository *AnalyticsRepository) GetGraphsOneCategories(title string) ([]string, error) {
	var categories []string
	id, _ := analyticsRepository.getProjectId(title)
	rows, err := analyticsRepository.db.Query("SELECT id FROM issues WHERE projectid = $1", id)
	if err != nil {
		return nil, fmt.Errorf("GetGraphsOneData: %w", err)
	}
	defer rows.Close()

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
	rows, err := analyticsRepository.db.Query("SELECT type, COUNT(*) AS count FROM issues "+
		"WHERE projectid = $1 GROUP BY TYPE ORDER BY TYPE", id)
	if err != nil {
		return nil, fmt.Errorf("GetGraphsFourData: %w", err)
	}
	defer rows.Close()

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
		return nil, fmt.Errorf("GetGraphsFourCategories: %w", err)
	}
	defer rows.Close()

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

func (analyticsRepository *AnalyticsRepository) GetGraphsFiveData(title string) ([]models.GraphFiveData, error) {
	var graph []models.GraphFiveData
	id, _ := analyticsRepository.getProjectId(title)
	rows, err := analyticsRepository.db.Query("SELECT priority, COUNT(*) AS count FROM issues "+
		"WHERE projectid = $1 GROUP BY PRIORITY ORDER BY PRIORITY", id)
	if err != nil {
		return nil, fmt.Errorf("GetGraphsFiveData: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var tmp models.GraphFiveData
		err := rows.Scan(&tmp.Priority, &tmp.Count)
		if err != nil {
			return nil, fmt.Errorf("GetGraphsFiveData: %w", err)
		}
		graph = append(graph, tmp)
	}
	return graph, nil
}

func (analyticsRepository *AnalyticsRepository) GetGraphsFiveCategories(title string) ([]string, error) {
	var categories []string
	id, _ := analyticsRepository.getProjectId(title)
	rows, err := analyticsRepository.db.Query("SELECT priority FROM issues WHERE projectid = $1 GROUP BY priority",
		id)
	if err != nil {
		return nil, fmt.Errorf("GetGraphsFiveCategories: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var tmp string
		err := rows.Scan(&tmp)
		if err != nil {
			return nil, fmt.Errorf("GetGraphsFiveCategories: %w", err)
		}
		categories = append(categories, tmp)
	}
	return categories, nil
}

func (analyticsRepository *AnalyticsRepository) GetGraphsSixData(title string) ([]models.GraphSixData, error) {
	var graph []models.GraphSixData
	id, _ := analyticsRepository.getProjectId(title)
	rows, err := analyticsRepository.db.Query("SELECT priority, COUNT(*) AS count FROM issues"+
		" WHERE projectid = $1 and status IN ('Closed', 'Resolved') GROUP BY PRIORITY ORDER BY PRIORITY", id)
	if err != nil {
		return nil, fmt.Errorf("GetGraphsSixData: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var tmp models.GraphSixData
		err := rows.Scan(&tmp.Priority, &tmp.Count)
		if err != nil {
			return nil, fmt.Errorf("GetGraphsSixData: %w", err)
		}
		graph = append(graph, tmp)
	}
	return graph, nil
}

func (analyticsRepository *AnalyticsRepository) GetGraphsSixCategories(title string) ([]string, error) {
	var categories []string
	id, _ := analyticsRepository.getProjectId(title)
	rows, err := analyticsRepository.db.Query("SELECT priority FROM issues "+
		"WHERE projectid = $1 and status IN ('Closed', 'Resolved') GROUP BY priority", id)
	if err != nil {
		return nil, fmt.Errorf("GetGraphsSixCategories: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var tmp string
		err := rows.Scan(&tmp)
		if err != nil {
			return nil, fmt.Errorf("GetGraphsSixCategories: %w", err)
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
