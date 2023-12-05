package transformer

import (
	"Jira-analyzer/jiraConnector/models"
	"time"
)

type Transformer struct {
}

func CreateNewTransformer() *Transformer {
	return &Transformer{}
}

func (dataTransformer *Transformer) TransformData(issues []models.Issue) []models.TransformedIssue {
	var transformedIssues []models.TransformedIssue
	for _, issue := range issues {

		createdTime, _ := time.Parse("2006-01-02T15:04:05.999-0700", issue.Fields.CreatedTime)
		closedTime, _ := time.Parse("2006-01-02T15:04:05.999-0700", issue.Fields.ClosedTime)
		updatedTime, _ := time.Parse("2006-01-02T15:04:05.999-0700", issue.Fields.UpdatedTime)
		timeSpent := closedTime.Sub(closedTime)
		transformedIssues = append(transformedIssues, models.TransformedIssue{
			Project:     issue.Fields.Project.Name,
			Author:      issue.Fields.Creator.Name,
			Assignee:    issue.Fields.AssigneeName.Name,
			Key:         issue.Key,
			Summary:     issue.Fields.Summary,
			Description: issue.Fields.Description,
			Type:        issue.Fields.Type.Name,
			Priority:    issue.Fields.Priority.Name,
			Status:      issue.Fields.Status.Name,
			CreatedTime: createdTime,
			ClosedTime:  closedTime,
			UpdatedTime: updatedTime,
			Timespent:   int64(timeSpent),
		})
	}
	return transformedIssues
}
