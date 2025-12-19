package domain

//go:generate easyjson -all -no_std_marshalers match.go

import (
	"time"

	"github.com/google/uuid"
)

type Match struct {
	ID        uuid.UUID  `json:"id" db:"id"`
	User1ID   uuid.UUID  `json:"user1_id" db:"user1_id"`
	User2ID   uuid.UUID  `json:"user2_id" db:"user2_id"`
	IsActive  bool       `json:"is_active" db:"is_active"`
	MatchedAt time.Time  `json:"matched_at" db:"matched_at"`
	ExpiresAt *time.Time `json:"expires_at" db:"expires_at"` // NULL если кто-то написал сообщение (матч активен навсегда)
}
