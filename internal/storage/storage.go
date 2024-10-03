package storage

import (
	"github.com/hddskull/urlShorty/config"
	"os"
	"path/filepath"
)

type Storage interface {
	Save(u string) (string, error)
	Get(id string) (string, error)
}

var Current Storage = newFileStorage()

// SetupStorage call in main to init var Current and create storage file
func SetupStorage() {
	//create storage file
	if config.StorageFileName == "" {
		config.StorageFileName = config.DefaultFileStoragePath
	}

	_, err := os.OpenFile(config.StorageFileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	//if file opened - return
	if err == nil {
		return
	}
	//else try to create dir
	if os.IsNotExist(err) {
		dir := filepath.Dir(config.StorageFileName)
		err = os.MkdirAll(dir, 0777)
		if err != nil {
			panic(err)
		}

		_, err = os.OpenFile(config.StorageFileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			panic(err)
		}

		return
	}

	panic(err)
}
