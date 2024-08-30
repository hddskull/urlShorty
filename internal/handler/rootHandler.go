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
		RootGetHandler(w, r)
	} else if r.Method == http.MethodPost {
		RootPostHandler(w, r)
	} else {
		http.Error(w, "invalid method", http.StatusBadRequest)
		return
	}
}

func RootGetHandler(w http.ResponseWriter, r *http.Request) {
	arr := strings.Split(r.URL.Path, "/")
	id := arr[len(arr)-1]

	url, err := storage.TempStorage.Get(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Println(url)

	w.Header().Add("Content-Type", "text/plain")
	w.Header().Add("Location", url)
	w.WriteHeader(http.StatusTemporaryRedirect)
}

func RootPostHandler(w http.ResponseWriter, r *http.Request) {
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

	fullID := fmt.Sprint("http://localhost:8080/", id)

	w.Header().Add("Content-Type", "text/plain")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(fullID))
}
