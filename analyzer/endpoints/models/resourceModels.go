package models

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

type responseStrucr struct {
	Links ListOfReferens `json:"_links"`
	//придумать как предавать info в случае разных запросов
}
