package psql

import (
	"database/sql"
	"github.com/hermeschat/engine/config"
	_ "github.com/lib/pq"
)
var db *sql.DB

func New() (*sql.DB, error) {
	if db != nil {
		return db, nil
	}
	if err := db.Ping(); err == nil {
		return db, nil
	}
	db, err := sql.Open("postgres", config.PostgresURI())
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}

func Migrate(db *sql.DB) error {
	tables := []string{
		`CREATE TABLE IF NOT EXISTS users ();`,
		`CREATE TABLE IF NOT EXISTS channels ();`,
		`CREATE TABLE IF NOT EXISTS messages ();`,
	}
	for _, t := range tables {
		_, err := db.Exec(t)
		if err != nil {
			return err
		}
	}
	return nil
}