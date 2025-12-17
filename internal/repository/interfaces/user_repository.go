package interfaces

import (
	"context"

	domain "github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/domain/entities"
	"github.com/google/uuid"
)

type UserRepository interface {
	Create(ctx context.Context, user *domain.User) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error)
	GetByEmail(ctx context.Context, email string) (*domain.User, error)
	GetByPhone(ctx context.Context, phone string) (*domain.User, error)
	Update(ctx context.Context, user *domain.User) error
	UpdateLastActive(ctx context.Context, userID uuid.UUID) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetUsersByIDs(ctx context.Context, ids []uuid.UUID) ([]domain.User, error)
	GetUsersForFeed(ctx context.Context, userID uuid.UUID, limit, offset int) ([]domain.User, error)
}

type UserPhotoRepository interface {
	Create(ctx context.Context, photo *domain.UserPhoto) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.UserPhoto, error)
	GetByUserID(ctx context.Context, userID uuid.UUID) ([]domain.UserPhoto, error)
	Update(ctx context.Context, photo *domain.UserPhoto) error
	Delete(ctx context.Context, id uuid.UUID) error
	UpdateDisplayOrder(ctx context.Context, userID uuid.UUID, photos []domain.UserPhoto) error
}

type UserPreferenceRepository interface {
	Create(ctx context.Context, preference *domain.UserPreference) error
	GetByUserID(ctx context.Context, userID uuid.UUID) (*domain.UserPreference, error)
	Update(ctx context.Context, preference *domain.UserPreference) error
	Delete(ctx context.Context, userID uuid.UUID) error
	GetInterests(ctx context.Context, userID uuid.UUID) ([]domain.Interest, error)
	UpdateInterests(ctx context.Context, userID uuid.UUID, interest []domain.Interest) error
}
