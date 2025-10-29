package interfaces

import (
	"context"

	domain "github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/domain/entities"
	"github.com/google/uuid"
)

type MatchRepository interface {
	Create(ctx context.Context, match *domain.Match) error
	GetByUsers(ctx context.Context, user1ID, user2ID uuid.UUID) (*domain.Match, error)
	GetUserMatches(ctx context.Context, userID uuid.UUID, limit, offset int) ([]domain.Match, error)
	UpdateActive(ctx context.Context, user1ID, user2ID uuid.UUID, isActive bool) error
	Delete(ctx context.Context, user1ID, user2ID uuid.UUID) error
	CheckMutualLike(ctx context.Context, user1ID, user2ID uuid.UUID) (bool, error)
}
