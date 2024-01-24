package v1

import (
	"fmt"
	"github.com/DinozvrrDan/jira-analyzer/connector/pkg/logger"
	"net/http"
)

func errorWriter(w http.ResponseWriter, log *logger.Logger, message string, status int) {
	fmt.Println(message)
	log.Log(logger.ERROR, message)
	w.WriteHeader(status)
}
