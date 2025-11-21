package service

import (
	"testing"
	"time"

	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/domain/entities"
	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/domain/errors"
	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/tests/mocks"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"golang.org/x/crypto/bcrypt"
)

func TestAuthService_Register_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mocks.NewMockUserRepository(ctrl)
	mockSessionRepo := mocks.NewMockSessionRepository(ctrl)

	service := NewAuthService(mockUserRepo, mockSessionRepo, nil)

	email := "test@example.com"
	password := "123456"

	// Настройка ожиданий
	mockUserRepo.EXPECT().
		GetByEmail(email).
		Return(nil, nil) // Пользователь не существует

	mockUserRepo.EXPECT().
		Create(gomock.Any()).
		DoAndReturn(func(user *domain.User) error {
			// Проверяем что пароль захэширован
			err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
			assert.NoError(t, err)
			assert.Equal(t, email, user.Email)
			return nil
		})

	mockSessionRepo.EXPECT().
		Set(gomock.Any()).
		Return(nil)

	// Вызов метода
	user, session, err := service.Register(email, password, password)

	// Проверки
	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.NotNil(t, session)
	assert.Equal(t, email, user.Email)
	assert.NotEqual(t, password, user.Password) // Пароль должен быть захэширован
}

func TestAuthService_Register_UserAlreadyExists(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mocks.NewMockUserRepository(ctrl)
	mockSessionRepo := mocks.NewMockSessionRepository(ctrl)

	service := NewAuthService(mockUserRepo, mockSessionRepo, nil)

	email := "existing@example.com"
	existingUser := &domain.User{Email: email}

	// Настройка ожиданий
	mockUserRepo.EXPECT().
		GetByEmail(email).
		Return(existingUser, nil) // Пользователь уже существует

	// Вызов метода
	user, session, err := service.Register(email, "123456", "123456")

	// Проверки
	assert.Error(t, err)
	assert.Equal(t, errors.ErrUserAlreadyExists, err)
	assert.Nil(t, user)
	assert.Nil(t, session)
}

func TestAuthService_Register_PasswordMismatch(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mocks.NewMockUserRepository(ctrl)
	mockSessionRepo := mocks.NewMockSessionRepository(ctrl)

	service := NewAuthService(mockUserRepo, mockSessionRepo, nil)

	// Вызов метода с разными паролями
	user, session, err := service.Register("test@example.com", "123456", "different")

	// Проверки
	assert.Error(t, err)
	assert.Equal(t, errors.ErrPasswordsDontMatch, err)
	assert.Nil(t, user)
	assert.Nil(t, session)
}

func TestAuthService_Register_ShortPassword(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mocks.NewMockUserRepository(ctrl)
	mockSessionRepo := mocks.NewMockSessionRepository(ctrl)

	service := NewAuthService(mockUserRepo, mockSessionRepo, nil)

	// Вызов метода с коротким паролем
	user, session, err := service.Register("test@example.com", "123", "123")

	// Проверки
	assert.Error(t, err)
	assert.Equal(t, errors.ErrPasswordTooShort, err)
	assert.Nil(t, user)
	assert.Nil(t, session)
}

func TestAuthService_Login_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mocks.NewMockUserRepository(ctrl)
	mockSessionRepo := mocks.NewMockSessionRepository(ctrl)

	service := NewAuthService(mockUserRepo, mockSessionRepo, nil)

	email := "test@example.com"
	password := "123456"

	// Хэшируем пароль как это делает bcrypt
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	existingUser := &domain.User{
		ID:       uuid.New(),
		Email:    email,
		Password: string(hashedPassword),
	}

	// Настройка ожиданий
	mockUserRepo.EXPECT().
		GetByEmail(email).
		Return(existingUser, nil)

	mockUserRepo.EXPECT().
		UpdateLastActive(existingUser.ID).
		Return(nil)

	mockSessionRepo.EXPECT().
		Set(gomock.Any()).
		Return(nil)

	// Вызов метода
	user, session, err := service.Login(email, password)

	// Проверки
	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.NotNil(t, session)
	assert.Equal(t, email, user.Email)
}

func TestAuthService_Login_InvalidCredentials(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mocks.NewMockUserRepository(ctrl)
	mockSessionRepo := mocks.NewMockSessionRepository(ctrl)

	service := NewAuthService(mockUserRepo, mockSessionRepo, nil)

	email := "test@example.com"

	// Настройка ожиданий - пользователь не найден
	mockUserRepo.EXPECT().
		GetByEmail(email).
		Return(nil, nil)

	// Вызов метода
	user, session, err := service.Login(email, "wrongpassword")

	// Проверки
	assert.Error(t, err)
	assert.Equal(t, errors.ErrInvalidCredentials, err)
	assert.Nil(t, user)
	assert.Nil(t, session)
}

func TestAuthService_Login_WrongPassword(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mocks.NewMockUserRepository(ctrl)
	mockSessionRepo := mocks.NewMockSessionRepository(ctrl)

	service := NewAuthService(mockUserRepo, mockSessionRepo, nil)

	email := "test@example.com"
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("correctpassword"), bcrypt.DefaultCost)
	existingUser := &domain.User{
		ID:       uuid.New(),
		Email:    email,
		Password: string(hashedPassword),
	}

	// Настройка ожиданий
	mockUserRepo.EXPECT().
		GetByEmail(email).
		Return(existingUser, nil)

	// Вызов метода с неправильным паролем
	user, session, err := service.Login(email, "wrongpassword")

	// Проверки
	assert.Error(t, err)
	assert.Equal(t, errors.ErrInvalidCredentials, err)
	assert.Nil(t, user)
	assert.Nil(t, session)
}

