package repository

import (
	"database/sql"
	"github.com/DinozvrrDan/jira-analyzer/backend/resource/internal/repository/postrege-sql"
)

type Repositories struct {
	ResourceRepository IResourceRepository
}

func NewRepositories(db *sql.DB) *Repositories {
	return &Repositories{
		ResourceRepository: postrege_sql.NewResourceRepository(db),
	}
}
