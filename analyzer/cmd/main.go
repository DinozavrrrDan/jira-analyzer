package main

import "Jira-analyzer/jiraConnector/connector"

func main() {
	connector.GetProjects()
	connector.GetProjectInfo("AAR")
}
