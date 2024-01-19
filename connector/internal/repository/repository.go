package repository

import (
	"database/sql"
	postrege_sql "github.com/DinozvrrDan/jira-analyzer/connector/internal/repository/postrege-sql"
)

type Repositories struct {
	ConnectorRepository IConnectorRepository
}

func NewRepositories(db *sql.DB) *Repositories {
	return &Repositories{
		ConnectorRepository: postrege_sql.NewConnectorRepository(db),
	}
}
