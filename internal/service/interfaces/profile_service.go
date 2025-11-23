package service

import (
	"context"
	"mime/multipart"

	domain "github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/domain/entities"
	"github.com/google/uuid"
)

type ProfileService interface {
	// Profile
	GetProfile(ctx context.Context, userID uuid.UUID) (*domain.ProfileResponse, error)
	UpdateProfileInfo(ctx context.Context, userID uuid.UUID, updateData *domain.ProfileUpdateRequest) error

	// Preferences
	UpdatePreferences(ctx context.Context, userID uuid.UUID, updateData *domain.PreferencesUpdateRequest) error
	UpdateInterests(ctx context.Context, userID uuid.UUID, interests []domain.Interest) error

	// Photos
	UploadPhotos(ctx context.Context, userID uuid.UUID, photos []*multipart.FileHeader) ([]domain.UserPhoto, error)
	DeletePhoto(ctx context.Context, userID uuid.UUID, photoID uuid.UUID) error
	SetPrimaryPhoto(ctx context.Context, userID uuid.UUID, photoID uuid.UUID) error
	ReorderPhotos(ctx context.Context, userID uuid.UUID, photoIDs []uuid.UUID) error

	// Валидация
	ValidateProfileUpdate(ctx context.Context, updateData *domain.ProfileUpdateRequest) error
	ValidatePreferencesUpdate(ctx context.Context, updateData *domain.PreferencesUpdateRequest) error
}
