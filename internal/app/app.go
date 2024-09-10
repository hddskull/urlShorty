package app

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/hddskull/urlShorty/config"
	"github.com/hddskull/urlShorty/internal/handler"
)

func Start() {
	config.ConfigureNetAddress()

	r := chi.NewRouter()

	r.Route("/", func(r chi.Router) {
		r.Post("/", handler.RootPostHandler)
		r.Get("/{id}", handler.RootGetHandler)
	})

	err := http.ListenAndServe(config.LaunchAdr.String(), r)

	if err != nil {
		panic(err)
	}
}
