package storage

import (
	"math/rand"

	"github.com/hddskull/urlShorty/tools/custom"
)

const (
	charset   = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	keyLength = 6
)

type TemporaryStorage struct {
	urls map[string]string
}

func newTemporaryStorage() *TemporaryStorage {
	return &TemporaryStorage{
		urls: make(map[string]string),
	}
}

var TempStorage Storage = newTemporaryStorage()

func (ts TemporaryStorage) Save(u string) (string, error) {
	if u == "" {
		return "", custom.ErrEmptyURL
	}

	id := generateShortKey()
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

func generateShortKey() string {
	shortKey := make([]byte, keyLength)
	for i := range shortKey {
		shortKey[i] = charset[rand.Intn(len(charset))]
	}
	return string(shortKey)
}
