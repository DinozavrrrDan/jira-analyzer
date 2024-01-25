package repository

import "github.com/DinozvrrDan/jira-analyzer/backend/analytics/internal/models"

type IAnalyticsRepository interface {
	GetGraphsOneData(title string) ([]models.GraphOneData, error)
	GetGraphsOneCategories(title string) ([]string, error)
	GetGraphsFourData(title string) ([]models.GraphFourData, error)
	GetGraphsFourCategories(title string) ([]string, error)
	GetGraphsSixData(title string) ([]models.GraphSixData, error)
	GetGraphsSixCategories(title string) ([]string, error)
	GetGraphsFiveData(title string) ([]models.GraphFiveData, error)
	GetGraphsFiveCategories(title string) ([]string, error)
}
