package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/chat-service/repository"
	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/chat-service/server"
	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/chat-service/service"
	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/config"
	gServer "github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/handler/grpc"
	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/repository/implementations/postgres"
	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/pkg/logger"
	pgxConn "github.com/go-park-mail-ru/2025_2_RoadTo200/backend/pkg/postgres"
	redisConn "github.com/go-park-mail-ru/2025_2_RoadTo200/backend/pkg/redis"
	pb "github.com/go-park-mail-ru/2025_2_RoadTo200/backend/proto/chat"
)

func main() {
	// Load config
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Initialize logger
	loggerInst := logger.New(&cfg.Logger)
	loggerInst.Info("🚀 Starting Chat Service...")

	// Connect to PostgreSQL
	ctx := context.Background()
	pgPool, err := pgxConn.NewConnect(ctx, &cfg.Postgres)
	if err != nil {
		loggerInst.Fatal(fmt.Errorf("failed to connect to PostgreSQL: %w", err))
	}
	defer pgPool.Close()
	loggerInst.Info("✅ PostgreSQL connected")

	// Connect to Redis (for Pub/Sub)
	redisClient, err := redisConn.NewPubSubClient(&cfg.Redis)
	if err != nil {
		loggerInst.Fatal(fmt.Errorf("failed to connect to Redis Pub/Sub: %w", err))
	}
	defer redisClient.Close()
	loggerInst.Info("✅ Redis Pub/Sub connected")

	// Initialize repositories
	messageRepo := repository.NewMessageRepository(pgPool, loggerInst)
	matchRepo := postgres.NewMatchRepository(pgPool)

	// Initialize chat service
	chatService := service.NewChatService(messageRepo, matchRepo, redisClient, loggerInst)

	// Create gRPC server
	chatServer := server.NewChatServer(chatService, loggerInst)
	grpcServer := gServer.NewGrpcServer(cfg.Port)
	pb.RegisterChatServiceServer(grpcServer, chatServer)

	listener, err := net.Listen("tcp", fmt.Sprintf(":%s", cfg.ChatService.Port))
	if err != nil {
		loggerInst.Fatal(fmt.Errorf("failed to listen: %w", err))
	}

	loggerInst.Info(fmt.Sprintf("Chat Service listening on port %s", cfg.ChatService.Port))

	// Graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	go func() {
		if err := grpcServer.Serve(listener); err != nil {
			loggerInst.Fatal(fmt.Errorf("failed to serve: %w", err))
		}
	}()

	<-stop
	loggerInst.Info("Shutting down Chat Service...")
	grpcServer.GracefulStop()
	loggerInst.Info("Chat Service stopped")
}
