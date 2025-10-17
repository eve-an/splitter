package feature

import (
	"context"
	"errors"
	"fmt"
	"maps"
	"slices"

	dbsqlc "github.com/eve-an/splitter/internal/db/sqlc"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	ErrFeatureNotFound      = errors.New("feature not found")
	ErrFeatureAlreadyExists = errors.New("feature already exists")
)

type postgresFeatureRepository struct {
	pool    *pgxpool.Pool
	queries *dbsqlc.Queries
}

var _ FeatureRepository = (*postgresFeatureRepository)(nil)

func NewPostgresFeatureRepository(pool *pgxpool.Pool, queries *dbsqlc.Queries) *postgresFeatureRepository {
	return &postgresFeatureRepository{
		pool:    pool,
		queries: queries,
	}
}

func (p *postgresFeatureRepository) List(ctx context.Context) ([]*Feature, error) {
	rows, err := p.queries.ListFeatures(ctx)
	if err != nil {
		return nil, fmt.Errorf("selecting features: %w", err)
	}

	featureMap := make(map[int32]*Feature, len(rows))
	for _, r := range rows {
		f, ok := featureMap[r.FeatureID]
		if !ok {
			f, err = mapFeatureRow(r.FeatureID, r.FeatureName, r.FeatureDescription, r.FeatureActive)
			if err != nil {
				return nil, fmt.Errorf("mapping feature: %w", err)
			}
			featureMap[r.FeatureID] = f
		}

		if !r.VariantID.Valid || !r.VariantName.Valid {
			continue
		}

		variant, err := mapVariantRow(r.VariantID, r.VariantName, r.VariantWeight)
		if err != nil {
			return nil, fmt.Errorf("mapping variant: %w", err)
		}

		if err := f.AddVariant(&variant); err != nil {
			return nil, fmt.Errorf("adding variant to feature: %w", err)
		}
	}

	return slices.Collect(maps.Values(featureMap)), nil
}

// GetByID implements FeatureRepository.
func (p *postgresFeatureRepository) GetByID(ctx context.Context, id int32) (*Feature, error) {
	rows, err := p.queries.GetFeature(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("selecting feature by id: %w", err)
	}

	if len(rows) == 0 {
		return nil, ErrFeatureNotFound
	}

	featureMap := make(map[int32]*Feature)
	for _, r := range rows {
		f, ok := featureMap[r.FeatureID]
		if !ok {
			f, err = mapFeatureRow(r.FeatureID, r.FeatureName, r.FeatureDescription, r.FeatureActive)
			if err != nil {
				return nil, fmt.Errorf("mapping feature: %w", err)
			}
			featureMap[r.FeatureID] = f
		}

		if !r.VariantID.Valid || !r.VariantName.Valid {
			continue
		}

		variant, err := mapVariantRow(r.VariantID, r.VariantName, r.VariantWeight)
		if err != nil {
			return nil, fmt.Errorf("mapping variant: %w", err)
		}

		if err := f.AddVariant(&variant); err != nil {
			return nil, fmt.Errorf("adding variant to feature: %w", err)
		}
	}

	if len(featureMap) == 0 {
		return nil, ErrFeatureNotFound
	}

	for _, feature := range featureMap {
		return feature, nil
	}

	return nil, ErrFeatureNotFound
}

// Create implements FeatureRepository.
func (p *postgresFeatureRepository) Create(ctx context.Context, feature *Feature) error {
	tx, err := p.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}
	defer func() { _ = tx.Rollback(ctx) }()

	queries := p.queries.WithTx(tx)

	featureID, err := queries.InsertFeature(ctx, dbsqlc.InsertFeatureParams{
		Name:        feature.Name,
		Description: textParam(feature.Descritption),
		Active:      feature.Active,
	})
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return ErrFeatureAlreadyExists
		}

		return fmt.Errorf("inserting feature: %w", err)
	}
	feature.ID = featureID
	featureIDParam := pgInt4FromInt32(featureID)

	for i := range feature.Variants {
		variant := &feature.Variants[i]

		variantID, err := queries.InsertVariant(ctx, dbsqlc.InsertVariantParams{
			FeatureID: featureIDParam,
			Name:      variant.Name,
			Weight:    int32(variant.Weight),
		})
		if err != nil {
			return fmt.Errorf("inserting variant %s: %w", variant.Name, err)
		}

		variant.ID = variantID
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit transaction: %w", err)
	}

	return nil
}

