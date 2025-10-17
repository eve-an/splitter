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
	mux.HandleFunc("GET /api/v1/features/{featureID}", featureHandler.GetFeature)
	mux.HandleFunc("DELETE /api/v1/features/{featureID}", featureHandler.DeleteFeature)
	mux.HandleFunc("POST /api/v1/features", featureHandler.CreateFeature)
	mux.HandleFunc("PUT /api/v1/features/{featureID}", featureHandler.UpdateFeature)
	mux.HandleFunc("GET /api/v1/features/{featureID}/events", featureHandler.ListFeatureEvents)
	mux.HandleFunc("POST /api/v1/features/{featureID}/events", featureHandler.RecordFeatureEvent)

	return chain(mux,
		recoveryMiddleware(logger), // runs first
		stripTrailingSlash,
		withTraceIDMiddleware,
		loggingMiddleware(logger),
		corsMiddleware, // runs last
	)
}
