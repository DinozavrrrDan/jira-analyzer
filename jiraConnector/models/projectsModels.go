package models

type Project struct {
	Id        int    `json:"Id"`
	Key       string `json:"Key"`
	Name      string `json:"Name"`
	Link      string `json:"Link"`
	Existence bool   `json:"Existence"`
}

type JiraProject struct {
	Key  string `json:"key"`
	Name string `json:"name"`
	Link string `json:"self"`
}

type Projects struct {
	Projects []Project `json:"projects"`
}

type Page struct {
	CurrentPageNumber  int `json:"currentPage"`
	TotalPageCount     int `json:"pageCount"`
	TotalProjectsCount int `json:"projectsCount"`
}
