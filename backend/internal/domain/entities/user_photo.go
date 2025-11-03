package domain

import (
	"time"

	"github.com/google/uuid"
)

// IsPrimary проверяет является ли фото основным (первым в порядке)
func (p UserPhoto) IsPrimary() bool {
	return p.DisplayOrder == 0
}

type UserPhoto struct {
	ID           uuid.UUID `json:"id" db:"id"`
	UserID       uuid.UUID `json:"user_id" db:"user_id"`
	PhotoURL     string    `json:"photo_url" db:"photo_url"`
	DisplayOrder int       `json:"display_order" db:"display_order"`
	IsApproved   bool      `json:"is_approved" db:"is_approved"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
}

// ToResponse преобразует в формат для ответа
func (p UserPhoto) ToResponse() map[string]interface{} {
	return map[string]interface{}{
		"id":            p.ID,
		"user_id":       p.UserID,
		"photo_url":     p.PhotoURL,
		"display_order": p.DisplayOrder,
		"is_primary":    p.IsPrimary(),
		"is_approved":   p.IsApproved,
		"created_at":    p.CreatedAt,
	}
}
