package connector

import (
	"Jira-analyzer/common/configReader"
	"Jira-analyzer/common/logger"
	"Jira-analyzer/jiraConnector/models"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"io"
	"math"
	"net/http"
	"strconv"
	"sync"
)

type Connector struct {
	logger            *logger.Logger
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
func (connector *Connector) GetProjectIssues(projectName string) ([]models.Issue, error) {
	isRequestNotCompleted := true
	timeUntilNewRequest := connector.configReader.GetMinTimeSleep()

	var issues []models.Issue

	for isRequestNotCompleted && timeUntilNewRequest <= connector.configReader.GetMaxTimeSleep() {
		httpClient := &http.Client{}
		response, err := httpClient.Get(connector.jiraRepositoryUrl + "/rest/api/2/search?jql=project=" +
			projectName + "&expand=changelog&startAt=0&maxResults=1")

		if err != nil || response.StatusCode != http.StatusOK {
			connector.logger.Log(logger.ERROR, "Error with get response from: "+projectName)
			return []models.Issue{}, fmt.Errorf("error with get response from: " + projectName)
		}

		body, err := io.ReadAll(response.Body)

		var issueResponse models.IssuesList
		err = json.Unmarshal(body, &issueResponse)

		if err != nil {
			connector.logger.Log(logger.ERROR, err.Error())

			return []models.Issue{}, err
		}

		counterOfIssues := issueResponse.IssuesCount

		if counterOfIssues == 0 {
			return []models.Issue{}, fmt.Errorf("error: no issues")
		}

		issues, timeUntilNewRequest, isRequestNotCompleted = connector.threadsFunc(counterOfIssues,
			httpClient, projectName, timeUntilNewRequest)

	}

	if timeUntilNewRequest > connector.configReader.GetMaxTimeSleep() {
		connector.logger.Log(logger.ERROR, "Error, too much time!")

		return []models.Issue{}, fmt.Errorf("error, too much time")
	}

	return issues, nil
}

func (connector *Connector) threadsFunc(counterOfIssues int, httpClient *http.Client,
	projectName string, timeUntilNewRequest int) ([]models.Issue, int, bool) {
	var issues []models.Issue

	counterOfThreads := connector.configReader.GetThreadCount()
	issueInOneRequest := connector.configReader.GetIssuesOnOneRequest()

	channelError := make(chan models.Issue)
	waitGroup := sync.WaitGroup{}
	mutex := sync.Mutex{}
	isError := false

	for i := 0; i < counterOfThreads; i++ {
		waitGroup.Add(1)

		go func(currentThreadNumber int) {
			defer waitGroup.Done()
			select {
			case <-channelError:
				connector.logger.Log(logger.ERROR, "Error while reading issues in thread")

				return
			default:
				startIndex := currentThreadNumber*(counterOfIssues/counterOfThreads) + 1
				numberOfRequests := int(math.Ceil(float64(counterOfIssues) / float64(counterOfThreads*issueInOneRequest)))

				for j := 0; j < numberOfRequests; j++ {
					startAt := startIndex + j*issueInOneRequest

					if startAt < counterOfIssues {
						response, errResponse := httpClient.Get(connector.jiraRepositoryUrl +
							"/rest/api/2/search?jql=project=" + projectName +
							"&expand=changelog&startAt=" + strconv.Itoa(startAt) +
							"&maxResults=" + strconv.Itoa(issueInOneRequest))

						body, errRead := io.ReadAll(response.Body)

						if errRead != nil || errResponse != nil {
							isError = true

							close(channelError)

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

	return issues, timeUntilNewRequest, isError
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
func (connector *Connector) GetProjects(limit int, page int, search string) ([]models.Project, models.Page, error) {
	httpClient := &http.Client{}

	response, err := httpClient.Get(connector.jiraRepositoryUrl + "/rest/api/2/project")

	if err != nil || response.StatusCode != http.StatusOK {
		connector.logger.Log(logger.ERROR, "Error with get response from about projects ")

		return []models.Project{}, models.Page{}, err
	}

	body, err := io.ReadAll(response.Body)

	if err != nil {
		connector.logger.Log(logger.ERROR, err.Error())

		return []models.Project{}, models.Page{}, err
	}

	var jiraProjects []models.JiraProject
	err = json.Unmarshal(body, &jiraProjects) //получаем информацию через сериализацию

	if err != nil {
		connector.logger.Log(logger.ERROR, err.Error())

		return []models.Project{}, models.Page{}, err
	}

	var projects []models.Project

	counterOfProjects := 0

	//Получение информации о определенном колчичестве проектов
	for _, element := range jiraProjects {
		if filterBySearch(element.Name, search) {
			counterOfProjects++

			projects = append(projects, models.Project{
				Existence: false,
				Id:        0,
				Name:      element.Name,
				Link:      element.Link,
				Key:       element.Key,
			})
		}
	}

	//обрезка проектов по странице

	startIndexOfProject := limit * (page - 1)
	endIndexOfProject := limit * page

	if endIndexOfProject >= len(projects) {
		endIndexOfProject = len(projects)
	}

	return projects[startIndexOfProject:endIndexOfProject],
		models.Page{
			CurrentPageNumber:  page,
			TotalPageCount:     counterOfProjects / limit,
			TotalProjectsCount: counterOfProjects,
		}, nil
}

func filterBySearch(projectName, search string) bool {
	return strings.Contains(strings.ToLower(projectName), strings.ToLower(search))
}
