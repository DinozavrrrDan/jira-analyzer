package main

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
)

func main() {

	connectorTarget, err := url.Parse(fmt.Sprintf("http://localhost:8003"))
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

	err = http.ListenAndServe("localhost:8000", gatewayMux)

}
