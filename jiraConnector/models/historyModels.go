package models

type IssueHistories struct {
	Changelog struct {
		Histories []History `json:"histories"`
	} `json:"changelog"`
}

type History struct {
	Author struct {
		Name string `json:"name"`
	} `json:"author"`
	ChangeTime    string `json:"created"`
	StatusChanges []StatusChange
}

type StatusChange struct {
	Field      string `json:"field"`
	FromStatus string `json:"fromString"`
	ToStatus   string `json:"toString"`
}
