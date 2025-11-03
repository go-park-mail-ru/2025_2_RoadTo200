package interfaces

import (
	domain "github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/domain/entities"
	"github.com/google/uuid"
)

type UserRepository interface {
	Create(user *domain.User) error
	GetByID(id uuid.UUID) (*domain.User, error)
	GetByEmail(email string) (*domain.User, error)
	GetByPhone(phone string) (*domain.User, error)
	Update(user *domain.User) error
	UpdateLastActive(userID uuid.UUID) error
	Delete(id uuid.UUID) error
	GetUsersByIDs(ids []uuid.UUID) ([]domain.User, error)
	GetUsersForFeed(userID uuid.UUID, limit, offset int) ([]domain.User, error)
}

type UserPhotoRepository interface {
	Create(photo *domain.UserPhoto) error
	GetByID(id uuid.UUID) (*domain.UserPhoto, error)
	GetByUserID(userID uuid.UUID) ([]domain.UserPhoto, error)
	Update(photo *domain.UserPhoto) error
	Delete(id uuid.UUID) error
	UpdateDisplayOrder(userID uuid.UUID, photos []domain.UserPhoto) error
}

type UserPreferenceRepository interface {
	Create(preference *domain.UserPreference) error
	GetByUserID(userID uuid.UUID) (*domain.UserPreference, error)
	Update(preference *domain.UserPreference) error
	Delete(userID uuid.UUID) error
}
