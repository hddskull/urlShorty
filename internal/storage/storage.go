package storage

import (
	"github.com/hddskull/urlShorty/config"
	"github.com/hddskull/urlShorty/internal/utils"
)

type Storage interface {
	Setup() error
	Save(u string) (string, error)
	Get(id string) (string, error)
	Ping() error
	Close() error
}

var Current Storage

// SetupStorage call in main to init var Current and create storage file
func SetupStorage() {

	if config.DBCredentials != "" {
		Current = newPostgresStorage()
		err := Current.Setup()
		if err != nil {
			panic(err)
		}
		utils.SugaredLogger.Debugln("Current Storage type: postgres")
	} else if config.StorageFileName != "" {
		Current = newFileStorage()
		err := Current.Setup()
		if err != nil {
			panic(err)
		}
		utils.SugaredLogger.Debugln("Current Storage type: file")
	} else {
		Current = NewTemporaryStorage()
		err := Current.Setup()
		if err != nil {
			panic(err)
		}
		utils.SugaredLogger.Debugln("Current Storage type: RAM memory")
	}

}
