package feature

import (
	"context"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
)

type dbEvent struct {
	ID        int64     `db:"id"`
	FeatureID int64     `db:"feature_id"`
	UserID    string    `db:"user_id"`
	Variant   string    `db:"variant"`
	Type      string    `db:"event_type"`
	CreatedAt time.Time `db:"created_at"`
}

func mapDBEvent(e *dbEvent) *Event {
	return &Event{
		ID:        e.ID,
		FeatureID: e.FeatureID,
		UserID:    e.UserID,
		Variant:   e.Variant,
		Type:      e.Type,
		CreatedAt: e.CreatedAt,
	}
}

type postgresEventRepository struct {
	db *sqlx.DB
}

var _ EventRepository = (*postgresEventRepository)(nil)

func NewPostgresEventRepository(db *sqlx.DB) *postgresEventRepository {
	return &postgresEventRepository{db: db}
}

func (p *postgresEventRepository) Create(ctx context.Context, event *Event) error {
	query := `
		INSERT INTO events (feature_id, user_id, variant, event_type)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at`

	if err := p.db.QueryRowxContext(ctx, query, event.FeatureID, event.UserID, event.Variant, event.Type).Scan(&event.ID, &event.CreatedAt); err != nil {
		return fmt.Errorf("inserting event: %w", err)
	}

	return nil
}

func (p *postgresEventRepository) ListByFeatureID(ctx context.Context, featureID int64) ([]*Event, error) {
	query := `
		SELECT
			id,
			feature_id,
			user_id,
			variant,
			event_type,
			created_at
		FROM events
		WHERE feature_id = $1
		ORDER BY created_at DESC`

	dbEvents := []dbEvent{}
	if err := p.db.SelectContext(ctx, &dbEvents, query, featureID); err != nil {
		return nil, fmt.Errorf("selecting events by feature id: %w", err)
	}

	events := make([]*Event, 0, len(dbEvents))
	for i := range dbEvents {
		dbEvent := dbEvents[i]
		events = append(events, mapDBEvent(&dbEvent))
	}

	return events, nil
}
