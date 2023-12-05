package analytic

import (
	"Jira-analyzer/analyzer/models"
	"Jira-analyzer/common/logger"
	"fmt"
	"log"
)

func (analyticServer *AnalyticServer) GraphFive(projectId int64) []models.GraphData {

	rows, err := analyticServer.database.Query(
		"SELECT" +
			" priority," +
			" COUNT(*) AS task_count" +
			" FROM issues" +
			" GROUP BY priority" +
			" CASE" +
			" WHEN priority = 'Critical' THEN 1" +
			" WHEN priority = 'Blocker' THEN 2" +
			" WHEN priority = 'Major' THEN 3" +
			" WHEN priority = 'Minor' THEN 4" +
			" END;",
	)

	if err != nil {
		analyticServer.logger.Log(logger.ERROR, fmt.Sprintf("GraphFive: %v", err))
		return nil
	}

	var result []models.GraphData

	for rows.Next() {
		var entry models.GraphData

		err := rows.Scan(&entry.PriorityType, &entry.Amount)
		if err != nil {
			analyticServer.logger.Log(logger.ERROR, err.Error())
			log.Fatal(err)
		}

		result = append(result, entry)
	}

	analyticServer.logger.Log(logger.INFO, "Successfully GraphFive")
	return result
}
