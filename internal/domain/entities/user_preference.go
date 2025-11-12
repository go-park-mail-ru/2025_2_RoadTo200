package domain

import (
	"time"

	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/domain/constants"
	"github.com/google/uuid"
)

type UserPreference struct {
	UserID       uuid.UUID                  `json:"user_id" db:"user_id"`
	ShowGender   constants.GenderPreference `json:"show_gender" db:"show_gender"`
	AgeMin       int                        `json:"age_min" db:"age_min"`
	AgeMax       int                        `json:"age_max" db:"age_max"`
	MaxDistance  int                        `json:"max_distance" db:"max_distance"`
	GlobalSearch bool                       `json:"global_search" db:"global_search"`
	CreatedAt    time.Time                  `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time                  `json:"updated_at" db:"updated_at"`
}
