package service

import (
	"context"
	"testing"
	"time"

	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/core-service/domain/constants"
	domain "github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/core-service/domain/entities"
	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/domain/errors"
	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/tests/mocks"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestFeedService_GetFeed_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mocks.NewMockUserRepository(ctrl)
	mockPrefRepo := mocks.NewMockUserPreferenceRepository(ctrl)
	mockPhotoRepo := mocks.NewMockUserPhotoRepository(ctrl)
	mockLogger := mocks.NewMockLogger()

	service := NewFeedService(mockUserRepo, mockPrefRepo, mockPhotoRepo, mockLogger)

	ctx := context.Background()
	userID := uuid.New()
	feedUserID := uuid.New()

	bio := "Test bio"
	birthDate := time.Date(time.Now().Year()-25, time.Now().Month(), time.Now().Day()-1, 0, 0, 0, 0, time.UTC)
	users := []domain.User{
		{
			ID:        feedUserID,
			Name:      "Test User",
			Gender:    constants.GenderMale,
			BirthDate: birthDate,
			Bio:       &bio,
		},
	}

	interests := []domain.Interest{
		{
			UserID: feedUserID,
			Theme:  "sports",
		},
	}

	photos := []domain.UserPhoto{
		{
			ID:         uuid.New(),
			UserID:     feedUserID,
			PhotoURL:   "http://example.com/photo.jpg",
			IsApproved: true,
		},
	}

	// Expectations
	mockUserRepo.EXPECT().
		GetUsersForFeed(ctx, userID, 15, 0).
		Return(users, nil)

	mockPrefRepo.EXPECT().
		GetInterests(ctx, feedUserID).
		Return(interests, nil)

	mockPhotoRepo.EXPECT().
		GetByUserID(ctx, feedUserID).
		Return(photos, nil)

	mockUserRepo.EXPECT().
		UpdateLastActive(ctx, userID).
		Return(nil).
		AnyTimes() // может быть вызван асинхронно

	// Execute
	feedUsers, err := service.GetFeed(ctx, userID, 15, 0)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, feedUsers)
	assert.Equal(t, 1, len(feedUsers))
	assert.Equal(t, feedUserID.String(), feedUsers[0].ID)
	assert.Equal(t, "Test User", feedUsers[0].Name)
	assert.Equal(t, 25, feedUsers[0].Age)
	assert.Equal(t, "male", feedUsers[0].Gender)
	assert.Equal(t, 1, len(feedUsers[0].Images))
	assert.Equal(t, 1, len(feedUsers[0].Interests))
}

func TestFeedService_GetFeed_EmptyResult(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mocks.NewMockUserRepository(ctrl)
	mockPrefRepo := mocks.NewMockUserPreferenceRepository(ctrl)
	mockPhotoRepo := mocks.NewMockUserPhotoRepository(ctrl)
	mockLogger := mocks.NewMockLogger()

	service := NewFeedService(mockUserRepo, mockPrefRepo, mockPhotoRepo, mockLogger)

	ctx := context.Background()
	userID := uuid.New()

	// Expectations
	mockUserRepo.EXPECT().
		GetUsersForFeed(ctx, userID, 15, 0).
		Return([]domain.User{}, nil)

	mockUserRepo.EXPECT().
		UpdateLastActive(ctx, userID).
		Return(nil).
		AnyTimes()

	// Execute
	feedUsers, err := service.GetFeed(ctx, userID, 15, 0)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, feedUsers)
	assert.Equal(t, 0, len(feedUsers))
}

func TestFeedService_GetFeed_GetUsersError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mocks.NewMockUserRepository(ctrl)
	mockPrefRepo := mocks.NewMockUserPreferenceRepository(ctrl)
	mockPhotoRepo := mocks.NewMockUserPhotoRepository(ctrl)
	mockLogger := mocks.NewMockLogger()

	service := NewFeedService(mockUserRepo, mockPrefRepo, mockPhotoRepo, mockLogger)

	ctx := context.Background()
	userID := uuid.New()
	expectedErr := errors.ErrInternalError

	// Expectations
	mockUserRepo.EXPECT().
		GetUsersForFeed(ctx, userID, 15, 0).
		Return(nil, expectedErr)

	// Execute
	feedUsers, err := service.GetFeed(ctx, userID, 15, 0)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, expectedErr, err)
	assert.Nil(t, feedUsers)
}

