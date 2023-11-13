package models

type IssuesList struct {
	IssuesCount int     `json:"total"`
	Issues      []Issue `json:"issues"`
}

type Issue struct {
	Key    string      `json:"key"`
	Fields IssueFields `json:"fields"`
}

/*
Возвращаемые проекты должны содержать следующие поля:
ключ проекта
имя проекта
url проекта
*/
type Project struct {
	Key  string `json:"key"`
	Name string `json:"name"`
	Link string `json:"link"`
}

type JiraProject struct {
	Key  string `json:"key"`
	Name string `json:"name"`
	Link string `json:"self"`
}

type Projects struct {
	Projects []Project `json:"projects"`
	Page     Page      `json:"page"`
}

type Page struct {
	TotalPageCount     int `json:"totalPageCount"`
	CurrentPageNumber  int `json:"currentPageNumber"`
	TotalProjectsCount int `json:"totalProjectsCount"`
}

type IssueFields struct {
	Summary string `json:"summary"`
	Type    struct {
		Name string `json:"name"`
	} `json:"issuetype"`
	Status struct {
		Name string `json:"name"`
	} `json:"status"`
	Priority struct {
		Name string `json:"name"`
	} `json:"priority"`
	Creator struct {
		Name string `json:"name"`
	} `json:"creator"`
	Project struct {
		Name string `json:"name"`
	} `json:"project"`
	Description  string `json:"description"`
	AssigneeName struct {
		Name string `json:"name"`
	} `json:"assignee"`
	CreatedTime string `json:"created"`
	UpdatedTime string `json:"updated"`
	ClosedTime  string `json:"resolutiondate"`
}
