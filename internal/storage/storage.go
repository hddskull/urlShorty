package storage

import (
	"context"
	"github.com/hddskull/urlShorty/config"
	"github.com/hddskull/urlShorty/internal/model"
	"github.com/hddskull/urlShorty/internal/utils"
)

type Storage interface {
	Setup() error
	Close() error
	Save(ctx context.Context, u string) (string, error)
	SaveBatch(ctx context.Context, arr []model.StorageModel) error
	Get(ctx context.Context, id string) (string, error)
	Ping(ctx context.Context) error
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
