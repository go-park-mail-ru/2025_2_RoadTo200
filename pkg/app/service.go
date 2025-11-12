package app

import (
	service "github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/service/implementations"
)

func (a *App) initServices() {
	// Инициализация сервисов с зависимостями
	a.services = &Services{
		Auth: service.NewAuthService(a.repositories.User, a.repositories.Session, a.logger),
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
