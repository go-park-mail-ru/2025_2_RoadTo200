package app

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/config"
	coreServer "github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/core-service/server"
	gServer "github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/handler/grpc"
	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/logger"
	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/repository/implementations/minio"
	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/repository/implementations/postgres"
	serviceImpl "github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/service/implementations"
	minioConn "github.com/go-park-mail-ru/2025_2_RoadTo200/backend/pkg/minio"
	pgxConn "github.com/go-park-mail-ru/2025_2_RoadTo200/backend/pkg/postgres"
	redisConn "github.com/go-park-mail-ru/2025_2_RoadTo200/backend/pkg/redis"
	pb "github.com/go-park-mail-ru/2025_2_RoadTo200/backend/proto/core"
)

func Run() {
	// Load configuration
	cfg, err := config.NewConfig()
	if err != nil {
		panic(fmt.Sprintf("failed to load config: %v", err))
	}

	// Initialize logger
	loggerInst := logger.New(&cfg.Logger)
	loggerInst.Info("🚀 Starting Core Service...")

	// Connect to PostgreSQL
	ctx := context.Background()
	pgPool, err := pgxConn.NewConnect(ctx, &cfg.Postgres)
	if err != nil {
		loggerInst.Fatal(fmt.Errorf("failed to connect to PostgreSQL: %w", err))
	}
	defer pgPool.Close()
	loggerInst.Info("✅ PostgreSQL connected")

	// Connect to Redis
	redisPool, err := redisConn.NewConnection(&cfg.Redis)
	if err != nil {
		loggerInst.Fatal(fmt.Errorf("failed to connect to Redis: %w", err))
	}
	defer redisPool.Close()
	loggerInst.Info("✅ Redis connected")

	// Connect to Redis Pub/Sub for notifications
	redisPubSub, err := redisConn.NewPubSubClient(&cfg.Redis)
	if err != nil {
		loggerInst.Fatal(fmt.Errorf("failed to connect to Redis Pub/Sub: %w", err))
	}
	defer redisPubSub.Close()
	loggerInst.Info("✅ Redis Pub/Sub connected")

	// Connect to MinIO
	minioClient, err := minioConn.NewMinioPool(&cfg.MinIO)
	if err != nil {
		loggerInst.Fatal(fmt.Errorf("failed to connect to MinIO: %w", err))
	}
	loggerInst.Info("✅ MinIO connected")

	// Initialize repositories
	userRepo := postgres.NewUserRepository(pgPool, loggerInst)
	photoRepo := postgres.NewUserPhotoRepository(pgPool, loggerInst)
	preferenceRepo := postgres.NewUserPreferenceRepository(pgPool, loggerInst)
	swipeRepo := postgres.NewSwipeRepository(pgPool)
	matchRepo := postgres.NewMatchRepository(pgPool)
	storageRepo := minio.NewStorageRepository(minioClient, &cfg.MinIO)
	strikeRepo := postgres.NewStrikeRepository(pgPool)
	notificationRepo := postgres.NewNotificationRepository(pgPool)

	// Initialize services
	profileService := serviceImpl.NewProfileService(userRepo, photoRepo, preferenceRepo, swipeRepo, matchRepo, storageRepo, loggerInst)
	feedService := serviceImpl.NewFeedService(userRepo, preferenceRepo, photoRepo, loggerInst)
	notificationService := serviceImpl.NewNotificationService(notificationRepo, redisPubSub, loggerInst)
	swipeService := serviceImpl.NewSwipeService(swipeRepo, matchRepo, userRepo, notificationService, loggerInst)
	matchService := serviceImpl.NewMatchService(matchRepo, userRepo, swipeRepo, photoRepo, loggerInst)
	strikeServie := serviceImpl.NewStrikeService(strikeRepo, userRepo, loggerInst)

	// Create gRPC server
	coreServiceServer := coreServer.NewCoreServer(profileService, feedService, swipeService, matchService, strikeServie, loggerInst)
	grpcServer := gServer.NewGrpcServer(cfg.Port, loggerInst)

	pb.RegisterCoreServiceServer(grpcServer, coreServiceServer)

	// Start listening
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", cfg.CoreService.Port))
	if err != nil {
		loggerInst.Fatal(fmt.Errorf("failed to listen: %w", err))
	}

	loggerInst.Info(fmt.Sprintf("Core Service listening on port %s", cfg.CoreService.Port))

	// Graceful shutdown
	go func() {
		if err := grpcServer.Serve(lis); err != nil {
			loggerInst.Fatal(fmt.Errorf("failed to serve: %w", err))
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	loggerInst.Info("Shutting down Core Service...")
	grpcServer.GracefulStop()
	loggerInst.Info("Core Service stopped")
}
