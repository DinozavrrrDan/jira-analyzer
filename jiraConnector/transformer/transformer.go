package transformer

import (
	"Jira-analyzer/jiraConnector/models"
	"fmt"
)

func TrasformData(issues models.IssuesList) {
	var trIss []models.TransformedIssue
	trIss = append(trIss, models.TransformedIssue{
		Project:     issues.Issues[0].Fields.Project.Name,
		Author:      issues.Issues[0].Fields.Creator.Name,
		Assignee:    issues.Issues[0].Fields.AssigneeName.Name,
		Key:         issues.Issues[0].Key,
		Summary:     issues.Issues[0].Fields.Summary,
		Description: issues.Issues[0].Fields.Description,
		Type:        issues.Issues[0].Fields.Type.Name,
		Priority:    issues.Issues[0].Fields.Priority.Name,
		Status:      issues.Issues[0].Fields.Status.Name,
	})

	fmt.Println("1: Project:     " + trIss[0].Project)
	fmt.Println("2: Author:      " + trIss[0].Author)
	fmt.Println("3: Assignee:    " + trIss[0].Assignee)
	fmt.Println("4: Key:         " + trIss[0].Key)
	fmt.Println("5: Summary:     " + trIss[0].Summary)
	fmt.Println("6: Description: " + trIss[0].Description)
	fmt.Println("7: Type:        " + trIss[0].Type)
	fmt.Println("8: Priority:    " + trIss[0].Priority)
	fmt.Println("9: Status:      " + trIss[0].Status)
}
