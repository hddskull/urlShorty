package storage

import (
	"context"
	"github.com/hddskull/urlShorty/internal/model"
	"github.com/hddskull/urlShorty/internal/utils"
	"github.com/hddskull/urlShorty/tools/custom"
)

type TemporaryStorage struct {
	urls map[string]string
}

func NewTemporaryStorage() *TemporaryStorage {
	return &TemporaryStorage{
		urls: make(map[string]string),
	}
}

// Storage interface
var _ Storage = NewTemporaryStorage()

func (ts TemporaryStorage) Setup() error {
	return nil
}

func (ts TemporaryStorage) Close() error {
	return custom.ErrFuncUnsupported
}

func (ts TemporaryStorage) Save(ctx context.Context, u string) (string, error) {
	if u == "" {
		return "", custom.ErrEmptyURL
	}

	id := utils.GenerateShortKey()
	ts.urls[id] = u

	return id, nil
}

func (ts TemporaryStorage) SaveBatch(ctx context.Context, arr []model.StorageModel) error {

	for _, v := range arr {
		ts.urls[v.UUID] = v.OriginalURL
	}

	return nil
}

func (ts TemporaryStorage) Get(ctx context.Context, id string) (string, error) {
	if id == "" {
		return "", custom.ErrEmptyURL
	}

	url, ok := ts.urls[id]
	if !ok {
		return "", custom.NoURLBy(id)
	}

	return url, nil
}

func (ts TemporaryStorage) GetUserURLs(ctx context.Context) (*[]model.UserURLModel, error) {
	return nil, custom.ErrFuncUnsupported
}

func (ts TemporaryStorage) BatchMarkDeleted(sessionID string, shortURLs ...string) {
}

func (ts TemporaryStorage) Ping(ctx context.Context) error {
	return custom.ErrFuncUnsupported
}
