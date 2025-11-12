package interfaces

import (
	"time"

	domain "github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/domain/entities"
)

type SessionRepository interface {
	Set(session *domain.Session) error
	Get(sessionID string) (*domain.Session, error)
	Delete(sessionID string) error
	Exists(sessionID string) (bool, error)
	SetWithExpiry(session *domain.Session, expiry time.Duration) error
}
