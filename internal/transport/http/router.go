package http

import (
	"log/slog"
	"net/http"

	"github.com/fm2901/taskflow/internal/transport/http/handlers"
	"github.com/fm2901/taskflow/internal/transport/http/middleware"
)

func NewRouter(logger *slog.Logger) http.Handler {
	mux := http.NewServeMux()

	healthHandler := handlers.NewHealthHandler()

	mux.HandleFunc("GET /health", healthHandler.Health)
	mux.HandleFunc("GET /ready", healthHandler.Ready)

	var handler http.Handler = mux

	handler = middleware.Recovery(logger)(handler)
	handler = middleware.Logging(logger)(handler)
	handler = middleware.RequestID(handler)

	return handler
}