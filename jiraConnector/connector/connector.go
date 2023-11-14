package connector

import (
	"Jira-analyzer/jiraConnector/configReader"
	"Jira-analyzer/jiraConnector/logger"
	"Jira-analyzer/jiraConnector/models"
	"Jira-analyzer/jiraConnector/transformer"
	"encoding/json"
	"io"
	"net/http"
)

type Connector struct {
	logger            *logger.JiraLogger
	configReader      *configReader.ConfigRaeder
	jiraRepositoryUrl string
}

func CreateNewJiraConnector() *Connector {
	newReader := configReader.CreateNewConfigReader()
	return &Connector{
		logger:            logger.CreateNewLogger(),
		configReader:      newReader,
		jiraRepositoryUrl: newReader.GetJiraUrl(),
	}
}

/*
В случае удачного выполнения запроса должен быть возвращен JSON,
который содержит массив проектов и общее количество страниц при
данном параметре limit
*/
func (connector *Connector) GetProjectIssues(projectName string) {
	httpClient := &http.Client{}

	//temp
	projectName = "ACE" //Просто для примера имя

	response, err := httpClient.Get(connector.jiraRepositoryUrl + "/rest/api/2/search?jql=project=" + projectName + "&expand=changelog&startAt=0&maxResults=1")
	if err != nil || response.StatusCode != http.StatusOK {
		connector.logger.Log(logger.ERROR, "Error with get response from: ")
		return
	}

	body, err := io.ReadAll(response.Body)
	var issueResponce models.IssuesList
	err = json.Unmarshal(body, &issueResponce)
	if err != nil {
		connector.logger.Log(logger.ERROR, " ")
		return
	}

	transformer.TrasformData(issueResponce)

}

/*
Параметр limit - сколько всего проектов необходимо вернуть
Параметр page - порядковый номер страницы, который необходимо
вернуть
Параметр search - фильтр, который накладывается на название и ключ
*/
func (connector *Connector) GetProjects(limit int, page int, search string) models.Projects {
	httpClient := &http.Client{}
	resp, err := httpClient.Get(connector.jiraRepositoryUrl + "/rest/api/2/project")
	if err != nil {
		connector.logger.Log(logger.ERROR, "Error with get response from about projects ")
		return models.Projects{}
	}

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		connector.logger.Log(logger.ERROR, " ")
		return models.Projects{}
	}

	var jiraProjects []models.JiraProject
	err = json.Unmarshal(body, &jiraProjects) //получаем информацию через сериализацию

	if err != nil {
		connector.logger.Log(logger.ERROR, " ")
		return models.Projects{}

	}
	var projects []models.Project

	counterOfProjects := 0

	//Получение информации о определенном колчичестве проектов
	for _, element := range jiraProjects {
		//Понять зачем search
		counterOfProjects++
		projects = append(projects, models.Project{
			Name: element.Name,
			Link: element.Link,
			Key:  element.Key,
		})
	}

	//обрезка проектов по странице

	startIndexOfProject := limit * (page - 1)
	endIndexOfProject := limit * page
	//подумать над косяками

	return models.Projects{
		Projects: projects[startIndexOfProject:endIndexOfProject],
		Page: models.Page{
			TotalPageCount:     int(counterOfProjects / limit),
			CurrentPageNumber:  page,
			TotalProjectsCount: counterOfProjects,
		},
	}
}
