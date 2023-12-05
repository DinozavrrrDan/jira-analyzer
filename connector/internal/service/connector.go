package service

import (
	"connector/config"
	"connector/internal/entities"
	"connector/pkg/logger"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
)

type ConnectorService struct {
	log *logger.Logger
	cfg *config.Reader
	url string
}

func NewConnectorService(jiraUrl string, log *logger.Logger, cfg *config.Reader) *ConnectorService {
	return &ConnectorService{
		log: log,
		cfg: cfg,
		url: jiraUrl,
	}
}

/*
В случае удачного выполнения запроса должен быть возвращен JSON,
который содержит массив проектов и общее количество страниц при
данном параметре limit
*/
func (connector *ConnectorService) GetProjectIssues(projectName string) ([]entities.Issue, error) {
	isRequestNotCompleted := true
	timeUntilNewRequest := connector.cfg.GetMinTimeSleep()

	var issues []entities.Issue

	for isRequestNotCompleted && timeUntilNewRequest <= connector.cfg.GetMaxTimeSleep() {
		httpClient := &http.Client{}
		response, err := httpClient.Get(connector.url + "/rest/api/2/search?jql=project=" +
			projectName + "&expand=changelog&startAt=0&maxResults=1")

		if err != nil || response.StatusCode != http.StatusOK {
			connector.log.Log(logger.ERROR, "Error with get response from: "+projectName)
			return []entities.Issue{}, fmt.Errorf("error with get response from: " + projectName)
		}

		body, err := io.ReadAll(response.Body)

		var issueResponse entities.IssuesList
		err = json.Unmarshal(body, &issueResponse)

		if err != nil {
			connector.log.Log(logger.ERROR, err.Error())

			return []entities.Issue{}, err
		}

		counterOfIssues := issueResponse.IssuesCount

		if counterOfIssues == 0 {
			return []entities.Issue{}, fmt.Errorf("error: no issues")
		}

		issues, timeUntilNewRequest, isRequestNotCompleted = connector.threadsFunc(counterOfIssues,
			httpClient, projectName, timeUntilNewRequest)

	}

	if timeUntilNewRequest > connector.cfg.GetMaxTimeSleep() {
		connector.log.Log(logger.ERROR, "Error, too much time!")

		return []entities.Issue{}, fmt.Errorf("error, too much time")
	}

	return issues, nil
}

func (connector *ConnectorService) threadsFunc(counterOfIssues int, httpClient *http.Client,
	projectName string, timeUntilNewRequest int) ([]entities.Issue, int, bool) {
	var issues []entities.Issue

	counterOfThreads := connector.cfg.GetThreadCount()
	issueInOneRequest := connector.cfg.GetIssuesOnOneRequest()

	channelError := make(chan entities.Issue)
	waitGroup := sync.WaitGroup{}
	mutex := sync.Mutex{}
	isError := false

	for i := 0; i < counterOfThreads; i++ {
		waitGroup.Add(1)

		go func(currentThreadNumber int) {
			defer waitGroup.Done()
			select {
			case <-channelError:
				connector.log.Log(logger.ERROR, "Error while reading issues in thread")

				return
			default:
				startIndex := currentThreadNumber*(counterOfIssues/counterOfThreads) + 1
				numberOfRequests := int(math.Ceil(float64(counterOfIssues) / float64(counterOfThreads*issueInOneRequest)))

				for j := 0; j < numberOfRequests; j++ {
					startAt := startIndex + j*issueInOneRequest

					if startAt < counterOfIssues {
						response, errResponse := httpClient.Get(connector.url +
							"/rest/api/2/search?jql=project=" + projectName +
							"&expand=changelog&startAt=" + strconv.Itoa(startAt) +
							"&maxResults=" + strconv.Itoa(issueInOneRequest))

						body, errRead := io.ReadAll(response.Body)

						if errRead != nil || errResponse != nil {
							isError = true

							close(channelError)

							return
						}

						var issueResponse entities.IssuesList
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

func (connector *ConnectorService) increaseTimeUntilNewRequest(timeUntilNewRequest int, projectName string) int {
	timeMultiplier := 2.0

	time.Sleep(time.Millisecond * time.Duration(timeUntilNewRequest))
	newTimeUntilNewRequest := int(math.Ceil(float64(timeUntilNewRequest) * timeMultiplier))
	connector.log.Log(logger.INFO, "Can`t download issues from project \""+
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
func (connector *ConnectorService) GetProjects(limit int, page int, search string) ([]entities.Project, entities.Page, error) {
	httpClient := &http.Client{}

	response, err := httpClient.Get(connector.url + "/rest/api/2/project")

	if err != nil || response.StatusCode != http.StatusOK {
		connector.log.Log(logger.ERROR, "Error with get response from about projects ")

		return []entities.Project{}, entities.Page{}, err
	}

	body, err := io.ReadAll(response.Body)

	if err != nil {
		connector.log.Log(logger.ERROR, err.Error())

		return []entities.Project{}, entities.Page{}, err
	}

	var jiraProjects []entities.JiraProject
	err = json.Unmarshal(body, &jiraProjects) //получаем информацию через сериализацию

	if err != nil {
		connector.log.Log(logger.ERROR, err.Error())

		return []entities.Project{}, entities.Page{}, err
	}

	var projects []entities.Project

	counterOfProjects := 0

	//Получение информации о определенном колчичестве проектов
	for _, element := range jiraProjects {
		if filterBySearch(element.Name, search) {
			counterOfProjects++

			projects = append(projects, entities.Project{
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
		entities.Page{
			CurrentPageNumber:  page,
			TotalPageCount:     counterOfProjects / limit,
			TotalProjectsCount: counterOfProjects,
		}, nil
}

func filterBySearch(projectName, search string) bool {
	return strings.Contains(strings.ToLower(projectName), strings.ToLower(search))
}
