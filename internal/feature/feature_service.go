package feature

import (
	"context"
	"fmt"
)

type Repository interface {
	GetByName(ctx context.Context, name string) (*Feature, error)
	List(ctx context.Context) ([]*Feature, error)
}

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) GetFeatureByName(ctx context.Context, name string) (*Feature, error) {
	feature, err := s.repo.GetByName(ctx, name)
	if err != nil {
		return nil, fmt.Errorf("get feature by name: %w", err)
	}

	return feature, nil
}

func (s *Service) ListFeatures(ctx context.Context) ([]*Feature, error) {
	features, err := s.repo.List(ctx)
	if err != nil {
		return nil, fmt.Errorf("list features: %w", err)
	}

	return features, nil
}
