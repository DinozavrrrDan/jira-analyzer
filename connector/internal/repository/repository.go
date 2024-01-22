package repository

import (
	"database/sql"
	postregesql "github.com/DinozvrrDan/jira-analyzer/connector/internal/repository/postrege-sql"
)

type Repositories struct {
	ConnectorRepository IConnectorRepository
}

func NewRepositories(db *sql.DB) *Repositories {
	return &Repositories{
		ConnectorRepository: postregesql.NewConnectorRepository(db),
	}
}
