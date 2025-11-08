package dto

import (
	"time"

	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/domain/constants"
	"github.com/google/uuid"
)

// TODO: Разбить на составляющие
// UpdateProfileRequest общий запрос на изменение профиля
type UpdateProfileRequest struct {
	// Для updateInfo
	Name      string           `json:"name,omitempty"`
	Phone     *string          `json:"phone,omitempty"`
	BirthDate *time.Time       `json:"birth_date,omitempty"`
	Gender    constants.Gender `json:"gender,omitempty"`
	Bio       *string          `json:"bio,omitempty"`
	Artist    *string          `json:"artist,omitempty"`
	Quote     *string          `json:"quote,omitempty"`
	Latitude  *float64         `json:"latitude,omitempty"`
	Longitude *float64         `json:"longitude,omitempty"`

	// Для updatePreferences
	ShowGender   constants.GenderPreference `json:"show_gender,omitempty"`
	AgeMin       int                        `json:"age_min,omitempty"`
	AgeMax       int                        `json:"age_max,omitempty"`
	MaxDistance  int                        `json:"max_distance,omitempty"`
	GlobalSearch bool                       `json:"global_search,omitempty"`

	// Для deletePhoto и setPrimaryPhoto
	PhotoID uuid.UUID `json:"photo_id,omitempty"`

	// Для reorderPhotos
	//PhotoIDs []uuid.UUID `json:"photo_ids,omitempty"`
}

// ProfileResponse ответ профиля
type ProfileResponse struct {
	User        interface{} `json:"user"`
	Preferences interface{} `json:"preferences,omitempty"`
	Photos      interface{} `json:"photos,omitempty"`
}
