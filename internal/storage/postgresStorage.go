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
		uuid TEXT PRIMARY KEY,
		shortURL TEXT NOT NULL,
		originalURL TEXT NOT NULL UNIQUE,
		session_id TEXT,
		created_at TIMESTAMP NOT NULL DEFAULT now(),
	    updated_at TIMESTAMP NOT NULL DEFAULT now()
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
	sessionID, ok := model.SessionIDFromContext(ctx)
	utils.SugaredLogger.Debugln("sessionID:", sessionID, "| ok:", ok)

	if !ok {
		return "", custom.ErrNoSessionID
	}

	//create model
	newModel, err := model.NewFileStorageModel(u, "", sessionID)
	if err != nil {
		return "", err
	}

	//create transaction
	tx, err := dbConnection.Begin()
	if err != nil {
		return "", err
	}

	//NewQuery
	query := `
		INSERT INTO urls (uuid, shortURL, originalURL, session_id)
			VALUES ($1, $2, $3, $4)
			ON CONFLICT (originalURL)
			DO UPDATE SET
				updated_at = now(),
				session_id = excluded.session_id
			RETURNING (created_at = updated_at) as is_new, shortURL
	`

	row := dbConnection.QueryRowContext(ctx, query, newModel.UUID, newModel.ShortURL, newModel.OriginalURL, newModel.SessionID)
	var isNew bool
	var shortURL string
	err = row.Scan(&isNew, &shortURL)
	utils.SugaredLogger.Debugln("Save(): isNew:", isNew, "|shortURL:", shortURL, "|err:", err)

	if err != nil {
		tx.Rollback()
		return "", err
	}
	if !isNew {
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
	sessionID, ok := model.SessionIDFromContext(ctx)
	utils.SugaredLogger.Debugln("sessionID:", sessionID, "| ok:", ok)

	if !ok {
		return custom.ErrNoSessionID
	}

	//create transaction
	tx, err := dbConnection.Begin()
	if err != nil {
		return err
	}

	//query
	query := "INSERT INTO urls (uuid, shortURL, originalURL, session_id) VALUES ($1, $2, $3, $4);"

	//batch query
	for _, v := range arr {
		_, err = tx.ExecContext(ctx, query, v.UUID, v.ShortURL, v.OriginalURL, sessionID)
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

func (ps *PostgresStorage) GetUserURLs(ctx context.Context) (*[]model.UserURLModel, error) {
	sessionID, ok := model.SessionIDFromContext(ctx)
	utils.SugaredLogger.Debugln("sessionID:", sessionID, "| ok:", ok)
	if !ok {
		return nil, custom.ErrNoSessionID
	}

	query := "SELECT shorturl, originalURL FROM urls WHERE session_id = $1;"

	rows, err := dbConnection.QueryContext(ctx, query, sessionID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	urlsSlice := make([]model.UserURLModel, 0)
	for rows.Next() {
		var (
			shortURL    string
			originalURL string
		)
		if err = rows.Scan(&shortURL, &originalURL); err != nil {
			return nil, err
		}
		urlsModel := model.NewUserURLModel(shortURL, originalURL)
		urlsSlice = append(urlsSlice, *urlsModel)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return &urlsSlice, nil
}

func (ps *PostgresStorage) Ping(ctx context.Context) error {
	err := dbConnection.Ping()
	if err != nil {
		utils.SugaredLogger.Errorln(err)
	}
	return err
}
