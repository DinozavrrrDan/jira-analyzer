package models

import "time"

type ListOfReferens struct {
	Issues struct {
		Href string `json:"href"`
	} `json:"issues"`
	Projects struct {
		Href string `json:"href"`
	} `json:"projects"`
	Histories struct {
		Href string `json:"href"`
	} `json:"histories"`
	Self struct {
		Href string `json:"href"`
	} `json:"self"`
}

type ResponseStrucrt struct {
	Links ListOfReferens `json:"_links"`
	//придумать как предавать info в случае разных запросов
	Message string `json:"message"`
	Name    string `json:"name"`
	Status  bool   `json:"status"`
}

type IssueInfo struct {
	IssueID  int `json:"id"`
	Projects struct {
		ProjectID int    `json:"projectID"`
		Key       string `json:"key"`
		Name      string `json:"name"`
		Url       string `json:"url"`
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
	Author      string    `json:"creator"`
	Assignee    string    `json:"assignee"`
	Timespent   int       `json:"timespent"`
}

type HistoryInfo struct {
	IssueID    int       `json:"issueID"`
	AuthorID   int       `json:"authorID"`
	ChangeTime time.Time `json:"changeTime"`
	FromStatus string    `json:"fromStatus"`
	ToStatus   string    `json:"toStatus"`
}

type ProjectInfo struct {
	ProjectID int    `json:"projectID"`
	Title     string `json:"title"`
}
