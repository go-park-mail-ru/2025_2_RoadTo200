package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/domain/constants"
	domain "github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/domain/entities"
	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/handler/dto"
	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/tests/mocks"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFeedService_GetFeed_Success(t *testing.T) {
	// Setup
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userRepo := mocks.NewMockUserRepository(ctrl)
	prefRepo := mocks.NewMockUserPreferenceRepository(ctrl)
	photoRepo := mocks.NewMockUserPhotoRepository(ctrl)
	logger := mocks.NewMockLogger()

	service := NewFeedService(userRepo, prefRepo, photoRepo, logger)

	ctx := context.Background()
	userID := uuid.New()
	limit := 10
	offset := 0

	// Test data - используем фиксированную дату для предсказуемого расчета возраста
	now := time.Now()
	testUsers := []domain.User{
		{
			ID:        uuid.New(),
			Name:      "John Doe",
			BirthDate: time.Date(now.Year()-25, now.Month(), now.Day(), 0, 0, 0, 0, now.Location()),
			Gender:    constants.GenderMale,
			Bio:       stringPtr("Software developer"),
		},
		{
			ID:        uuid.New(),
			Name:      "Jane Smith",
			BirthDate: time.Date(now.Year()-30, now.Month(), now.Day(), 0, 0, 0, 0, now.Location()),
			Gender:    constants.GenderFemale,
			Bio:       stringPtr("Designer"),
		},
	}

	interests := []domain.Interest{
		{Theme: "music"},
		{Theme: "sports"},
		{Theme: "travel"},
	}

	photos := []domain.UserPhoto{
		{PhotoURL: "http://example.com/photo1.jpg", IsApproved: true},
		{PhotoURL: "http://example.com/photo2.jpg", IsApproved: true},
	}

	// Mock expectations
	userRepo.EXPECT().
		GetUsersForFeed(ctx, userID, limit, offset).
		Return(testUsers, nil)

	// Expectations for first user
	prefRepo.EXPECT().
		GetInterests(ctx, testUsers[0].ID).
		Return(interests, nil)
	photoRepo.EXPECT().
		GetByUserID(ctx, testUsers[0].ID).
		Return(photos, nil)

	// Expectations for second user
	prefRepo.EXPECT().
		GetInterests(ctx, testUsers[1].ID).
		Return(interests, nil)
	photoRepo.EXPECT().
		GetByUserID(ctx, testUsers[1].ID).
		Return(photos, nil)

	// UpdateLastActive вызывается в горутине, поэтому используем AnyTimes
	userRepo.EXPECT().
		UpdateLastActive(gomock.Any(), gomock.Any()).
		Return(nil).
		AnyTimes()

	// Execute
	result, err := service.GetFeed(ctx, userID, limit, offset)

	// Assert
	require.NoError(t, err)
	require.Len(t, result, 2)

	// Verify first user
	assert.Equal(t, testUsers[0].ID.String(), result[0].ID)
	assert.Equal(t, testUsers[0].Name, result[0].Name)
	assert.Equal(t, 24, result[0].Age)
	assert.Equal(t, "male", result[0].Gender)
	assert.Equal(t, "Software developer", result[0].Description)
	assert.Equal(t, []string{"http://example.com/photo1.jpg", "http://example.com/photo2.jpg"}, result[0].Images)
	assert.Equal(t, 2, result[0].PhotosCount)
	assert.Len(t, result[0].Interests, 3)

	// Verify second user
	assert.Equal(t, 30, result[1].Age)

	// Verify logger was called
	assert.True(t, logger.WasCalled("Trace"))
	assert.True(t, logger.WasCalled("Debugf"))
}

func TestFeedService_GetFeed_UserRepoError(t *testing.T) {
	// Setup
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userRepo := mocks.NewMockUserRepository(ctrl)
	prefRepo := mocks.NewMockUserPreferenceRepository(ctrl)
	photoRepo := mocks.NewMockUserPhotoRepository(ctrl)
	logger := mocks.NewMockLogger()

	service := NewFeedService(userRepo, prefRepo, photoRepo, logger)

	ctx := context.Background()
	userID := uuid.New()
	expectedErr := errors.New("database error")

	// Mock expectations
	userRepo.EXPECT().
		GetUsersForFeed(ctx, userID, 15, 0).
		Return([]domain.User{}, expectedErr)

	// Execute
	result, err := service.GetFeed(ctx, userID, 0, -1) // Invalid params that get defaulted

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, expectedErr, err)
	assert.True(t, logger.WasCalled("Errorf"))
}

