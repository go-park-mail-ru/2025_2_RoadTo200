package repository

import (
	//"time"

	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/domain/entities"
)

type UserRepository interface {
	Create(user *domain.User) error
	FindByEmail(email string) (*domain.User, error)
	FindByID(id string) (*domain.User, error)
	Exists(email string) (bool, error)
}

type SessionRepository interface {
	Create(session *domain.Session) error
	FindByToken(token string) (*domain.Session, error)
	Delete(token string) error
	DeleteExpired() error
}
