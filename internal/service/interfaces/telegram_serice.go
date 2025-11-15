package service

import (
	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/domain/entities"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type TelegramService interface {
	// SendNewTicketNotification отправляет уведомление о новом тикете в Telegram
	SendNewTicketNotification(ticket *domain.Report) (int, error)
	// UpdateTicketMessage updates message with new status
	UpdateTicketMessage(chatID int64, messageID int, ticket *domain.Report) error

	// ProcessCallback handles button clicks
	ProcessCallback(callbackData string) (ticketID string, action string, err error)
	GetBot() *tgbotapi.BotAPI // Добавляем этот метод
}
