package main

import (
	"Jira-analyzer/jiraConnector/apiServer"
)

func main() {
	apiServer.CreateNewApiServer().StartServer()
}
