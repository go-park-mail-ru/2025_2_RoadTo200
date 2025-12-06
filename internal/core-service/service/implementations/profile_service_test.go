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

func TestProfileService_GetProfile_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mocks.NewMockUserRepository(ctrl)
	mockPhotoRepo := mocks.NewMockUserPhotoRepository(ctrl)
	mockPrefRepo := mocks.NewMockUserPreferenceRepository(ctrl)
	mockFileStorage := mocks.NewMockFileStorage(ctrl)
	mockLogger := mocks.NewMockLogger()

	service := NewProfileService(mockUserRepo, mockPhotoRepo, mockPrefRepo, mockFileStorage, mockLogger)

	ctx := context.Background()
	userID := uuid.New()

	user := &domain.User{
		ID:        userID,
		Name:      "Test User",
		Email:     "test@example.com",
		Gender:    constants.GenderMale,
		BirthDate: time.Now().AddDate(-25, 0, 0),
	}

	preferences := &domain.UserPreference{
		UserID:      userID,
		AgeMin:      18,
		AgeMax:      35,
		MaxDistance: 50,
	}

	photos := []domain.UserPhoto{
		{
			ID:           uuid.New(),
			UserID:       userID,
			PhotoURL:     "http://example.com/photo.jpg",
			DisplayOrder: 0,
			IsApproved:   true,
		},
	}

	interests := []domain.Interest{
		{
			UserID: userID,
			Theme:  constants.InterestTypeWorkout,
		},
	}

	// Expectations
	mockUserRepo.EXPECT().
		GetByID(ctx, userID).
		Return(user, nil)

	mockPrefRepo.EXPECT().
		GetByUserID(ctx, userID).
		Return(preferences, nil)

	mockPhotoRepo.EXPECT().
		GetByUserID(ctx, userID).
		Return(photos, nil)

	mockPrefRepo.EXPECT().
		GetInterests(ctx, userID).
		Return(interests, nil)

	mockUserRepo.EXPECT().
		UpdateLastActive(ctx, userID).
		Return(nil).
		AnyTimes() // async call

	// Execute
	profile, err := service.GetProfile(ctx, userID)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, profile)
	assert.Equal(t, user, profile.User)
	assert.Equal(t, preferences, profile.Preferences)
	assert.Equal(t, 1, len(profile.Photos))
	assert.Equal(t, 1, len(profile.Interests))
}

func TestProfileService_GetProfile_UserNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mocks.NewMockUserRepository(ctrl)
	mockPhotoRepo := mocks.NewMockUserPhotoRepository(ctrl)
	mockPrefRepo := mocks.NewMockUserPreferenceRepository(ctrl)
	mockFileStorage := mocks.NewMockFileStorage(ctrl)
	mockLogger := mocks.NewMockLogger()

	service := NewProfileService(mockUserRepo, mockPhotoRepo, mockPrefRepo, mockFileStorage, mockLogger)

	ctx := context.Background()
	userID := uuid.New()

	// Expectations
	mockUserRepo.EXPECT().
		GetByID(ctx, userID).
		Return(nil, errors.ErrUserNotFound)

	// Execute
	profile, err := service.GetProfile(ctx, userID)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, errors.ErrUserNotFound, err)
	assert.Nil(t, profile)
}

func TestProfileService_UpdateProfileInfo_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mocks.NewMockUserRepository(ctrl)
	mockPhotoRepo := mocks.NewMockUserPhotoRepository(ctrl)
	mockPrefRepo := mocks.NewMockUserPreferenceRepository(ctrl)
	mockFileStorage := mocks.NewMockFileStorage(ctrl)
	mockLogger := mocks.NewMockLogger()

	service := NewProfileService(mockUserRepo, mockPhotoRepo, mockPrefRepo, mockFileStorage, mockLogger)

	ctx := context.Background()
	userID := uuid.New()

	user := &domain.User{
		ID:   userID,
		Name: "Test User",
	}

	updateData := &domain.ProfileUpdateRequest{
		Name: "Updated Name",
	}

	// Expectations
	mockUserRepo.EXPECT().
		GetByID(ctx, userID).
		Return(user, nil)

	mockUserRepo.EXPECT().
		Update(ctx, gomock.Any()).
		DoAndReturn(func(ctx context.Context, u *domain.User) error {
			assert.Equal(t, "Updated Name", u.Name)
			return nil
		})

	// Execute
	err := service.UpdateProfileInfo(ctx, userID, updateData)

	// Assert
	assert.NoError(t, err)
}

