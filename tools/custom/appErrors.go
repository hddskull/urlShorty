package custom

import (
	"errors"
	"fmt"
)

var ErrNoServerAddress = errors.New("no server address")
var ErrInvalidAddressPattern = errors.New("invalid host:port")

var ErrEmptyURL = errors.New("empty url")

func NoURLBy(id string) error {
	return fmt.Errorf("no url by id: %s", id)
}
