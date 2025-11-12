package interfaces

import (
	domain "github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/domain/entities"
	"github.com/google/uuid"
)

type SwipeRepository interface {
	Create(swipe *domain.Swipe) error
	GetBySwiperAndTarget(swiperID, targetID uuid.UUID) (*domain.Swipe, error)
	GetSwipesBySwiper(swiperID uuid.UUID, limit, offset int) ([]domain.Swipe, error)
	GetSwipesByTarget(targetID uuid.UUID, limit, offset int) ([]domain.Swipe, error)
	Exists(swiperID, targetID uuid.UUID) (bool, error)
	GetSwipesStats(userID uuid.UUID) (likesCount, dislikesCount, superLikesCount int, err error)
	GetMutualLikes(userID uuid.UUID) ([]uuid.UUID, error)
}
