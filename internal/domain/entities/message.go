package domain

import (
	"time"

	"github.com/google/uuid"
)

type Message struct {
	ID          uuid.UUID `json:"id" db:"id"`
	MatchID     uuid.UUID `json:"match_id" db:"match_id"`
	SenderID    uuid.UUID `json:"sender_id" db:"sender_id"`
	MessageText string    `json:"message_text" db:"message_text"`
	Status      int       `json:"status" db:"status"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
}

type Conversation struct {
	MatchID         uuid.UUID `json:"match_id"`
	OtherUserID     uuid.UUID `json:"other_user_id"`
	OtherUserName   string    `json:"other_user_name"`
	LastMessage     string    `json:"last_message"`
	LastMessageTime time.Time `json:"last_message_time"`
	UnreadCount     int       `json:"unread_count"`
}
