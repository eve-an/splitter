package db

import (
	"context"
	"fmt"
	"time"

	"github.com/eve-an/splitter/internal/config"

	_ "github.com/jackc/pgx/stdlib"
	"github.com/jmoiron/sqlx"
)

// New opens a PostgreSQL connection using the provided configuration.
func New(cfg config.DatabaseConfig) (*sqlx.DB, error) {
	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("invalid database config: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	db, err := sqlx.ConnectContext(ctx, "pgx", cfg.URL)
	if err != nil {
		return nil, fmt.Errorf("connect to database: %w", err)
	}

	if cfg.MaxOpenConns > 0 {
		db.SetMaxOpenConns(cfg.MaxOpenConns)
	}

	if cfg.MaxIdleConns > 0 {
		db.SetMaxIdleConns(cfg.MaxIdleConns)
	}

	if cfg.ConnMaxLifetime > 0 {
		db.SetConnMaxLifetime(cfg.ConnMaxLifetime)
	} else {
		db.SetConnMaxLifetime(30 * time.Minute)
	}

	return db, nil
}
