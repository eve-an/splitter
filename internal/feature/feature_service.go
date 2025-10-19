package feature

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/eve-an/splitter/internal/cache"
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
	featureRepo  FeatureRepository
	featureCache cache.Cache[*Feature]
	eventRepo    EventRepository
}

func NewService(
	featureRepo FeatureRepository,
	eventRepo EventRepository,
	featureCache cache.Cache[*Feature],
) *Service {
	return &Service{
		featureRepo:  featureRepo,
		featureCache: featureCache,
		eventRepo:    eventRepo,
	}
}

func (s *Service) GetFeature(ctx context.Context, id int32) (*Feature, error) {
	featureKey := strconv.Itoa(int(id))
	if feature, ok := s.featureCache.Get(featureKey); ok {
		return feature, nil
	}

	feature, err := s.featureRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get feature: %w", err)
	}

	s.featureCache.Set(featureKey, feature, 1*time.Minute)

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
	events, err := s.eventRepo.ListByFeatureID(ctx, featureID)
	if err != nil {
		return nil, fmt.Errorf("list events: %w", err)
	}

	return events, nil
}

func (s *Service) DeleteFeature(ctx context.Context, id int32) error {
	if err := s.featureRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("delete feature: %w", err)
	}

	return nil
}
