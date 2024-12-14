package storage

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/hddskull/urlShorty/config"
	"github.com/hddskull/urlShorty/internal/model"
	"github.com/hddskull/urlShorty/internal/utils"
	"github.com/hddskull/urlShorty/tools/custom"
	_ "github.com/lib/pq"
)

var dbConnection *sql.DB

type PostgresStorage struct {
}

func newPostgresStorage() *PostgresStorage {
	return &PostgresStorage{}
}

// Storage interface
var _ Storage = newPostgresStorage()

func (ps *PostgresStorage) Setup() error {
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

func (ps *PostgresStorage) Close() error {
	err := dbConnection.Close()
	if err != nil {
		utils.SugaredLogger.Errorln(err)
	}
	return err
}

func (ps *PostgresStorage) Save(ctx context.Context, u string) (string, error) {
	utils.SugaredLogger.Debugln("Save() called")
	//create model
	newModel, err := model.NewFileStorageModel(u, "")
	if err != nil {
		return "", err
	}

	//create transaction
	tx, err := dbConnection.Begin()
	if err != nil {
		return "", err
	}

	//write query
	//query := "INSERT INTO urls (uuid, shortURL, originalURL) VALUES ($1, $2, $3) ON CONFLICT DO NOTHING;"
	//res, err := tx.ExecContext(ctx, query, newModel.UUID, newModel.ShortURL, newModel.OriginalURL)

	//NewQuery
	query := `
		INSERT INTO urls (uuid, shortURL, originalURL)
			VALUES ($1, $2, $3)
			ON CONFLICT (originalURL) 
			DO UPDATE SET
    			originalURL=EXCLUDED.originalURL
			RETURNING shortURL != $2, shortURL
	`

	row := dbConnection.QueryRowContext(ctx, query, newModel.UUID, newModel.ShortURL, newModel.OriginalURL)
	var isConflict bool
	var shortURL string
	err = row.Scan(&isConflict, &shortURL)
	utils.SugaredLogger.Debugln("Save(): isConflict:", isConflict, "|shortURL:", shortURL, "|err:", err)

	if err != nil {
		tx.Rollback()
		return "", err
	}
	if isConflict {
		tx.Rollback()
		conflictErr := custom.NewUniqueViolationError(fmt.Errorf("duplicate of %s", newModel.OriginalURL), shortURL)
		return "", conflictErr
	}

	//commit
	err = tx.Commit()
	if err != nil {
		return "", err
	}

	return newModel.ShortURL, nil
}

func (ps *PostgresStorage) SaveBatch(ctx context.Context, arr []model.StorageModel) error {

	//create transaction
	tx, err := dbConnection.Begin()
	if err != nil {
		return err
	}

	//query
	query := "INSERT INTO urls (uuid, shortURL, originalURL) VALUES ($1, $2, $3);"

	//batch query
	for _, v := range arr {
		_, err = tx.ExecContext(ctx, query, v.UUID, v.ShortURL, v.OriginalURL)
		if err != nil {
			//on error roll back
			tx.Rollback()
			return err
		}
	}

	//commit
	err = tx.Commit()
	if err != nil {
		return err
	}

	//if transaction successful return models
	return nil
}

func (ps *PostgresStorage) Get(ctx context.Context, id string) (string, error) {
	query := "SELECT originalURL FROM urls WHERE shortURL = $1;"
	row := dbConnection.QueryRowContext(ctx, query, id)
	var originalURL string
	err := row.Scan(&originalURL)
	if err != nil {
		return "", err
	}
	return originalURL, nil
}

func (ps *PostgresStorage) Ping(ctx context.Context) error {
	err := dbConnection.Ping()
	if err != nil {
		utils.SugaredLogger.Errorln(err)
	}
	return err
}

func handleUniqueViolation(ctx context.Context, originalURL string) error {
	query := "SELECT shortURL FROM urls WHERE originalURL = $1;"
	row := dbConnection.QueryRowContext(ctx, query, originalURL)
	var shortURL string
	err := row.Scan(&shortURL)
	if err != nil {
		return err
	}

	return custom.NewUniqueViolationError(fmt.Errorf("duplicate of %s", originalURL), shortURL)
}
