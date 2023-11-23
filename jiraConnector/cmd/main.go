package main

import (
	//	"Jira-analyzer/jiraConnector/apiServer"
	"Jira-analyzer/jiraConnector/connector"
)

func main() {
	//	apiServer.CreateNewApiServer().StartServer()
	con := connector.CreateNewJiraConnector()
	con.GetProjects(1, 1, "")
}
