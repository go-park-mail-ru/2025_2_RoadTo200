package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/domain/constants"
	domain "github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/domain/entities"
	expectation "github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/domain/errors"
	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/tests/mocks"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

func TestAuthService_Register_Success(t *testing.T) {
	// Setup
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userRepo := mocks.NewMockUserRepository(ctrl)
	sessionRepo := mocks.NewMockSessionRepository(ctrl)
	logger := mocks.NewMockLogger()

	service := NewAuthService(userRepo, sessionRepo, logger)

	ctx := context.Background()
	email := "test@example.com"
	password := "password123"
	passwordConfirm := "password123"

	// Mock expectations
	userRepo.EXPECT().
		GetByEmail(ctx, email).
		Return(nil, nil) // User doesn't exist

	userRepo.EXPECT().
		Create(ctx, gomock.Any()).
		DoAndReturn(func(ctx context.Context, user *domain.User) error {
			// Verify user data
			assert.Equal(t, email, user.Email)
			assert.Equal(t, constants.Gender("male"), user.Gender)
			assert.True(t, user.IsVerified)
			// Verify password is hashed
			err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
			assert.NoError(t, err)
			return nil
		})

	sessionRepo.EXPECT().
		Set(ctx, gomock.Any()).
		DoAndReturn(func(ctx context.Context, session *domain.Session) error {
			assert.Equal(t, email, session.UserEmail)
			assert.NotEmpty(t, session.Token)
			assert.False(t, session.IsExpired())
			return nil
		})

	// Execute
	user, session, err := service.Register(ctx, email, password, passwordConfirm)

	// Assert
	require.NoError(t, err)
	require.NotNil(t, user)
	require.NotNil(t, session)
	assert.Equal(t, email, user.Email)
	assert.Equal(t, email, user.Name)
	assert.Equal(t, constants.Gender("male"), user.Gender)
	assert.NotEmpty(t, session.Token)
	assert.Equal(t, email, session.UserEmail)
}

func TestAuthService_Register_PasswordsDontMatch(t *testing.T) {
	// Setup
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userRepo := mocks.NewMockUserRepository(ctrl)
	sessionRepo := mocks.NewMockSessionRepository(ctrl)
	logger := mocks.NewMockLogger()

	service := NewAuthService(userRepo, sessionRepo, logger)

	ctx := context.Background()
	email := "test@example.com"
	password := "password123"
	passwordConfirm := "differentpassword"

	// Execute
	user, session, err := service.Register(ctx, email, password, passwordConfirm)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, expectation.ErrPasswordsDontMatch, err)
	assert.Nil(t, user)
	assert.Nil(t, session)
}

func TestAuthService_Register_PasswordTooShort(t *testing.T) {
	// Setup
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userRepo := mocks.NewMockUserRepository(ctrl)
	sessionRepo := mocks.NewMockSessionRepository(ctrl)
	logger := mocks.NewMockLogger()

	service := NewAuthService(userRepo, sessionRepo, logger)

	ctx := context.Background()
	email := "test@example.com"
	password := "123"
	passwordConfirm := "123"

	// Execute
	user, session, err := service.Register(ctx, email, password, passwordConfirm)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, expectation.ErrPasswordTooShort, err)
	assert.Nil(t, user)
	assert.Nil(t, session)
}

func TestAuthService_Register_UserAlreadyExists(t *testing.T) {
	// Setup
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userRepo := mocks.NewMockUserRepository(ctrl)
	sessionRepo := mocks.NewMockSessionRepository(ctrl)
	logger := mocks.NewMockLogger()

	service := NewAuthService(userRepo, sessionRepo, logger)

	ctx := context.Background()
	email := "existing@example.com"
	password := "password123"
	passwordConfirm := "password123"

	existingUser := &domain.User{
		ID:    uuid.New(),
		Email: email,
	}

	// Mock expectations
	userRepo.EXPECT().
		GetByEmail(ctx, email).
		Return(existingUser, nil)

	// Execute
	user, session, err := service.Register(ctx, email, password, passwordConfirm)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, expectation.ErrUserAlreadyExists, err)
	assert.Nil(t, user)
	assert.Nil(t, session)
}

