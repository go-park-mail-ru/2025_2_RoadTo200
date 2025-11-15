package app

import (
	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/domain/constants"
	service "github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/service/implementations"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/google/uuid"
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
		// Создаем заглушку если телеграм не работает
		//telegramService = &service.TelegramServiceStub{}
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
		Swipe:    service.NewSwipeService(a.repositories.Swipe, a.repositories.Match),
		Match:    service.NewMatchService(a.repositories.Match, a.repositories.User, a.repositories.Swipe, a.repositories.Photo, a.logger),
		Report:   service.NewSupportService(a.repositories.Report, a.logger, telegramService),
		Telegram: telegramService,
	}
}

// startTelegramListener запускает обработчик callback'ов из Telegram
func (a *App) startTelegramListener() {
	if !a.config.Telegram.Enabled {
		a.logger.Info("Telegram listener disabled")
		return
	}

	a.logger.Info("Starting Telegram callback listener...")

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	a.logger.Info("Telegram callback listener started")
	updates := a.services.Telegram.GetBot().GetUpdatesChan(u)
	a.logger.Debugf("Updates: %v", updates)
	go func() {
		for update := range updates {
			if update.CallbackQuery != nil {
				a.handleTelegramCallback(update.CallbackQuery)
			}
		}
	}()
}

// handleTelegramCallback обрабатывает нажатия кнопок в Telegram
func (a *App) handleTelegramCallback(callback *tgbotapi.CallbackQuery) {
	a.logger.Debugf("Received Telegram callback: %s", callback.Data)

	// Обрабатываем callback данные
	ticketID, action, err := a.services.Telegram.ProcessCallback(callback.Data)
	if err != nil {
		a.logger.Errorf("Invalid callback data: %v", err)
		return
	}

	// Конвертируем action в статус
	var newStatus constants.ReportStatus
	switch action {
	case "work":
		newStatus = constants.ReportStatusWork
	case "close":
		newStatus = constants.ReportStatusClosed
	default:
		a.logger.Errorf("Unknown action: %s", action)
		return
	}

	// Обновляем статус в БД
	ticketUUID, err := uuid.Parse(ticketID)
	if err != nil {
		a.logger.Errorf("Invalid ticket ID: %v", err)
		return
	}
	a.logger.Infof("Updating ticket %v", ticketUUID)
	err = a.repositories.Report.UpdateStatus(ticketUUID, newStatus)
	if err != nil {
		a.logger.Errorf("Failed to update ticket status: %v", err)
		return
	}

	a.logger.Info("Get Report")
	// Получаем обновленный тикет
	ticket, err := a.repositories.Report.GetById(ticketUUID)
	if err != nil {
		a.logger.Errorf("Failed to get updated ticket: %v", err)
		return
	}

	a.logger.Infof("Updating ticket Message %v", ticket)
	// Обновляем сообщение в Telegram
	err = a.services.Telegram.UpdateTicketMessage(
		callback.Message.Chat.ID,
		callback.Message.MessageID,
		ticket,
	)
	if err != nil {
		a.logger.Errorf("Failed to update telegram message: %v", err)
		return
	}

	a.logger.Info("Callback")
	// Отправляем ответ на callback (убирает "часики" у кнопки)
	callbackConfig := tgbotapi.NewCallback(callback.ID, "✅ Статус обновлен!")
	if _, err := a.services.Telegram.GetBot().Request(callbackConfig); err != nil {
		a.logger.Errorf("Failed to send callback answer: %v", err)
	}

	a.logger.Infof("Ticket %s status updated to %s", ticketID, newStatus)
}
