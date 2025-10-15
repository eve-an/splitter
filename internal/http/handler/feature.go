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

	Ok(w, features)
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
	case errors.Is(err, feature.ErrMaximumWeightExceeded),
		errors.Is(err, feature.ErrVariantAlreadyExist),
		errors.Is(err, feature.ErrEventFeatureIDRequired),
		errors.Is(err, feature.ErrEventTypeRequired):
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

func parseFeatureID(w http.ResponseWriter, r *http.Request) (int64, bool) {
	featureIDValue := r.PathValue("featureID")
	if featureIDValue == "" {
		Error(w, http.StatusBadRequest, "missing feature id")
		return 0, false
	}

	id, err := strconv.ParseInt(featureIDValue, 10, 64)
	if err != nil || id <= 0 {
		Error(w, http.StatusBadRequest, "invalid feature id", featureIDValue)
		return 0, false
	}

	return id, true
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

func buildFeatureFromRequest(req *featureRequest) (*feature.Feature, error) {
	variants := make([]feature.Variant, 0, len(req.Variants))
	for _, v := range req.Variants {
		variant, err := feature.NewVariant(v.Name, v.Weight)
		if err != nil {
			return nil, err
		}
		variants = append(variants, variant)
	}

	domainVariants, err := feature.NewVariants(variants...)
	if err != nil {
		return nil, err
	}

	return feature.NewFeature(req.Name, req.Description, req.Active, &domainVariants)
}