// Update implements FeatureRepository.
func (p *postgresFeatureRepository) Update(ctx context.Context, feature *Feature) (err error) {
	tx, err := p.pool.BeginTx(ctx, pgx.TxOptions{})
	defer func() {
		if err != nil {
			_ = tx.Rollback(ctx)
		}
	}()

	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}

	queries := p.queries.WithTx(tx)

	if err := queries.UpdateFeature(ctx, dbsqlc.UpdateFeatureParams{
		Name:        feature.Name,
		Description: textParam(feature.Descritption),
		Active:      feature.Active,
		ID:          feature.ID,
	}); err != nil {
		return fmt.Errorf("updating feature: %w", err)
	}

	if err := queries.DeleteVariantsByFeature(ctx, pgInt4FromInt32(feature.ID)); err != nil {
		return fmt.Errorf("deleting existing variants: %w", err)
	}

	featureIDParam := pgInt4FromInt32(feature.ID)
	for i := range feature.Variants {
		variant := &feature.Variants[i]

		variantID, err := queries.InsertVariant(ctx, dbsqlc.InsertVariantParams{
			FeatureID: featureIDParam,
			Name:      variant.Name,
			Weight:    int32(variant.Weight),
		})
		if err != nil {
			return fmt.Errorf("inserting variant %s: %w", variant.Name, err)
		}

		variant.ID = variantID
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit transaction: %w", err)
	}

	return nil
}

func (r *postgresFeatureRepository) Delete(ctx context.Context, id int32) (err error) {
	tx, err := r.pool.BeginTx(ctx, pgx.TxOptions{})
	defer func() {
		if err != nil {
			_ = tx.Rollback(ctx)
		}
	}()

	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}

	queries := r.queries.WithTx(tx)

	if err := queries.DeleteVariantsByFeature(ctx, pgInt4FromInt32(id)); err != nil {
		return fmt.Errorf("deleting existing variants: %w", err)
	}

	if err := queries.DeleteFeature(ctx, id); err != nil {
		return fmt.Errorf("deleting feature: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit transaction: %w", err)
	}

	return nil
}

func mapFeatureRow(id int32, name string, description pgtype.Text, active bool) (*Feature, error) {
	feature, err := NewFeature(name, textToString(description), active, &Variants{})
	if err != nil {
		return nil, err
	}

	feature.ID = id

	return feature, nil
}

func mapVariantRow(id pgtype.Int4, name pgtype.Text, weight pgtype.Int4) (Variant, error) {
	if !id.Valid {
		return Variant{}, errors.New("variant id is null")
	}
	if !name.Valid {
		return Variant{}, errors.New("variant name is null")
	}
	if !weight.Valid {
		return Variant{}, errors.New("variant weight is null")
	}

	converted, err := uint8FromInt32(weight.Int32)
	if err != nil {
		return Variant{}, err
	}

	variant, err := NewVariant(name.String, converted)
	if err != nil {
		return Variant{}, err
	}

	variant.ID = id.Int32

	return variant, nil
}

func textToString(value pgtype.Text) string {
	if !value.Valid {
		return ""
	}

	return value.String
}

func textParam(value string) pgtype.Text {
	return pgtype.Text{
		String: value,
		Valid:  true,
	}
}

func uint8FromInt32(value int32) (uint8, error) {
	if value < 0 || value > 255 {
		return 0, fmt.Errorf("value %d cannot be represented as uint8", value)
	}

	return uint8(value), nil
}

func pgInt4FromInt32(value int32) pgtype.Int4 {
	return pgtype.Int4{
		Int32: value,
		Valid: true,
	}
}