func TestAuthService_Register_GetByEmailError(t *testing.T) {
	// Setup
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userRepo := mocks.NewMockUserRepository(ctrl)
	sessionRepo := mocks.NewMockSessionRepository(ctrl)
	logger := mocks.NewMockLogger()

	service := NewAuthService(userRepo, sessionRepo, logger)

	ctx := context.Background()
	email := "test@example.com"
	password := "password123"
	passwordConfirm := "password123"

	expectedErr := errors.New("database error")

	// Mock expectations
	userRepo.EXPECT().
		GetByEmail(ctx, email).
		Return(nil, expectedErr)

	// Execute
	user, session, err := service.Register(ctx, email, password, passwordConfirm)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, expectedErr, err)
	assert.Nil(t, user)
	assert.Nil(t, session)
}

func TestAuthService_Register_CreateUserError(t *testing.T) {
	// Setup
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userRepo := mocks.NewMockUserRepository(ctrl)
	sessionRepo := mocks.NewMockSessionRepository(ctrl)
	logger := mocks.NewMockLogger()

	service := NewAuthService(userRepo, sessionRepo, logger)

	ctx := context.Background()
	email := "test@example.com"
	password := "password123"
	passwordConfirm := "password123"

	expectedErr := errors.New("create user error")

	// Mock expectations
	userRepo.EXPECT().
		GetByEmail(ctx, email).
		Return(nil, nil)

	userRepo.EXPECT().
		Create(ctx, gomock.Any()).
		Return(expectedErr)

	// Execute
	user, session, err := service.Register(ctx, email, password, passwordConfirm)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, expectedErr, err)
	assert.Nil(t, user)
	assert.Nil(t, session)
}

func TestAuthService_Register_SetSessionError(t *testing.T) {
	// Setup
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userRepo := mocks.NewMockUserRepository(ctrl)
	sessionRepo := mocks.NewMockSessionRepository(ctrl)
	logger := mocks.NewMockLogger()

	service := NewAuthService(userRepo, sessionRepo, logger)

	ctx := context.Background()
	email := "test@example.com"
	password := "password123"
	passwordConfirm := "password123"

	expectedErr := errors.New("session error")

	// Mock expectations
	userRepo.EXPECT().
		GetByEmail(ctx, email).
		Return(nil, nil)

	userRepo.EXPECT().
		Create(ctx, gomock.Any()).
		Return(nil)

	sessionRepo.EXPECT().
		Set(ctx, gomock.Any()).
		Return(expectedErr)

	// Execute
	user, session, err := service.Register(ctx, email, password, passwordConfirm)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, expectedErr, err)
	assert.Nil(t, user)
	assert.Nil(t, session)
}

func TestAuthService_Login_Success(t *testing.T) {
	// Setup
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userRepo := mocks.NewMockUserRepository(ctrl)
	sessionRepo := mocks.NewMockSessionRepository(ctrl)
	logger := mocks.NewMockLogger()

	service := NewAuthService(userRepo, sessionRepo, logger)

	ctx := context.Background()
	email := "test@example.com"
	password := "password123"

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	existingUser := &domain.User{
		ID:       uuid.New(),
		Email:    email,
		Password: string(hashedPassword),
		Name:     "Test User",
	}

	// Mock expectations
	userRepo.EXPECT().
		GetByEmail(ctx, email).
		Return(existingUser, nil)

	userRepo.EXPECT().
		UpdateLastActive(ctx, existingUser.ID).
		Return(nil)

	sessionRepo.EXPECT().
		Set(ctx, gomock.Any()).
		DoAndReturn(func(ctx context.Context, session *domain.Session) error {
			assert.Equal(t, email, session.UserEmail)
			assert.NotEmpty(t, session.Token)
			return nil
		})

	// Execute
	user, session, err := service.Login(ctx, email, password)

	// Assert
	require.NoError(t, err)
	require.NotNil(t, user)
	require.NotNil(t, session)
	assert.Equal(t, existingUser.ID, user.ID)
	assert.Equal(t, email, user.Email)
	assert.Equal(t, email, session.UserEmail)
	assert.NotEmpty(t, session.Token)
}

func TestAuthService_Login_UserNotFound(t *testing.T) {
	// Setup
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userRepo := mocks.NewMockUserRepository(ctrl)
	sessionRepo := mocks.NewMockSessionRepository(ctrl)
	logger := mocks.NewMockLogger()

	service := NewAuthService(userRepo, sessionRepo, logger)

	ctx := context.Background()
	email := "nonexistent@example.com"
	password := "password123"

	// Mock expectations
	userRepo.EXPECT().
		GetByEmail(ctx, email).
		Return(nil, nil)

	// Execute
	user, session, err := service.Login(ctx, email, password)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, expectation.ErrInvalidCredentials, err)
	assert.Nil(t, user)
	assert.Nil(t, session)
}

