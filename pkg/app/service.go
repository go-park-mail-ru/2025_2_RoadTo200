package app

import (
	"fmt"

	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/gateway/adapters"
	grpcPkg "github.com/go-park-mail-ru/2025_2_RoadTo200/backend/pkg/grpc"
)

func (a *App) initServices() {
	// Initialize gRPC client for Auth Service
	authClient, _, err := grpcPkg.NewAuthClient(fmt.Sprintf("localhost:%s", a.config.AuthService.Port))
	if err != nil {
		a.logger.Fatal(fmt.Errorf("failed to create auth client: %w", err))
	}

	// Initialize gRPC client for Core Service
	coreClient, _, err := grpcPkg.NewCoreClient(fmt.Sprintf("localhost:%s", a.config.CoreService.Port))
	if err != nil {
		a.logger.Fatal(fmt.Errorf("failed to create core client: %w", err))
	}

	// Инициализация сервисов с зависимостями
	a.services = &Services{
		Auth:    adapters.NewAuthServiceAdapter(authClient),
		Profile: adapters.NewProfileServiceAdapter(coreClient),
		Feed:    adapters.NewFeedServiceAdapter(coreClient),
		Swipe:   adapters.NewSwipeServiceAdapter(coreClient),
		Match:   adapters.NewMatchServiceAdapter(coreClient),
		Strike:  adapters.NewStrikeServiceAdapter(coreClient),
	}
}
