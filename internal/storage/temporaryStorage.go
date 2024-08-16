package storage

import (
	"errors"
	"fmt"
	"math/rand"
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
		return "", errors.New("empty url")
	}

	id := generateShortKey()
	ts.urls[id] = u

	return id, nil
}

func (ts TemporaryStorage) Get(id string) (string, error) {
	if id == "" {
		return "", errors.New("empty url")
	}

	url, ok := ts.urls[id]
	if !ok {
		return "", fmt.Errorf("no url by id: %s", id)
	}

	return url, nil
}

func generateShortKey() string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	const keyLength = 6

	// rand.Seed(time.Now().UnixNano())
	shortKey := make([]byte, keyLength)
	for i := range shortKey {
		shortKey[i] = charset[rand.Intn(len(charset))]
	}
	return string(shortKey)
}
