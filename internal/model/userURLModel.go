package model

import (
	"fmt"
	"github.com/hddskull/urlShorty/config"
)

type UserURLModel struct {
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

func NewUserURLModel(shortURL, originalURL string) *UserURLModel {
	return &UserURLModel{
		ShortURL:    fmt.Sprint("http://", config.Address.BaseURL, "/", shortURL),
		OriginalURL: originalURL,
	}
}
