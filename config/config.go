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
	serverAddressKey = "SERVER_ADDRESS"
	baseURLKey       = "BASE_URL"
	defaultAddress   = "localhost:8080"
)

type appConfig struct {
	ServerAddress string
	BaseURL       string
}

var Address = new(appConfig)

func Setup() {
	launchENV, redirectEnv := getEnv()
	launchFlag, redirectFlag := getFlags()

	Address.ServerAddress, Address.BaseURL = launchENV, redirectEnv

	if launchENV == "" {
		Address.ServerAddress = launchFlag
	}

	if redirectEnv == "" {
		Address.BaseURL = redirectFlag
	}

	utils.SugaredLogger.Infow("configuration setup finished", "launch address:", Address.ServerAddress, "redirect address:", Address.BaseURL)
}

func getEnv() (string, string) {
	lEnv := os.Getenv(serverAddressKey)
	rEnv := os.Getenv(baseURLKey)

	l, err := validateAddress(lEnv)
	if err != nil {
		utils.SugaredLogger.Debugln("environment variable", serverAddressKey, ":", err)
	}

	r, err := validateAddress(rEnv)
	if err != nil {
		utils.SugaredLogger.Debugln("environment variable", baseURLKey, ":", err)
	}

	return l, r
}

func getFlags() (string, string) {
	lFlag := flag.String("a", defaultAddress, "start server - host:port")
	rFlag := flag.String("b", defaultAddress, "redirect url - host:port")
	flag.Parse()

	l, err := validateAddress(*lFlag)
	if err != nil {
		utils.SugaredLogger.Debugln("launch flag error:", err)
	}

	r, err := validateAddress(*rFlag)
	if err != nil {
		utils.SugaredLogger.Debugln("redirect flag error:", err)
	}

	return l, r
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
