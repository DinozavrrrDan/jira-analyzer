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
	getProjects()
	getProjectInfo()
}

func getProjectInfo() {
	httpClient := &http.Client{}
	projectName := "AAR" //Просто для примера имя
	responce, err := httpClient.Get("http://issues.apache.org/jira/rest/api/2/search?jql=project=" + projectName + "&expand=changelog&startAt=0&maxResults=1")

	if err != nil {
		fmt.Print(err) //заменю на логирование
		return
	}

	defer responce.Body.Close()

	body, err := io.ReadAll(responce.Body)

	var issueResponce models.IssuesList
	err = json.Unmarshal(body, &issueResponce)

	transformer.TrasformData(issueResponce)

}

func getProjects() {
	httpClient := &http.Client{}
	resp, err := httpClient.Get("http://issues.apache.org/jira/rest/api/2/project")

	if err != nil {
		fmt.Print(err)
		return
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Print(err)
		return
	}
	var jiraProjects []models.JiraProject
	err = json.Unmarshal(body, &jiraProjects) //получаем информацию через сериализацию

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
