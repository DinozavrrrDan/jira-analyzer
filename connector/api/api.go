package api

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

type ListOfReferences struct {
	Issues    Link `json:"issues"`
	Projects  Link `json:"projects"`
	Histories Link `json:"histories"`
	Self      Link `json:"self"`
}

type Link struct {
	Href string `json:"href"`
}
