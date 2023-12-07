package main

import endpoints "Jira-analyzer/analyzer/endpoints/resource"

func main() {
	err := endpoints.CreateNewResourceHandler().StartServer()
	if err != nil {
		panic(err)
	}
}
