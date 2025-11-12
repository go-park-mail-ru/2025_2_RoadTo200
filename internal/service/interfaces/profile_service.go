package service

import (
	"mime/multipart"

	domain "github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/domain/entities"
	"github.com/google/uuid"
)

type ProfileService interface {
	// Profile
	GetProfile(userID uuid.UUID) (*domain.ProfileResponse, error)
	UpdateProfileInfo(userID uuid.UUID, updateData *domain.ProfileUpdateRequest) error

	// Preferences
	UpdatePreferences(userID uuid.UUID, updateData *domain.PreferencesUpdateRequest) error
	UpdateInterests(uuid.UUID, []domain.Interest) error

	// Photos
	UploadPhotos(userID uuid.UUID, photos []*multipart.FileHeader) ([]domain.UserPhoto, error)
	DeletePhoto(userID uuid.UUID, photoID uuid.UUID) error
	SetPrimaryPhoto(userID uuid.UUID, photoID uuid.UUID) error
	ReorderPhotos(userID uuid.UUID, photoIDs []uuid.UUID) error

	// Валидация
	ValidateProfileUpdate(updateData *domain.ProfileUpdateRequest) error
	ValidatePreferencesUpdate(updateData *domain.PreferencesUpdateRequest) error
}
