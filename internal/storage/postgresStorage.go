package storage

import (
	"context"
	"database/sql"
	"errors"
	"github.com/hddskull/urlShorty/config"
	"github.com/hddskull/urlShorty/internal/utils"
)

var dbConnection *sql.DB

type PostgresStorage struct {
}

func newPostgresStorage() *PostgresStorage {
	return &PostgresStorage{}
}

// Storage interface
var _ Storage = newPostgresStorage()

func (ps PostgresStorage) Setup() error {
	var err error
	//establish connection
	dbConnection, err = sql.Open("postgres", config.DBCredentials)
	if err != nil {
		return err
	}

	//Ping
	err = dbConnection.Ping()
	if err != nil {
		return err
	}

	//Create table
	ctx := context.Background()
	tableQuery := `
	CREATE TABLE IF NOT EXISTS urls (
		uuid UUID PRIMARY KEY,
		shortURL TEXT NOT NULL,
		originalURL TEXT NOT NULL
	);`
	_, err = dbConnection.ExecContext(ctx, tableQuery)
	if err != nil {
		return err
	}

	// Debug: insert test data into DB here

	return nil
}

func (ps PostgresStorage) Save(u string) (string, error) {
	//TODO implement SAVE() method
	return "", errors.New("implementation needed")
}

func (ps PostgresStorage) Get(id string) (string, error) {
	//TODO implement GET() method
	return "", errors.New("implementation needed")
}

func (ps PostgresStorage) Ping() error {
	//fmt.Println("\n\n", dbConnection, "\n")
	err := dbConnection.Ping()
	if err != nil {
		utils.SugaredLogger.Errorln(err)
	}
	return err
}

func (ps PostgresStorage) Close() error {
	err := dbConnection.Close()
	if err != nil {
		utils.SugaredLogger.Errorln(err)
	}
	return err
}
