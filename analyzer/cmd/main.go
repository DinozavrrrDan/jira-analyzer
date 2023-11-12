package main

import "Jira-analyzer/jiraConnector/connector"

func main() {
	jiraConnector := connector.CreateNewJiraConnector()
	jiraConnector.GetProjects()
	jiraConnector.GetProjectIssues("ACE")
}
