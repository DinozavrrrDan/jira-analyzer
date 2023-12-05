package service

import (
	"connector/internal/entities"
	"time"
)

type TransformerService struct {
}

func NewTransformerService() *TransformerService {
	return &TransformerService{}
}

func (dataTransformer *TransformerService) TransformData(issues []entities.Issue) []entities.TransformedIssue {
	var transformedIssues []entities.TransformedIssue
	for _, issue := range issues {

		createdTime, _ := time.Parse("2006-01-02T15:04:05.999-0700", issue.Fields.CreatedTime)
		closedTime, _ := time.Parse("2006-01-02T15:04:05.999-0700", issue.Fields.ClosedTime)
		updatedTime, _ := time.Parse("2006-01-02T15:04:05.999-0700", issue.Fields.UpdatedTime)
		timeSpent := closedTime.Sub(closedTime)
		transformedIssues = append(transformedIssues, entities.TransformedIssue{
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
