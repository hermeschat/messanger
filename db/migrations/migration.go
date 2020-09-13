package migration

import (
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/golang-migrate/migrate/v4"
	"github.com/hermeschat/engine/config"
	"github.com/hermeschat/engine/db"
	"time"
)

func New() (*migrate.Migrate, error) {
	prv ,err := db.NewSQLProvider()
	if err != nil {
		return nil, err
	}
	_db, err := prv.DB()
	if err != nil {
		return nil, err
	}
	drv, err := postgres.WithInstance(_db, &postgres.Config{
		MigrationsTable:  "migrations",
		DatabaseName:     "hermes",
		SchemaName:       "public",
		StatementTimeout: time.Millisecond*100,
	})
	if err != nil {
		return nil, err
	}
	dbName := config.C.GetString("database.name")
	m, err := migrate.NewWithDatabaseInstance("file://./db/migrations",dbName , drv)
	if err != nil {
		return nil, err
	}
	return m, nil
}