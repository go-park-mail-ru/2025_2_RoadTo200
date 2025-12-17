package redis

import (
	"context"
	"fmt"
	"os"
	"strconv"

	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/config"
	"github.com/redis/go-redis/v9"
)

// NewPubSubClient creates a new Redis client for Pub/Sub
func NewPubSubClient(cfg *config.RedisConfig) (*redis.Client, error) {
	host := os.Getenv(cfg.Host)         // ❌ Это ищет env переменную с именем "localhost"
	sport := os.Getenv(cfg.Port)        // ❌ Это ищет env переменную с именем "5435"
	password := os.Getenv(cfg.Password) // ❌ Это ищет env переменную с именем "password"
	base := os.Getenv(cfg.Base)         // ❌ Это ищет env переменную с именем "dating_app"

	// Преобразуем порт и базу в числа
	port, err := strconv.Atoi(sport)
	if err != nil {
		return nil, fmt.Errorf("failed to convert port to int: %w", err)
	}

	baseNum, err := strconv.Atoi(base)
	if err != nil {
		return nil, fmt.Errorf("failed to convert base to int: %w", err)
	}

	address := fmt.Sprintf("%s:%d", host, port)

	client := redis.NewClient(&redis.Options{
		Addr:     address,
		Password: password,
		DB:       baseNum, // Use default DB for Pub/Sub
	})

	// Test connection
	ctx := context.Background()
	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	return client, nil
}
