package app

import (
	"compress/gzip"
	"github.com/go-chi/chi/v5/middleware"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/hddskull/urlShorty/config"
	"github.com/hddskull/urlShorty/internal/handler/api/shorten"
	"github.com/hddskull/urlShorty/internal/handler/root"
	customMiddleware "github.com/hddskull/urlShorty/internal/middleware"
)

func Start() {

	r := chi.NewRouter()
	r.Use(middleware.AllowContentEncoding("gzip"))
	r.Use(middleware.Compress(gzip.DefaultCompression, "application/json", "text/html"))
	r.Use(customMiddleware.WithLogging)
	//r.Use(customMiddleware.WithGzip)

	r.Route("/", func(r chi.Router) {
		r.Post("/", root.RootPostHandler)
		r.Get("/{id}", root.RootGetHandler)
	})
	r.Route("/api", func(r chi.Router) {
		r.Post("/shorten", shorten.PostHandler)
	})

	err := http.ListenAndServe(config.Address.ServerAddress, r)

	if err != nil {
		panic(err)
	}

}
