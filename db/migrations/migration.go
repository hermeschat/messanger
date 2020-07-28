package migration

import (
	"github.com/golang-migrate/migrate/v4/database/mysql"
	"github.com/golang-migrate/migrate/v4"
	"github.com/hermeschat/engine/config"
	"github.com/hermeschat/engine/db"
)

func New() (*migrate.Migrate, error) {
	prv ,err := db.NewSQLProvider()
	if err != nil {
		return nil, err
	}
	db, err := prv.DB()
	if err != nil {
		return nil, err
	}
	drv, err := mysql.WithInstance(db, &mysql.Config{})
	if err != nil {
		return nil, err
	}
	dbName, err := config.C.GetString("database.name")
	if err != nil {
	    return nil, err
	}
	m, err := migrate.NewWithDatabaseInstance("./db/migrations",dbName , drv)
	if err != nil {
		return nil, err
	}
	return m, nil
}