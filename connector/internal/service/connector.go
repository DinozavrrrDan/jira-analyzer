package service

import (
	"encoding/json"
	"fmt"
	"github.com/DinozvrrDan/jira-analyzer/connector/config"
	"github.com/DinozvrrDan/jira-analyzer/connector/internal/models"
	"github.com/DinozvrrDan/jira-analyzer/connector/pkg/logger"
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
	cfg *config.Config
	url string
}

func NewConnectorService(jiraUrl string, log *logger.Logger, cfg *config.Config) *ConnectorService {
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

const jiraRequestPart1 = "/rest/api/2/search?jql=project="
const jiraRequestPart2 = "&expand=changelog&startAt=0&maxResults=1"

func (connector *ConnectorService) GetProjectIssues(projectName string) ([]models.Issue, error) {
	isRequestNotCompleted := true
	timeUntilNewRequest := connector.cfg.Connector.MinTimeSleep

	var issues []models.Issue

	for isRequestNotCompleted && timeUntilNewRequest <= connector.cfg.Connector.MaxTimeSleep {
		httpClient := &http.Client{}

		response, err := httpClient.Get(connector.url + jiraRequestPart1 +
			projectName + jiraRequestPart2)

		if err != nil || response.StatusCode != http.StatusOK {
			connector.log.Log(logger.ERROR, "Error with get response from: "+projectName)
			return nil, fmt.Errorf("error with get response from: " + projectName)
		}

		body, err := io.ReadAll(response.Body)

		var issueResponse models.IssuesList
		err = json.Unmarshal(body, &issueResponse)

		if err != nil {
			connector.log.Log(logger.ERROR, err.Error())

			return nil, err
		}

		counterOfIssues := issueResponse.IssuesCount

		if counterOfIssues == 0 {
			return nil, fmt.Errorf("error: no issues")
		}

		issues, timeUntilNewRequest, isRequestNotCompleted = connector.threadsFunc(counterOfIssues,
			httpClient, projectName, timeUntilNewRequest)

	}

	if timeUntilNewRequest > connector.cfg.Connector.MaxTimeSleep {
		connector.log.Log(logger.ERROR, "error, too much time!")

		return []models.Issue{}, fmt.Errorf("error, too much time")
	}

	return issues, nil
}

func (connector *ConnectorService) threadsFunc(counterOfIssues int, httpClient *http.Client,
	projectName string, timeUntilNewRequest int) ([]models.Issue, int, bool) {
	var issues []models.Issue

	counterOfThreads := connector.cfg.Connector.ThreadCount
	issueInOneRequest := connector.cfg.Connector.IssueInRequest

	channelError := make(chan struct{})
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

						var issueResponse models.IssuesList
						err := json.Unmarshal(body, &issueResponse)

						if err != nil {
							isError = true
							close(channelError)
							return
						}
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
func (connector *ConnectorService) GetProjects(limit int, page int, search string) ([]models.Project, models.Page, error) {
	httpClient := &http.Client{}

	response, err := httpClient.Get(connector.url + "/rest/api/2/project")

	if err != nil || response.StatusCode != http.StatusOK {
		connector.log.Log(logger.ERROR, "Error with get response from about projects ")

		return []models.Project{}, models.Page{}, err
	}

	body, err := io.ReadAll(response.Body)

	if err != nil {
		connector.log.Log(logger.ERROR, err.Error())

		return []models.Project{}, models.Page{}, err
	}

	var jiraProjects []models.JiraProject
	err = json.Unmarshal(body, &jiraProjects) //получаем информацию через сериализацию
	if err != nil {
		connector.log.Log(logger.ERROR, err.Error())

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
