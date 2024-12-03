package root

import (
	"compress/gzip"
	"errors"
	"fmt"
	"github.com/hddskull/urlShorty/internal/storage"
	"github.com/hddskull/urlShorty/tools/custom"
	"io"
	"net/http"
	"strings"

	"github.com/hddskull/urlShorty/config"
)

func GetHandler(w http.ResponseWriter, r *http.Request) {
	arr := strings.Split(r.URL.Path, "/")
	id := arr[len(arr)-1]

	url, err := storage.Current.Get(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Add("Content-Type", "text/plain")
	w.Header().Add("Location", url)
	w.WriteHeader(http.StatusTemporaryRedirect)
}

func PostHandler(w http.ResponseWriter, r *http.Request) {
	reader := r.Body

	if r.Header.Get(`Content-Encoding`) == `gzip` {
		gz, err := gzip.NewReader(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		reader = gz
	}
	defer reader.Close()

	bodyB, err := io.ReadAll(reader)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	bodyS := string(bodyB)
	id, err := storage.Current.Save(bodyS)
	if err != nil {
		var uvError *custom.UniqueViolationError
		if errors.As(err, &uvError) {
			fullID := fmt.Sprint("http://", config.Address.BaseURL, "/", uvError.ShortURL)
			http.Error(w, fullID, http.StatusConflict)
			return
		}

		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	fullID := fmt.Sprint("http://", config.Address.BaseURL, "/", id)

	w.Header().Add("Content-Type", "text/plain")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(fullID))
}

func PingHandler(w http.ResponseWriter, r *http.Request) {
	err := storage.Current.Ping()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
}
