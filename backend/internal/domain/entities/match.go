package domain

import (
	"time"

	"github.com/google/uuid"
)

type Match struct {
	User1ID   uuid.UUID `json:"user1_id" db:"user1_id"`
	User2ID   uuid.UUID `json:"user2_id" db:"user2_id"`
	IsActive  bool      `json:"is_active" db:"is_active"`
	MatchedAt time.Time `json:"matched_at" db:"matched_at"`
}
