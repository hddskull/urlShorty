package handler

import (
	"fmt"
	"io"
	"net/http"
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
	id := r.URL.Query().Get("id")

	fmt.Println(id)

	url, err := getURL(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	http.Redirect(w, r, url, http.StatusTemporaryRedirect)

	// w.Header().Add("Content-Type", "text/plain")
	// w.WriteHeader(http.StatusTemporaryRedirect)
	// w.Write([]byte(url))
}

func getURL(id string) (string, error) {
	if id == "EwHXdJfB" {
		return "https://practicum.yandex.ru/", nil
	}

	return "", fmt.Errorf("no url found from id: %s", id)
}

func handlePost(w http.ResponseWriter, r *http.Request) {
	bodyB, err := io.ReadAll(r.Body)
	defer r.Body.Close()

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	bodyS := string(bodyB)
	processURL(bodyS)

	w.Header().Add("Content-Type", "text/plain")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("success"))
}

func processURL(b string) {

	fmt.Println(b)

}
