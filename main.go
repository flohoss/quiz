package main

import (
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/flohoss/quiz/config"
	handlers "github.com/flohoss/quiz/handler"
)

func setupRouter() *echo.Echo {
	e := echo.New()

	e.HideBanner = true
	e.HidePort = true

	e.Use(middleware.Recover())
	e.Use(middleware.CORS())
	e.Use(middleware.GzipWithConfig(middleware.GzipConfig{
		Skipper: func(c echo.Context) bool {
			return strings.Contains(c.Path(), "events")
		},
	}))

	return e
}

func main() {
	e := setupRouter()
	config.New()

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: config.GetLogLevel(),
	}))
	slog.SetDefault(logger)

	handlers.SetupRouter(e)

	slog.Info("Starting server", "url", fmt.Sprintf("http://%s", config.GetServer()))
	slog.Error(e.Start(config.GetServer()).Error())
}
