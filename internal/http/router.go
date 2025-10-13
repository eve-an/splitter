package http

import (
	"log/slog"
	"net/http"

	"github.com/eve-an/splitter/internal/http/handler"
)

func newRouter(logger *slog.Logger) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/api/v1/features", handler.ListAllFeatures)

	return chain(mux,
		recoveryMiddleware(logger), // runs first
		withRequestIDMiddleware,
		loggingMiddleware(logger), // runs last
	)
}