func TestAuthService_Login_GetByEmailError(t *testing.T) {
	// Setup
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userRepo := mocks.NewMockUserRepository(ctrl)
	sessionRepo := mocks.NewMockSessionRepository(ctrl)
	logger := mocks.NewMockLogger()

	service := NewAuthService(userRepo, sessionRepo, logger)

	ctx := context.Background()
	email := "test@example.com"
	password := "password123"

	expectedErr := errors.New("database error")

	// Mock expectations
	userRepo.EXPECT().
		GetByEmail(ctx, email).
		Return(nil, expectedErr)

	// Execute
	user, session, err := service.Login(ctx, email, password)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, expectedErr, err)
	assert.Nil(t, user)
	assert.Nil(t, session)
}

func TestAuthService_Login_InvalidPassword(t *testing.T) {
	// Setup
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userRepo := mocks.NewMockUserRepository(ctrl)
	sessionRepo := mocks.NewMockSessionRepository(ctrl)
	logger := mocks.NewMockLogger()

	service := NewAuthService(userRepo, sessionRepo, logger)

	ctx := context.Background()
	email := "test@example.com"
	password := "wrongpassword"

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("correctpassword"), bcrypt.DefaultCost)
	existingUser := &domain.User{
		ID:       uuid.New(),
		Email:    email,
		Password: string(hashedPassword),
	}

	// Mock expectations
	userRepo.EXPECT().
		GetByEmail(ctx, email).
		Return(existingUser, nil)

	// Execute
	user, session, err := service.Login(ctx, email, password)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, expectation.ErrInvalidCredentials, err)
	assert.Nil(t, user)
	assert.Nil(t, session)
}

func TestAuthService_Login_UpdateLastActiveError(t *testing.T) {
	// Setup
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userRepo := mocks.NewMockUserRepository(ctrl)
	sessionRepo := mocks.NewMockSessionRepository(ctrl)
	logger := mocks.NewMockLogger()

	service := NewAuthService(userRepo, sessionRepo, logger)

	ctx := context.Background()
	email := "test@example.com"
	password := "password123"

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	existingUser := &domain.User{
		ID:       uuid.New(),
		Email:    email,
		Password: string(hashedPassword),
	}

	expectedErr := errors.New("update error")

	// Mock expectations
	userRepo.EXPECT().
		GetByEmail(ctx, email).
		Return(existingUser, nil)

	userRepo.EXPECT().
		UpdateLastActive(ctx, existingUser.ID).
		Return(expectedErr)

	// Execute
	user, session, err := service.Login(ctx, email, password)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, expectedErr, err)
	assert.Nil(t, user)
	assert.Nil(t, session)
}

func TestAuthService_Login_SetSessionError(t *testing.T) {
	// Setup
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userRepo := mocks.NewMockUserRepository(ctrl)
	sessionRepo := mocks.NewMockSessionRepository(ctrl)
	logger := mocks.NewMockLogger()

	service := NewAuthService(userRepo, sessionRepo, logger)

	ctx := context.Background()
	email := "test@example.com"
	password := "password123"

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	existingUser := &domain.User{
		ID:       uuid.New(),
		Email:    email,
		Password: string(hashedPassword),
	}

	expectedErr := errors.New("session error")

	// Mock expectations
	userRepo.EXPECT().
		GetByEmail(ctx, email).
		Return(existingUser, nil)

	userRepo.EXPECT().
		UpdateLastActive(ctx, existingUser.ID).
		Return(nil)

	sessionRepo.EXPECT().
		Set(ctx, gomock.Any()).
		Return(expectedErr)

	// Execute
	user, session, err := service.Login(ctx, email, password)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, expectedErr, err)
	assert.Nil(t, user)
	assert.Nil(t, session)
}

func TestAuthService_Logout_Success(t *testing.T) {
	// Setup
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userRepo := mocks.NewMockUserRepository(ctrl)
	sessionRepo := mocks.NewMockSessionRepository(ctrl)
	logger := mocks.NewMockLogger()

	service := NewAuthService(userRepo, sessionRepo, logger)

	ctx := context.Background()
	token := "test-token"

	session := &domain.Session{
		Token:     token,
		UserEmail: "test@example.com",
		ExpiresAt: time.Now().Add(time.Hour),
	}

	// Mock expectations
	sessionRepo.EXPECT().
		Get(ctx, token).
		Return(session, nil)

	sessionRepo.EXPECT().
		Delete(ctx, token).
		Return(nil)

	// Execute
	err := service.Logout(ctx, token)

	// Assert
	assert.NoError(t, err)
}

