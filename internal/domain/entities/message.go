package domain

import (
	"time"

	"github.com/google/uuid"
)

// Message - database entity
type Message struct {
	ID         uuid.UUID `json:"id" db:"id"`
	MatchID    uuid.UUID `json:"match_id" db:"match_id"`
	SenderID   uuid.UUID `json:"sender_id" db:"sender_id"`
	ReceiverID uuid.UUID `json:"receiver_id" db:"receiver_id"`
	Content    string    `json:"content" db:"content"`
	IsRead     bool      `json:"is_read" db:"is_read"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
}

// ChatMessage - WebSocket message format
type ChatMessage struct {
	Type      string    `json:"type"` // "message", "typing", "read", "error"
	MessageID uuid.UUID `json:"message_id,omitempty"`
	MatchID   uuid.UUID `json:"match_id"`
	SenderID  uuid.UUID `json:"sender_id,omitempty"`
	Content   string    `json:"content,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	Error     string    `json:"error,omitempty"`
}

type LastMessage struct {
	MatchID   uuid.UUID `json:"message_id,omitempty"`
	Content   string    `json:"content,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
}

type UnreadCount struct {
	MatchID uuid.UUID `json:"message_id,omitempty"`
	Amount  int       `json:"content,omitempty"`
}

// Conversation - chat preview for conversations list
type Conversation struct {
	MatchID         uuid.UUID `json:"match_id"`
	OtherUserID     uuid.UUID `json:"other_user_id"`
	OtherUserName   string    `json:"other_user_name"`
	OtherUserPhoto  string    `json:"other_user_photo"`
	LastMessage     string    `json:"last_message"`
	LastMessageTime time.Time `json:"last_message_time"`
	UnreadCount     int       `json:"unread_count"`
}
