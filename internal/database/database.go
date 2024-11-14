package database

type Database interface {
	ConnectDB(creds string) error
	CloseDB() error
	Ping() error
}

var Current Database = NewPostgres()
