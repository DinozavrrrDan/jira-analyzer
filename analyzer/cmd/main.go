package main

import "Jira-analyzer/jiraConnector/connector"

func main() {
	jiraConnector := connector.NewConnector()
	jiraConnector.GetProjects()
	jiraConnector.GetProjectInfo("ACE")
}
