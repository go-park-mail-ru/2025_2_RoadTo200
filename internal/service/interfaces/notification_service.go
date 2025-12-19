package service

import (
	"context"

	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/domain/constants"
	domain "github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/domain/entities"
	"github.com/google/uuid"
)

type NotificationService interface {
	// SendNotification создает уведомление в БД и отправляет через WebSocket
	SendNotification(ctx context.Context, userID uuid.UUID, notificationType constants.NotificationType, fromUserID *uuid.UUID, matchID *uuid.UUID) error
	// GetNotifications возвращает список уведомлений пользователя
	GetNotifications(ctx context.Context, userID uuid.UUID, limit, offset int) ([]domain.Notification, error)
	// MarkAsRead помечает уведомление как прочитанное
	MarkAsRead(ctx context.Context, userID, notificationID uuid.UUID) error
	// DeleteUserNotifications удаляет все уведомления (лайк, суперлайк, мэтч) между двумя пользователями
	DeleteUserNotifications(ctx context.Context, user1ID, user2ID uuid.UUID) error
}
