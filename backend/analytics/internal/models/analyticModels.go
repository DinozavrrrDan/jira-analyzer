package models

<<<<<<< HEAD:analyzer/models/analyticModels.go
type GraphData struct {
	PriorityType string `json:"priority_type"`
	Amount       int    `json:"count"`
=======
type GraphFourData struct {
	Type  string `json:"type"`
	Count int    `json:"count"`
}

type GraphFour struct {
	GraphFourData []GraphFourData `json:"graphFourData"`
	Categories    []string        `json:"categories"`
}

type GraphOneData struct {
	Id       int    `json:"id"`
	OpenTime uint64 `json:"openTime"`
}

type GraphOne struct {
	GraphOneData []GraphOneData `json:"graphOneData"`
	Categories   []string       `json:"categories"`
}

type CompareGraphsOne struct {
	Count     int        `json:"count"`
	GraphsOne []GraphOne `json:"graphsOne"`
>>>>>>> 6a15cb1650a9c1e304607e8f9b48d77b20ebf674:backend/analytics/internal/models/analyticModels.go
}