func TestFeedService_GetFeed_InterestsError(t *testing.T) {
	// Setup
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userRepo := mocks.NewMockUserRepository(ctrl)
	prefRepo := mocks.NewMockUserPreferenceRepository(ctrl)
	photoRepo := mocks.NewMockUserPhotoRepository(ctrl)
	logger := mocks.NewMockLogger()

	service := NewFeedService(userRepo, prefRepo, photoRepo, logger)

	ctx := context.Background()
	userID := uuid.New()

	testUser := domain.User{
		ID:        uuid.New(),
		Name:      "John Doe",
		BirthDate: time.Date(time.Now().Year()-25, time.Now().Month(), time.Now().Day(), 0, 0, 0, 0, time.Local),
		Gender:    constants.GenderMale,
		Bio:       stringPtr("Test bio"),
	}

	// Mock expectations
	userRepo.EXPECT().
		GetUsersForFeed(ctx, userID, 15, 0).
		Return([]domain.User{testUser}, nil)

	interestsErr := errors.New("failed to get interests")
	prefRepo.EXPECT().
		GetInterests(ctx, testUser.ID).
		Return([]domain.Interest{}, interestsErr)

	// UpdateLastActive вызывается в горутине, используем AnyTimes
	userRepo.EXPECT().
		UpdateLastActive(gomock.Any(), gomock.Any()).
		Return(nil).
		AnyTimes()

	// Execute
	result, err := service.GetFeed(ctx, userID, 15, 0)

	// Assert
	assert.NoError(t, err)
	assert.Len(t, result, 0) // Пользователь пропущен из-за ошибки интересов
	assert.True(t, logger.WasCalled("Warnf"))
}

func TestFeedService_GetFeed_PhotoRepoError(t *testing.T) {
	// Setup
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userRepo := mocks.NewMockUserRepository(ctrl)
	prefRepo := mocks.NewMockUserPreferenceRepository(ctrl)
	photoRepo := mocks.NewMockUserPhotoRepository(ctrl)
	logger := mocks.NewMockLogger()

	service := NewFeedService(userRepo, prefRepo, photoRepo, logger)

	ctx := context.Background()
	userID := uuid.New()

	testUser := domain.User{
		ID:        uuid.New(),
		Name:      "John Doe",
		BirthDate: time.Date(time.Now().Year()-25, time.Now().Month(), time.Now().Day(), 0, 0, 0, 0, time.Local),
		Gender:    constants.GenderMale,
		Bio:       stringPtr("Test bio"),
	}

	// Mock expectations
	userRepo.EXPECT().
		GetUsersForFeed(ctx, userID, 15, 0).
		Return([]domain.User{testUser}, nil)

	interests := []domain.Interest{{Theme: "music"}}
	prefRepo.EXPECT().
		GetInterests(ctx, testUser.ID).
		Return(interests, nil)

	photoErr := errors.New("failed to get photos")
	photoRepo.EXPECT().
		GetByUserID(ctx, testUser.ID).
		Return([]domain.UserPhoto{}, photoErr)

	// UpdateLastActive вызывается в горутине, используем AnyTimes
	userRepo.EXPECT().
		UpdateLastActive(gomock.Any(), gomock.Any()).
		Return(nil).
		AnyTimes()

	// Execute
	result, err := service.GetFeed(ctx, userID, 15, 0)

	// Assert
	assert.NoError(t, err)
	assert.Empty(t, result) // User should be skipped due to photo error
	assert.True(t, logger.WasCalled("Warnf"))
}

