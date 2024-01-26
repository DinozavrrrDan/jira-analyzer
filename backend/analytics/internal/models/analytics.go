package models

import "time"

type ListOfReferences struct {
	Issues   Link `json:"issues"`
	Projects Link `json:"projects"`
	Graphs   Link `json:"graphs"`
	Self     Link `json:"self"`
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
	ID      int `json:"id"`
	Project struct {
		ID    int    `json:"projectId"`
		Title string `json:"title"`
	} `json:"project"`
	Key         string    `json:"key"`
	CreatedTime time.Time `json:"createdTime"`
	ClosedTime  time.Time `json:"closedTime"`
	UpdatedTime time.Time `json:"updatedTime"`
	Summary     string    `json:"summary"`
	Description string    `json:"description"`
	Type        string    `json:"type"`
	Priority    string    `json:"priority"`
	Status      bool      `json:"status"`
	Author      string    `json:"author"`
	Assignee    string    `json:"assignee"`
	TimeSpent   int       `json:"timeSpent"`
}

type HistoryInfo struct {
	IssueID    int       `json:"issueId"`
	AuthorID   int       `json:"authorId"`
	ChangeTime time.Time `json:"changeTime"`
	FromStatus string    `json:"fromStatus"`
	ToStatus   string    `json:"toStatus"`
}

type ProjectInfo struct {
	ID    int    `json:"projectId"`
	Title string `json:"title"`
}
