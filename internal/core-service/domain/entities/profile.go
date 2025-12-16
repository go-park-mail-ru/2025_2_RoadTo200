package domain

import (
	"time"

	cons "github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/auth-service/domain/constants"
	domain "github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/auth-service/domain/entities"
	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/core-service/domain/constants"
)

// TODO: remove user

// ProfileResponse представляет полный ответ профиля для фронта
type ProfileResponse struct {
	User        *domain.User    `json:"user"`
	Preferences *UserPreference `json:"preferences,omitempty"`
	Photos      []UserPhoto     `json:"photos,omitempty"`
	Interests   []Interest      `json:"interests,omitempty"`
}

// ProfileUpdateRequest запрос на обновление профиля
type ProfileUpdateRequest struct {
	Name      string      `json:"name,omitempty"`
	Phone     *string     `json:"phone,omitempty"`
	BirthDate *time.Time  `json:"birth_date,omitempty"`
	Gender    cons.Gender `json:"gender,omitempty"`
	Bio       *string     `json:"bio,omitempty"`
	City      *string     `json:"city,omitempty"`
	Artist    *string     `json:"artist,omitempty"`
	Quote     *string     `json:"quote,omitempty"`
	Latitude  *float64    `json:"latitude,omitempty"`
	Longitude *float64    `json:"longitude,omitempty"`
}

// PreferencesUpdateRequest запрос на обновление предпочтений
type PreferencesUpdateRequest struct {
	ShowGender   constants.GenderPreference `json:"show_gender,omitempty"`
	AgeMin       int                        `json:"age_min,omitempty"`
	AgeMax       int                        `json:"age_max,omitempty"`
	MaxDistance  int                        `json:"max_distance,omitempty"`
	GlobalSearch bool                       `json:"global_search,omitempty"`
}
