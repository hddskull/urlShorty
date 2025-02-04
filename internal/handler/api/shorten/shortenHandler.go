package shorten

import (
	"compress/gzip"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/hddskull/urlShorty/internal/storage"
	"net/http"

	"github.com/hddskull/urlShorty/config"
	"github.com/hddskull/urlShorty/internal/utils"
	"github.com/hddskull/urlShorty/tools/custom"
)

type (
	requestPostModel struct {
		URL string `json:"url"`
	}
	responsePostModel struct {
		Result string `json:"result"`
	}
)

func PostHandler(w http.ResponseWriter, r *http.Request) {
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

	reqModel := requestPostModel{}
	if err := json.NewDecoder(reader).Decode(&reqModel); err != nil {
		utils.SugaredLogger.Debugln("PostHandler decoding error:", err)
		custom.JSONError(w, err, http.StatusBadRequest)
		return
	}

	if reqModel.URL == "" {
		utils.SugaredLogger.Debugln("Save() empty arg:", custom.ErrEmptyURL)
		custom.JSONError(w, custom.ErrEmptyURL, http.StatusBadRequest)
	}

	id, err := storage.Current.Save(r.Context(), reqModel.URL)
	if err != nil {
		utils.SugaredLogger.Debugln("PostHandler saving to storage error:", err)
		var uvError *custom.UniqueViolationError
		if errors.As(err, &uvError) {
			resModel := responsePostModel{
				Result: fmt.Sprint("http://", config.Address.BaseURL, "/", uvError.ShortURL),
			}
			w.WriteHeader(http.StatusConflict)
			json.NewEncoder(w).Encode(resModel)
		} else {
			custom.JSONError(w, err, http.StatusBadRequest)
		}

		return
	}

	fullID := fmt.Sprint("http://", config.Address.BaseURL, "/", id)

	resModel := responsePostModel{
		Result: fullID,
	}

	w.WriteHeader(http.StatusCreated)

	err = json.NewEncoder(w).Encode(resModel)
	if err != nil {
		utils.SugaredLogger.Debugln("PostHandler encoding error:", err)
		custom.JSONError(w, err, http.StatusInternalServerError)
		return
	}
}
