package service

import (
	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/domain/entities"
	"github.com/google/uuid"
)

// AuthService defines authentication service contract
type AuthService interface {
	// Register creates new user
	Register(email, password, passwordConfirm string) (*domain.User, *domain.Session, error)

	// Login authenticates user
	Login(email, password string) (*domain.User, *domain.Session, error)

	// Logout terminates user session
	Logout(token string) error

	// ValidateSession checks session validity
	ValidateSession(token string) (*domain.User, error)

	// GetUserByID returns user by ID
	GetUserByID(userID uuid.UUID) (*domain.User, error)
}
