package main

import endpoints "Jira-analyzer/analyzer/endpoints/resource"

func main() {
	endpoints.CreateNewResourceHandler().StartServer()
}
