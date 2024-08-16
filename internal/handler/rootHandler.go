package handler

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/hddskull/urlShorty/internal/storage"
)

func RootHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		handleGet(w, r)
	} else if r.Method == http.MethodPost {
		handlePost(w, r)
	} else {
		http.Error(w, "invalid method", http.StatusBadRequest)
		return
	}
}

func handleGet(w http.ResponseWriter, r *http.Request) {
	// id := r.URL.Query().Get("id")

	arr := strings.Split(r.URL.Path, "/")
	id := arr[len(arr)-1]
	fmt.Println(id)

	url, err := storage.TempStorage.Get(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Add("Content-Type", "text/plain")
	w.Header().Add("Location", url)
	w.WriteHeader(http.StatusTemporaryRedirect)
	// http.Redirect(w, r, url, http.StatusTemporaryRedirect)

	// w.Header().Add("Content-Type", "text/plain")
	// w.WriteHeader(http.StatusTemporaryRedirect)
	w.Write([]byte(url))
}

func handlePost(w http.ResponseWriter, r *http.Request) {
	bodyB, err := io.ReadAll(r.Body)
	defer r.Body.Close()

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	bodyS := string(bodyB)
	id, err := storage.TempStorage.Save(bodyS)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Add("Content-Type", "text/plain")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(id))
}
