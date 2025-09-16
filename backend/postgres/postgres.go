package postgres

import (
	"context"
	"fmt"
	"wildberies/L0/backend/internal/config"

	"github.com/jackc/pgx/v5/pgxpool"
)

func NewConn(ctx context.Context, cfg *config.Config) (*pgxpool.Pool, error) {
	connStr := fmt.Sprintf("%s://%s:%s@%s:%d/%s", cfg.PostgresConfig.Driver, cfg.PostgresConfig.User,
		cfg.PostgresConfig.Password, cfg.PostgresConfig.Host, cfg.PostgresConfig.Port, cfg.PostgresConfig.DBName)

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
