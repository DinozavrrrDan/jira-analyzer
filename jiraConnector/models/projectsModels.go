package models

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
