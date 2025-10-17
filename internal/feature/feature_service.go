package feature

import (
	"context"
	"errors"
	"fmt"
)

var (
	ErrInvalidFeatureID = errors.New("invalid feature id")
	ErrEventsRepoUnset  = errors.New("event repository not configured")
)

type FeatureRepository interface {
	GetByID(ctx context.Context, id int32) (*Feature, error)
	List(ctx context.Context) ([]*Feature, error)
	Create(ctx context.Context, feature *Feature) error
	Update(ctx context.Context, feature *Feature) error
	Delete(ctx context.Context, id int32) error
}

type EventRepository interface {
	Create(ctx context.Context, event *Event) error
	ListByFeatureID(ctx context.Context, featureID int32) ([]*Event, error)
}

type Service struct {
	featureRepo FeatureRepository
	eventRepo   EventRepository
}

func NewService(featureRepo FeatureRepository, eventRepo EventRepository) *Service {
	return &Service{
		featureRepo: featureRepo,
		eventRepo:   eventRepo,
	}
}

func (s *Service) GetFeature(ctx context.Context, id int32) (*Feature, error) {
	if id <= 0 {
		return nil, fmt.Errorf("get feature: %w %d", ErrInvalidFeatureID, id)
	}

	feature, err := s.featureRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get feature: %w", err)
	}

	return feature, nil
}

func (s *Service) ListFeatures(ctx context.Context) ([]*Feature, error) {
	features, err := s.featureRepo.List(ctx)
	if err != nil {
		return nil, fmt.Errorf("list features: %w", err)
	}

	return features, nil
}

func (s *Service) CreateFeature(ctx context.Context, feature *Feature) error {
	if err := feature.Validate(); err != nil {
		return fmt.Errorf("validate feature: %w", err)
	}

	if err := s.featureRepo.Create(ctx, feature); err != nil {
		return fmt.Errorf("create feature: %w", err)
	}

	return nil
}

func (s *Service) UpdateFeature(ctx context.Context, feature *Feature) error {
	if feature.ID <= 0 {
		return fmt.Errorf("update feature: %w %d", ErrInvalidFeatureID, feature.ID)
	}

	if err := feature.Validate(); err != nil {
		return fmt.Errorf("validate feature: %w", err)
	}

	if err := s.featureRepo.Update(ctx, feature); err != nil {
		return fmt.Errorf("update feature: %w", err)
	}

	return nil
}

func (s *Service) RecordEvent(ctx context.Context, event *Event) error {
	if err := event.Validate(); err != nil {
		return fmt.Errorf("validate event: %w", err)
	}

	if err := s.eventRepo.Create(ctx, event); err != nil {
		return fmt.Errorf("create event: %w", err)
	}

	return nil
}

func (s *Service) ListEventsByFeature(ctx context.Context, featureID int32) ([]*Event, error) {
	if featureID <= 0 {
		return nil, fmt.Errorf("list events: %w %d", ErrInvalidFeatureID, featureID)
	}

	events, err := s.eventRepo.ListByFeatureID(ctx, featureID)
	if err != nil {
		return nil, fmt.Errorf("list events: %w", err)
	}

	return events, nil
}

func (s *Service) DeleteFeature(ctx context.Context, id int32) error {
	if id <= 0 {
		return fmt.Errorf("delete feature: %w %d", ErrInvalidFeatureID, id)
	}

	if err := s.featureRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("delete feature: %w", err)
	}

	return nil
}
