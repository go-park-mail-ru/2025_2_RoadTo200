package app

import (
	"fmt"

	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/gateway/adapters"
	service "github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/service/implementations"
	grpcPkg "github.com/go-park-mail-ru/2025_2_RoadTo200/backend/pkg/grpc"
)

func (a *App) initServices() {
	// Initialize gRPC client for Auth Service
	authClient, _, err := grpcPkg.NewAuthClient(fmt.Sprintf("localhost:%s", a.config.AuthService.Port))
	if err != nil {
		a.logger.Fatal(fmt.Errorf("failed to create auth client: %w", err))
	}

	// Инициализация сервисов с зависимостями
	a.services = &Services{
		Auth: adapters.NewAuthServiceAdapter(authClient),
		Feed: service.NewFeedService(a.repositories.User, a.repositories.Preference, a.repositories.Photo, a.logger),
		Profile: service.NewProfileService(
			a.repositories.User,
			a.repositories.Photo,
			a.repositories.Preference,
			a.repositories.Storage,
			a.logger,
		),
		Swipe: service.NewSwipeService(a.repositories.Swipe, a.repositories.Match),
		Match: service.NewMatchService(a.repositories.Match, a.repositories.User, a.repositories.Swipe, a.repositories.Photo, a.logger),
	}
}