func TestFeedService_GetFeed_OnlyApprovedPhotos(t *testing.T) {
	// Setup
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userRepo := mocks.NewMockUserRepository(ctrl)
	prefRepo := mocks.NewMockUserPreferenceRepository(ctrl)
	photoRepo := mocks.NewMockUserPhotoRepository(ctrl)
	logger := mocks.NewMockLogger()

	service := NewFeedService(userRepo, prefRepo, photoRepo, logger)

	ctx := context.Background()
	userID := uuid.New()

	testUser := domain.User{
		ID:        uuid.New(),
		Name:      "John Doe",
		BirthDate: time.Date(time.Now().Year()-25, time.Now().Month(), time.Now().Day(), 0, 0, 0, 0, time.Local),
		Gender:    constants.GenderMale,
		Bio:       stringPtr("Test bio"),
	}

	mixedPhotos := []domain.UserPhoto{
		{PhotoURL: "approved1.jpg", IsApproved: true},
		{PhotoURL: "not_approved.jpg", IsApproved: false},
		{PhotoURL: "approved2.jpg", IsApproved: true},
		{PhotoURL: "not_approved2.jpg", IsApproved: false},
	}

	// Mock expectations
	userRepo.EXPECT().
		GetUsersForFeed(ctx, userID, 15, 0).
		Return([]domain.User{testUser}, nil)

	prefRepo.EXPECT().
		GetInterests(ctx, testUser.ID).
		Return([]domain.Interest{}, nil)

	photoRepo.EXPECT().
		GetByUserID(ctx, testUser.ID).
		Return(mixedPhotos, nil)

	// UpdateLastActive вызывается в горутине, используем AnyTimes
	userRepo.EXPECT().
		UpdateLastActive(gomock.Any(), gomock.Any()).
		Return(nil).
		AnyTimes()

	// Execute
	result, err := service.GetFeed(ctx, userID, 15, 0)

	// Assert
	assert.NoError(t, err)
	assert.Len(t, result, 1)

	// Only approved photos should be included
	expectedPhotos := []string{"approved1.jpg", "approved2.jpg"}
	assert.Equal(t, expectedPhotos, result[0].Images)
	assert.Equal(t, 2, result[0].PhotosCount)
}

func TestFeedService_GetFeed_ParameterValidation(t *testing.T) {
	testCases := []struct {
		name           string
		limit          int
		offset         int
		expectedLimit  int
		expectedOffset int
	}{
		{"Zero limit", 0, 0, 15, 0},
		{"Negative limit", -5, 0, 15, 0},
		{"Limit too large", 100, 0, 15, 0},
		{"Negative offset", 10, -5, 10, 0},
		{"Valid parameters", 25, 10, 25, 10},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			userRepo := mocks.NewMockUserRepository(ctrl)
			prefRepo := mocks.NewMockUserPreferenceRepository(ctrl)
			photoRepo := mocks.NewMockUserPhotoRepository(ctrl)
			logger := mocks.NewMockLogger()

			service := NewFeedService(userRepo, prefRepo, photoRepo, logger)

			ctx := context.Background()
			userID := uuid.New()

			// Mock expectations
			userRepo.EXPECT().
				GetUsersForFeed(ctx, userID, tc.expectedLimit, tc.expectedOffset).
				Return([]domain.User{}, nil)

			// UpdateLastActive вызывается в горутине, используем AnyTimes
			userRepo.EXPECT().
				UpdateLastActive(gomock.Any(), gomock.Any()).
				Return(nil).
				AnyTimes()

			// Execute
			result, err := service.GetFeed(ctx, userID, tc.limit, tc.offset)

			// Assert
			assert.NoError(t, err)
			assert.Empty(t, result)
		})
	}
}

