package database

import (
	"database/sql"
	"github.com/hddskull/urlShorty/internal/utils"
)

type Postgres struct {
}

func NewPostgres() *Postgres {
	return &Postgres{}
}

var _ Database = NewPostgres()

// ConnectDB uses creds (credentials) to establish a connection: host, user, password, db name, etc.
func (p Postgres) ConnectDB(creds string) (*sql.DB, error) {

	db, err := sql.Open("postgres", creds)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}

func (p Postgres) Ping() error {
	err := dbConnection.Ping()
	if err != nil {
		utils.SugaredLogger.Debugln("Ping():", err)
	}
	return err
}

func (p Postgres) CloseDB() error {
	return dbConnection.Close()
}
