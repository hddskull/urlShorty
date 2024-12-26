package root

import (
	"compress/gzip"
	"errors"
	"fmt"
	"github.com/hddskull/urlShorty/internal/storage"
	"github.com/hddskull/urlShorty/internal/utils"
	"github.com/hddskull/urlShorty/tools/custom"
	"io"
	"net/http"
	"strings"

	"github.com/hddskull/urlShorty/config"
)

func GetHandler(w http.ResponseWriter, r *http.Request) {
	arr := strings.Split(r.URL.Path, "/")
	id := arr[len(arr)-1]

	url, err := storage.Current.Get(r.Context(), id)
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
	utils.SugaredLogger.Debugln("PostHandler()")

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
	if bodyS == "" {
		utils.SugaredLogger.Debugln("Save() empty arg:", custom.ErrEmptyURL)
		custom.JSONError(w, custom.ErrEmptyURL, http.StatusBadRequest)
	}

	id, err := storage.Current.Save(r.Context(), bodyS)
	utils.SugaredLogger.Debugf("storage.Current.Save(...) ID: %v|| err: %v", id, err)
	if err != nil {
		var uvError *custom.UniqueViolationError
		if errors.As(err, &uvError) {
			fullID := fmt.Sprint("http://", config.Address.BaseURL, "/", uvError.ShortURL)
			w.Header().Add("Content-Type", "text/plain")
			w.WriteHeader(http.StatusConflict)
			w.Write([]byte(fullID))
			return
		}

		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	utils.SugaredLogger.Debugln("(root)/ PostHandler() successfully saved", bodyS)
	fullID := fmt.Sprint("http://", config.Address.BaseURL, "/", id)

	w.Header().Add("Content-Type", "text/plain")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(fullID))
}

func PingHandler(w http.ResponseWriter, r *http.Request) {
	err := storage.Current.Ping(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
}
