package resource

import (
	"Jira-analyzer/jiraConnector/configReader"
	"Jira-analyzer/jiraConnector/logger"
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type ResourseHandler struct {
	configReader *configReader.ConfigReader
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

func (resourseHandler *ResourseHandler) HandlePostIssue(responseWriter http.ResponseWriter, request *http.Request) {
	body, err := io.ReadAll(request.Body)
	if err != nil {
		return
	}

	var requestDataIssue models.DataInfoIssue
	err = json.Unmarshal(body, &requestDataIssue)

	if err != nil {
		resourseHandler.logger.Log(logger.ERROR, "error when encoding request")
		return
	}
	var statusCode int
	//id, err := PutIssueIntoDB функция которая будет помещать узел в БД БОРЯ
	if err != nil {
		//как-то напишем об ошибке
		//statusCode = http.Status - подобрать верный статус
	} else {
		//statusCode = http.Status - подобрать верный статус
	}
	statusCode = statusCode + 1 // заглушка
	response, err := json.Marshal(models.ResponseStrucrt{})
	if err != nil {

	}
	_, err = responseWriter.Write(response)
	if err != nil {

	}
}

func (resourseHandler *ResourseHandler) HandlePostHistory(responseWriter http.ResponseWriter, request *http.Request) {
	body, err := io.ReadAll(request.Body)
	if err != nil {
		return
	}

	var requestDataIssue models.DataInfoHistory
	err = json.Unmarshal(body, &requestDataIssue)

	if err != nil {
		resourseHandler.logger.Log(logger.ERROR, "error when encoding request")
		return
	}
	var statusCode int
	//id, err := PutIssueIntoDB функция которая будет помещать узел в БД БОРЯ
	if err != nil {
		//как-то напишем об ошибке
		//statusCode = http.Status - подобрать верный статус
	} else {
		//statusCode = http.Status - подобрать верный статус
	}
	statusCode = statusCode + 1 // заглушка
	response, err := json.Marshal(models.ResponseStrucrt{})
	if err != nil {

	}
	_, err = responseWriter.Write(response)
	if err != nil {

	}
}

func (resourseHandler *ResourseHandler) HandlePostProject(responseWriter http.ResponseWriter, request *http.Request) {
	body, err := io.ReadAll(request.Body)
	if err != nil {
		return
	}

	var requestDataIssue models.DataInfoProject
	err = json.Unmarshal(body, &requestDataIssue)

	if err != nil {
		resourseHandler.logger.Log(logger.ERROR, "error when encoding request")
		return
	}
	var statusCode int
	//id, err := PutIssueIntoDB функция которая будет помещать узел в БД БОРЯ
	if err != nil {
		//как-то напишем об ошибке
		//statusCode = http.Status - подобрать верный статус
	} else {
		//statusCode = http.Status - подобрать верный статус
	}
	statusCode = statusCode + 1 // заглушка
	response, err := json.Marshal(models.ResponseStrucrt{})
	if err != nil {

	}
	_, err = responseWriter.Write(response)
	if err != nil {

	}
}
