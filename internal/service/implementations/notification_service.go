package service

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/domain/constants"
	domain "github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/domain/entities"
	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/logger"
	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/repository/interfaces"
	serviceInterfaces "github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/service/interfaces"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

var _ serviceInterfaces.NotificationService = (*NotificationService)(nil)

type NotificationService struct {
	notificationRepo interfaces.NotificationRepository
	redisClient      *redis.Client
	logger           logger.Log
}

func NewNotificationService(
	notificationRepo interfaces.NotificationRepository,
	redisClient *redis.Client,
	logger logger.Log,
) *NotificationService {
	return &NotificationService{
		notificationRepo: notificationRepo,
		redisClient:      redisClient,
		logger:           logger,
	}
}

// SendNotification создает уведомление в БД и отправляет через WebSocket
func (s *NotificationService) SendNotification(ctx context.Context, userID uuid.UUID, notificationType constants.NotificationType, fromUserID *uuid.UUID, matchID *uuid.UUID) error {
	s.logger.Tracef("NotificationService.SendNotification: userID=%s, type=%s", userID, notificationType)

	// Создаем уведомление
	notification := &domain.Notification{
		UserID:     userID,
		Type:       notificationType,
		FromUserID: fromUserID,
		MatchID:    matchID,
		IsRead:     false,
	}

	// Сохраняем в БД
	if err := s.notificationRepo.Create(ctx, notification); err != nil {
		s.logger.Errorf("Failed to create notification: %v", err)
		return fmt.Errorf("failed to create notification: %w", err)
	}

	// Публикуем в Redis для отправки через WebSocket
	if err := s.publishNotification(ctx, notification); err != nil {
		s.logger.Warnf("Failed to publish notification to Redis: %v", err)
		// Не возвращаем ошибку, уведомление уже сохранено в БД
	}

	return nil
}

// publishNotification публикует уведомление в Redis канал для WebSocket
func (s *NotificationService) publishNotification(ctx context.Context, notification *domain.Notification) error {
	// Skip if Redis client is not configured
	if s.redisClient == nil {
		return nil
	}

	// Создаем WebSocket сообщение
	// Убеждаемся, что from_user_id всегда заполнен
	var fromUserID uuid.UUID
	if notification.FromUserID != nil {
		fromUserID = *notification.FromUserID
	} else {
		// Если from_user_id не заполнен (не должно быть, но на всякий случай)
		s.logger.Warnf("Notification without from_user_id: %s", notification.ID)
		fromUserID = uuid.Nil
	}

	notifMsg := domain.NotificationMessage{
		Type:       "notification",
		ID:         notification.ID,
		NotifType:  string(notification.Type),
		FromUserID: fromUserID,
		CreatedAt:  notification.CreatedAt,
	}

	data, err := json.Marshal(notifMsg)
	if err != nil {
		return fmt.Errorf("failed to marshal notification: %w", err)
	}

	// Публикуем в канал для пользователя
	channel := fmt.Sprintf("notification:user:%s", notification.UserID.String())
	if err := s.redisClient.Publish(ctx, channel, data).Err(); err != nil {
		return fmt.Errorf("failed to publish to channel %s: %w", channel, err)
	}

	s.logger.Debugf("Published notification to channel %s: type=%s", channel, notification.Type)
	return nil
}

// GetNotifications возвращает список уведомлений пользователя
func (s *NotificationService) GetNotifications(ctx context.Context, userID uuid.UUID, limit, offset int) ([]domain.Notification, error) {
	s.logger.Tracef("NotificationService.GetNotifications: userID=%s, limit=%d, offset=%d", userID, limit, offset)

	notifications, err := s.notificationRepo.GetByUserID(ctx, userID, limit, offset)
	if err != nil {
		s.logger.Errorf("Failed to get notifications: %v", err)
		return nil, fmt.Errorf("failed to get notifications: %w", err)
	}

	return notifications, nil
}

// MarkAsRead помечает уведомление как прочитанное
func (s *NotificationService) MarkAsRead(ctx context.Context, userID, notificationID uuid.UUID) error {
	s.logger.Tracef("NotificationService.MarkAsRead: userID=%s, notificationID=%s", userID, notificationID)

	// Получаем уведомление для проверки принадлежности пользователю
	// Используем GetByUserID с большим лимитом, чтобы найти нужное уведомление
	notifications, err := s.notificationRepo.GetByUserID(ctx, userID, 1000, 0)
	if err != nil {
		return fmt.Errorf("failed to get notifications: %w", err)
	}

	// Проверяем, что уведомление существует и принадлежит пользователю
	found := false
	for _, notif := range notifications {
		if notif.ID == notificationID {
			found = true
			break
		}
	}

	if !found {
		return fmt.Errorf("notification not found or access denied")
	}

	// Помечаем как прочитанное
	if err := s.notificationRepo.MarkAsRead(ctx, notificationID); err != nil {
		s.logger.Errorf("Failed to mark notification as read: %v", err)
		return fmt.Errorf("failed to mark notification as read: %w", err)
	}

	return nil
}
