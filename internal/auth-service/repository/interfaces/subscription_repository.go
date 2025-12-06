package interfaces

import (
	"context"

	domain "github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/auth-service/domain/entities"
	"github.com/google/uuid"
)

type SubscriptionRepository interface {
	Create(ctx context.Context, subscription *domain.Subscription) error
	GetByUserID(ctx context.Context, userID uuid.UUID) (*domain.Subscription, error)
	Update(ctx context.Context, subscription *domain.Subscription) error
	Delete(ctx context.Context, userID uuid.UUID) error
	GetActiveSubscription(ctx context.Context, userID uuid.UUID) (*domain.Subscription, error)
}
