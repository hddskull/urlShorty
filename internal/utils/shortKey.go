package utils

import (
	"math/rand"
	"os/exec"
)

const (
	charset   = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	keyLength = 6

	UUIDCommand = "uuidgen"
)

func GenerateShortKey() string {
	shortKey := make([]byte, keyLength)
	for i := range shortKey {
		shortKey[i] = charset[rand.Intn(len(charset))]
	}
	return string(shortKey)
}

func GenerateUUID() (string, error) {
	newUUID, err := exec.Command(UUIDCommand).Output()
	if err != nil {
		return "", err
	}
	strUUID := string(newUUID)
	strUUID = strUUID[:len(strUUID)-1]
	return strUUID, nil
}
