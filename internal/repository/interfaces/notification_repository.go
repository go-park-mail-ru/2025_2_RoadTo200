package interfaces

import (
	"context"

	domain "github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/domain/entities"
	"github.com/google/uuid"
)

type NotificationRepository interface {
	Create(ctx context.Context, notification *domain.Notification) error
	GetByUserID(ctx context.Context, userID uuid.UUID, limit, offset int) ([]domain.Notification, error)
	GetUnreadCount(ctx context.Context, userID uuid.UUID) (int, error)
	MarkAsRead(ctx context.Context, notificationID uuid.UUID) error
	MarkAllAsRead(ctx context.Context, userID uuid.UUID) error
	DeleteUserNotifications(ctx context.Context, user1ID, user2ID uuid.UUID) error
}
