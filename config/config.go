package config

import (
	"flag"
	"os"
	"strconv"
	"strings"

	"github.com/hddskull/urlShorty/internal/utils"
	"github.com/hddskull/urlShorty/tools/custom"
)

const (
	//Environment property keys
	serverAddressKey   = "SERVER_ADDRESS"
	baseURLKey         = "BASE_URL"
	fileStoragePathKey = "FILE_STORAGE_PATH"
	dbCredentialsKey   = "DATABASE_DSN"

	//defaults
	defaultAddress         = "localhost:8080"
	DefaultFileStoragePath = "internal/storage/someStorage.json"
	defaultDBCredentials   = "host=localhost port=5432 user=postgres password=password dbname=urlshorty sslmode=disable"
)

type appConfig struct {
	ServerAddress string
	BaseURL       string
}

var (
	Address         appConfig
	StorageFileName string
	DBCredentials   string
)

func Setup() {
	//set default values
	config := appConfig{
		ServerAddress: defaultAddress,
		BaseURL:       defaultAddress,
	}
	StorageFileName = DefaultFileStoragePath
	DBCredentials = defaultDBCredentials

	getConfigFromFlags(&config)

	getConfigFromEnv(&config)

	Address = config

	utils.SugaredLogger.Infoln("configuration setup finished")
	utils.SugaredLogger.Debugln("launch address:  ", Address.ServerAddress)
	utils.SugaredLogger.Debugln("redirect address:", Address.BaseURL)
	utils.SugaredLogger.Debugln("StorageFileName: ", StorageFileName)
	utils.SugaredLogger.Debugln("DBCredentials: ", DBCredentials)
}

// getConfigFromFlags will parse flags ["a", "b", "f"]
func getConfigFromFlags(config *appConfig) {
	//flag "a" - server launch address
	flag.Func("a", "start server - host:port", func(s string) error {
		formattedAddr, err := validateAddress(s)
		if err != nil {
			utils.SugaredLogger.Debugf("flag %s error: %s", "a", err)
			return err
		}
		config.ServerAddress = formattedAddr
		return nil
	})

	//flag "b" - redirect address
	flag.Func("b", "redirect url - host:port", func(s string) error {
		formattedAddr, err := validateAddress(s)
		if err != nil {
			utils.SugaredLogger.Debugf("flag %s error: %s", "b", err)
			return err
		}
		config.BaseURL = formattedAddr
		return nil
	})

	//flag "f" - file storage path
	flag.Func("f", "/path/to/file.extension", func(s string) error {
		if s == "" {
			err := custom.ErrEmptyPath
			utils.SugaredLogger.Debugf("flag %s error: %s", "b", err)
			return err
		}
		StorageFileName = s
		return nil
	})

	//flag "d" - credentials for db connect
	flag.Func("d", "credentials for db connect", func(s string) error {
		if s == "" {
			err := custom.ErrEmptyPath
			utils.SugaredLogger.Debugf("flag %s error: %s", "d", err)
			return err
		}
		DBCredentials = s
		return nil
	})

	flag.Parse()
}

// getConfigFromEnv will parse environment values, if values aren't empty will overwrite default values
func getConfigFromEnv(config *appConfig) {
	//server launch
	launchEnv, err := validateAddress(os.Getenv(serverAddressKey))
	if err != nil {
		utils.SugaredLogger.Debugf("env %s err: %s", baseURLKey, err)
	}
	if launchEnv != "" {
		config.ServerAddress = launchEnv
	}

	//redirect url
	redirectEnv, err := validateAddress(os.Getenv(baseURLKey))
	if err != nil {
		utils.SugaredLogger.Debugf("env %s err: %s", baseURLKey, err)
	}
	if redirectEnv != "" {
		config.BaseURL = redirectEnv
	}

	//storage path
	storagePathEnv := os.Getenv(fileStoragePathKey)
	if storagePathEnv == "" {
		utils.SugaredLogger.Debugf("env %s err: %s", fileStoragePathKey, custom.ErrEmptyEnvVar)
	} else {
		StorageFileName = storagePathEnv
	}

	//database credentials
	dbCredentialsEnv := os.Getenv(dbCredentialsKey)
	if dbCredentialsEnv == "" {
		utils.SugaredLogger.Debugf("env %s err: %s", fileStoragePathKey, custom.ErrEmptyEnvVar)
	} else {
		StorageFileName = dbCredentialsEnv
	}
}

func validateAddress(adr string) (string, error) {

	if adr == "" {
		return "", custom.ErrNoServerAddress
	}

	if i := strings.Index(adr, "https://"); i != -1 {
		adr = adr[(len("https://")):]
	}

	if i := strings.Index(adr, "http://"); i != -1 {
		adr = adr[(len("http://")):]
	}

	vals := strings.Split(adr, ":")

	if len(vals) != 2 {
		return "", custom.ErrInvalidAddressPattern
	}

	_, err := strconv.Atoi(vals[1])

	if err != nil {
		return "", err
	}

	return adr, nil
}
