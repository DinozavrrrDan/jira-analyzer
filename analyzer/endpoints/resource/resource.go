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

func (resourseHandler *ResourseHandler) HandleGetIssue(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		resourseHandler.logger.Log(logger.ERROR, err.Error())
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	issue, err := GetIssueInfoByID(id)
	if err != nil {
		resourseHandler.logger.Log(logger.ERROR, err.Error())
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	project, err := GetProjectInfoByID(issue.ProjectID)
	if err != nil {
		resourseHandler.logger.Log(logger.ERROR, err.Error())
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	resourseHandler.logger.Log(logger.INFO, "HandleGetIssue successfully")
	rw.WriteHeader(http.StatusOK)
}

func (resourseHandler *ResourseHandler) HandleGetHistory(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		resourseHandler.logger.Log(logger.ERROR, err.Error())
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	history, err := GetAllHistoryInfoByIssueID(id)
	if err != nil {
		resourseHandler.logger.Log(logger.ERROR, err.Error())
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	resourseHandler.logger.Log(logger.INFO, "HandleGetHistory successfully")
	rw.WriteHeader(http.StatusOK)
}

func (resourseHandler *ResourseHandler) HandleGetProject(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		resourseHandler.logger.Log(logger.ERROR, err.Error())
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	project, err := GetProjectInfoByID(id)
	if err != nil {
		resourseHandler.logger.Log(logger.ERROR, err.Error())
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	resourseHandler.logger.Log(logger.INFO, "HandleGetProject successfully")
	rw.WriteHeader(http.StatusOK)
}

func (resourseHandler *ResourseHandler) HandlePostIssue(rw http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		resourseHandler.logger.Log(logger.ERROR, err.Error())
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	var requestDataIssue models.IssueInfo
	err = json.Unmarshal(body, &requestDataIssue)

	if err != nil {
		resourseHandler.logger.Log(logger.ERROR, err.Error())
		http.Error(rw, err.Error(), http.StatusBadRequest)
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
	_, err = rw.Write(response)
	if err != nil {

	}
}

func (resourseHandler *ResourseHandler) HandlePostHistory(rw http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		resourseHandler.logger.Log(logger.ERROR, err.Error())
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	var requestDataIssue models.HistoryInfo
	err = json.Unmarshal(body, &requestDataIssue)

	if err != nil {
		resourseHandler.logger.Log(logger.ERROR, err.Error())
		http.Error(rw, err.Error(), http.StatusBadRequest)
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
		resourseHandler.logger.Log(logger.ERROR, err.Error())
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}
	_, err = responseWriter.Write(response)
	if err != nil {
		resourseHandler.logger.Log(logger.ERROR, err.Error())
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}
}

func (resourseHandler *ResourseHandler) HandlePostProject(rw http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		resourseHandler.logger.Log(logger.ERROR, err.Error())
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	var requestDataIssue models.ProjectInfo
	err = json.Unmarshal(body, &requestDataIssue)

	if err != nil {
		resourseHandler.logger.Log(logger.ERROR, err.Error())
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}
	var statusCode int
	//id, err := PutIssueIntoDB функция которая будет помещать узел в БД БОРЯ
	if err != nil {
		if err != nil {
			resourseHandler.logger.Log(logger.ERROR, err.Error())
			http.Error(rw, err.Error(), http.StatusBadRequest)
			return
		}
	} else {
		//statusCode = http.Status - подобрать верный статус
	}
	statusCode = statusCode + 1 // заглушка
	response, err := json.Marshal(models.ResponseStrucrt{})
	if err != nil {
		resourseHandler.logger.Log(logger.ERROR, err.Error())
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}
	_, err = rw.Write(response)
	if err != nil {
		resourseHandler.logger.Log(logger.ERROR, err.Error())
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}
}
