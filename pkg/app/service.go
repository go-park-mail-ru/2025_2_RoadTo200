package app

import (
	"fmt"

	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/gateway/adapters"
	grpc "github.com/go-park-mail-ru/2025_2_RoadTo200/backend/pkg/grpc"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
)

func (a *App) initServices() error {
	// Метрики для gRPC клиента
	grpcMetrics := grpc_prometheus.NewClientMetrics()
	a.metrics.MustRegister(grpcMetrics)

	// Initialize gRPC client for Auth Service Client
	authClient, err := grpc.NewAuthClient(fmt.Sprintf("localhost:%s", a.config.AuthService.Port))
	if err != nil {
		a.logger.Fatal(fmt.Errorf("failed to create auth client: %w", err))
	}

	// Initialize gRPC client for	// Core Service Client
	coreClient, err := grpc.NewCoreClient(fmt.Sprintf("localhost:%s", a.config.CoreService.Port))
	if err != nil {
		a.logger.Fatal(fmt.Errorf("failed to create core client: %w", err))
	}

	// Инициализация сервисов с зависимостями
	// Chat Service Client
	chatClient, err := grpc.NewChatClient(fmt.Sprintf("localhost:%s", a.config.ChatService.Port))
	if err != nil {
		a.logger.Fatal("Failed to create chat client", err)
	}
	// Note: we should probably close chatConn on shutdown, but for now we rely on app exit.

	a.services = &Services{
		Auth:    adapters.NewAuthServiceAdapter(authClient),
		Feed:    adapters.NewFeedServiceAdapter(coreClient),
		Profile: adapters.NewProfileServiceAdapter(coreClient),
		Swipe:   adapters.NewSwipeServiceAdapter(coreClient),
		Match:   adapters.NewMatchServiceAdapter(coreClient),
		Chat:    adapters.NewChatServiceAdapter(chatClient, a.logger),
		Strike:  adapters.NewStrikeServiceAdapter(coreClient),
	}
	return nil
}
