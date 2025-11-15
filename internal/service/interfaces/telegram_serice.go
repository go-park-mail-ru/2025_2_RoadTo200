package service

import (
	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/handler/dto"
)

type TelegramService interface {
	// SendNewTicketNotification отправляет уведомление о новом тикете в Telegram
	SendNewTicketNotification(ticket *dto.SupportTicketResponse) error
}
