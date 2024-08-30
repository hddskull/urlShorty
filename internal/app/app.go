package app

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/hddskull/urlShorty/internal/handler"
)

func Start() {
	r := chi.NewRouter()

	r.Route("/", func(r chi.Router) {
		r.Post("/", handler.RootPostHandler)
		r.Get("/{id}", handler.RootGetHandler)
	})

	err := http.ListenAndServe(":8080", r)

	if err != nil {
		panic(err)
	}
}
