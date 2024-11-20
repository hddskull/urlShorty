package model

import "github.com/hddskull/urlShorty/internal/utils"

type StorageModel struct {
	UUID        string `json:"uuid"`
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

func NewFileStorageModel(originalURL string) (*StorageModel, error) {
	//create uuid
	uuid, err := utils.GenerateUUID()
	if err != nil {
		return nil, err
	}

	//create shortURL
	shortURL := utils.GenerateShortKey()

	return &StorageModel{
		UUID:        uuid,
		ShortURL:    shortURL,
		OriginalURL: originalURL,
	}, nil
}
