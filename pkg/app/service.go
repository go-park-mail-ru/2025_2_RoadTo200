package app

import (
	service "github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/service/implementations"
)

func (a *App) initServices() {
	// Инициализация сервисов с зависимостями
	telegramService, err := service.NewTelegramService(
		a.config.Telegram.BotToken,
		a.logger,
		a.config.Telegram.ChatID,
		a.config.Telegram.Enabled,
	)
	if err != nil {
		a.logger.Fatalf("Failed to initialize Telegram service: %v", err)
	} else {
		a.logger.Info("Telegram service initialized successfully")
	}

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
		Swipe:  service.NewSwipeService(a.repositories.Swipe, a.repositories.Match),
		Match:  service.NewMatchService(a.repositories.Match, a.repositories.User, a.repositories.Swipe, a.repositories.Photo, a.logger),
		Report: service.NewSupportService(a.repositories.Report, a.repositories.Screen, a.repositories.Storage, a.logger, telegramService),
	}
}
