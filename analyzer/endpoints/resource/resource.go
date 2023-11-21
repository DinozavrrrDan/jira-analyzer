package resource

import (
	"Jira-analyzer/jiraConnector/configReader"
	"Jira-analyzer/jiraConnector/logger"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type ResourseHandler struct {
	configReader *configReader.ConfigRaeder
	logger       *logger.JiraLogger
}

func CreateNewResourseHandler() *ResourseHandler {
	return &ResourseHandler{
		configReader: configReader.CreateNewConfigReader(),
		logger:       logger.CreateNewLogger(),
	}
}

func (resourseHandler *ResourseHandler) HandleGetIssue(responseWriter http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		resourseHandler.logger.Log(logger.INFO, "Invalid ID!") //Подумать над уровнем логирования
		return
	}
	//Фунция которую реализует БОРЯ
	id = id + 1 //заглушка
	//issue, err := ВОТ ТУТ ДОЛЖНА БЫТЬ ДЛЯ ISSUES
	if err != nil {
		return
	}
	//работа с полученной инфой и формирование ответа
}

func (resourseHandler *ResourseHandler) HandleGetHistory(responseWriter http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		resourseHandler.logger.Log(logger.INFO, "Invalid ID!") //Подумать над уровнем логирования
		return
	}
	//Фунция которую реализует БОРЯ
	id = id + 1 //заглушка
	//issue, err := ВОТ ТУТ ДОЛЖНА БЫТЬ ДЛЯ ИСТОРИИ
	if err != nil {
		return
	}
	//работа с полученной инфой и формирование ответа
}

func (resourseHandler *ResourseHandler) HandleGetProject(responseWriter http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		resourseHandler.logger.Log(logger.INFO, "Invalid ID!") //Подумать над уровнем логирования
		return
	}
	//Фунция которую реализует БОРЯ
	id = id + 1 //заглушка
	//issue, err := ВОТ ТУТ ДОЛЖНА БЫТЬ ДЛЯ ПРОЕКТОВ
	if err != nil {
		return
	}
	//работа с полученной инфой и формирование ответа
}
