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
	Id    int    `json:"projectId"`
	Title string `json:"title"`
}
