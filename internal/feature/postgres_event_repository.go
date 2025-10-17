package feature

import (
	"context"
	"fmt"
	"time"

	dbsqlc "github.com/eve-an/splitter/internal/db/sqlc"

	"github.com/jackc/pgx/v5/pgtype"
)

type postgresEventRepository struct {
	queries *dbsqlc.Queries
}

var _ EventRepository = (*postgresEventRepository)(nil)

func NewPostgresEventRepository(queries *dbsqlc.Queries) *postgresEventRepository {
	return &postgresEventRepository{queries: queries}
}

func (p *postgresEventRepository) Create(ctx context.Context, event *Event) error {
	inserted, err := p.queries.InsertEvent(ctx, dbsqlc.InsertEventParams{
		FeatureID: pgInt4FromInt32(event.FeatureID),
		UserID:    textParam(event.UserID),
		Variant:   textParam(event.Variant),
		EventType: textParam(event.Type),
	})
	if err != nil {
		return fmt.Errorf("inserting event: %w", err)
	}

	applyEvent(inserted, event)

	return nil
}

func (p *postgresEventRepository) ListByFeatureID(ctx context.Context, featureID int32) ([]*Event, error) {
	rows, err := p.queries.ListEventsByFeatureID(ctx, pgInt4FromInt32(featureID))
	if err != nil {
		return nil, fmt.Errorf("selecting events by feature id: %w", err)
	}

	events := make([]*Event, 0, len(rows))
	for _, row := range rows {
		e := &Event{}
		applyEvent(row, e)
		events = append(events, e)
	}

	return events, nil
}

func applyEvent(dbEvent dbsqlc.Event, target *Event) {
	target.ID = dbEvent.ID
	if dbEvent.FeatureID.Valid {
		target.FeatureID = dbEvent.FeatureID.Int32
	} else {
		target.FeatureID = 0
	}
	target.UserID = textToString(dbEvent.UserID)
	target.Variant = textToString(dbEvent.Variant)
	target.Type = textToString(dbEvent.EventType)
	target.CreatedAt = timestamptzToTime(dbEvent.CreatedAt)
}

func timestamptzToTime(value pgtype.Timestamptz) time.Time {
	if !value.Valid {
		return time.Time{}
	}

	return value.Time
}
