package shorten

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/hddskull/urlShorty/config"
	"github.com/hddskull/urlShorty/internal/storage"
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

	reqModel := requestPostModel{}
	if err := json.NewDecoder(r.Body).Decode(&reqModel); err != nil {
		utils.SugaredLogger.Debugln("PostHandler decoding error:", err)
		formattedError := custom.ErrorResponseModel{Message: err.Error()}
		custom.JSONError(w, formattedError, http.StatusBadRequest)
		return
	}

	id, err := storage.TempStorage.Save(reqModel.URL)
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
