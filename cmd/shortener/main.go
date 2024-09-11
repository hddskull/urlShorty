package main

import (
	"github.com/hddskull/urlShorty/config"
	"github.com/hddskull/urlShorty/internal/app"
)

func main() {
	config.Setup()
	app.Start()
}
