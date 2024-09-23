package errors

import (
	"errors"
	"fmt"
)

var NoServerAddress = errors.New("no server address")
var InvalidAddressPattern = errors.New("invalid host:port")

var EmptyURL = errors.New("empty url")

func NoURLBy(id string) error {
	return fmt.Errorf("no url by id: %s", id)
}
