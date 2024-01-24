package service

import (
	"fmt"
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

		createdTime, _ := time.Parse("2006-01-02T15:04:05.999-0700", issue.Fields.CreatedTime)
		closedTime, _ := time.Parse("2006-01-02T15:04:05.999-0700", issue.Fields.ClosedTime)
		updatedTime, _ := time.Parse("2006-01-02T15:04:05.999-0700", issue.Fields.UpdatedTime)
		if len(issue.Fields.ClosedTime) == 0 {
			closedTime = createdTime.Add(time.Hour * 24 * 365 * 10)
		}
		timespent := closedTime.Sub(createdTime)
		fmt.Print(timespent.Seconds())
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
<<<<<<< HEAD:jiraConnector/transformer/transformer.go
			Timespent:   int64(timeSpent),
=======
			TimeSpent:   int64(timespent.Seconds()),
>>>>>>> 6a15cb1650a9c1e304607e8f9b48d77b20ebf674:connector/internal/service/transformer.go
		})
	}
	return transformedIssues
}
