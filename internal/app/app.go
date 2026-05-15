package app

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os/signal"
	"syscall"

	"github.com/fm2901/taskflow/internal/config"
	httptransport "github.com/fm2901/taskflow/internal/transport/http"
)

type App struct {
	cfg    config.Config
	logger *slog.Logger
}

func New(cfg config.Config, logger *slog.Logger) *App {
	return &App{
		cfg:    cfg,
		logger: logger,
	}
}

func (a *App) Run(ctx context.Context) error {
	ctx, stop := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	router := httptransport.NewRouter(a.logger)

	server := &http.Server{
		Addr:         a.cfg.HTTP.Addr,
		Handler:      router,
		ReadTimeout:  a.cfg.HTTP.ReadTimeout,
		WriteTimeout: a.cfg.HTTP.WriteTimeout,
		IdleTimeout:  a.cfg.HTTP.IdleTimeout,
	}

	errCh := make(chan error, 1)

	go func() {
		a.logger.Info("starting HTTP server", slog.String("addr", a.cfg.HTTP.Addr))

		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errCh <- err
			return
		}

		errCh <- nil
	}()

	select {
	case err := <-errCh:
		return err

	case <-ctx.Done():
		a.logger.Info("shutdown signal received")

		shutdownCtx, cancel := context.WithTimeout(context.Background(), a.cfg.HTTP.ShutdownTimeout)
		defer cancel()

		if err := server.Shutdown(shutdownCtx); err != nil {
			return err
		}

		a.logger.Info("HTTP server stopped gracefully")

		return nil
	}
}