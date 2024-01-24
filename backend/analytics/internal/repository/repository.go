package repository

import (
	"database/sql"
	"github.com/DinozvrrDan/jira-analyzer/backend/analytics/internal/repository/postrege-sql"
)

type Repositories struct {
	AnalyticsRepository IAnalyticsRepository
}

func NewRepositories(db *sql.DB) *Repositories {
	return &Repositories{
		AnalyticsRepository: postrege_sql.NewAnalyticsRepository(db),
	}
}