func TestFeedService_GetFeed_InvalidLimit(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mocks.NewMockUserRepository(ctrl)
	mockPrefRepo := mocks.NewMockUserPreferenceRepository(ctrl)
	mockPhotoRepo := mocks.NewMockUserPhotoRepository(ctrl)
	mockLogger := mocks.NewMockLogger()

	service := NewFeedService(mockUserRepo, mockPrefRepo, mockPhotoRepo, mockLogger)

	ctx := context.Background()
	userID := uuid.New()

	// Expectations - limit должен быть нормализован до 15
	mockUserRepo.EXPECT().
		GetUsersForFeed(ctx, userID, 15, 0).
		Return([]domain.User{}, nil)

	mockUserRepo.EXPECT().
		UpdateLastActive(ctx, userID).
		Return(nil).
		AnyTimes()

	// Execute with invalid limit
	feedUsers, err := service.GetFeed(ctx, userID, -1, 0)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, feedUsers)
}

func TestFeedService_GetFeed_InvalidOffset(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mocks.NewMockUserRepository(ctrl)
	mockPrefRepo := mocks.NewMockUserPreferenceRepository(ctrl)
	mockPhotoRepo := mocks.NewMockUserPhotoRepository(ctrl)
	mockLogger := mocks.NewMockLogger()

	service := NewFeedService(mockUserRepo, mockPrefRepo, mockPhotoRepo, mockLogger)

	ctx := context.Background()
	userID := uuid.New()

	// Expectations - offset должен быть нормализован до 0
	mockUserRepo.EXPECT().
		GetUsersForFeed(ctx, userID, 15, 0).
		Return([]domain.User{}, nil)

	mockUserRepo.EXPECT().
		UpdateLastActive(ctx, userID).
		Return(nil).
		AnyTimes()

	// Execute with invalid offset
	feedUsers, err := service.GetFeed(ctx, userID, 15, -10)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, feedUsers)
}

func TestFeedService_GetFeed_GetPhotosError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mocks.NewMockUserRepository(ctrl)
	mockPrefRepo := mocks.NewMockUserPreferenceRepository(ctrl)
	mockPhotoRepo := mocks.NewMockUserPhotoRepository(ctrl)
	mockLogger := mocks.NewMockLogger()

	service := NewFeedService(mockUserRepo, mockPrefRepo, mockPhotoRepo, mockLogger)

	ctx := context.Background()
	userID := uuid.New()
	feedUserID := uuid.New()

	users := []domain.User{
		{
			ID:        feedUserID,
			Name:      "Test User",
			Gender:    constants.GenderMale,
			BirthDate: time.Now().AddDate(-25, 0, 0),
		},
	}

	// Expectations
	mockUserRepo.EXPECT().
		GetUsersForFeed(ctx, userID, 15, 0).
		Return(users, nil)

	mockPrefRepo.EXPECT().
		GetInterests(ctx, feedUserID).
		Return([]domain.Interest{}, nil)

	mockPhotoRepo.EXPECT().
		GetByUserID(ctx, feedUserID).
		Return(nil, errors.ErrInternalError)

	mockUserRepo.EXPECT().
		UpdateLastActive(ctx, userID).
		Return(nil).
		AnyTimes()

	// Execute
	feedUsers, err := service.GetFeed(ctx, userID, 15, 0)

	// Assert
	assert.NoError(t, err) // Сервис должен пропустить пользователя с ошибкой
	assert.NotNil(t, feedUsers)
	assert.Equal(t, 0, len(feedUsers)) // Пользователь пропущен из-за ошибки
}
