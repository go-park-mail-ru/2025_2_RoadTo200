package service

import (
	"context"
	"time"

	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/auth-service/domain/entities"
	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/auth-service/repository/interfaces"
	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/domain/errors"
	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/pkg/logger"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	userRepo    interfaces.UserRepository
	sessionRepo interfaces.SessionRepository
	logger      logger.Log
}

func NewAuthService(userRepo interfaces.UserRepository, sessionRepo interfaces.SessionRepository, l logger.Log) *AuthService {
	return &AuthService{
		userRepo:    userRepo,
		sessionRepo: sessionRepo,
		logger:      l,
	}
}

func (s *AuthService) Register(ctx context.Context, email, password, passwordConfirm string) (*domain.User, *domain.Session, error) {
	s.logger.Trace("AuthService.Register")
	// Валидация
	if password != passwordConfirm {
		return nil, nil, errors.ErrPasswordsDontMatch
	}

	if len(password) < 6 {
		return nil, nil, errors.ErrPasswordTooShort
	}

	// Проверка существования пользователя
	existingUser, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		s.logger.Errorf("Get by email error: %s", err)
		return nil, nil, err
	}
	if existingUser != nil {
		return nil, nil, errors.ErrUserAlreadyExists
	}

	// Хэширование пароля
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		s.logger.Errorf("hashing password error: %s", err)
		return nil, nil, err
	}

	// Создание пользователя
	user := &domain.User{
		ID:         uuid.New(),
		Email:      email,
		Password:   string(hashedPassword),
		Name:       email,  // можно генерировать или оставить пустым
		Gender:     "male", // ← пустая строка вместо NULL
		IsVerified: true,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		s.logger.Errorf("Create user error: %s", err)
		return nil, nil, err
	}

	// Создание сессии
	session := domain.NewSession(user.Email, 3600*time.Second)
	session.Token = uuid.New().String()

	if err := s.sessionRepo.Set(ctx, session); err != nil {
		s.logger.Errorf("Set session error: %s", err)
		return nil, nil, err
	}

	return user, session, nil
}

func (s *AuthService) Login(ctx context.Context, email, password string) (*domain.User, *domain.Session, error) {
	s.logger.Trace("AuthService.Login")

	user, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		s.logger.Errorf("Get by email error: %s", err)
		return nil, nil, err
	}
	if user == nil {
		s.logger.Warn("User not found")
		return nil, nil, errors.ErrInvalidCredentials
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		s.logger.Errorf("CompareHashAndPassword error: %s", err)
		return nil, nil, errors.ErrInvalidCredentials
	}

	// Обновляем last_active
	if err := s.userRepo.UpdateLastActive(ctx, user.ID); err != nil {
		s.logger.Errorf("UpdateLastActive error: %s", err)
		return nil, nil, err
	}

	// Создаем новую сессию
	session := domain.NewSession(user.Email, 3600*time.Second)
	session.Token = uuid.New().String()

	if err := s.sessionRepo.Set(ctx, session); err != nil {
		s.logger.Errorf("Set session error: %s", err)
		return nil, nil, err
	}

	return user, session, nil
}

func (s *AuthService) Logout(ctx context.Context, token string) error {
	s.logger.Trace("AuthService.Logout")

	session, err := s.sessionRepo.Get(ctx, token)
	if err != nil {
		s.logger.Errorf("Get session error: %s", err)
		return err
	}
	if session == nil {
		s.logger.Warn("Session not found")
		return errors.ErrSessionNotFound
	}

	return s.sessionRepo.Delete(ctx, session.Token)
}

func (s *AuthService) ValidateSession(ctx context.Context, token string) (*domain.User, error) {
	session, err := s.sessionRepo.Get(ctx, token)
	if err != nil {
		s.logger.Errorf("Get session error: %s", err)
		return nil, err
	}
	if session == nil {
		s.logger.Warn("Session not found")
		return nil, errors.ErrSessionNotFound
	}

	if session.IsExpired() {
		s.logger.Warn("Session is expired")
		go s.sessionRepo.Delete(ctx, session.Token)
		return nil, errors.ErrSessionExpired
	}

	user, err := s.userRepo.GetByEmail(ctx, session.UserEmail)
	if err != nil {
		s.logger.Errorf("Get by email error: %s", err)
		return nil, err
	}
	if user == nil {
		s.logger.Warn("User not found")
		return nil, errors.ErrUserNotFound
	}

	return user, nil
}

func (s *AuthService) GetUserByID(ctx context.Context, userID uuid.UUID) (*domain.User, error) {
	return s.userRepo.GetByID(ctx, userID)
}
