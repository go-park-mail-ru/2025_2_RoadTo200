package domain

import (
	"time"

	"github.com/google/uuid"
)

type UserPhoto struct {
	ID           uuid.UUID `json:"id" db:"id"`
	UserID       uuid.UUID `json:"user_id" db:"user_id"`
	PhotoURL     string    `json:"photo_url" db:"photo_url"`
	DisplayOrder int       `json:"display_order" db:"display_order"`
	IsApproved   bool      `json:"is_approved" db:"is_approved"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
}
