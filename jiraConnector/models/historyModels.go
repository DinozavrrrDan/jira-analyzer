package models

import "time"

type IssueHistories struct {
	Changelog struct {
		Histories []History `json:"histories"`
	} `json:"changelog"`
}

type History struct {
	Author struct {
		Name string `json:"name"`
	} `json:"author"`
	ChangeTime    time.Time `json:"changeTime"`
	StatusChanges []StatusChange
}

type StatusChange struct {
	FromStatus string `json:"fromStatus"`
	ToStatus   string `json:"toStatus"`
}
