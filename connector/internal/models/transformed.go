package models

import "time"

type TransformedIssue struct {
	Project     string    `json:"Project"`
	Author      string    `json:"Author"`
	Assignee    string    `json:"Assignee"`
	Key         string    `json:"Key"`
	Summary     string    `json:"Summary"`
	Description string    `json:"Description"`
	Type        string    `json:"Type"`
	Priority    string    `json:"Priority"`
	Status      string    `json:"Status"`
	CreatedTime time.Time `json:"CreatedTime"`
	ClosedTime  time.Time `json:"ClosedTime"`
	UpdatedTime time.Time `json:"UpdatedTime"`
	TimeSpent   int64     `json:"TimeSpent"`
}
