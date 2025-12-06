package service

import (
	"context"
	"testing"
	"time"

	domain "github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/chat-service/domain/entities"
	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/domain/errors"
	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/tests/mocks"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestMatchService_GetUserMatches_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockMatchRepo := mocks.NewMockMatchRepository(ctrl)
	mockUserRepo := mocks.NewMockUserRepository(ctrl)
	mockSwipeRepo := mocks.NewMockSwipeRepository(ctrl)
	mockPhotoRepo := mocks.NewMockUserPhotoRepository(ctrl)
	mockLogger := mocks.NewMockLogger()

	service := NewMatchService(mockMatchRepo, mockUserRepo, mockSwipeRepo, mockPhotoRepo, mockLogger)

	ctx := context.Background()
	userID := uuid.New()
	otherUserID := uuid.New()
	matchID := uuid.New()

	matches := []domain.Match{
		{
			ID:        matchID,
			User1ID:   userID,
			User2ID:   otherUserID,
			IsActive:  true,
			MatchedAt: time.Now(),
		},
	}

	users := []domain.User{
		{
			ID:        otherUserID,
			Name:      "Test User",
			BirthDate: time.Now().AddDate(-25, 0, 0),
		},
	}

	photos := []domain.UserPhoto{
		{
			ID:         uuid.New(),
			UserID:     otherUserID,
			PhotoURL:   "http://example.com/photo.jpg",
			IsApproved: true,
		},
	}

	// Expectations
	mockMatchRepo.EXPECT().
		GetUserMatches(ctx, userID, 20, 0).
		Return(matches, nil)

	mockUserRepo.EXPECT().
		GetUsersByIDs(ctx, []uuid.UUID{otherUserID}).
		Return(users, nil)

	mockPhotoRepo.EXPECT().
		GetByUserID(ctx, otherUserID).
		Return(photos, nil)

	// Execute
	response, err := service.GetUserMatches(ctx, userID, 20, 0)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, 1, len(response.Matches))
	assert.Equal(t, 20, response.Limit)
	assert.Equal(t, 0, response.Offset)
}

func TestMatchService_GetUserMatches_NoMatches(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockMatchRepo := mocks.NewMockMatchRepository(ctrl)
	mockUserRepo := mocks.NewMockUserRepository(ctrl)
	mockSwipeRepo := mocks.NewMockSwipeRepository(ctrl)
	mockPhotoRepo := mocks.NewMockUserPhotoRepository(ctrl)
	mockLogger := mocks.NewMockLogger()

	service := NewMatchService(mockMatchRepo, mockUserRepo, mockSwipeRepo, mockPhotoRepo, mockLogger)

	ctx := context.Background()
	userID := uuid.New()

	// Expectations
	mockMatchRepo.EXPECT().
		GetUserMatches(ctx, userID, 20, 0).
		Return([]domain.Match{}, nil)

	// Execute
	response, err := service.GetUserMatches(ctx, userID, 20, 0)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, 0, len(response.Matches))
	assert.Equal(t, 0, response.Total)
}

func TestMatchService_GetUserMatches_GetMatchesError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockMatchRepo := mocks.NewMockMatchRepository(ctrl)
	mockUserRepo := mocks.NewMockUserRepository(ctrl)
	mockSwipeRepo := mocks.NewMockSwipeRepository(ctrl)
	mockPhotoRepo := mocks.NewMockUserPhotoRepository(ctrl)
	mockLogger := mocks.NewMockLogger()

	service := NewMatchService(mockMatchRepo, mockUserRepo, mockSwipeRepo, mockPhotoRepo, mockLogger)

	ctx := context.Background()
	userID := uuid.New()
	expectedErr := errors.ErrInternalError

	// Expectations
	mockMatchRepo.EXPECT().
		GetUserMatches(ctx, userID, 20, 0).
		Return(nil, expectedErr)

	// Execute
	response, err := service.GetUserMatches(ctx, userID, 20, 0)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, expectedErr, err)
	assert.Nil(t, response)
}

func TestMatchService_Unmatch_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockMatchRepo := mocks.NewMockMatchRepository(ctrl)
	mockUserRepo := mocks.NewMockUserRepository(ctrl)
	mockSwipeRepo := mocks.NewMockSwipeRepository(ctrl)
	mockPhotoRepo := mocks.NewMockUserPhotoRepository(ctrl)
	mockLogger := mocks.NewMockLogger()

	service := NewMatchService(mockMatchRepo, mockUserRepo, mockSwipeRepo, mockPhotoRepo, mockLogger)

	ctx := context.Background()
	userID := uuid.New()
	targetUserID := uuid.New()

	match := &domain.Match{
		ID:        uuid.New(),
		User1ID:   userID,
		User2ID:   targetUserID,
		IsActive:  true,
		MatchedAt: time.Now(),
	}

	// Expectations
	mockMatchRepo.EXPECT().
		GetByUsers(ctx, userID, targetUserID).
		Return(match, nil)

	mockMatchRepo.EXPECT().
		UpdateActive(ctx, userID, targetUserID, false).
		Return(nil)

	// Execute
	err := service.Unmatch(ctx, userID, targetUserID)

	// Assert
	assert.NoError(t, err)
}

func TestMatchService_Unmatch_MatchNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockMatchRepo := mocks.NewMockMatchRepository(ctrl)
	mockUserRepo := mocks.NewMockUserRepository(ctrl)
	mockSwipeRepo := mocks.NewMockSwipeRepository(ctrl)
	mockPhotoRepo := mocks.NewMockUserPhotoRepository(ctrl)
	mockLogger := mocks.NewMockLogger()

	service := NewMatchService(mockMatchRepo, mockUserRepo, mockSwipeRepo, mockPhotoRepo, mockLogger)

	ctx := context.Background()
	userID := uuid.New()
	targetUserID := uuid.New()

	// Expectations
	mockMatchRepo.EXPECT().
		GetByUsers(ctx, userID, targetUserID).
		Return(nil, nil)

	// Execute
	err := service.Unmatch(ctx, userID, targetUserID)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, errors.ErrMatchNotFound, err)
}

func TestMatchService_Unmatch_GetByUsersError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockMatchRepo := mocks.NewMockMatchRepository(ctrl)
	mockUserRepo := mocks.NewMockUserRepository(ctrl)
	mockSwipeRepo := mocks.NewMockSwipeRepository(ctrl)
	mockPhotoRepo := mocks.NewMockUserPhotoRepository(ctrl)
	mockLogger := mocks.NewMockLogger()

	service := NewMatchService(mockMatchRepo, mockUserRepo, mockSwipeRepo, mockPhotoRepo, mockLogger)

	ctx := context.Background()
	userID := uuid.New()
	targetUserID := uuid.New()
	expectedErr := errors.ErrInternalError

	// Expectations
	mockMatchRepo.EXPECT().
		GetByUsers(ctx, userID, targetUserID).
		Return(nil, expectedErr)

	// Execute
	err := service.Unmatch(ctx, userID, targetUserID)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, expectedErr, err)
}
