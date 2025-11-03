package postgres

import (
	"context"
	"fmt"
	"os"
	"strconv"

	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/config"
	"github.com/jackc/pgx/v5/pgxpool"
)

func NewConnect(ctx context.Context, cfg *config.PostgresConfig) (*pgxpool.Pool, error) {
	host := os.Getenv(cfg.Host)
	sport := os.Getenv(cfg.Port)
	user := os.Getenv(cfg.User)
	password := os.Getenv(cfg.Password)
	base := os.Getenv(cfg.Base)
	port, err := strconv.Atoi(sport)
	if err != nil {
		return nil, err
	}

	poolConfig, err := pgxpool.ParseConfig(
		fmt.Sprintf("postgres://%s:%s@%s:%d/%s", user, password, host, port, base))
	if err != nil {
		return nil, fmt.Errorf("failed to parse pool config: %w", err)
	}

	poolConfig.MaxConns = cfg.MaxConns
	poolConfig.MinConns = cfg.MinConns
	poolConfig.MaxConnLifetime = cfg.MaxLife
	poolConfig.MaxConnIdleTime = cfg.MaxIdle
	poolConfig.HealthCheckPeriod = cfg.HealthCheckInterval

	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create connection pool: %w", err)
	}

	// Test connection
	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return pool, nil
}
