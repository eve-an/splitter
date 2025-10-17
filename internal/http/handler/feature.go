package handler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"

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
		logger:     logger,
		featureSvc: featureSvc,
	}
}

func (f *Feature) ListFeatures(w http.ResponseWriter, r *http.Request) {
	features, err := f.featureSvc.ListFeatures(r.Context())
	if err != nil {
		f.respondError(w, err, "failed to list features")
		return
	}

	apiFeatures := make([]featureResponse, len(features))
	for i, feature := range features {
		apiFeatures[i] = mapFeatureResponse(feature)
	}

	Ok(w, apiFeatures)
}

func (f *Feature) GetFeature(w http.ResponseWriter, r *http.Request) {
	id, ok := parseFeatureID(w, r)
	if !ok {
		return
	}

	feat, err := f.featureSvc.GetFeature(r.Context(), id)
	if err != nil {
		f.respondError(w, err, fmt.Sprintf("failed to get feature by id %d", id))
		return
	}

	Ok(w, feat)
}

func (f *Feature) CreateFeature(w http.ResponseWriter, r *http.Request) {
	req, ok := decodeFeatureRequest(w, r)
	if !ok {
		return
	}

	domainFeature, err := buildFeatureFromRequest(req)
	if err != nil {
		f.respondError(w, err, "failed to build feature")
		return
	}

	if err := f.featureSvc.CreateFeature(r.Context(), domainFeature); err != nil {
		f.respondError(w, err, "failed to create feature")
		return
	}

	writeJSON(w, http.StatusCreated, domainFeature)
}

func (f *Feature) UpdateFeature(w http.ResponseWriter, r *http.Request) {
	id, ok := parseFeatureID(w, r)
	if !ok {
		return
	}

	req, ok := decodeFeatureRequest(w, r)
	if !ok {
		return
	}

	domainFeature, err := buildFeatureFromRequest(req)
	if err != nil {
		f.respondError(w, err, "failed to build feature")
		return
	}
	domainFeature.ID = id

	if err := f.featureSvc.UpdateFeature(r.Context(), domainFeature); err != nil {
		f.respondError(w, err, fmt.Sprintf("failed to update feature %d", id))
		return
	}

	Ok(w, domainFeature)
}

func (f *Feature) DeleteFeature(w http.ResponseWriter, r *http.Request) {
	id, ok := parseFeatureID(w, r)
	if !ok {
		return
	}

	if err := f.featureSvc.DeleteFeature(r.Context(), id); err != nil {
		f.respondError(w, err, fmt.Sprintf("failed to delete feature %d", id))
		return
	}

	Ok(w, nil)
}

func (f *Feature) ListFeatureEvents(w http.ResponseWriter, r *http.Request) {
	id, ok := parseFeatureID(w, r)
	if !ok {
		return
	}

	events, err := f.featureSvc.ListEventsByFeature(r.Context(), id)
	if err != nil {
		f.respondError(w, err, fmt.Sprintf("failed to list events for feature %d", id))
		return
	}

	Ok(w, events)
}

func (f *Feature) RecordFeatureEvent(w http.ResponseWriter, r *http.Request) {
	id, ok := parseFeatureID(w, r)
	if !ok {
		return
	}

	defer r.Body.Close() // nolint: errcheck

	var req eventRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		Error(w, http.StatusBadRequest, "invalid event payload")
		return
	}

	event, err := feature.NewEvent(id, req.UserID, req.Variant, req.Type)
	if err != nil {
		f.respondError(w, err, "failed to create event")
		return
	}

	if err := f.featureSvc.RecordEvent(r.Context(), event); err != nil {
		f.respondError(w, err, fmt.Sprintf("failed to record event for feature %d", id))
		return
	}

	writeJSON(w, http.StatusCreated, event)
}

func (f *Feature) respondError(w http.ResponseWriter, err error, msg string) {
	f.logger.Error(msg, "error", err)

	switch {
	case errors.Is(err, context.DeadlineExceeded):
		Error(w, http.StatusGatewayTimeout, "deadline exceeded")
	case
		errors.Is(err, feature.ErrMaximumWeightExceeded),
		errors.Is(err, feature.ErrVariantAlreadyExist),
		errors.Is(err, feature.ErrEventFeatureIDRequired),
		errors.Is(err, feature.ErrEventTypeRequired),
		errors.Is(err, feature.ErrFeatureAlreadyExists):
		Error(w, http.StatusBadRequest, err.Error())
	case errors.Is(err, feature.ErrFeatureNotFound):
		Error(w, http.StatusNotFound, "feature not found")
	case errors.Is(err, feature.ErrInvalidFeatureID):
		Error(w, http.StatusBadRequest, "invalid feature id")
	case errors.Is(err, feature.ErrEventsRepoUnset):
		Error(w, http.StatusInternalServerError, "event repository not configured")
	default:
		Error(w, http.StatusInternalServerError, "unexpected error")
	}
}

func parseFeatureID(w http.ResponseWriter, r *http.Request) (int32, bool) {
	featureIDValue := r.PathValue("featureID")
	if featureIDValue == "" {
		Error(w, http.StatusBadRequest, "missing feature id")
		return 0, false
	}

	id, err := strconv.ParseInt(featureIDValue, 10, 32)
	if err != nil || id <= 0 {
		Error(w, http.StatusBadRequest, "invalid feature id", featureIDValue)
		return 0, false
	}

	return int32(id), true
}

func decodeFeatureRequest(w http.ResponseWriter, r *http.Request) (*featureRequest, bool) {
	defer r.Body.Close() // nolint: errcheck

	var req featureRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		Error(w, http.StatusBadRequest, "invalid feature payload")
		return nil, false
	}

	return &req, true
}
