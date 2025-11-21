package interfaces

import (
	"context"
	"time"

	domain "github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/domain/entities"
)

type SessionRepository interface {
	Set(ctx context.Context, session *domain.Session) error
	Get(ctx context.Context, sessionID string) (*domain.Session, error)
	Delete(ctx context.Context, sessionID string) error
	Exists(ctx context.Context, sessionID string) (bool, error)
	SetWithExpiry(ctx context.Context, session *domain.Session, expiry time.Duration) error
}
