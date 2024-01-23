package service

import (
	"github.com/DinozvrrDan/jira-analyzer/connector/internal/models"
	"time"
)

type TransformerService struct {
}

func NewTransformerService() *TransformerService {
	return &TransformerService{}
}

func (dataTransformer *TransformerService) TransformData(issues []models.Issue) []models.TransformedIssue {
	var transformedIssues []models.TransformedIssue
	for _, issue := range issues {

		createdTime, _ := time.Parse(time.DateTime, issue.Fields.CreatedTime)
		closedTime, _ := time.Parse(time.DateTime, issue.Fields.ClosedTime)
		updatedTime, _ := time.Parse(time.DateTime, issue.Fields.UpdatedTime)

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
<<<<<<< HEAD:connector/internal/service/transformer.go
			TimeSpent:   int64(timeSpent),
=======
			Timespent:   int64(timeSpent),
>>>>>>> 77eb92a4ac4e65cecd3e27102698a8f5679caf6c:jiraConnector/transformer/transformer.go
		})
	}
	return transformedIssues
}
