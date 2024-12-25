package urls

import (
	"encoding/json"
	"github.com/hddskull/urlShorty/internal/storage"
	"github.com/hddskull/urlShorty/tools/custom"
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
