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
