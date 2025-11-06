package service

import (
	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/domain/entities"
	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/repository/interfaces"
	"github.com/google/uuid"
)

type AuthService interface {
	NewAuthService(userRepo interfaces.UserRepository, sessionRepo interfaces.SessionRepository)
	Register(email, password, passwordConfirm string) (*domain.User, *domain.Session, error)
	Login(email, password string) (*domain.User, *domain.Session, error)
	Logout(token string) error
	ValidateSession(token string) (*domain.User, error)
	GetUserByID(userID uuid.UUID) (*domain.User, error)
}

type RegisterRequest struct {
	Email           string `json:"email" binding:"required"`
	Password        string `json:"password" binding:"required"`
	PasswordConfirm string `json:"passwordConfirm" binding:"required"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}
