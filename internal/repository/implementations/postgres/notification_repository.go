package postgres

import (
	"context"

	domain "github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/domain/entities"
	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/repository/interfaces"
	"github.com/google/uuid"
)

var _ interfaces.NotificationRepository = (*NotificationRepository)(nil)

type NotificationRepository struct {
	pool interfaces.PgxIface
}

func NewNotificationRepository(pool interfaces.PgxIface) *NotificationRepository {
	return &NotificationRepository{pool: pool}
}

func (r *NotificationRepository) Create(ctx context.Context, notification *domain.Notification) error {
	query := `
		INSERT INTO notification (user_id, type, from_user_id, match_id, is_read)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at`

	err := r.pool.QueryRow(ctx, query,
		notification.UserID, notification.Type, notification.FromUserID, notification.MatchID, notification.IsRead).
		Scan(&notification.ID, &notification.CreatedAt)

	return err
}

func (r *NotificationRepository) GetByUserID(ctx context.Context, userID uuid.UUID, limit, offset int) ([]domain.Notification, error) {
	query := `
		SELECT id, user_id, type, from_user_id, match_id, is_read, created_at
		FROM notification
		WHERE user_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3`

	rows, err := r.pool.Query(ctx, query, userID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var notifications []domain.Notification
	for rows.Next() {
		var notification domain.Notification
		err := rows.Scan(
			&notification.ID, &notification.UserID, &notification.Type,
			&notification.FromUserID, &notification.MatchID, &notification.IsRead, &notification.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		notifications = append(notifications, notification)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return notifications, nil
}

func (r *NotificationRepository) GetUnreadCount(ctx context.Context, userID uuid.UUID) (int, error) {
	query := `SELECT COUNT(*) FROM notification WHERE user_id = $1 AND is_read = false`

	var count int
	err := r.pool.QueryRow(ctx, query, userID).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (r *NotificationRepository) MarkAsRead(ctx context.Context, notificationID uuid.UUID) error {
	query := `UPDATE notification SET is_read = true WHERE id = $1`

	_, err := r.pool.Exec(ctx, query, notificationID)
	return err
}

func (r *NotificationRepository) MarkAllAsRead(ctx context.Context, userID uuid.UUID) error {
	query := `UPDATE notification SET is_read = true WHERE user_id = $1 AND is_read = false`

	_, err := r.pool.Exec(ctx, query, userID)
	return err
}
