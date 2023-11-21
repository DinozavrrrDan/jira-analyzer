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

type DataInfoIssue struct {
	IssueID  string `json:"id"`
	Projects struct {
		ProjectID string `json:"projectID"`
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
	Status      string    `json:"status"`
	Creator     string    `json:"creator"`
	Assignee    string    `json:"assignee"`
	Timespent   int64     `json:"timespent"`
}

type DataInfoHistory struct {
	IssueID    int    `json:"issueID"`
	AuthorID   int    `json:"authorID"`
	ChangeTime uint64 `json:"changeTime"`
	FromStatus string `json:"fromStatus"`
	ToStatus   string `json:"toStatus"`
}

type DataInfoProject struct {
	ProjectID string `json:"projectID"`
	Title     int    `json:"title"`
}
