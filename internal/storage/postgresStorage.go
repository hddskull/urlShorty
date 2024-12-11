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
	"sync"
)

var dbConnection *sql.DB

type PostgresStorage struct {
	*sync.Mutex
}

func newPostgresStorage() *PostgresStorage {
	return &PostgresStorage{&sync.Mutex{}}
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
	//check that url isn't empty
	if u == "" {
		utils.SugaredLogger.Debugln("Save() empty arg:", custom.ErrEmptyURL)
		return "", custom.ErrEmptyURL
	}

	//create model
	newModel, err := model.NewFileStorageModel(u, "")
	if err != nil {
		return "", err
	}

	ps.Lock()
	defer ps.Unlock()

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
		WITH insert_attempt AS (
			INSERT INTO urls (uuid, shortURL, originalURL)
			VALUES ($1, $2, $3)
			ON CONFLICT (originalURL) DO NOTHING
			RETURNING shortURL
		)
		SELECT
			CASE
				WHEN EXISTS (SELECT 1 FROM insert_attempt) THEN false
				ELSE true
			END AS conflict,
			COALESCE((SELECT shortURL FROM insert_attempt), (SELECT shortURL FROM urls WHERE originalURL = $3)) AS shortURL;
	`

	//INTO conflict, _short_url;
	//
	//RAISE NOTICE 'Conflict: %, Short URL: %', conflict, _short_url;

	row := dbConnection.QueryRowContext(ctx, query, newModel.UUID, newModel.ShortURL, newModel.OriginalURL)
	var isConflict bool
	var shortURL string
	err = row.Scan(&isConflict, &shortURL)
	utils.SugaredLogger.Debugln("scanned row| isConflict:", isConflict, "|shortURL:", shortURL, "|err:", err)

	if err != nil {
		tx.Rollback()
		return "", err
	}
	if isConflict {
		tx.Rollback()
		conflictErr := custom.NewUniqueViolationError(fmt.Errorf("duplicate of %s", newModel.OriginalURL), shortURL)
		return "", conflictErr
	}
	//affRows, err := res.RowsAffected()
	//utils.SugaredLogger.Debugln("RowsAffected", affRows)

	//if err != nil {
	//	tx.Rollback()
	//	return "", err
	//}

	//conflict on INSERT
	//if affRows == 0 {
	//	tx.Rollback()
	//	return "", handleUniqueViolation(newModel.OriginalURL)
	//}

	//if err != nil {
	//	//on error roll back
	//	tx.Rollback()
	//
	//	var pqErr *pq.Error
	//	if errors.As(err, &pqErr); pqErr.Code == pgerrcode.UniqueViolation {
	//		return "", handleUniqueViolation(newModel.OriginalURL, pqErr)
	//	}
	//
	//	return "", err
	//}

	//commit
	err = tx.Commit()
	if err != nil {
		return "", err
	}

	return newModel.ShortURL, nil
}

func (ps *PostgresStorage) SaveBatch(ctx context.Context, arr []model.StorageModel) ([]model.StorageModel, error) {

	ps.Lock()
	defer ps.Unlock()

	//create transaction
	tx, err := dbConnection.Begin()
	if err != nil {
		return nil, err
	}

	//query
	query := "INSERT INTO urls (uuid, shortURL, originalURL) VALUES ($1, $2, $3);"

	//batch query
	for _, v := range arr {
		_, err = tx.ExecContext(ctx, query, v.UUID, v.ShortURL, v.OriginalURL)
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

func handleUniqueViolation(originalURL string) error {
	query := "SELECT shortURL FROM urls WHERE originalURL = $1;"
	row := dbConnection.QueryRowContext(context.Background(), query, originalURL)
	var shortURL string
	err := row.Scan(&shortURL)
	if err != nil {
		return err
	}

	return custom.NewUniqueViolationError(fmt.Errorf("duplicate of %s", originalURL), shortURL)
}
