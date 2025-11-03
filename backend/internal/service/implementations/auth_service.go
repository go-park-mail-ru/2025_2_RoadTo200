package service

import (
	"time"

	//"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/domain"
	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/domain/entities"
	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/domain/errors"
	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/repository/interfaces"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	userRepo    interfaces.UserRepository
	sessionRepo interfaces.SessionRepository
}

func NewAuthService(userRepo interfaces.UserRepository, sessionRepo interfaces.SessionRepository) *AuthService {
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
	existingUser, err := s.userRepo.GetByEmail(email)
	if err != nil {
		return nil, nil, err
	}
	if existingUser != nil {
		return nil, nil, errors.ErrUserAlreadyExists
	}

	// Хэширование пароля
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, nil, err
	}

	// Создание пользователя
	user := &domain.User{
		ID:         uuid.New(),
		Email:      email,
		Password:   string(hashedPassword),
		Name:       "",     // можно генерировать или оставить пустым
		Gender:     "male", // ← пустая строка вместо NULL
		IsVerified: false,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, nil, err
	}

	// Создание сессии
	session := domain.NewSession(user.Email, 3600*time.Second)
	session.Token = uuid.New().String()

	if err := s.sessionRepo.Set(session); err != nil {
		return nil, nil, err
	}

	return user, session, nil
}

func (s *AuthService) Login(email, password string) (*domain.User, *domain.Session, error) {
	user, err := s.userRepo.GetByEmail(email)
	if err != nil {
		return nil, nil, err
	}
	if user == nil {
		return nil, nil, errors.ErrInvalidCredentials
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, nil, errors.ErrInvalidCredentials
	}

	// Обновляем last_active
	if err := s.userRepo.UpdateLastActive(user.ID); err != nil {
		return nil, nil, err
	}

	// Создаем новую сессию
	session := domain.NewSession(user.Email, 3600*time.Second)
	session.Token = uuid.New().String()

	if err := s.sessionRepo.Set(session); err != nil {
		return nil, nil, err
	}

	return user, session, nil
}

func (s *AuthService) Logout(token string) error {
	session, err := s.sessionRepo.Get(token)
	if err != nil {
		return err
	}
	if session == nil {
		return errors.ErrSessionNotFound
	}

	return s.sessionRepo.Delete(session.Token)
}

func (s *AuthService) ValidateSession(token string) (*domain.User, error) {
	session, err := s.sessionRepo.Get(token)
	if err != nil {
		return nil, err
	}
	if session == nil {
		return nil, errors.ErrSessionNotFound
	}

	if session.IsExpired() {
		go s.sessionRepo.Delete(session.Token)
		return nil, errors.ErrSessionExpired
	}

	user, err := s.userRepo.GetByEmail(session.UserEmail)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.ErrUserNotFound
	}

	return user, nil
}

func (s *AuthService) GetUserByID(userID uuid.UUID) (*domain.User, error) {
	return s.userRepo.GetByID(userID)
}
