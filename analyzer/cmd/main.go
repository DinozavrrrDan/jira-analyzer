package main

import (
	"Jira-analyzer/jiraConnector/connector"
	"fmt"
)

func main() {
	jiraConnector := connector.CreateNewJiraConnector()
	projects := jiraConnector.GetProjects(5, 1, "")
	fmt.Print(projects)
	//jiraConnector.GetProjectIssues("ACE")
}
