package feature

import (
	"context"
	"errors"
	"fmt"
	"maps"
	"slices"
	"time"

	"github.com/jmoiron/sqlx"
)

var ErrFeatureNotFound = errors.New("feature not found")

type dbFeature struct {
	FeatureID        int64     `db:"feature_id"`
	FeatureName      string    `db:"feature_name"`
	FeatureDesc      string    `db:"feature_description"`
	FeatureActive    bool      `db:"feature_active"`
	FeatureCreatedAt time.Time `db:"feature_created_at"`
	VariantID        int64     `db:"variant_id"`
	VariantName      string    `db:"variant_name"`
	VariantWeight    int64     `db:"variant_weight"`
}

func mapFeature(dbFeature *dbFeature) (*Feature, error) {
	feature, err := NewFeature(dbFeature.FeatureName, dbFeature.FeatureDesc, dbFeature.FeatureActive, &Variants{})
	if err != nil {
		return nil, err
	}

	feature.ID = dbFeature.FeatureID

	return feature, nil
}

func mapVariant(dbVariant *dbFeature) (Variant, error) {
	variant, err := NewVariant(dbVariant.VariantName, uint8(dbVariant.VariantWeight))
	if err != nil {
		return Variant{}, err
	}

	variant.ID = dbVariant.VariantID

	return variant, nil
}

type postgresFeatureRepository struct {
	db *sqlx.DB
}

var _ FeatureRepository = (*postgresFeatureRepository)(nil)

func NewPostgresFeatureRepository(db *sqlx.DB) *postgresFeatureRepository {
	return &postgresFeatureRepository{db: db}
}

func (p *postgresFeatureRepository) List(ctx context.Context) ([]*Feature, error) {
	query := `
		SELECT
			f.id          AS feature_id,
			f.name        AS feature_name,
			f.description AS feature_description,
			f.active      AS feature_active,
			f.created_at  AS feature_created_at,
			v.id          AS variant_id,
			v.name        AS variant_name,
			v.weight      AS variant_weight
		FROM features f
		LEFT JOIN variants v ON f.id = v.feature_id;`

	rows := []dbFeature{}
	if err := p.db.SelectContext(ctx, &rows, query); err != nil {
		return nil, fmt.Errorf("selecting features: %w", err)
	}

	var err error
	featureMap := make(map[int64]*Feature, len(rows))
	for _, r := range rows {
		f, ok := featureMap[r.FeatureID]
		if !ok {
			f, err = mapFeature(&r)
			if err != nil {
				return nil, fmt.Errorf("mapping feature: %w", err)
			}
			featureMap[r.FeatureID] = f
		}

		if r.VariantID != 0 {
			variant, err := mapVariant(&r)
			if err != nil {
				return nil, fmt.Errorf("mapping variant: %w", err)
			}

			if err := f.AddVariant(&variant); err != nil {
				return nil, fmt.Errorf("adding variant to feature: %w", err)
			}
		}
	}

	return slices.Collect(maps.Values(featureMap)), nil
}

// GetByID implements FeatureRepository.
func (p *postgresFeatureRepository) GetByID(ctx context.Context, id int64) (*Feature, error) {
	query := `
		SELECT
			f.id          AS feature_id,
			f.name        AS feature_name,
			f.description AS feature_description,
			f.active      AS feature_active,
			f.created_at  AS feature_created_at,
			v.id          AS variant_id,
			v.name        AS variant_name,
			v.weight      AS variant_weight
		FROM features f
		LEFT JOIN variants v ON f.id = v.feature_id
		WHERE f.id = $1;`

	rows := []dbFeature{}
	if err := p.db.SelectContext(ctx, &rows, query, id); err != nil {
		return nil, fmt.Errorf("selecting feature by id: %w", err)
	}

	if len(rows) == 0 {
		return nil, ErrFeatureNotFound
	}

	var err error
	featureMap := make(map[int64]*Feature)
	for _, r := range rows {
		f, ok := featureMap[r.FeatureID]
		if !ok {
			f, err = mapFeature(&r)
			if err != nil {
				return nil, fmt.Errorf("mapping feature: %w", err)
			}
			featureMap[r.FeatureID] = f
		}

		if r.VariantID != 0 {
			variant, err := mapVariant(&r)
			if err != nil {
				return nil, fmt.Errorf("mapping variant: %w", err)
			}

			if err := f.AddVariant(&variant); err != nil {
				return nil, fmt.Errorf("adding variant to feature: %w", err)
			}
		}
	}

	for _, feature := range featureMap {
		return feature, nil
	}

	return nil, ErrFeatureNotFound
}

// Create implements FeatureRepository.
func (p *postgresFeatureRepository) Create(ctx context.Context, feature *Feature) error {
	tx, err := p.db.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}
	defer func() { _ = tx.Rollback() }()

	var featureID int64
	query := `
		INSERT INTO features (name, description, active)
		VALUES ($1, $2, $3)
		RETURNING id`

	if err := tx.GetContext(ctx, &featureID, query, feature.Name, feature.Descritption, feature.Active); err != nil {
		return fmt.Errorf("inserting feature: %w", err)
	}
	feature.ID = featureID

	for i := range feature.Variants {
		variant := &feature.Variants[i]
		variantQuery := `
			INSERT INTO variants (feature_id, name, weight)
			VALUES ($1, $2, $3)
			RETURNING id`

		if err := tx.GetContext(ctx, &variant.ID, variantQuery, featureID, variant.Name, variant.Weight); err != nil {
			return fmt.Errorf("inserting variant %s: %w", variant.Name, err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit transaction: %w", err)
	}

	return nil
}

// Update implements FeatureRepository.
func (p *postgresFeatureRepository) Update(ctx context.Context, feature *Feature) error {
	tx, err := p.db.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}
	defer func() { _ = tx.Rollback() }()

	updateFeatureQuery := `
		UPDATE features
		SET name = $1, description = $2, active = $3
		WHERE id = $4`

	if _, err := tx.ExecContext(ctx, updateFeatureQuery, feature.Name, feature.Descritption, feature.Active, feature.ID); err != nil {
		return fmt.Errorf("updating feature: %w", err)
	}

	deleteVariantsQuery := `DELETE FROM variants WHERE feature_id = $1`
	if _, err := tx.ExecContext(ctx, deleteVariantsQuery, feature.ID); err != nil {
		return fmt.Errorf("deleting existing variants: %w", err)
	}

	for i := range feature.Variants {
		variant := &feature.Variants[i]
		variantQuery := `
			INSERT INTO variants (feature_id, name, weight)
			VALUES ($1, $2, $3)
			RETURNING id`

		if err := tx.GetContext(ctx, &variant.ID, variantQuery, feature.ID, variant.Name, variant.Weight); err != nil {
			return fmt.Errorf("inserting variant %s: %w", variant.Name, err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit transaction: %w", err)
	}

	return nil
}
