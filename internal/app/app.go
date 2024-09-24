package app

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/hddskull/urlShorty/config"
	"github.com/hddskull/urlShorty/internal/handler"
	"github.com/hddskull/urlShorty/internal/middleware"
)

func Start() {

	r := chi.NewRouter()
	r.Use(middleware.WithLogging)
	r.Route("/", func(r chi.Router) {
		r.Post("/", handler.RootPostHandler)
		r.Get("/{id}", handler.RootGetHandler)
	})

	err := http.ListenAndServe(config.Address.ServerAddress, r)

	if err != nil {
		panic(err)
	}

}
