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
	serverAddressKey   = "SERVER_ADDRESS"
	baseURLKey         = "BASE_URL"
	fileStoragePathKey = "FILE_STORAGE_PATH"

	defaultAddress         = "localhost:8080"
	defaultFileStoragePath = "internal/storage/someStorage.json"
)

type appConfig struct {
	ServerAddress string
	BaseURL       string
}

var (
	Address         appConfig
	StorageFileName string
)

func Setup() {
	//set default values
	config := appConfig{
		ServerAddress: defaultAddress,
		BaseURL:       defaultAddress,
	}
	StorageFileName = defaultFileStoragePath
	//utils.SugaredLogger.Infof("default \n\t\tServerAddress: %v \n\t\tBaseURL: %v \n\t\tStorageFileName: %v\n\n", config.ServerAddress, config.BaseURL, StorageFileName)

	getConfigFromFlags(&config)
	//utils.SugaredLogger.Infof("after flags \n\t\tServerAddress: %v \n\t\tBaseURL: %v \n\t\tStorageFileName: %v\n\n", config.ServerAddress, config.BaseURL, StorageFileName)

	getConfigFromEnv(&config)
	//utils.SugaredLogger.Infof("after enc \n\t\tServerAddress: %v \n\t\tBaseURL: %v \n\t\tStorageFileName: %v\n\n", config.ServerAddress, config.BaseURL, StorageFileName)

	Address = config

	utils.SugaredLogger.Infoln("configuration setup finished")
	utils.SugaredLogger.Debugln("launch address:  ", Address.ServerAddress)
	utils.SugaredLogger.Debugln("redirect address:", Address.BaseURL)
	utils.SugaredLogger.Debugln("StorageFileName: ", StorageFileName)

	//launchENV, redirectEnv := getEnv()
	//launchFlag, redirectFlag := getFlags()
	//
	//Address.ServerAddress, Address.BaseURL = launchENV, redirectEnv
	//
	//if launchENV == "" {
	//	Address.ServerAddress = launchFlag
	//}
	//
	//if redirectEnv == "" {
	//	Address.BaseURL = redirectFlag
	//}
	//
	//setFileStoragePath()

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

	flag.Parse()

}

// getConfigFromEnv will parse environment values, if values aren't empty will overwrite default values
func getConfigFromEnv(config *appConfig) {
	//server launch
	launchEnv, err := validateAddress(os.Getenv(serverAddressKey))
	if err != nil {
		utils.SugaredLogger.Debugf("environment value %s err: %s", launchEnv, err)
	}
	if launchEnv != "" {
		config.ServerAddress = launchEnv
	}

	//redirect url
	redirectEnv, err := validateAddress(os.Getenv(baseURLKey))
	if err != nil {
		utils.SugaredLogger.Debugf("environment value %s err: %s", redirectEnv, err)
	}
	if redirectEnv != "" {
		config.BaseURL = redirectEnv
	}

	storagePathEnv := os.Getenv(fileStoragePathKey)
	if redirectEnv != "" {
		StorageFileName = storagePathEnv
	} else {
		utils.SugaredLogger.Debugf("environment value %s err: %s", storagePathEnv, err)
	}
}

//func getEnv() (string, string) {
//	lEnv := os.Getenv(serverAddressKey)
//	rEnv := os.Getenv(baseURLKey)
//
//	l, err := validateAddress(lEnv)
//	if err != nil {
//		utils.SugaredLogger.Debugln("environment variable", serverAddressKey, ":", err)
//	}
//
//	r, err := validateAddress(rEnv)
//	if err != nil {
//		utils.SugaredLogger.Debugln("environment variable", baseURLKey, ":", err)
//	}
//
//	return l, r
//}

//func getFlags() (string, string) {
//	lFlag := flag.String("a", defaultAddress, "start server - host:port")
//	rFlag := flag.String("b", defaultAddress, "redirect url - host:port")
//	flag.Parse()
//
//	l, err := validateAddress(*lFlag)
//	if err != nil {
//		utils.SugaredLogger.Debugln("launch flag error:", err)
//	}
//
//	r, err := validateAddress(*rFlag)
//	if err != nil {
//		utils.SugaredLogger.Debugln("redirect flag error:", err)
//	}
//
//	return l, r
//}

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

//func setFileStoragePath() {
//	filePathEnv := os.Getenv(fileStoragePathKey)
//	StorageFileName = filePathEnv
//
//	if filePathEnv != "" {
//		return
//	}
//	filePathFlag := flag.String("f", defaultFileStoragePath, "/path/to/file.extension")
//	flag.Parse()
//	StorageFileName = *filePathFlag
//}
