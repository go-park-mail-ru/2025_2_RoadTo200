package domain

import (
	"time"

	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/domain/constants"
	"github.com/google/uuid"
)

type Notification struct {
	ID         uuid.UUID                  `json:"id" db:"id"`
	UserID     uuid.UUID                  `json:"user_id" db:"user_id"`
	Type       constants.NotificationType `json:"type" db:"type"`
	FromUserID *uuid.UUID                 `json:"from_user_id,omitempty" db:"from_user_id"`
	MatchID    *uuid.UUID                 `json:"match_id,omitempty" db:"match_id"`
	IsRead     bool                       `json:"is_read" db:"is_read"`
	CreatedAt  time.Time                  `json:"created_at" db:"created_at"`
}

// NotificationMessage - WebSocket message format for notifications
type NotificationMessage struct {
	Type       string    `json:"type"` // "notification"
	ID         uuid.UUID `json:"id"`
	NotifType  string    `json:"notif_type"`   // "match", "super_like", "like"
	FromUserID uuid.UUID `json:"from_user_id"` // Всегда заполнен, даже для матчей
	CreatedAt  time.Time `json:"created_at"`
}
