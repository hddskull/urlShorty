package model

import (
	"context"
	"github.com/hddskull/urlShorty/internal/utils"
	"github.com/hddskull/urlShorty/tools/custom"
)

type StorageModel struct {
	UUID        string `json:"uuid"`
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
	SessionID   string `json:"session_id"`
}

type key string

var SessionIDKey key = "sessionID"

func NewFileStorageModel(originalURL, correlationID, sessionID string) (*StorageModel, error) {
	//create uuid
	var err error
	if correlationID == "" {
		correlationID, err = utils.GenerateUUID()
	}
	if sessionID == "" {
		return nil, custom.ErrNoSessionID
	}

	if err != nil {
		return nil, err
	}

	//create shortURL
	shortURL := utils.GenerateShortKey()

	return &StorageModel{
		UUID:        correlationID,
		ShortURL:    shortURL,
		OriginalURL: originalURL,
		SessionID:   sessionID,
	}, nil
}

func NewContextWithSessionID(ctx context.Context, sessionID string) context.Context {
	return context.WithValue(ctx, SessionIDKey, sessionID)
}

func SessionIDFromContext(ctx context.Context) (string, bool) {
	sessionID, ok := ctx.Value(SessionIDKey).(string)
	return sessionID, ok
}
