package models

type ResponseStruct struct {
	Links    ListOfReferences `json:"_links"`
	Info     interface{}      `json:"data"`
	Message  string           `json:"message"`
	Name     string           `json:"name"`
	PageInfo Page             `json:"pageInfo"`
	Status   bool             `json:"status"`
}

type ListOfReferences struct {
	Issues   Link `json:"issues"`
	Projects Link `json:"projects"`
	Graphs   Link `json:"graphs"`
	Self     Link `json:"self"`
}

type Link struct {
	Href string `json:"href"`
}
