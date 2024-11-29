package storage

import (
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

func (ts TemporaryStorage) Save(u string) (string, error) {
	if u == "" {
		return "", custom.ErrEmptyURL
	}

	id := utils.GenerateShortKey()
	ts.urls[id] = u

	return id, nil
}

func (ts TemporaryStorage) SaveBatch(arr []model.StorageModel) ([]model.StorageModel, error) {

	for _, v := range arr {
		ts.urls[v.UUID] = v.OriginalURL
	}

	return arr, nil
}

func (ts TemporaryStorage) Get(id string) (string, error) {
	if id == "" {
		return "", custom.ErrEmptyURL
	}

	url, ok := ts.urls[id]
	if !ok {
		return "", custom.NoURLBy(id)
	}

	return url, nil
}

func (ts TemporaryStorage) Ping() error {
	return custom.ErrFuncUnsupported
}

func (ts TemporaryStorage) Close() error {
	return custom.ErrFuncUnsupported
}
