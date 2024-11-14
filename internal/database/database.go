package database

import (
	"database/sql"
	"github.com/hddskull/urlShorty/config"
)

type Database interface {
	ConnectDB(creds string) (*sql.DB, error)
	CloseDB() error
	Ping() error
}

var Current Database = NewPostgres()
var dbConnection *sql.DB

func SaveConnection(db *sql.DB) {
	dbConnection = db
}

func Start() {
	creds := config.DBCredentials
	db, err := Current.ConnectDB(creds)
	if err != nil {
		panic(err)
	}
	SaveConnection(db)
}
