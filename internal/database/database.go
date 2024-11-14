package database

import "database/sql"

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
