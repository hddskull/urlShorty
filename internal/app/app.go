package app

import (
	"github.com/go-chi/chi/v5"
	"github.com/hddskull/urlShorty/config"
	"github.com/hddskull/urlShorty/internal/handler/api/shorten"
	"github.com/hddskull/urlShorty/internal/handler/api/shorten/batch"
	"github.com/hddskull/urlShorty/internal/handler/api/user/urls"
	"github.com/hddskull/urlShorty/internal/handler/root"
	customMiddleware "github.com/hddskull/urlShorty/internal/middleware"
	"net/http"
)

func Start() {

	r := chi.NewRouter()
	r.Use(customMiddleware.WithLogging)
	r.Use(customMiddleware.CompressResponseGzip)
	r.Use(customMiddleware.WithJWT)

	r.Route("/", func(r chi.Router) {
		r.Post("/", root.PostHandler)
		r.Get("/{id}", root.GetHandler)
		r.Get("/ping", root.PingHandler)
	})
	r.Route("/api", func(r chi.Router) {
		r.Post("/shorten", shorten.PostHandler)
		r.Post("/shorten/batch", batch.BatchHandler)
		r.Get("/user/urls", urls.GetHandler)
		r.Delete("/user/urls", urls.DeleteHandler)
	})

	err := http.ListenAndServe(config.Address.ServerAddress, r)

	if err != nil {
		panic(err)
	}

}
