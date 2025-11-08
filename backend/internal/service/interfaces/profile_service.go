package service

import (
	"mime/multipart"
	"time"

	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/domain/constants"
	domain "github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/domain/entities"
	"github.com/google/uuid"
)

type ProfileService interface {
	// Profile
	GetProfile(userID uuid.UUID) (*domain.ProfileResponse, error)
	UpdateProfileInfo(userID uuid.UUID, updateData *domain.ProfileUpdateRequest) error

	// Preferences
	UpdatePreferences(userID uuid.UUID, updateData *domain.PreferencesUpdateRequest) error

	// Photos
	UploadPhotos(userID uuid.UUID, photos []*multipart.FileHeader) ([]domain.UserPhoto, error)
	DeletePhoto(userID uuid.UUID, photoID uuid.UUID) error
	SetPrimaryPhoto(userID uuid.UUID, photoID uuid.UUID) error
	ReorderPhotos(userID uuid.UUID, photoIDs []uuid.UUID) error

	// Валидация
	ValidateProfileUpdate(updateData *domain.ProfileUpdateRequest) error
	ValidatePreferencesUpdate(updateData *domain.PreferencesUpdateRequest) error
}

// TODO: Разбить на составляющие
// UpdateProfileRequest общий запрос на изменение профиля
type UpdateProfileRequest struct {
	Action string `json:"action"` // "updateInfo", "updatePreferences", "deletePhoto", "setPrimaryPhoto"

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
	PhotoIDs []uuid.UUID `json:"photo_ids,omitempty"`
}

// ProfileResponse ответ профиля
type ProfileResponse struct {
	User        interface{} `json:"user"`
	Preferences interface{} `json:"preferences,omitempty"`
	Photos      interface{} `json:"photos,omitempty"`
}

// UploadPhotosResponse ответ загрузки фото
type UploadPhotosResponse struct {
	Photos []interface{} `json:"photos"`
}

// SuccessResponse общий успешный ответ
type SuccessResponse struct {
	Message string `json:"message"`
}

// ErrorResponse ответ с ошибкой
type ErrorResponse struct {
	Error string `json:"error"`
}
