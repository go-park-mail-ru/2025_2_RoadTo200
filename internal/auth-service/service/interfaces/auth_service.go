package service

import (
	"context"

	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/auth-service/domain/entities"
	"github.com/google/uuid"
)

// AuthService defines authentication service contract
type AuthService interface {
	// Register creates new user
	Register(ctx context.Context, email, password, passwordConfirm string) (*domain.User, *domain.Session, error)

	// Login authenticates user
	Login(ctx context.Context, email, password string) (*domain.User, *domain.Session, error)

	// Logout terminates user session
	Logout(ctx context.Context, token string) error

	// ValidateSession checks session validity
	ValidateSession(ctx context.Context, token string) (*domain.User, error)

	// GetUserByID returns user by ID
	GetUserByID(ctx context.Context, userID uuid.UUID) (*domain.User, error)
}
