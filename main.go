package main

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/flohoss/quiz/config"
	"github.com/flohoss/quiz/handlers"
)

func main() {
	config.New()

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: config.GetLogLevel(),
	}))
	slog.SetDefault(logger)

	e := handlers.SetupRouter()

	slog.Info("Starting server", "url", fmt.Sprintf("http://%s", config.GetServer()))
	slog.Error(e.Start(config.GetServer()).Error())
}
