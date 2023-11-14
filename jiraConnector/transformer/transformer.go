package transformer

import (
	"Jira-analyzer/jiraConnector/models"
	"fmt"
	"time"
)

func TrasformData(issues models.IssuesList) {

	createdTime, _ := time.Parse("2006-01-02T15:04:05.999-0700", issues.Issues[0].Fields.CreatedTime)
	closedTime, _ := time.Parse("2006-01-02T15:04:05.999-0700", issues.Issues[0].Fields.ClosedTime)
	updatedTime, _ := time.Parse("2006-01-02T15:04:05.999-0700", issues.Issues[0].Fields.UpdatedTime)

	var transormedIssue []models.TransformedIssue
	transormedIssue = append(transormedIssue, models.TransformedIssue{
		Project:     issues.Issues[0].Fields.Project.Name,
		Author:      issues.Issues[0].Fields.Creator.Name,
		Assignee:    issues.Issues[0].Fields.AssigneeName.Name,
		Key:         issues.Issues[0].Key,
		Summary:     issues.Issues[0].Fields.Summary,
		Description: issues.Issues[0].Fields.Description,
		Type:        issues.Issues[0].Fields.Type.Name,
		Priority:    issues.Issues[0].Fields.Priority.Name,
		Status:      issues.Issues[0].Fields.Status.Name,
		CreatedTime: createdTime,
		ClosedTime:  closedTime,
		UpdatedTime: updatedTime,
	})
	fmt.Println("1: Project:     " + transormedIssue[0].Project)
	fmt.Println("2: Author:      " + transormedIssue[0].Author)
	fmt.Println("3: Assignee:    " + transormedIssue[0].Assignee)
	fmt.Println("4: Key:         " + transormedIssue[0].Key)
	fmt.Println("5: Summary:     " + transormedIssue[0].Summary)
	fmt.Println("6: Description: " + transormedIssue[0].Description)
	fmt.Println("7: Type:        " + transormedIssue[0].Type)
	fmt.Println("8: Priority:    " + transormedIssue[0].Priority)
	fmt.Println("9: Status:      " + transormedIssue[0].Status)
	fmt.Println("10: Status:      " + transormedIssue[0].CreatedTime.String())
	fmt.Println("11: Status:      " + transormedIssue[0].ClosedTime.String())
	fmt.Println("12: Status:      " + transormedIssue[0].UpdatedTime.String())
}
