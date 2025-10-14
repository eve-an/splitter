package http

import (
	"log/slog"
	"net/http"

	"github.com/eve-an/splitter/internal/http/handler"
)

func NewRouter(
	logger *slog.Logger,
	featureHandler *handler.Feature,
) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /api/v1/features", featureHandler.ListFeatures)

	return chain(mux,
		recoveryMiddleware(logger), // runs first
		stripTrailingSlash,
		withRequestIDMiddleware,
		loggingMiddleware(logger), // runs last
	)
}
