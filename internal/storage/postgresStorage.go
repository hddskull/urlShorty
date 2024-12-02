package storage

import (
	"context"
	"database/sql"
	"github.com/hddskull/urlShorty/config"
	"github.com/hddskull/urlShorty/internal/model"
	"github.com/hddskull/urlShorty/internal/utils"
	"github.com/hddskull/urlShorty/tools/custom"
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
	//check that url isn't empty
	if u == "" {
		utils.SugaredLogger.Debugln("Save() empty arg:", custom.ErrEmptyURL)
		return "", custom.ErrEmptyURL
	}

	//check if url is already saved
	existingUUID, err := ps.checkIfExists(u)
	if err != nil {
		utils.SugaredLogger.Debugln("checkExistence() err:", err)
		return "", err
	}

	//if exists return uuid
	if existingUUID != "" {
		utils.SugaredLogger.Debugln("Save() url already saved:", existingUUID, u)
		return existingUUID, nil
	}

	//create model
	newModel, err := model.NewFileStorageModel(u, "")
	if err != nil {
		return "", err
	}

	//write query
	query := "INSERT INTO urls (uuid, shortURL, originalURL) VALUES ($1, $2, $3);"
	_, err = dbConnection.Exec(query, newModel.UUID, newModel.ShortURL, newModel.OriginalURL)
	if err != nil {
		return "", err
	}

	return newModel.ShortURL, nil
}

func (ps PostgresStorage) SaveBatch(arr []model.StorageModel) ([]model.StorageModel, error) {

	//create transaction
	tx, err := dbConnection.Begin()
	if err != nil {
		return nil, err
	}

	//query
	query := "INSERT INTO urls (uuid, shortURL, originalURL) VALUES ($1, $2, $3);"

	//batch query
	for _, v := range arr {
		_, err = tx.Exec(query, v.UUID, v.ShortURL, v.OriginalURL)
		if err != nil {
			//on error roll back
			tx.Rollback()
			return nil, err
		}
	}

	//commit
	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	//if transaction successful return models
	return arr, nil
}

func (ps PostgresStorage) Get(id string) (string, error) {
	query := "SELECT originalURL FROM urls WHERE shortURL = $1;"
	row := dbConnection.QueryRowContext(context.Background(), query, id)
	var originalURL string
	err := row.Scan(&originalURL)
	if err != nil {
		return "", err
	}
	return originalURL, nil
}

func (ps PostgresStorage) Ping() error {
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

// Supporting methods
func (ps PostgresStorage) checkIfExists(origin string) (string, error) {
	query := "SELECT shortURL FROM urls WHERE originalURL = $1;"
	row := dbConnection.QueryRowContext(context.Background(), query, origin)
	var shortURL string
	err := row.Scan(&shortURL)
	if err != nil && err.Error() != "sql: no rows in result set" {
		return "", err
	}

	return shortURL, nil
}
