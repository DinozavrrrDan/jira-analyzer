package models

import "time"

type ListOfReferences struct {
	Issues    Link `json:"issues"`
	Projects  Link `json:"projects"`
	Histories Link `json:"histories"`
	Self      Link `json:"self"`
}

type Link struct {
	Href string `json:"href"`
}

type Project struct {
	Id        int    `json:"Id"`
	Key       string `json:"Key"`
	Name      string `json:"Name"`
	Link      string `json:"Link"`
	Existence bool   `json:"Existence"`
}

type ResponseStruct struct {
	Links    ListOfReferences `json:"_links"`
	Info     interface{}      `json:"data"`
	Message  string           `json:"message"`
	Name     string           `json:"name"`
	PageInfo Page             `json:"pageInfo"`
	Status   bool             `json:"status"`
}

type Page struct {
	CurrentPageNumber  int `json:"currentPage"`
	TotalPageCount     int `json:"pageCount"`
	TotalProjectsCount int `json:"projectsCount"`
}

type IssueInfo struct {
	Id          int         `json:"id"`
	Project     ProjectInfo `json:"project"`
	Key         string      `json:"key"`
	CreatedTime time.Time   `json:"createdTime"`
	ClosedTime  time.Time   `json:"closedTime"`
	UpdatedTime time.Time   `json:"updatedTime"`
	Summary     string      `json:"summary"`
	Description string      `json:"description"`
	Type        string      `json:"type"`
	Priority    string      `json:"priority"`
	Status      bool        `json:"status"`
	Author      string      `json:"author"`
	Assignee    string      `json:"assignee"`
	TimeSpent   int         `json:"timeSpent"`
}

type ProjectInfo struct {
	Id                  int    `json:"projectId"`
	Title               string `json:"title"`
	IssueCount          int    `json:"issueCount"`
	ClosedIssuesCount   int    `json:"closedIssuesCount"`
	OpenedIssuesCount   int    `json:"openedIssuesCount"`
	ResolvedIssuesCount int    `json:"resolvedIssuesCount"`
	ReopenedIssueCount  int    `json:"reopenedIssueCount"`
	ProgressIssuesCount int    `json:"progressIssuesCount"`
<<<<<<< HEAD
=======
	AverageTime         int    `json:"averageTime"`
>>>>>>> 6a15cb1650a9c1e304607e8f9b48d77b20ebf674
}
