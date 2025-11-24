package app

import (
	"fmt"

	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/gateway/adapters"
	grpc "github.com/go-park-mail-ru/2025_2_RoadTo200/backend/pkg/grpc"
)

func (a *App) initServices() {
	// Initialize gRPC client for Auth Service Client
	authClient, _, err := grpc.NewAuthClient(fmt.Sprintf("localhost:%s", a.config.AuthService.Port))
	if err != nil {
		a.logger.Fatal(fmt.Errorf("failed to create auth client: %w", err))
	}

	// Initialize gRPC client for	// Core Service Client
	coreClient, _, err := grpc.NewCoreClient(fmt.Sprintf("localhost:%s", a.config.CoreService.Port))
	if err != nil {
		a.logger.Fatal(fmt.Errorf("failed to create core client: %w", err))
	}

	// Инициализация сервисов с зависимостями
	// Chat Service Client
	chatClient, _, err := grpc.NewChatClient(fmt.Sprintf("localhost:%s", a.config.ChatService.Port))
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
}
