package postgres

import (
	"context"
	"fmt"
	"wildberies/L0/backend/internal/config"

	"github.com/jackc/pgx/v5/pgxpool"
)

func NewConn(ctx context.Context, cfg *config.PostgresConfig) (*pgxpool.Pool, error) {
	connStr := fmt.Sprintf("%s://%s:%s@%s:%d/%s", cfg.Driver, cfg.User,
		cfg.Password, cfg.Host, cfg.Port, cfg.DBName)

	pool, err := pgxpool.New(ctx, connStr)
	if err != nil {
		return nil, fmt.Errorf("connecting to database: %w", err)
	}

	err = pool.Ping(ctx)
	if err != nil {
		return nil, fmt.Errorf("pinging database: %w", err)
	}

	return pool, nil
}
