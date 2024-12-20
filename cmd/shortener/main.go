package main

import (
	"github.com/hddskull/urlShorty/config"
	"github.com/hddskull/urlShorty/internal/app"
	"github.com/hddskull/urlShorty/internal/storage"
)

func main() {
	config.Setup()
	storage.SetupStorage()
	app.Start()
}
