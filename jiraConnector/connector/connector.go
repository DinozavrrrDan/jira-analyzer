package connector

import (
	"Jira-analyzer/jiraConnector/models"
	"Jira-analyzer/jiraConnector/transformer"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func main() {
	GetProjects()
	GetProjectInfo("AAR")
}

func GetProjectInfo(projectName string) {
	httpClient := &http.Client{}

	//temp
	projectName = "AAR" //Просто для примера имя

	//Вынести в конфиг
	jiraRepositoryUrl := "http://issues.apache.org/jira"

	response, err := httpClient.Get(jiraRepositoryUrl + "/rest/api/2/search?jql=project=" + projectName + "&expand=changelog&startAt=0&maxResults=1")
	if err != nil || response.StatusCode != http.StatusOK {
		fmt.Print(err) //заменю на логирование
		return
	}

	body, err := io.ReadAll(response.Body)
	var issueResponce models.IssuesList
	err = json.Unmarshal(body, &issueResponce)
	if err != nil {
		fmt.Print(err) //заменю на логирование
		return
	}

	transformer.TrasformData(issueResponce)

}

func GetProjects() {
	httpClient := &http.Client{}

	//Вынести в конфиг
	jiraRepositoryUrl := "http://issues.apache.org/jira"

	resp, err := httpClient.Get(jiraRepositoryUrl + "/rest/api/2/project")
	if err != nil {
		fmt.Print(err)
		return
	}

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		fmt.Print(err)
		return
	}
	var jiraProjects []models.JiraProject
	err = json.Unmarshal(body, &jiraProjects) //получаем информацию через сериализацию
	if err != nil {
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
