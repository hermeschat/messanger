package db

import (
	"database/sql"
	"fmt"
	"github.com/hermeschat/engine/config"
	"github.com/hermeschat/engine/transport/psql"
)

type SQLProvider interface {
	DB() (*sql.DB, error)
}

func NewSQLProvider() (SQLProvider, error) {
	config.C.SetDefault("database.type", "sqlite3")
	dbType := config.C.GetString("database.type")
	switch dbType {
	case "psql":
		return &psql.Postgres{}, nil
	default:
		return SQLProvider(nil), fmt.Errorf("%s is not supported as a database provider", dbType)
	}
}
