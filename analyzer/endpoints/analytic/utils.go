package analytic

import (
	"Jira-analyzer/analyzer/models"
	"Jira-analyzer/common/logger"
	"fmt"
)

func (analyticServer *AnalyticServer) GraphFive(projectId int64) ([]models.GraphData, error) {
	var result []models.GraphData

	rows, err := analyticServer.database.Query("SELECT" +
		" priority," +
		" COUNT(*) AS task_count" +
		" FROM issues" +
		" GROUP BY priority" +
		" CASE" +
		" WHEN priority = 'Critical' THEN 1" +
		" WHEN priority = 'Blocker' THEN 2" +
		" WHEN priority = 'Major' THEN 3" +
		" WHEN priority = 'Minor' THEN 4" +
		" END;")
	if err != nil {
		return nil, fmt.Errorf("GraphFive: select project info %d: %v", projectId, err)
	}
	defer rows.Close()

	for rows.Next() {
		var projectInfo models.GraphData
		if err := rows.Scan(&projectInfo.PriorityType, &projectInfo.Amount); err != nil {
			return nil, fmt.Errorf("GraphFive with projectId %d: %v", projectId, err)
		}
		result = append(result, projectInfo)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("GraphFive with projectId %d: %v", projectId, err)
	}

	analyticServer.logger.Log(logger.INFO, "Successfully GraphFive")
	return result, nil
}