func TestProfileService_UpdateInterests_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mocks.NewMockUserRepository(ctrl)
	mockPhotoRepo := mocks.NewMockUserPhotoRepository(ctrl)
	mockPrefRepo := mocks.NewMockUserPreferenceRepository(ctrl)
	mockFileStorage := mocks.NewMockFileStorage(ctrl)
	mockLogger := mocks.NewMockLogger()

	service := NewProfileService(mockUserRepo, mockPhotoRepo, mockPrefRepo, mockFileStorage, mockLogger)

	ctx := context.Background()
	userID := uuid.New()

	// Use valid interest types from constants
	interests := []domain.Interest{
		{
			UserID: userID,
			Theme:  constants.InterestTypeWorkout,
		},
		{
			UserID: userID,
			Theme:  constants.InterestTypeFun,
		},
	}

	// Expectations
	mockPrefRepo.EXPECT().
		UpdateInterests(ctx, userID, interests).
		Return(nil)

	// Execute
	err := service.UpdateInterests(ctx, userID, interests)

	// Assert
	assert.NoError(t, err)
}

func TestProfileService_DeletePhoto_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mocks.NewMockUserRepository(ctrl)
	mockPhotoRepo := mocks.NewMockUserPhotoRepository(ctrl)
	mockPrefRepo := mocks.NewMockUserPreferenceRepository(ctrl)
	mockFileStorage := mocks.NewMockFileStorage(ctrl)
	mockLogger := mocks.NewMockLogger()

	service := NewProfileService(mockUserRepo, mockPhotoRepo, mockPrefRepo, mockFileStorage, mockLogger)

	ctx := context.Background()
	userID := uuid.New()
	photoID := uuid.New()

	photo := &domain.UserPhoto{
		ID:       photoID,
		UserID:   userID,
		PhotoURL: "http://example.com/photo.jpg",
	}

	// Expectations
	mockPhotoRepo.EXPECT().
		GetByID(ctx, photoID).
		Return(photo, nil)

	mockFileStorage.EXPECT().
		DeleteByURL(ctx, photo.PhotoURL).
		Return(nil)

	mockPhotoRepo.EXPECT().
		Delete(ctx, photoID).
		Return(nil)

	mockPhotoRepo.EXPECT().
		GetByUserID(ctx, userID).
		Return([]domain.UserPhoto{}, nil)

	mockPhotoRepo.EXPECT().
		UpdateDisplayOrder(ctx, userID, gomock.Any()).
		Return(nil)

	// Execute
	err := service.DeletePhoto(ctx, userID, photoID)

	// Assert
	assert.NoError(t, err)
}

func TestProfileService_DeletePhoto_PhotoNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mocks.NewMockUserRepository(ctrl)
	mockPhotoRepo := mocks.NewMockUserPhotoRepository(ctrl)
	mockPrefRepo := mocks.NewMockUserPreferenceRepository(ctrl)
	mockFileStorage := mocks.NewMockFileStorage(ctrl)
	mockLogger := mocks.NewMockLogger()

	service := NewProfileService(mockUserRepo, mockPhotoRepo, mockPrefRepo, mockFileStorage, mockLogger)

	ctx := context.Background()
	userID := uuid.New()
	photoID := uuid.New()

	// Expectations
	mockPhotoRepo.EXPECT().
		GetByID(ctx, photoID).
		Return(nil, errors.ErrPhotoNotFound)

	// Execute
	err := service.DeletePhoto(ctx, userID, photoID)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, errors.ErrPhotoNotFound, err)
}

func TestProfileService_UpdatePreferences_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mocks.NewMockUserRepository(ctrl)
	mockPhotoRepo := mocks.NewMockUserPhotoRepository(ctrl)
	mockPrefRepo := mocks.NewMockUserPreferenceRepository(ctrl)
	mockFileStorage := mocks.NewMockFileStorage(ctrl)
	mockLogger := mocks.NewMockLogger()

	service := NewProfileService(mockUserRepo, mockPhotoRepo, mockPrefRepo, mockFileStorage, mockLogger)

	ctx := context.Background()
	userID := uuid.New()

	updateData := &domain.PreferencesUpdateRequest{
		AgeMin:      18,
		AgeMax:      35,
		MaxDistance: 50,
	}

	// Expectations
	mockPrefRepo.EXPECT().
		GetByUserID(ctx, userID).
		Return(nil, nil) // Return nil to trigger Create instead of Update

	mockPrefRepo.EXPECT().
		Create(ctx, gomock.Any()).
		Return(nil)

	// Execute
	err := service.UpdatePreferences(ctx, userID, updateData)

	// Assert
	assert.NoError(t, err)
}
