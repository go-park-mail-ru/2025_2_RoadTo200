package app

import (
	"fmt"

	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/gateway/adapters"
	payment_service "github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/service/implementations"
	grpc "github.com/go-park-mail-ru/2025_2_RoadTo200/backend/pkg/grpc"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
)

func (a *App) initServices() error {
	// Метрики для gRPC клиента
	grpcMetrics := grpc_prometheus.NewClientMetrics()
	a.metrics.MustRegister(grpcMetrics)

	// Initialize gRPC client for Auth Service Client
	authClient, err := grpc.NewAuthClient(fmt.Sprintf("localhost:%s", a.config.AuthService.Port), grpcMetrics)
	if err != nil {
		a.logger.Fatal(fmt.Errorf("failed to create auth client: %w", err))
	}

	// Initialize gRPC client for	// Core Service Client
	coreClient, err := grpc.NewCoreClient(fmt.Sprintf("localhost:%s", a.config.CoreService.Port), grpcMetrics)
	if err != nil {
		a.logger.Fatal(fmt.Errorf("failed to create core client: %w", err))
	}

	// Инициализация сервисов с зависимостями
	// Chat Service Client
	chatClient, err := grpc.NewChatClient(fmt.Sprintf("localhost:%s", a.config.ChatService.Port), grpcMetrics)
	if err != nil {
		a.logger.Fatal("Failed to create chat client", err)
	}
	// Note: we should probably close chatConn on shutdown, but for now we rely on app exit.

	// NotificationService создается локально (не через gRPC)
	notificationService := payment_service.NewNotificationService(
		a.repositories.Notification,
		a.resources.RedisPubSub,
		a.logger,
	)

	a.services = &Services{
		Auth:    adapters.NewAuthServiceAdapter(authClient),
		Feed:    adapters.NewFeedServiceAdapter(coreClient),
		Profile: adapters.NewProfileServiceAdapter(coreClient),
		Swipe:   adapters.NewSwipeServiceAdapter(coreClient),
		Match:   adapters.NewMatchServiceAdapter(coreClient),
		Chat:    adapters.NewChatServiceAdapter(chatClient, a.logger),
		Strike:  adapters.NewStrikeServiceAdapter(coreClient),
		Payment: payment_service.NewPaymentService(
			a.repositories.User,
			a.repositories.Subscription,
			a.config,
			a.logger,
		),
		Notification: notificationService,
	}
	return nil
}
