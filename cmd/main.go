package main

import (
	"github.com/nekizz/init-project/config"
	"github.com/nekizz/init-project/internal/server"
)

func main() {
	//Load config
	cfg, err := config.Load()
	if err != nil {
		panic(err)
	}

	// Initialize Echo HTTP server
	e := server.New(cfg)

	//init service

	server.Start(e)
}
