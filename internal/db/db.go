package db

import (
	"context"
	"fmt"
	"time"

	"github.com/eve-an/splitter/internal/config"
	dbsqlc "github.com/eve-an/splitter/internal/db/sqlc"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Client struct {
	Pool    *pgxpool.Pool
	Queries *dbsqlc.Queries
}

func New(cfg config.DatabaseConfig) (*Client, error) {
	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("invalid database config: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	poolConfig, err := pgxpool.ParseConfig(cfg.URL)
	if err != nil {
		return nil, fmt.Errorf("parse database config: %w", err)
	}

	if cfg.MaxOpenConns > 0 {
		poolConfig.MaxConns = int32(cfg.MaxOpenConns)
	}

	if cfg.MaxIdleConns > 0 {
		poolConfig.MinConns = int32(cfg.MaxIdleConns)
	}

	if cfg.ConnMaxLifetime > 0 {
		poolConfig.MaxConnLifetime = cfg.ConnMaxLifetime
	} else {
		poolConfig.MaxConnLifetime = 30 * time.Minute
	}

	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		return nil, fmt.Errorf("connect to database: %w", err)
	}

	return &Client{
		Pool:    pool,
		Queries: dbsqlc.New(pool),
	}, nil
}

func (c *Client) Close() error {
	if c == nil || c.Pool == nil {
		return nil
	}

	c.Pool.Close()
	return nil
}
