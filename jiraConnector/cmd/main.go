package main

import (
	"Jira-analyzer/jiraConnector/apiServer"
)

func main() {
	err := apiServer.CreateNewApiServer().StartServer()
	if err != nil {
		panic(err)
	}
}
