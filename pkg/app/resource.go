package app

import (
	"context"
	"fmt"

	minio_connect "github.com/go-park-mail-ru/2025_2_RoadTo200/backend/pkg/minio"
	postgres_connect "github.com/go-park-mail-ru/2025_2_RoadTo200/backend/pkg/postgres"
	redis_connect "github.com/go-park-mail-ru/2025_2_RoadTo200/backend/pkg/redis"
)

func (a *App) initResources() error {
	a.resources = &Resources{}

	if err := a.initPostgres(context.Background()); err != nil {
		return err
	}

	if err := a.initRedis(); err != nil {
		return err
	}

	if err := a.initMinIO(); err != nil {
		return err
	}

	a.logger.Info("✅ All connections established")
	return nil
}

func (a *App) initPostgres(ctx context.Context) error {
	a.logger.Info("🔌 Connecting to PostgreSQL...")

	psgPool, err := postgres_connect.NewConnect(ctx, &a.config.Postgres)
	if err != nil {
		return fmt.Errorf("failed to connect to PostgreSQL: %w", err)
	}

	a.resources.Postgres = psgPool
	a.logger.Info("✅ Postgres connected successfully")
	return nil
}

func (a *App) initRedis() error {
	a.logger.Info("🔌 Connecting to Redis...")

	redisPool, err := redis_connect.NewConnection(&a.config.Redis)
	if err != nil {
		return fmt.Errorf("failed to connect to Redis: %w", err)
	}

	a.resources.Redis = redisPool
	a.logger.Info("✅ Redis connected successfully")
	return nil
}

func (a *App) initMinIO() error {
	a.logger.Info("🔌 Connecting to MinIO...")

	minioPool, err := minio_connect.NewMinioPool(&a.config.MinIO)
	if err != nil {
		return fmt.Errorf("failed to connect to MinIO: %w", err)
	}

	a.resources.MinIO = minioPool
	a.logger.Info("✅ MinIO connected successfully")
	return nil
}
