package app

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/hddskull/urlShorty/config"
	"github.com/hddskull/urlShorty/internal/handler/api/shorten"
	"github.com/hddskull/urlShorty/internal/handler/root"
	customMiddleware "github.com/hddskull/urlShorty/internal/middleware"
)

func Start() {

	r := chi.NewRouter()
	r.Use(customMiddleware.WithLogging)
	r.Use(customMiddleware.CompressResponseGzip)

	r.Route("/", func(r chi.Router) {
		r.Post("/", root.PostHandler)
		r.Get("/{id}", root.GetHandler)
	})
	r.Route("/api", func(r chi.Router) {
		r.Post("/shorten", shorten.PostHandler)
	})

	err := http.ListenAndServe(config.Address.ServerAddress, r)

	if err != nil {
		panic(err)
	}

}
