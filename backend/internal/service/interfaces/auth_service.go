package service

import (
	"time"

	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/domain/entities"
	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/domain/errors"
	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/repository/interfaces"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	userRepo    repository.UserRepository
	sessionRepo repository.SessionRepository
}

func NewAuthService(userRepo repository.UserRepository, sessionRepo repository.SessionRepository) *AuthService {
	return &AuthService{
		userRepo:    userRepo,
		sessionRepo: sessionRepo,
	}
}

func (s *AuthService) Register(email, password, passwordConfirm string) (*domain.User, *domain.Session, error) {
	// Валидация
	if password != passwordConfirm {
		return nil, nil, errors.ErrPasswordsDontMatch
	}

	if len(password) < 6 {
		return nil, nil, errors.ErrPasswordTooShort
	}

	// Проверка существования пользователя
	exists, err := s.userRepo.Exists(email)
	if err != nil {
		return nil, nil, err
	}
	if exists {
		return nil, nil, errors.ErrUserAlreadyExists
	}

	// Хэширование пароля
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, nil, err
	}

	// Создание пользователя
	user := domain.NewUser(email, string(hashedPassword))
	user.ID = uuid.New().String()

	if err := s.userRepo.Create(user); err != nil {
		return nil, nil, err
	}

	// Создание сессии
	session := domain.NewSession(email, 3600*time.Second)
	session.Token = uuid.New().String()

	if err := s.sessionRepo.Create(session); err != nil {
		return nil, nil, err
	}

	return user, session, nil
}

func (s *AuthService) Login(email, password string) (*domain.User, *domain.Session, error) {
	user, err := s.userRepo.FindByEmail(email)
	if err != nil {
		return nil, nil, errors.ErrInvalidCredentials
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, nil, errors.ErrInvalidCredentials
	}

	session := domain.NewSession(email, 3600*time.Second)
	session.Token = uuid.New().String()

	if err := s.sessionRepo.Create(session); err != nil {
		return nil, nil, err
	}

	return user, session, nil
}

func (s *AuthService) Logout(token string) error {
	return s.sessionRepo.Delete(token)
}

func (s *AuthService) ValidateSession(token string) (*domain.User, error) {
	session, err := s.sessionRepo.FindByToken(token)
	if err != nil {
		return nil, err
	}

	user, err := s.userRepo.FindByEmail(session.UserEmail)
	if err != nil {
		return nil, errors.ErrUserNotFound
	}

	return user, nil
}