func TestAuthService_Logout_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mocks.NewMockUserRepository(ctrl)
	mockSessionRepo := mocks.NewMockSessionRepository(ctrl)

	service := NewAuthService(mockUserRepo, mockSessionRepo, nil)

	token := "valid-session-token"
	session := &domain.Session{
		Token: token,
	}

	// Настройка ожиданий
	mockSessionRepo.EXPECT().
		Get(token).
		Return(session, nil)

	mockSessionRepo.EXPECT().
		Delete(token).
		Return(nil)

	// Вызов метода
	err := service.Logout(token)

	// Проверки
	assert.NoError(t, err)
}

func TestAuthService_Logout_SessionNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mocks.NewMockUserRepository(ctrl)
	mockSessionRepo := mocks.NewMockSessionRepository(ctrl)

	service := NewAuthService(mockUserRepo, mockSessionRepo, nil)

	token := "non-existent-token"

	// Настройка ожиданий
	mockSessionRepo.EXPECT().
		Get(token).
		Return(nil, nil)

	// Вызов метода
	err := service.Logout(token)

	// Проверки
	assert.Error(t, err)
	assert.Equal(t, errors.ErrSessionNotFound, err)
}

func TestAuthService_Logout_DeleteError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mocks.NewMockUserRepository(ctrl)
	mockSessionRepo := mocks.NewMockSessionRepository(ctrl)

	service := NewAuthService(mockUserRepo, mockSessionRepo, nil)

	token := "valid-token"
	session := &domain.Session{
		Token: token,
	}
	expectedErr := errors.ErrInternalError

	// Настройка ожиданий
	mockSessionRepo.EXPECT().
		Get(token).
		Return(session, nil)

	mockSessionRepo.EXPECT().
		Delete(token).
		Return(expectedErr)

	// Вызов метода
	err := service.Logout(token)

	// Проверки
	assert.Error(t, err)
	assert.Equal(t, expectedErr, err)
}

func TestAuthService_ValidateSession_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mocks.NewMockUserRepository(ctrl)
	mockSessionRepo := mocks.NewMockSessionRepository(ctrl)

	service := NewAuthService(mockUserRepo, mockSessionRepo, nil)

	token := "valid-token"
	userEmail := "test@example.com"
	userID := uuid.New()

	session := &domain.Session{
		Token:     token,
		UserEmail: userEmail,
		ExpiresAt: time.Now().Add(1 * time.Hour), // сессия не истекла
	}

	user := &domain.User{
		ID:    userID,
		Email: userEmail,
	}

	// Настройка ожиданий
	mockSessionRepo.EXPECT().
		Get(token).
		Return(session, nil)

	mockUserRepo.EXPECT().
		GetByEmail(userEmail).
		Return(user, nil)

	// Вызов метода
	resultUser, err := service.ValidateSession(token)

	// Проверки
	assert.NoError(t, err)
	assert.NotNil(t, resultUser)
	assert.Equal(t, userID, resultUser.ID)
	assert.Equal(t, userEmail, resultUser.Email)
}

func TestAuthService_ValidateSession_SessionNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mocks.NewMockUserRepository(ctrl)
	mockSessionRepo := mocks.NewMockSessionRepository(ctrl)

	service := NewAuthService(mockUserRepo, mockSessionRepo, nil)

	token := "non-existent-token"

	// Настройка ожиданий
	mockSessionRepo.EXPECT().
		Get(token).
		Return(nil, nil)

	// Вызов метода
	user, err := service.ValidateSession(token)

	// Проверки
	assert.Error(t, err)
	assert.Equal(t, errors.ErrSessionNotFound, err)
	assert.Nil(t, user)
}

// func TestAuthService_ValidateSession_Expired(t *testing.T) {
// 	ctrl := gomock.NewController(t)
// 	defer ctrl.Finish()

// 	mockUserRepo := mocks.NewMockUserRepository(ctrl)
// 	mockSessionRepo := mocks.NewMockSessionRepository(ctrl)

// 	service := service.NewAuthService(mockUserRepo, mockSessionRepo, nil)

// 	token := "expired-token"

// 	session := &domain.Session{
// 		Token:     token,
// 		UserEmail: "test@example.com",
// 		ExpiresAt: time.Now().Add(-1 * time.Hour), // сессия истекла
// 	}

// 	// Настройка ожиданий
// 	mockSessionRepo.EXPECT().
// 		Get(token).
// 		Return(session, nil)

// 	mockSessionRepo.EXPECT().
// 		Delete(token).
// 		Return(nil)

// 	// Вызов метода
// 	user, err := service.ValidateSession(token)

// 	// Проверки
// 	assert.Error(t, err)
// 	assert.Equal(t, errors.ErrSessionExpired, err)
// 	assert.Nil(t, user)
// }

func TestAuthService_ValidateSession_UserNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mocks.NewMockUserRepository(ctrl)
	mockSessionRepo := mocks.NewMockSessionRepository(ctrl)

	service := NewAuthService(mockUserRepo, mockSessionRepo, nil)

	token := "valid-token"
	userEmail := "deleted@example.com"

	session := &domain.Session{
		Token:     token,
		UserEmail: userEmail,
		ExpiresAt: time.Now().Add(1 * time.Hour),
	}

	// Настройка ожиданий
	mockSessionRepo.EXPECT().
		Get(token).
		Return(session, nil)

	mockUserRepo.EXPECT().
		GetByEmail(userEmail).
		Return(nil, nil)

	// Вызов метода
	user, err := service.ValidateSession(token)

	// Проверки
	assert.Error(t, err)
	assert.Equal(t, errors.ErrUserNotFound, err)
	assert.Nil(t, user)
}
