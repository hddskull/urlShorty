package storage

import (
	"context"
	"encoding/json"
	"github.com/hddskull/urlShorty/config"
	"github.com/hddskull/urlShorty/internal/model"
	"github.com/hddskull/urlShorty/internal/utils"
	"github.com/hddskull/urlShorty/tools/custom"
	"os"
	"path/filepath"
)

type FileStorage struct {
}

func newFileStorage() *FileStorage {
	return &FileStorage{}
}

// Storage interface
var _ Storage = newFileStorage()

func (fs FileStorage) Setup() error {
	_, err := os.OpenFile(config.StorageFileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	//if file opened - return
	if err == nil {
		return nil
	}
	//else try to create dir
	if os.IsNotExist(err) {
		dir := filepath.Dir(config.StorageFileName)
		err = os.MkdirAll(dir, 0777)
		if err != nil {
			return err
		}

		_, err = os.OpenFile(config.StorageFileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			return err
		}

		return nil
	}

	return err
}

func (fs FileStorage) Close() error {
	return custom.ErrFuncUnsupported
}

func (fs FileStorage) Save(ctx context.Context, u string) (string, error) {
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

	sessionID, ok := model.SessionIDFromContext(ctx)

	if !ok {
		return "", custom.ErrNoSessionID
	}

	storageModel, err := model.NewFileStorageModel(u, "", sessionID)
	if err != nil {
		return "", err
	}

	//save model
	err = fs.saveToFile(storageModel)
	if err != nil {
		return "", err
	}

	return storageModel.ShortURL, nil
}

func (fs FileStorage) SaveBatch(ctx context.Context, arr []model.StorageModel) error {
	err := fs.saveBatchToFile(&arr)
	if err != nil {
		return err
	}

	return nil
}

func (fs FileStorage) Get(ctx context.Context, id string) (string, error) {
	if id == "" {
		return "", custom.ErrEmptyURL
	}

	originalURL, err := fs.getFromFile(id)
	if err != nil {
		return "", err
	}

	return originalURL, nil
}

func (fs FileStorage) GetUserURLs(ctx context.Context) (*[]model.UserURLModel, error) {
	return nil, custom.ErrFuncUnsupported
}

func (fs FileStorage) Ping(ctx context.Context) error {
	return custom.ErrFuncUnsupported
}

func (fs FileStorage) BatchMarkDeleted(sessionID string, shortURLs ...string) {
}

// Supporting methods

func (fs FileStorage) readAllFromFile() ([]model.StorageModel, error) {
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

	modelSlice := []model.StorageModel{}
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

func (fs FileStorage) checkExistence(originalURL string) (*model.StorageModel, error) {
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

func (fs FileStorage) saveToFile(model *model.StorageModel) error {

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

func (fs FileStorage) saveBatchToFile(batch *[]model.StorageModel) error {
	//read from file
	modelSlice, err := fs.readAllFromFile()
	if err != nil {
		return err
	}

	//append new Data
	modelSlice = append(modelSlice, *batch...)

	//to json
	data, err := json.Marshal(modelSlice)
	if err != nil {
		return err
	}

	//write to file
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
