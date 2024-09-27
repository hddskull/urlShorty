package shorten

import (
	"encoding/json"
	"fmt"
	"github.com/hddskull/urlShorty/config"
	"github.com/hddskull/urlShorty/internal/storage"
	"github.com/hddskull/urlShorty/internal/utils"
	"net/http"
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

	reqModel := requestPostModel{}
	if err := json.NewDecoder(r.Body).Decode(&reqModel); err != nil {
		utils.SugaredLogger.Debugln("PostHandler decoding error:", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	id, err := storage.TempStorage.Save(reqModel.URL)
	if err != nil {
		utils.SugaredLogger.Debugln("PostHandler saving to storage error:", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	fullID := fmt.Sprint("http://", config.Address.BaseURL, "/", id)

	resModel := responsePostModel{
		Result: fullID,
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	err = json.NewEncoder(w).Encode(resModel)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}
}
