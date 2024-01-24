package models

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
}
