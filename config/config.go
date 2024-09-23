package config

import (
	//"errors"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/hddskull/urlShorty/tools/errors"
)

type appConfig struct {
	ServerAddress string
	BaseURL       string
}

var Address = new(appConfig)

func Setup() {
	launchENV, redirectEnv := getEnv()
	lFlag, rFlag := getFlags()

	Address.ServerAddress, Address.BaseURL = launchENV, redirectEnv

	if launchENV == "" {
		Address.ServerAddress = lFlag
	}

	if redirectEnv == "" {
		Address.BaseURL = rFlag
	}

	fmt.Println("launchENV:  ", launchENV, "redirectEnv:  ", redirectEnv)
	fmt.Println("lFlag: ", lFlag, "rFlag: ", rFlag)
	fmt.Println("Address.ServerAddress: ", Address.ServerAddress)
	fmt.Println("Address.BaseURL: ", Address.BaseURL)
}

func getEnv() (string, string) {
	lEnv := os.Getenv("SERVER_ADDRESS")
	rEnv := os.Getenv("BASE_URL")

	l, err := validateAddress(lEnv)
	if err != nil {
		fmt.Println(err)
	}

	r, err := validateAddress(rEnv)
	if err != nil {
		fmt.Println(err)
	}

	return l, r
}

func getFlags() (string, string) {
	//default
	def := "localhost:8080"

	lFlag := flag.String("a", def, "start server - host:port")
	rFlag := flag.String("b", def, "redirect url - host:port")
	flag.Parse()

	l, err := validateAddress(*lFlag)
	if err != nil {
		fmt.Println(err)
	}

	r, err := validateAddress(*rFlag)
	if err != nil {
		fmt.Println(err)
	}

	return l, r
}

func validateAddress(adr string) (string, error) {

	if adr == "" {
		return "", errors.NoServerAddress //errors.New("no server address")
	}

	if i := strings.Index(adr, "https://"); i != -1 {
		adr = adr[(len("https://")):]
	}

	if i := strings.Index(adr, "http://"); i != -1 {
		adr = adr[(len("http://")):]
	}

	vals := strings.Split(adr, ":")

	if len(vals) != 2 {
		return "", errors.InvalidAddressPattern
	}

	_, err := strconv.Atoi(vals[1])

	if err != nil {
		return "", err
	}

	return adr, nil
}

//
//// Properties
//var LaunchAdr = defaultNetAddress()
//var RedirectAdr = defaultNetAddress()
//
//func ConfigureNetAddress() {
//	flag.Var(LaunchAdr, "a", "Network address host:port")
//	flag.Var(RedirectAdr, "b", "Network address host:port")
//	flag.Parse()
//}
