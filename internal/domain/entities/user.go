package domain

import (
	"time"

	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/domain/constants"
	"github.com/google/uuid"
)

type User struct {
	ID         uuid.UUID        `json:"id" db:"id"`
	Email      string           `json:"email" db:"email"`
	Phone      *string          `json:"phone,omitempty" db:"phone"`
	Name       string           `json:"name" db:"name"`
	Password   string           `json:"-" db:"password"`
	BirthDate  time.Time        `json:"birth_date" db:"birth_date"`
	Gender     constants.Gender `json:"gender" db:"gender"`
	Bio        *string          `json:"bio,omitempty" db:"bio"`
	City       *string          `json:"city,omitempty" db:"city"`
	Artist     *string          `json:"artist,omitempty" db:"artist"`
	Quote      *string          `json:"quote,omitempty" db:"quote"`
	IsVerified bool             `json:"is_verified" db:"is_verified"`
	LastActive time.Time        `json:"last_active" db:"last_active"`
	CreatedAt  time.Time        `json:"created_at" db:"created_at"`
	UpdatedAt  time.Time        `json:"updated_at" db:"updated_at"`
}
