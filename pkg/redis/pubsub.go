package redis

import (
	"context"
	"fmt"

	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/config"
	"github.com/redis/go-redis/v9"
)

// NewPubSubClient creates a new Redis client for Pub/Sub
func NewPubSubClient(cfg *config.RedisConfig) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", cfg.Host, cfg.Port),
		Password: cfg.Password,
		DB:       0, // Use default DB for Pub/Sub
	})

	// Test connection
	ctx := context.Background()
	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	return client, nil
}
