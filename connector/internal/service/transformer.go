package service

import (
	"connector/internal/models"
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
			Timespent:   int(timeSpent),
		})
	}
	return transformedIssues
}
