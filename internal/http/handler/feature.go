package handler

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/eve-an/splitter/internal/feature"
)

type Feature struct {
	logger     *slog.Logger
	featureSvc *feature.Service
}

func NewFeatureHandler(
	logger *slog.Logger,
	featureSvc *feature.Service,
) *Feature {
	return &Feature{
		logger: logger,
	}
}

func (f *Feature) ListFeatures(w http.ResponseWriter, r *http.Request) {
	featureName := r.URL.Query().Get("feature_name")

	var features []*feature.Feature
	var err error
	if featureName != "" {
		feat, err := f.featureSvc.GetFeatureByName(r.Context(), featureName)
		if err != nil {
			f.respondError(w, err, "failed to get feature by name")
			return
		}
		features = []*feature.Feature{feat}
	} else {
		features, err = f.featureSvc.ListFeatures(r.Context())
		if err != nil {
			f.respondError(w, err, "failed to list features")
			return
		}
	}

	if err := json.NewEncoder(w).Encode(features); err != nil {
		f.respondError(w, err, "failed write features")
		return
	}
}

func (f *Feature) respondError(w http.ResponseWriter, err error, msg string) {
	f.logger.Error(msg, "error", err)

	switch {
	case errors.Is(err, context.DeadlineExceeded):
		writeError(w, http.StatusGatewayTimeout, "deadline exceeded")
	default:
		writeError(w, http.StatusInternalServerError, "unexpected error")
	}
}
