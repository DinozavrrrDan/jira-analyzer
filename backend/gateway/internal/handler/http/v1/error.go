package v1

import (
	"github.com/DinozvrrDan/jira-analyzer/backend/gateway/pkg/logger"
	"net/http"
)

func errorWriter(w http.ResponseWriter, log *logger.Logger, message string, status int) {
	log.Log(logger.ERROR, message)
	w.WriteHeader(status)
}
