package connector

import (
	"Jira-analyzer/jiraConnector/configReader"
	"Jira-analyzer/jiraConnector/logger"
	"Jira-analyzer/jiraConnector/models"
	"encoding/json"
	"time"

	"io"
	"math"
	"net/http"
	"strconv"
	"sync"
)

type Connector struct {
	logger            *logger.JiraLogger
	configReader      *configReader.ConfigReader
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
func (connector *Connector) GetProjectIssues(projectName string) []models.Issue {
	isRequestGompleteSuccsesfully := false
	timeUntilNewRequest := connector.configReader.GetMinTimeSleep()
	var issues []models.Issue

	for isRequestGompleteSuccsesfully || timeUntilNewRequest <= connector.configReader.GetMaxTimeSleep() {

		httpClient := &http.Client{}

		response, err := httpClient.Get(connector.jiraRepositoryUrl + "/rest/api/2/search?jql=project=" + projectName + "&expand=changelog&startAt=0&maxResults=1")
		if err != nil || response.StatusCode != http.StatusOK {
			connector.logger.Log(logger.ERROR, "Error with get response from: ")
			return []models.Issue{}
		}

		body, err := io.ReadAll(response.Body)

		var issueResponce models.IssuesList
		err = json.Unmarshal(body, &issueResponce)
		if err != nil {
			connector.logger.Log(logger.ERROR, " ")
			return []models.Issue{}
		}

		counterOfIssues := issueResponce.IssuesCount
		if counterOfIssues == 0 {
			return []models.Issue{}
		}

		issues, timeUntilNewRequest, isRequestGompleteSuccsesfully = connector.threadsFunc(counterOfIssues, httpClient, projectName, timeUntilNewRequest)

	}
	if timeUntilNewRequest > connector.configReader.GetMaxTimeSleep() {
		connector.logger.Log(logger.ERROR, "Error, too much time!")
		return []models.Issue{}
	}
	return issues
}

func (connector *Connector) threadsFunc(counterOfIssues int, httpClient *http.Client, projectName string, timeUntilNewRequest int) ([]models.Issue, int, bool) {
	var issues []models.Issue
	counterOfThreads := connector.configReader.GetThreadCount()
	issueInOneRequest := connector.configReader.GetIssusOnOneRequest()

	channelErrorr := make(chan models.Issue)
	waitGroup := sync.WaitGroup{}
	mutex := sync.Mutex{}
	isError := false
	for i := 0; i < counterOfThreads; i++ {
		waitGroup.Add(1)
		go func(currentThreadNumber int) {
			defer waitGroup.Done()
			select {
			case <-channelErrorr:
				connector.logger.Log(logger.ERROR, "Error while reading issues in thread")
				return
			default:
				startIndex := currentThreadNumber*(counterOfIssues/counterOfThreads) + 1
				numberOfRequests := int(math.Ceil(float64(counterOfIssues) / float64(counterOfThreads*issueInOneRequest)))

				for j := 0; j < numberOfRequests; j++ {

					startAt := startIndex + j*issueInOneRequest
					if startAt < counterOfIssues {

						response, errResponce := httpClient.Get(connector.jiraRepositoryUrl +
							"/rest/api/2/search?jql=project=" + projectName +
							"&expand=changelog&startAt=" + strconv.Itoa(startAt) +
							"&maxResults=" + strconv.Itoa(issueInOneRequest))

						body, errRead := io.ReadAll(response.Body)

						if errRead != nil || errResponce != nil {
							isError = true
							close(channelErrorr)
							return
						}
						var issueResponse models.IssuesList
						_ = json.Unmarshal(body, &issueResponse)

						mutex.Lock()
						issues = append(issues, issueResponse.Issues...)
						mutex.Unlock()
					}
				}
			}
		}(i)
	}
	waitGroup.Wait()
	if isError {
		timeUntilNewRequest = connector.increaseTimeUntilNewRequest(timeUntilNewRequest, projectName)
	}
	return issues, timeUntilNewRequest, !isError
}

func (connector *Connector) increaseTimeUntilNewRequest(timeUntilNewRequest int, projectName string) int {
	timeMultiplier := 2.0
	time.Sleep(time.Millisecond * time.Duration(timeUntilNewRequest))
	newTimeUntilNewRequest := int(math.Ceil(float64(timeUntilNewRequest) * timeMultiplier))
	connector.logger.Log(logger.INFO, "Can`t download issues from project \""+
		projectName+"\": waiting "+strconv.Itoa(timeUntilNewRequest)+" Millisecond")
	return newTimeUntilNewRequest
}

/*
Выгружает проекты
Параметр limit - сколько всего проектов необходимо вернуть
Параметр page - порядковый номер страницы, который необходимо
вернуть
Параметр search - фильтр, который накладывается на название и ключ
*/
func (connector *Connector) GetProjects(limit int, page int, search string) models.Projects {
	httpClient := &http.Client{}
	response, err := httpClient.Get(connector.jiraRepositoryUrl + "/rest/api/2/project")
	if err != nil || response.StatusCode != http.StatusOK {
		connector.logger.Log(logger.ERROR, "Error with get response from about projects ")
		return models.Projects{}
	}

	body, err := io.ReadAll(response.Body)

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
