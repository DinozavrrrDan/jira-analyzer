package transformer

import (
	"Jira-analyzer/jiraConnector/models"
	"fmt"
	"time"
)

type Transformer struct {
}

func CreateNewTransformer() *Transformer {
	return &Transformer{}
}

func (dataTransformer *Transformer) TrasformData(issues []models.Issue) []models.TransformedIssue {
	var transformedIssues []models.TransformedIssue
	for _, issue := range issues {

		createdTime, _ := time.Parse("2006-01-02T15:04:05.999-0700", issue.Fields.CreatedTime)
		closedTime, _ := time.Parse("2006-01-02T15:04:05.999-0700", issue.Fields.ClosedTime)
		updatedTime, _ := time.Parse("2006-01-02T15:04:05.999-0700", issue.Fields.UpdatedTime)
		//timeSpent := closedTime - closedTime;
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
			Timespent:   0,
		})
	}
	fmt.Println("1: Project:     " + transformedIssues[0].Project)
	fmt.Println("2: Author:      " + transformedIssues[0].Author)
	fmt.Println("3: Assignee:    " + transformedIssues[0].Assignee)
	fmt.Println("4: Key:         " + transformedIssues[0].Key)
	fmt.Println("5: Summary:     " + transformedIssues[0].Summary)
	fmt.Println("6: Description: " + transformedIssues[0].Description)
	fmt.Println("7: Type:        " + transformedIssues[0].Type)
	fmt.Println("8: Priority:    " + transformedIssues[0].Priority)
	fmt.Println("9: Status:      " + transformedIssues[0].Status)
	fmt.Println("10: CreatedTime:      " + transformedIssues[0].CreatedTime.String())
	fmt.Println("11: ClosedTime:      " + transformedIssues[0].ClosedTime.String())
	fmt.Println("12: UpdatedTime:      " + transformedIssues[0].UpdatedTime.String())
	return transformedIssues
}
