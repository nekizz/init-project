package server

import (
	"context"
	"fmt"
	"github.com/nekizz/init-project/pkg/server/middleware/secure"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
)

// Config represents server specific config
type Config struct {
	Stage        string
	Port         int
	ReadTimeout  int
	WriteTimeout int
	Debug        bool
	AllowOrigins []string
}

var (
	// DefaultConfig for the API server
	DefaultConfig = Config{
		Stage:        "development",
		Port:         8080,
		ReadTimeout:  10,
		WriteTimeout: 5,
		Debug:        true,
		AllowOrigins: []string{"*"},
	}

	version   = "dev" // sha1 revision used to build the server
	buildTime = "now" // when the server was built
)

func (c *Config) fillDefaults() {
	if c.Stage == "" {
		c.Stage = DefaultConfig.Stage
	}
	if c.Port == 0 {
		c.Port = DefaultConfig.Port
	}
	if c.ReadTimeout == 0 {
		c.ReadTimeout = DefaultConfig.ReadTimeout
	}
	if c.WriteTimeout == 0 {
		c.WriteTimeout = DefaultConfig.WriteTimeout
	}
	if c.AllowOrigins == nil && len(c.AllowOrigins) == 0 {
		c.AllowOrigins = DefaultConfig.AllowOrigins
	}
}

// New instantates new Echo server
func New(cfg *Config) *echo.Echo {
	cfg.fillDefaults()
	e := echo.New()
	e.Use(middleware.Logger(), middleware.Recover(), secure.Headers(),
		secure.CORS(&secure.Config{AllowOrigins: cfg.AllowOrigins}))
	e.GET("/", healthCheck)
	e.Validator = NewValidator()
	e.Debug = cfg.Debug
	if e.Debug {
		e.Logger.SetLevel(log.DEBUG)
		e.Use(secure.BodyDump())
	} else {
		e.Logger.SetLevel(log.ERROR)
	}
	e.Server.Addr = fmt.Sprintf(":%d", cfg.Port)
	e.Server.ReadTimeout = time.Duration(cfg.ReadTimeout) * time.Minute
	e.Server.WriteTimeout = time.Duration(cfg.WriteTimeout) * time.Minute

	return e
}

func healthCheck(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]interface{}{
		"status":     "ok",
		"version":    version,
		"build_time": buildTime,
	})
}

// Start starts echo server with graceful shutdown process
func Start(e *echo.Echo) {
	idleConnsClosed := make(chan struct{})
	go func() {
		sigint := make(chan os.Signal, 1)

		// interrupt signal sent from terminal
		signal.Notify(sigint, os.Interrupt)
		// sigterm signal sent from kubernetes
		signal.Notify(sigint, syscall.SIGTERM)

		sigrev := <-sigint
		e.Logger.Infof("signal received: %s", sigrev.String())

		// We received an interrupt signal, shut down.
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := e.Shutdown(ctx); err != nil {
			// Error from closing listeners, or context timeout:
			e.Logger.Errorf("http server Shutdown: %v", err)
		}
		close(idleConnsClosed)
	}()

	// Hide verbose logs
	e.HideBanner = true
	// e.HidePort = true
	if err := e.StartServer(e.Server); err != nil {
		if err == http.ErrServerClosed {
			e.Logger.Info("http server stopped")
		} else {
			e.Logger.Errorf("http server StartServer: %v", err)
		}
	}

	<-idleConnsClosed
}
