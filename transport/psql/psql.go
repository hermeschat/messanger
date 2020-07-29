package psql

import (
	"database/sql"
	"fmt"
	"github.com/hermeschat/engine/config"
	_ "github.com/lib/pq"
)

type Postgres struct {
	conn *sql.DB
}

func (p *Postgres) DB() (*sql.DB, error) {
	SetPostgresDefaultConfiguration()
	host := config.C.GetString("database.host")
	port := config.C.GetString("database.port")
	user := config.C.GetString("database.user")
	password := config.C.GetString("database.password")
	name := config.C.GetString("database.name")
	sslmode := config.C.GetString("database.sslmode")
	if p.conn == nil {
		conn, err := postgresConnect(host, port, user, password, name, sslmode)
		if err != nil {
			return nil, err
		}
		p.conn = conn
		return p.conn, nil
	}
	if err := p.conn.Ping(); err != nil {
		return nil, err
	}
	return p.conn, nil
}

func postgresConnect(host, port, user, password, name, sslmode string) (*sql.DB, error) {
	conString := postgresConnectionString(host, port, user, password, name, sslmode)
	conn, err := sql.Open("postgres", conString)
	if err != nil {
		return nil, fmt.Errorf("Error in openning postgres connection: %w", err)
	}

	err = conn.Ping()
	if err != nil {
		return nil, fmt.Errorf("Error in pinging postgres database: %w", err)
	}
	return conn, nil
}

func postgresConnectionString(host, port, user, password, name, sslmode string) string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s", host, port, user, password, name, sslmode)
}

func SetPostgresDefaultConfiguration() {
	config.C.SetDefault("database.host", "localhost")
	config.C.SetDefault("database.port", "5432")
	config.C.SetDefault("database.user", "hermes")
	config.C.SetDefault("database.password", "hermes")
	config.C.SetDefault("database.name", "hermes")
	config.C.SetDefault("database.sslmode", "disable")
}
