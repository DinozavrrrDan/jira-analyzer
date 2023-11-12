package connector

import (
	"Jira-analyzer/jiraConnector/configReader"
	"Jira-analyzer/jiraConnector/logger"
	"Jira-analyzer/jiraConnector/models"
	"Jira-analyzer/jiraConnector/transformer"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Connector struct {
	logger            *logger.JiraLogger
	configReader      *configReader.ConfigRaeder
	jiraRepositoryUrl string
}

func NewConnector() *Connector {
	newReader := configReader.CreateNewConfigReader()
	return &Connector{
		logger:            logger.CreateNewLogger(),
		configReader:      newReader,
		jiraRepositoryUrl: newReader.GetJiraUrl(),
	}
}

func (connector *Connector) GetProjectInfo(projectName string) {
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

func (connector *Connector) GetProjects() {
	httpClient := &http.Client{}

	resp, err := httpClient.Get(connector.jiraRepositoryUrl + "/rest/api/2/project")
	if err != nil {
		connector.logger.Log(logger.ERROR, "Error with get response from about projects ")
		return
	}

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		connector.logger.Log(logger.ERROR, " ")
		return
	}
	var jiraProjects []models.JiraProject
	err = json.Unmarshal(body, &jiraProjects) //получаем информацию через сериализацию
	if err != nil {
		connector.logger.Log(logger.ERROR, " ")
		return
	}
	var projects []models.Project

	counterOfProjects := 0

	//Получение информации о определенном колчичестве проектов
	for _, element := range jiraProjects {
		counterOfProjects++
		projects = append(projects, models.Project{
			Name: element.Name,
			Link: element.Link,
		})
		if counterOfProjects == 5 {
			break
		}
	}

	for i := 0; i < counterOfProjects; i++ {
		fmt.Println(projects[i].Name + "  \t:  " + projects[i].Link)
	}
}
