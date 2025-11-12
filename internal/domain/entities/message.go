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
