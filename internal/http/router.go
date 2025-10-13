package http

import (
	"log/slog"
	"net/http"

	"github.com/eve-an/splitter/internal/http/handler"
	"github.com/eve-an/splitter/internal/http/middleware"
)

func newRouter(logger *slog.Logger) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/api/v1/features", handler.ListAllFeatures)

	return middleware.Chain(mux,
		middleware.Recovery(logger), // runs first
		middleware.WithRequestID,
		middleware.Logging(logger), // runs last
	)
}
