package storage

import (
	"encoding/json"
	"github.com/hddskull/urlShorty/config"
	"github.com/hddskull/urlShorty/internal/utils"
	"github.com/hddskull/urlShorty/tools/custom"
	"os"
)

type fileStorageModel struct {
	UUID        string `json:"uuid"`
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

func newFileStorageModel(originalURL string) (*fileStorageModel, error) {
	//create uuid
	uuidBytes, err := utils.GenerateUUID()
	if err != nil {
		return nil, err
	}

	//create shortURL
	shortURL := utils.GenerateShortKey()

	return &fileStorageModel{
		UUID:        string(uuidBytes),
		ShortURL:    shortURL,
		OriginalURL: originalURL,
	}, nil
}

type FileStorage struct {
}

func newFileStorage() *FileStorage {
	return &FileStorage{}
}

var _ Storage = newFileStorage()

func (fs FileStorage) Save(u string) (string, error) {
	if u == "" {
		utils.SugaredLogger.Debugln("Save() empty arg:", custom.ErrEmptyURL)
		return "", custom.ErrEmptyURL
	}

	//check if url is already saved
	existingModel, err := fs.checkExistence(u)
	if err != nil {
		utils.SugaredLogger.Debugln("checkExistence() err:", err)
		return "", err
	}

	//if already saved exit func without error
	if existingModel != nil {
		utils.SugaredLogger.Debugln("Save() url already saved:", existingModel)
		return existingModel.ShortURL, nil
	}

	model, err := newFileStorageModel(u)
	if err != nil {
		return "", err
	}

	//save model
	err = fs.saveToFile(model)
	if err != nil {
		return "", err
	}

	return model.ShortURL, nil
}

func (fs FileStorage) Get(id string) (string, error) {
	if id == "" {
		return "", custom.ErrEmptyURL
	}

	originalURL, err := fs.getFromFile(id)
	if err != nil {
		return "", err
	}

	return originalURL, nil
}

func (fs FileStorage) readAllFromFile() ([]fileStorageModel, error) {
	filename := config.StorageFileName

	fileBytes, err := os.ReadFile(filename)
	if err != nil {
		utils.SugaredLogger.Debugln("readAllFromFile() couldn't read from file", filename)
		return nil, err
	}

	if len(fileBytes) == 0 {
		utils.SugaredLogger.Debugln("readAllFromFile() len(fileBytes) == 0")
		return nil, nil
	}

	modelSlice := []fileStorageModel{}
	err = json.Unmarshal(fileBytes, &modelSlice)
	if err != nil {
		utils.SugaredLogger.Debugln("readAllFromFile() couldn't unmarshal to slice", filename)
		return nil, err
	}

	return modelSlice, nil
}

func (fs FileStorage) getFromFile(id string) (string, error) {
	modelSlice, err := fs.readAllFromFile()
	if err != nil {
		return "", err
	}

	for _, model := range modelSlice {
		if model.ShortURL == id {
			return model.OriginalURL, nil
		}
	}

	return "", custom.NoURLBy(id)
}

func (fs FileStorage) checkExistence(originalURL string) (*fileStorageModel, error) {
	models, err := fs.readAllFromFile()
	if err != nil {
		utils.SugaredLogger.Debugln("checkExistence() readAllFromFile error", err)
		return nil, err
	}

	for _, m := range models {
		if m.OriginalURL == originalURL {
			return &m, nil
		}
	}

	return nil, nil
}

func (fs FileStorage) saveToFile(model *fileStorageModel) error {

	modelSlice, err := fs.readAllFromFile()
	if err != nil {
		return err
	}

	modelSlice = append(modelSlice, *model)

	data, err := json.Marshal(modelSlice)
	if err != nil {
		return err
	}

	filename := config.StorageFileName

	file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.Write(data)
	if err != nil {
		return err
	}

	return nil
}
