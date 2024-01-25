package models

type GraphOneData struct {
	SpentTime string `json:"spentTime"`
	Count     int    `json:"count"`
}

type GraphOne struct {
	GraphOneData []GraphOneData `json:"graphOneData"`
	Categories   []string       `json:"categories"`
}

type GraphFourData struct {
	Type  string `json:"type"`
	Count int    `json:"count"`
}

type GraphFour struct {
	GraphFourData []GraphFourData `json:"graphFourData"`
	Categories    []string        `json:"categories"`
}

type GraphFiveData struct {
	Priority string `json:"priority"`
	Count    int    `json:"count"`
}

type GraphFive struct {
	GraphFiveData []GraphFiveData `json:"graphFiveData"`
	Categories    []string        `json:"categories"`
}

type GraphSixData struct {
	Priority string `json:"priority"`
	Count    int    `json:"count"`
}

type GraphSix struct {
	GraphSixData []GraphSixData `json:"graphSixData"`
	Categories   []string       `json:"categories"`
}

type CompareGraphsOne struct {
	Counter    int        `json:"counter"`
	GraphsOne  []GraphOne `json:"graphsOne"`
	Categories []string   `json:"categories"`
}
