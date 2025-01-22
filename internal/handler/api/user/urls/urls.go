package urls

import (
	"compress/gzip"
	"encoding/json"
	"github.com/hddskull/urlShorty/internal/model"
	"github.com/hddskull/urlShorty/internal/storage"
	"github.com/hddskull/urlShorty/internal/utils"
	"github.com/hddskull/urlShorty/tools/custom"
	"io"
	"net/http"
)

func GetHandler(w http.ResponseWriter, r *http.Request) {
	urls, err := storage.Current.GetUserURLs(r.Context())
	if err != nil {
		custom.JSONError(w, err, http.StatusBadRequest)
	}
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(urls)
}

func DeleteHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	reader := r.Body

	if r.Header.Get(`Content-Encoding`) == `gzip` {
		gz, err := gzip.NewReader(r.Body)
		if err != nil {
			utils.SugaredLogger.Debugln("PostHandler gzip decompression error:", err)
			custom.JSONError(w, err, http.StatusInternalServerError)
		}
		reader = gz
	}
	defer reader.Close()

	body, err := io.ReadAll(r.Body)
	if err != nil {
		utils.SugaredLogger.Debugln("DeleteHandler reading body error:", err)
		custom.JSONError(w, err, http.StatusBadRequest)
	}

	var urls []string
	if err = json.Unmarshal(body, &urls); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = verifyDeleteURLs(urls...)
	if err != nil {
		custom.JSONError(w, err, http.StatusBadRequest)
		return
	}

	//call delete method
	sessionID, ok := model.SessionIDFromContext(r.Context())
	utils.SugaredLogger.Debugln("sessionID:", sessionID, "| ok:", ok)

	if !ok {
		utils.SugaredLogger.Debugln("err on ErrNoSessionID")
		w.WriteHeader(http.StatusForbidden)
		return //custom.ErrNoSessionID
	}
	go storage.Current.BatchMarkDeleted(sessionID, urls...)
	w.WriteHeader(http.StatusAccepted)
}

func verifyDeleteURLs(urls ...string) error {
	if len(urls) == 0 {
		return custom.ErrEmptyURL
	}
	return nil
}
