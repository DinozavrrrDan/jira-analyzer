package main

import (
	"Jira-analyzer/jiraConnector/apiServer"
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
)

func main() {
	apiServer := apiServer.CreateNewApiServer()
	apiServer.StrartServer()

	connectorTarget, err := url.Parse(fmt.Sprintf("http://%s:%d", "localhost", 8000))
	if err != nil {
		fmt.Printf("Error parsing target URL: %v\n", err)
	}

	gatewayMux := http.NewServeMux()
	fmt.Println("Create proxy for connector server.")
	gatewayMux.HandleFunc("/api/v1"+"/connector/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("Recieved request for connector proxy: %s", r.URL.Path)
		proxy := httputil.NewSingleHostReverseProxy(connectorTarget)
		proxy.ServeHTTP(w, r)
	})

}
