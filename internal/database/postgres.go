package database

import (
	"database/sql"
)

type Postgres struct {
	DB *sql.DB
}

func NewPostgres() *Postgres {
	return &Postgres{}
}

var _ Database = NewPostgres()

// ConnectDB uses creds (credentials) to establish a connection: host, user, password, db name, etc.
func (p Postgres) ConnectDB(creds string) error {

	db, err := sql.Open("postgres", creds)
	if err != nil {
		return err
	}

	err = db.Ping()
	if err != nil {
		return err
	}
	p.DB = db

	return nil
}

func (p Postgres) Ping() error {
	return p.DB.Ping()
}

func (p Postgres) CloseDB() error {
	return p.DB.Close()
}
