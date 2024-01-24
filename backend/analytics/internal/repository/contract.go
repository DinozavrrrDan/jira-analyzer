package repository

import "github.com/DinozvrrDan/jira-analyzer/backend/analytics/internal/models"

type IAnalyticsRepository interface {
	//функции для графики
	GetGraphsOneData(title string) ([]models.GraphOneData, error)
	GetGraphsOneCategories(title string) ([]string, error)
	GetGraphsFourData(title string) ([]models.GraphFourData, error)
	GetGraphsFourCategories(title string) ([]string, error)
}
