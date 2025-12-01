package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/auth-service/repository/implementations/postgres"
	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/auth-service/repository/implementations/redis"
	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/auth-service/server"
	serviceImpl "github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/auth-service/service/implementations"
	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/config"
	gServer "github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/handler/grpc"
	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/pkg/logger"
	pgxConn "github.com/go-park-mail-ru/2025_2_RoadTo200/backend/pkg/postgres"
	redisConn "github.com/go-park-mail-ru/2025_2_RoadTo200/backend/pkg/redis"
	pb "github.com/go-park-mail-ru/2025_2_RoadTo200/backend/proto/auth"
)

func main() {
	// Load config
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Initialize logger
	loggerInst := logger.New(&cfg.Logger)
	loggerInst.Info("🚀 Starting Auth Service...")

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

	// Initialize repositories
	userRepo := postgres.NewUserRepository(pgPool, loggerInst)
	sessionRepo := redis.NewSessionRepository(redisPool)

	// Initialize auth service
	authService := serviceImpl.NewAuthService(userRepo, sessionRepo, loggerInst)

	// Create gRPC server
	authServer := server.NewAuthServer(authService, loggerInst)
	grpcServer := gServer.NewGrpcServer(cfg.Port, loggerInst)
	pb.RegisterAuthServiceServer(grpcServer, authServer)

	listener, err := net.Listen("tcp", fmt.Sprintf(":%s", cfg.AuthService.Port))
	if err != nil {
		loggerInst.Fatal(fmt.Errorf("failed to listen: %w", err))
	}

	loggerInst.Info(fmt.Sprintf("Auth Service listening on port %s", cfg.AuthService.Port))

	// Graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	go func() {
		if err := grpcServer.Serve(listener); err != nil {
			loggerInst.Fatal(fmt.Errorf("failed to serve: %w", err))
		}
	}()

	<-stop
	loggerInst.Info("Shutting down Auth Service...")
	grpcServer.GracefulStop()
	loggerInst.Info("Auth Service stopped")
}
