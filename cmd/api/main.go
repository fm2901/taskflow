package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/fm2901/taskflow/internal/app"
	"github.com/fm2901/taskflow/internal/config"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))

	cfg := config.Load()

	application := app.New(cfg, logger)

	if err := application.Run(context.Background()); err != nil {
		logger.Error("application stopped with error", slog.Any("error", err))
		os.Exit(1)
	}
}