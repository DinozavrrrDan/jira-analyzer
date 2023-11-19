package resource

import (
	"Jira-analyzer/jiraConnector/configReader"
	"net/http"
)

type ResourseHandler struct {
	configReader *configReader.ConfigRaeder
}

func CreateNewResourseHandler() *ResourseHandler {
	return &ResourseHandler{
		configReader: configReader.CreateNewConfigReader(),
	}
}

func (resourseHandler *ResourseHandler) HandlerGetIssue(responseWriter http.ResponseWriter, request *http.Request) {

}

func (resourseHandler *ResourseHandler) HandlerGetHistory(responseWriter http.ResponseWriter, request *http.Request) {

}

func (resourseHandler *ResourseHandler) HandlerGetProject(responseWriter http.ResponseWriter, request *http.Request) {

}
