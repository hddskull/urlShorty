package app

import (
	"net/http"

	"github.com/hddskull/urlShorty/internal/handler"
)

func Start() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", handler.RootHandler)

	err := http.ListenAndServe(":8080", mux)

	if err != nil {
		panic(err)
	}
}
