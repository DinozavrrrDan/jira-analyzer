package analytic

import (
	"Jira-analyzer/analyzer/models"
	"Jira-analyzer/common/logger"
	"fmt"
)

func (analyticServer *AnalyticServer) GraphFive(projectId int64) ([]models.GraphData, error) {
	result := make([]models.GraphData, 0)

	const PriorityCount = "SELECT " +
		"priority, COUNT(*) AS task_count " +
		"FROM issues " +
		"GROUP BY priority "

	rows, err := analyticServer.database.Query(PriorityCount)
	if err != nil {
		return nil, fmt.Errorf("GraphFive: select project info %d: %w", projectId, err)
	}
	defer rows.Close()

	for rows.Next() {
		var projectInfo models.GraphData
		if err := rows.Scan(&projectInfo.PriorityType, &projectInfo.Amount); err != nil {
			return nil, fmt.Errorf("GraphFive with projectId %d: %w", projectId, err)
		}
		result = append(result, projectInfo)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("GraphFive with projectId %d: %w", projectId, err)
	}

	analyticServer.logger.Log(logger.INFO, "Successfully GraphFive")
	return result, nil
}