func TestFeedService_convertToFeedUser_Success(t *testing.T) {
	// Setup
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userRepo := mocks.NewMockUserRepository(ctrl)
	prefRepo := mocks.NewMockUserPreferenceRepository(ctrl)
	photoRepo := mocks.NewMockUserPhotoRepository(ctrl)
	logger := mocks.NewMockLogger()

	service := NewFeedService(userRepo, prefRepo, photoRepo, logger)

	ctx := context.Background()
	user := domain.User{
		ID:        uuid.New(),
		Name:      "Test User",
		BirthDate: time.Date(time.Now().Year()-25, time.Now().Month(), time.Now().Day(), 0, 0, 0, 0, time.Local),
		Gender:    constants.GenderMale,
		Bio:       stringPtr("Test bio"),
	}

	photos := []domain.UserPhoto{
		{PhotoURL: "photo1.jpg", IsApproved: true},
		{PhotoURL: "photo2.jpg", IsApproved: true},
	}

	// Mock expectations
	photoRepo.EXPECT().
		GetByUserID(ctx, user.ID).
		Return(photos, nil)

	// Execute
	result, err := service.convertToFeedUser(ctx, user)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, user.ID.String(), result.ID)
	assert.Equal(t, user.Name, result.Name)
	assert.Equal(t, 24, result.Age)
	assert.Equal(t, "male", result.Gender)
	assert.Equal(t, "Test bio", result.Description)
	assert.Equal(t, []string{"photo1.jpg", "photo2.jpg"}, result.Images)
	assert.Equal(t, 2, result.PhotosCount)
}

func TestFeedService_convertToFeedUser_NoBio(t *testing.T) {
	// Setup
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userRepo := mocks.NewMockUserRepository(ctrl)
	prefRepo := mocks.NewMockUserPreferenceRepository(ctrl)
	photoRepo := mocks.NewMockUserPhotoRepository(ctrl)
	logger := mocks.NewMockLogger()

	service := NewFeedService(userRepo, prefRepo, photoRepo, logger)

	ctx := context.Background()
	user := domain.User{
		ID:        uuid.New(),
		Name:      "Test User",
		BirthDate: time.Date(time.Now().Year()-25, time.Now().Month(), time.Now().Day(), 0, 0, 0, 0, time.Local),
		Gender:    constants.GenderMale,
		Bio:       nil, // No bio
	}

	photos := []domain.UserPhoto{
		{PhotoURL: "photo1.jpg", IsApproved: true},
	}

	// Mock expectations
	photoRepo.EXPECT().
		GetByUserID(ctx, user.ID).
		Return(photos, nil)

	// Execute
	result, err := service.convertToFeedUser(ctx, user)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, "Нет описания", result.Description)
}

func TestFeedService_convertToFeedUser_PhotoError(t *testing.T) {
	// Setup
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userRepo := mocks.NewMockUserRepository(ctrl)
	prefRepo := mocks.NewMockUserPreferenceRepository(ctrl)
	photoRepo := mocks.NewMockUserPhotoRepository(ctrl)
	logger := mocks.NewMockLogger()

	service := NewFeedService(userRepo, prefRepo, photoRepo, logger)

	ctx := context.Background()
	user := domain.User{
		ID:        uuid.New(),
		Name:      "Test User",
		BirthDate: time.Date(time.Now().Year()-25, time.Now().Month(), time.Now().Day(), 0, 0, 0, 0, time.Local),
		Gender:    constants.GenderMale,
		Bio:       stringPtr("Test bio"),
	}

	photoErr := errors.New("photo repository error")

	// Mock expectations
	photoRepo.EXPECT().
		GetByUserID(ctx, user.ID).
		Return([]domain.UserPhoto{}, photoErr)

	// Execute
	result, err := service.convertToFeedUser(ctx, user)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, dto.FeedUser{}, result)
}

// Helper functions tests
func TestGetGenderString(t *testing.T) {
	testCases := []struct {
		name     string
		gender   constants.Gender
		expected string
	}{
		{"Male", constants.GenderMale, "male"},
		{"Female", constants.GenderFemale, "female"},
		{"Empty", "", "not_specified"},
		{"Other", "other", "other"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := getGenderString(tc.gender)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestGetDescription(t *testing.T) {
	testCases := []struct {
		name     string
		bio      *string
		expected string
	}{
		{"With bio", stringPtr("Hello world"), "Hello world"},
		{"Nil bio", nil, "Нет описания"},
		{"Empty bio", stringPtr(""), "Нет описания"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := getDescription(tc.bio)
			assert.Equal(t, tc.expected, result)
		})
	}
}

// Helper function
func stringPtr(s string) *string {
	return &s
}
