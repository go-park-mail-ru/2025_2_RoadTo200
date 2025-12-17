package app

import (
	"context"
	"fmt"

	minio_connect "github.com/go-park-mail-ru/2025_2_RoadTo200/backend/pkg/minio"
	postgres_connect "github.com/go-park-mail-ru/2025_2_RoadTo200/backend/pkg/postgres"
	redis_connect "github.com/go-park-mail-ru/2025_2_RoadTo200/backend/pkg/redis"
)

func (a *App) initResources() error {
	// PostgreSQL
	a.logger.Info("🔌 Connecting to PostgreSQL...")
	postgresPool, err := postgres_connect.NewConnect(context.Background(), &a.config.Postgres)
	if err != nil {
		return fmt.Errorf("failed to connect to PostgreSQL: %w", err)
	}
	a.logger.Info("✅ Postgres connected successfully")

	// Redis (for sessions)
	a.logger.Info("🔌 Connecting to Redis...")
	redisPool, err := redis_connect.NewConnection(&a.config.Redis)
	if err != nil {
		return fmt.Errorf("failed to connect to Redis: %w", err)
	}
	a.logger.Info("✅ Redis connected successfully")

	// MinIO
	a.logger.Info("🔌 Connecting to MinIO...")
	minioClient, err := minio_connect.NewMinioPool(&a.config.MinIO)
	if err != nil {
		return fmt.Errorf("failed to connect to MinIO: %w", err)
	}
	a.logger.Info("✅ MinIO connected successfully")

	// Redis Pub/Sub for chat
	a.logger.Info("🔌 Connecting to Redis Pub/Sub...")
	redisPubSub, err := redis_connect.NewPubSubClient(&a.config.Redis)
	if err != nil {
		return fmt.Errorf("failed to connect to Redis Pub/Sub: %w", err)
	}
	a.logger.Info("✅ Redis Pub/Sub connected successfully")

	a.resources = &Resources{
		Postgres:    postgresPool,
		Redis:       redisPool,
		RedisPubSub: redisPubSub,
		MinIO:       minioClient,
	}

	a.logger.Info("✅ All connections established")
	return nil
}
