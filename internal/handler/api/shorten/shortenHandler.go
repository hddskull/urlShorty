package shorten

import (
	"compress/gzip"
	"encoding/json"
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
			formattedError := custom.ErrorResponseModel{Message: err.Error()}
			custom.JSONError(w, formattedError, http.StatusInternalServerError)
		}
		reader = gz
	}
	defer reader.Close()

	reqModel := requestPostModel{}
	if err := json.NewDecoder(reader).Decode(&reqModel); err != nil {
		utils.SugaredLogger.Debugln("PostHandler decoding error:", err)
		formattedError := custom.ErrorResponseModel{Message: err.Error()}
		custom.JSONError(w, formattedError, http.StatusBadRequest)
		return
	}

	id, err := storage.Current.Save(reqModel.URL)
	if err != nil {
		utils.SugaredLogger.Debugln("PostHandler saving to storage error:", err)
		formattedError := custom.ErrorResponseModel{Message: err.Error()}
		custom.JSONError(w, formattedError, http.StatusBadRequest)
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
		formattedError := custom.ErrorResponseModel{Message: err.Error()}
		custom.JSONError(w, formattedError, http.StatusBadRequest)
		return
	}
}