func TestAuthService_Logout_SessionNotFound(t *testing.T) {
	// Setup
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userRepo := mocks.NewMockUserRepository(ctrl)
	sessionRepo := mocks.NewMockSessionRepository(ctrl)
	logger := mocks.NewMockLogger()

	service := NewAuthService(userRepo, sessionRepo, logger)

	ctx := context.Background()
	token := "nonexistent-token"

	// Mock expectations
	sessionRepo.EXPECT().
		Get(ctx, token).
		Return(nil, nil)

	// Execute
	err := service.Logout(ctx, token)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, expectation.ErrSessionNotFound, err)
}

func TestAuthService_Logout_GetSessionError(t *testing.T) {
	// Setup
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userRepo := mocks.NewMockUserRepository(ctrl)
	sessionRepo := mocks.NewMockSessionRepository(ctrl)
	logger := mocks.NewMockLogger()

	service := NewAuthService(userRepo, sessionRepo, logger)

	ctx := context.Background()
	token := "test-token"

	expectedErr := errors.New("get session error")

	// Mock expectations
	sessionRepo.EXPECT().
		Get(ctx, token).
		Return(nil, expectedErr)

	// Execute
	err := service.Logout(ctx, token)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, expectedErr, err)
}

func TestAuthService_ValidateSession_Success(t *testing.T) {
	// Setup
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userRepo := mocks.NewMockUserRepository(ctrl)
	sessionRepo := mocks.NewMockSessionRepository(ctrl)
	logger := mocks.NewMockLogger()

	service := NewAuthService(userRepo, sessionRepo, logger)

	ctx := context.Background()
	token := "valid-token"
	email := "test@example.com"

	session := &domain.Session{
		Token:     token,
		UserEmail: email,
		ExpiresAt: time.Now().Add(time.Hour),
	}

	user := &domain.User{
		ID:    uuid.New(),
		Email: email,
		Name:  "Test User",
	}

	// Mock expectations
	sessionRepo.EXPECT().
		Get(ctx, token).
		Return(session, nil)

	userRepo.EXPECT().
		GetByEmail(ctx, email).
		Return(user, nil)

	// Execute
	resultUser, err := service.ValidateSession(ctx, token)

	// Assert
	require.NoError(t, err)
	require.NotNil(t, resultUser)
	assert.Equal(t, user.ID, resultUser.ID)
	assert.Equal(t, user.Email, resultUser.Email)
}

func TestAuthService_ValidateSession_SessionNotFound(t *testing.T) {
	// Setup
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userRepo := mocks.NewMockUserRepository(ctrl)
	sessionRepo := mocks.NewMockSessionRepository(ctrl)
	logger := mocks.NewMockLogger()

	service := NewAuthService(userRepo, sessionRepo, logger)

	ctx := context.Background()
	token := "invalid-token"

	// Mock expectations
	sessionRepo.EXPECT().
		Get(ctx, token).
		Return(nil, nil)

	// Execute
	user, err := service.ValidateSession(ctx, token)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, expectation.ErrSessionNotFound, err)
	assert.Nil(t, user)
}

func TestAuthService_ValidateSession_GetSessionError(t *testing.T) {
	// Setup
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userRepo := mocks.NewMockUserRepository(ctrl)
	sessionRepo := mocks.NewMockSessionRepository(ctrl)
	logger := mocks.NewMockLogger()

	service := NewAuthService(userRepo, sessionRepo, logger)

	ctx := context.Background()
	token := "test-token"

	expectedErr := errors.New("session error")

	// Mock expectations
	sessionRepo.EXPECT().
		Get(ctx, token).
		Return(nil, expectedErr)

	// Execute
	user, err := service.ValidateSession(ctx, token)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, expectedErr, err)
	assert.Nil(t, user)
}

