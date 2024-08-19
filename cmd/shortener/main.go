package main

import (
	"github.com/hddskull/urlShorty/internal/app"
)

// func start(mux *http.ServeMux) {
// 	err := http.ListenAndServe(":8080", mux)

// 	if err != nil {
// 		panic(err)
// 	}
// }

func main() {
	app.Start()
	// mux := http.NewServeMux()
	// mux.HandleFunc("/", handler.RootHandler)

	// start(mux)
}
