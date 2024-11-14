package app

import (
	"github.com/go-chi/chi/v5"
	"github.com/hddskull/urlShorty/config"
	"github.com/hddskull/urlShorty/internal/database"
	"github.com/hddskull/urlShorty/internal/handler/api/shorten"
	"github.com/hddskull/urlShorty/internal/handler/root"
	customMiddleware "github.com/hddskull/urlShorty/internal/middleware"
	_ "github.com/lib/pq"
	"net/http"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "password"
	dbname   = "urlshorty"
)

func Start() {

	creds := config.DBCredentials
	err := database.Current.ConnectDB(creds)
	if err != nil {
		panic(err)
	}
	defer database.Current.CloseDB()

	r := chi.NewRouter()
	r.Use(customMiddleware.WithLogging)
	r.Use(customMiddleware.CompressResponseGzip)

	r.Route("/", func(r chi.Router) {
		r.Post("/", root.PostHandler)
		r.Get("/{id}", root.GetHandler)
		r.Get("/ping", root.PingHandler)
	})
	r.Route("/api", func(r chi.Router) {
		r.Post("/shorten", shorten.PostHandler)
	})

	err = http.ListenAndServe(config.Address.ServerAddress, r)

	if err != nil {
		panic(err)
	}

}
