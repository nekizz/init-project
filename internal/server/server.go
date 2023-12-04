package server

import (
	"github.com/labstack/echo/v4"
	"github.com/nekizz/init-project/config"
	"github.com/nekizz/init-project/pkg/server"
)

// New initializes server
func New(cfg *config.Configuration) *echo.Echo {
	// Initialize HTTP server
	e := server.New(&server.Config{
		Stage:        cfg.Stage,
		Port:         cfg.Port,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
		AllowOrigins: cfg.AllowOrigins,
		Debug:        cfg.Debug,
	})

	return e
}

// Start starts the server
func Start(e *echo.Echo) {
	// Start the HTTP server
	server.Start(e)
}
