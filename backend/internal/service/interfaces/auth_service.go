package service

import (
	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/domain/entities"
	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/repository/interfaces"
	"github.com/google/uuid"
)

// AuthService defines authentication service contract
type AuthService interface {
	NewAuthService(userRepo interfaces.UserRepository, sessionRepo interfaces.SessionRepository)

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

// RegisterRequest represents registration request
// @Description Запрос для регистрации нового пользователя
type RegisterRequest struct {
	Email           string `json:"email" binding:"required" example:"user@example.com"`
	Password        string `json:"password" binding:"required" example:"securepassword123"`
	PasswordConfirm string `json:"passwordConfirm" binding:"required" example:"securepassword123"`
}

// LoginRequest represents login request
// @Description Запрос для аутентификации пользователя
type LoginRequest struct {
	Email    string `json:"email" binding:"required" example:"user@example.com"`
	Password string `json:"password" binding:"required" example:"securepassword123"`
}
