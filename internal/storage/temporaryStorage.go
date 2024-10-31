package storage

import (
	"github.com/hddskull/urlShorty/internal/utils"
	"github.com/hddskull/urlShorty/tools/custom"
)

type TemporaryStorage struct {
	urls map[string]string
}

// interface compliance check
var _ Storage = NewTemporaryStorage()

func NewTemporaryStorage() *TemporaryStorage {
	return &TemporaryStorage{
		urls: make(map[string]string),
	}
}

func (ts TemporaryStorage) Save(u string) (string, error) {
	if u == "" {
		return "", custom.ErrEmptyURL
	}

	id := utils.GenerateShortKey()
	ts.urls[id] = u

	return id, nil
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
