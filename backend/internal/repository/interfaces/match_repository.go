package interfaces

import (
	domain "github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/domain/entities"
	"github.com/google/uuid"
)

type MatchRepository interface {
	Create(match *domain.Match) error
	GetByUsers(user1ID, user2ID uuid.UUID) (*domain.Match, error)
	GetUserMatches(userID uuid.UUID, limit, offset int) ([]domain.Match, error)
	UpdateActive(user1ID, user2ID uuid.UUID, isActive bool) error
	Delete(user1ID, user2ID uuid.UUID) error
	CheckMutualLike(user1ID, user2ID uuid.UUID) (bool, error)
}
