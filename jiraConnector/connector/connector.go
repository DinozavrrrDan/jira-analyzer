package connector

import (
	"Jira-analyzer/jiraConnector/models"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type TransformedIssue struct {
	Project     string
	Author      string
	Assignee    string
	Key         string
	Summary     string
	Description string
	Type        string
	Priority    string
	Status      string
	//CreatedTime time.Time
	//ClosedTime  time.Time
	//UpdatedTime time.Time
	//Timespent   int64
}

func main() {
	getProjects()
	getProjectInfo()
}

func getProjectInfo() {
	httpClient := &http.Client{}
	projectName := "AAR"
	resp, err := httpClient.Get("http://issues.apache.org/jira/rest/api/2/search?jql=project=" + projectName + "&expand=changelog&startAt=0&maxResults=1")
	if err != nil {
		fmt.Print(err)
		return
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)

	var issueResp models.IssuesList
	err = json.Unmarshal(body, &issueResp)
	//fmt.Print(issueResp)

	var trIss []TransformedIssue
	trIss = append(trIss, TransformedIssue{
		Project:     issueResp.Issues[0].Fields.Project.Name,
		Author:      issueResp.Issues[0].Fields.Creator.Name,
		Assignee:    issueResp.Issues[0].Fields.AssigneeName.Name,
		Key:         issueResp.Issues[0].Key,
		Summary:     issueResp.Issues[0].Fields.Summary,
		Description: issueResp.Issues[0].Fields.Description,
		Type:        issueResp.Issues[0].Fields.Type.Name,
		Priority:    issueResp.Issues[0].Fields.Priority.Name,
		Status:      issueResp.Issues[0].Fields.Status.Name,
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
	//	fmt.Print(string(body))

	// if err != nil {
	// 	fmt.Print(err)
	// 	return
	// }

	// var issueFields IssueFields
	// err = json.Unmarshal(body, &issueFields)

}

func getProjects() {
	httpClient := &http.Client{}
	resp, err := httpClient.Get("http://issues.apache.org/jira/rest/api/2/project")

	if err != nil {
		fmt.Print(err)
		return
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Print(err)
		return
	}
	var jiraProjects []models.JiraProject
	err = json.Unmarshal(body, &jiraProjects) //получаем информацию через сериализацию

	var projects []models.Project

	counterOfProjects := 0

	for _, element := range jiraProjects {
		counterOfProjects++
		projects = append(projects, models.Project{
			Name: element.Name,
			Link: element.Link,
		})
		if counterOfProjects == 5 {
			break
		}
	}

	for i := 0; i < counterOfProjects; i++ {
		fmt.Println(projects[i].Name + "  \t:  " + projects[i].Link)
	}
}
