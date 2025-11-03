package dto

import (
	"time"

	"github.com/google/uuid"
)

type Conversation struct {
	MatchID         uuid.UUID `json:"match_id"`
	OtherUserID     uuid.UUID `json:"other_user_id"`
	OtherUserName   string    `json:"other_user_name"`
	LastMessage     string    `json:"last_message"`
	LastMessageTime time.Time `json:"last_message_time"`
	UnreadCount     int       `json:"unread_count"`
}
