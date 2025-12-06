package interfaces

import (
	"context"

	domain "github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/domain/entities"
	"github.com/google/uuid"
)

type SwipeRepository interface {
	Create(ctx context.Context, swipe *domain.Swipe) error
	GetBySwiperAndTarget(ctx context.Context, swiperID, targetID uuid.UUID) (*domain.Swipe, error)
	GetSwipesBySwiper(ctx context.Context, swiperID uuid.UUID, limit, offset int) ([]domain.Swipe, error)
	GetSwipesByTarget(ctx context.Context, targetID uuid.UUID, limit, offset int) ([]domain.Swipe, error)
	Exists(ctx context.Context, swiperID, targetID uuid.UUID) (bool, error)
	GetSwipesStats(ctx context.Context, userID uuid.UUID) (likesCount, dislikesCount, superLikesCount int, err error)
	GetMutualLikes(ctx context.Context, userID uuid.UUID) ([]uuid.UUID, error)
	GetUsersForFeed(ctx context.Context, userID uuid.UUID, limit, offset int) ([]uuid.UUID, error)
}
