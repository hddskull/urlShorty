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

type UniqueViolationError struct {
	Err      error
	ShortURL string
}

func (e *UniqueViolationError) Error() string {
	return fmt.Sprintf("%v; original: %v", e.Err, e.ShortURL)
}

func (e *UniqueViolationError) Unwrap() error {
	return e.Err
}

func NewUniqueViolationError(err error, shortURL string) error {
	return &UniqueViolationError{
		Err:      err,
		ShortURL: shortURL,
	}
}

type CookieError struct {
	Err        error
	HTTPStatus int
}

func (e *CookieError) Error() string {
	return fmt.Sprintf("%v; original: %v", e.Err, e.HTTPStatus)
}

func (e *CookieError) Unwrap() error {
	return e.Err
}

func NewCookieError(err error, httpStatus int) error {
	return &CookieError{
		Err:        err,
		HTTPStatus: httpStatus,
	}
}

var ErrNoServerAddress = errors.New("no server address")
var ErrInvalidAddressPattern = errors.New("invalid host:port")

var ErrEmptyURL = errors.New("empty url")
var ErrEmptyPath = errors.New("empty path")
var ErrEmptyEnvVar = errors.New("empty environment variable")
var ErrEmptyBatch = errors.New("empty batch")
var ErrInvalidBatch = errors.New("invalid model in batch")

var ErrFuncUnsupported = errors.New("current implementation doesn't support such function")

var ErrUnauthorized = errors.New("unauthorized")

var ErrUnknown = errors.New("unknown error")

func NoURLBy(id string) error {
	return fmt.Errorf("no url by id: %s", id)
}

// JSONError equivalent to http.Error(...), but content type is "application/json"
func JSONError(w http.ResponseWriter, err error, code int) {
	formattedError := ErrorResponseModel{Message: err.Error()}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(formattedError)
}