func TestAuthService_ValidateSession_SessionExpired(t *testing.T) {
	// Setup
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userRepo := mocks.NewMockUserRepository(ctrl)
	sessionRepo := mocks.NewMockSessionRepository(ctrl)
	logger := mocks.NewMockLogger()

	service := NewAuthService(userRepo, sessionRepo, logger)

	ctx := context.Background()
	token := "expired-token"
	email := "test@example.com"

	session := &domain.Session{
		Token:     token,
		UserEmail: email,
		ExpiresAt: time.Now().Add(-time.Hour), // Expired
	}

	deleteCalled := make(chan bool, 1)

	// Mock expectations
	sessionRepo.EXPECT().
		Get(ctx, token).
		Return(session, nil)

	sessionRepo.EXPECT().
		Delete(ctx, token).
		DoAndReturn(func(ctx context.Context, token string) error {
			deleteCalled <- true
			return nil
		})

	// Execute
	user, err := service.ValidateSession(ctx, token)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, expectation.ErrSessionExpired, err)
	assert.Nil(t, user)

	// Ждем вызова Delete (с таймаутом)
	select {
	case <-deleteCalled:
		// Успех - Delete был вызван
	case <-time.After(100 * time.Millisecond):
		t.Error("Delete was not called within timeout")
	}
}

func TestAuthService_ValidateSession_UserNotFound(t *testing.T) {
	// Setup
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userRepo := mocks.NewMockUserRepository(ctrl)
	sessionRepo := mocks.NewMockSessionRepository(ctrl)
	logger := mocks.NewMockLogger()

	service := NewAuthService(userRepo, sessionRepo, logger)

	ctx := context.Background()
	token := "valid-token"
	email := "nonexistent@example.com"

	session := &domain.Session{
		Token:     token,
		UserEmail: email,
		ExpiresAt: time.Now().Add(time.Hour),
	}

	// Mock expectations
	sessionRepo.EXPECT().
		Get(ctx, token).
		Return(session, nil)

	userRepo.EXPECT().
		GetByEmail(ctx, email).
		Return(nil, nil)

	// Execute
	user, err := service.ValidateSession(ctx, token)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, expectation.ErrUserNotFound, err)
	assert.Nil(t, user)
}

func TestAuthService_ValidateSession_GetUserError(t *testing.T) {
	// Setup
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userRepo := mocks.NewMockUserRepository(ctrl)
	sessionRepo := mocks.NewMockSessionRepository(ctrl)
	logger := mocks.NewMockLogger()

	service := NewAuthService(userRepo, sessionRepo, logger)

	ctx := context.Background()
	token := "valid-token"
	email := "test@example.com"

	session := &domain.Session{
		Token:     token,
		UserEmail: email,
		ExpiresAt: time.Now().Add(time.Hour),
	}

	expectedErr := errors.New("user error")

	// Mock expectations
	sessionRepo.EXPECT().
		Get(ctx, token).
		Return(session, nil)

	userRepo.EXPECT().
		GetByEmail(ctx, email).
		Return(nil, expectedErr)

	// Execute
	user, err := service.ValidateSession(ctx, token)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, expectedErr, err)
	assert.Nil(t, user)
}

func TestAuthService_GetUserByID_Success(t *testing.T) {
	// Setup
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userRepo := mocks.NewMockUserRepository(ctrl)
	sessionRepo := mocks.NewMockSessionRepository(ctrl)
	logger := mocks.NewMockLogger()

	service := NewAuthService(userRepo, sessionRepo, logger)

	ctx := context.Background()
	userID := uuid.New()

	expectedUser := &domain.User{
		ID:    userID,
		Email: "test@example.com",
		Name:  "Test User",
	}

	// Mock expectations
	userRepo.EXPECT().
		GetByID(ctx, userID).
		Return(expectedUser, nil)

	// Execute
	user, err := service.GetUserByID(ctx, userID)

	// Assert
	require.NoError(t, err)
	require.NotNil(t, user)
	assert.Equal(t, expectedUser.ID, user.ID)
	assert.Equal(t, expectedUser.Email, user.Email)
}

func TestAuthService_GetUserByID_Error(t *testing.T) {
	// Setup
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userRepo := mocks.NewMockUserRepository(ctrl)
	sessionRepo := mocks.NewMockSessionRepository(ctrl)
	logger := mocks.NewMockLogger()

	service := NewAuthService(userRepo, sessionRepo, logger)

	ctx := context.Background()
	userID := uuid.New()

	expectedErr := errors.New("user not found")

	// Mock expectations
	userRepo.EXPECT().
		GetByID(ctx, userID).
		Return(nil, expectedErr)

	// Execute
	user, err := service.GetUserByID(ctx, userID)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, expectedErr, err)
	assert.Nil(t, user)
}
