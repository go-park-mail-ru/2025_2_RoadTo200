package postgres

import (
	"context"
	"fmt"
	"os"

	"strconv"

	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/config"
	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/logger"
	"github.com/jackc/pgx/v4/pgxpool"
)

func NewConnect(ctx context.Context, cfg *config.PostgresConfig) (*pgxpool.Pool, error) {
	logger.Printf("🔍 PostgreSQL Config: Host=%s, Port=%s, User=%s, Base=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.Base)
	logger.Printf("🔍 PostgreSQL Password length: %d", len(cfg.Password))
	logger.Printf("🔍 PostgreSQL Port type: %T, value: '%s'", cfg.Port, cfg.Port)

	if cfg.Port == "" {
		return nil, fmt.Errorf("PostgreSQL port is empty")
	}

	// УБЕРИ эти строки - они пытаются получить значения из env переменных
	host := os.Getenv(cfg.Host)         // ❌ Это ищет env переменную с именем "localhost"
	sport := os.Getenv(cfg.Port)        // ❌ Это ищет env переменную с именем "5435"
	user := os.Getenv(cfg.User)         // ❌ Это ищет env переменную с именем "postgres"
	password := os.Getenv(cfg.Password) // ❌ Это ищет env переменную с именем "password"
	base := os.Getenv(cfg.Base)         // ❌ Это ищет env переменную с именем "dating_app"

	// Преобразуем порт в число
	port, err := strconv.Atoi(sport)
	if err != nil {
		return nil, fmt.Errorf("failed to convert port to int: %w", err)
	}

	poolConfig, err := pgxpool.ParseConfig(
		fmt.Sprintf("postgres://%s:%s@%s:%d/%s", user, password, host, port, base))
	if err != nil {
		return nil, fmt.Errorf("failed to parse pool config: %w", err)
	}

	logger.Printf("🔍 Pool Config: MaxConns=%d, MinConns=%d", cfg.MaxConns, cfg.MinConns)

	poolConfig.MaxConns = cfg.MaxConns
	poolConfig.MinConns = cfg.MinConns
	poolConfig.MaxConnLifetime = cfg.MaxLife
	poolConfig.MaxConnIdleTime = cfg.MaxIdle
	poolConfig.HealthCheckPeriod = cfg.HealthCheckInterval

	pool, err := pgxpool.ConnectConfig(ctx, poolConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create connection pool: %w", err)
	}

	// Test connection
	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	logger.Println("✅ PostgreSQL connected successfully")
	return pool, nil
}
