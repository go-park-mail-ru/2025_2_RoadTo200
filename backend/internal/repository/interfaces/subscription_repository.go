package interfaces

import (
	domain "github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/domain/entities"
	"github.com/google/uuid"
)

type SubscriptionRepository interface {
	Create(subscription *domain.Subscription) error
	GetByUserID(userID uuid.UUID) (*domain.Subscription, error)
	Update(subscription *domain.Subscription) error
	Delete(userID uuid.UUID) error
	GetActiveSubscription(userID uuid.UUID) (*domain.Subscription, error)
}
