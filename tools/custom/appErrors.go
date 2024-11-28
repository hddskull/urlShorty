package custom

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

type ErrorResponseModel struct {
	Message string `json:"message"`
}

var ErrNoServerAddress = errors.New("no server address")
var ErrInvalidAddressPattern = errors.New("invalid host:port")

var ErrEmptyURL = errors.New("empty url")
var ErrEmptyPath = errors.New("empty path")
var ErrEmptyEnvVar = errors.New("empty environment variable")
var ErrEmptyBatch = errors.New("empty batch")
var ErrInvalidBatch = errors.New("invalid model in batch")

var ErrFuncUnsupported = errors.New("current implementation doesn't support such function")

func NoURLBy(id string) error {
	return fmt.Errorf("no url by id: %s", id)
}

// JSONError equivalent to http.Error(...), but content type is "application/json"
func JSONError(w http.ResponseWriter, err interface{}, code int) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(err)
}
