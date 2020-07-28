package db

import (
	"app/config"
	"database/sql"
	"fmt"
)
type SQLProvider interface {
	DB() (*sql.DB, error)
}

func NewSQLProvider() (SQLProvider, error) {
	config.C.Set("database.type", "sqlite3")
	dbType, err := config.C.GetString("database.type")
	if err != nil {
	    return nil, fmt.Errorf("can't create DB provider: %w", err)
	}
	switch dbType {
	case "mysql":
		return &Mysql{}, nil
	case "postgres":
		return &Postgres{}, nil
	case "sqlite3":
		return &SQLite{}, nil
	default:
		return SQLProvider(nil), fmt.Errorf("%s is not supported as a database provider", dbType)
	}
}