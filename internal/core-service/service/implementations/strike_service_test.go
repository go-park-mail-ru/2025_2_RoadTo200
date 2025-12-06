package service

import (
	"context"
	"testing"
	"time"

	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/core-service/domain/constants"
	domain "github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/core-service/domain/entities"
	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/domain/errors"
	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/handler/dto"
	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/tests/mocks"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestStrikeService_CreateStrike_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStrikeRepo := mocks.NewMockStrikeRepository(ctrl)
	mockUserRepo := mocks.NewMockUserRepository(ctrl)
	mockLogger := mocks.NewMockLogger()

	service := NewStrikeService(mockStrikeRepo, mockUserRepo, mockLogger)

	ctx := context.Background()
	reporterID := uuid.New()
	targetUserID := uuid.New()

	strikeData := &dto.StrikeCreateRequest{
		ReporterID:   reporterID,
		TargetUserID: targetUserID,
		Type:         constants.StrikeTypeSpam,
		Reason:       "Spam messages",
	}

	reporter := &domain.User{ID: reporterID, Name: "Reporter"}
	target := &domain.User{ID: targetUserID, Name: "Target"}

	// Expectations
	mockUserRepo.EXPECT().
		GetByID(ctx, reporterID).
		Return(reporter, nil)

	mockUserRepo.EXPECT().
		GetByID(ctx, targetUserID).
		Return(target, nil)

	mockStrikeRepo.EXPECT().
		HasActiveStrikeFromUser(ctx, reporterID.String(), targetUserID.String()).
		Return(false, nil)

	mockStrikeRepo.EXPECT().
		CreateStrike(ctx, gomock.Any()).
		DoAndReturn(func(ctx context.Context, strike *domain.Strike) error {
			assert.Equal(t, reporterID, strike.ReporterID)
			assert.Equal(t, targetUserID, strike.TargetUserID)
			assert.Equal(t, constants.StrikeTypeSpam, strike.Type)
			assert.Equal(t, "Spam messages", strike.Reason)
			assert.Equal(t, constants.StrikeStatusPending, strike.Status)
			return nil
		})

	// Execute
	strike, err := service.CreateStrike(ctx, strikeData)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, strike)
}

func TestStrikeService_CreateStrike_SelfStrike(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStrikeRepo := mocks.NewMockStrikeRepository(ctrl)
	mockUserRepo := mocks.NewMockUserRepository(ctrl)
	mockLogger := mocks.NewMockLogger()

	service := NewStrikeService(mockStrikeRepo, mockUserRepo, mockLogger)

	ctx := context.Background()
	userID := uuid.New()

	strikeData := &dto.StrikeCreateRequest{
		ReporterID:   userID,
		TargetUserID: userID, // Same user!
		Type:         constants.StrikeTypeSpam,
		Reason:       "Test",
	}

	user := &domain.User{ID: userID, Name: "User"}

	// Expectations
	mockUserRepo.EXPECT().
		GetByID(ctx, userID).
		Return(user, nil).
		Times(2) // Called twice for reporter and target

	// Execute
	strike, err := service.CreateStrike(ctx, strikeData)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, errors.ErrSelfStrikeNotAllowed, err)
	assert.Nil(t, strike)
}

func TestStrikeService_CreateStrike_DuplicateStrike(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStrikeRepo := mocks.NewMockStrikeRepository(ctrl)
	mockUserRepo := mocks.NewMockUserRepository(ctrl)
	mockLogger := mocks.NewMockLogger()

	service := NewStrikeService(mockStrikeRepo, mockUserRepo, mockLogger)

	ctx := context.Background()
	reporterID := uuid.New()
	targetUserID := uuid.New()

	strikeData := &dto.StrikeCreateRequest{
		ReporterID:   reporterID,
		TargetUserID: targetUserID,
		Type:         constants.StrikeTypeSpam,
		Reason:       "Spam",
	}

	reporter := &domain.User{ID: reporterID}
	target := &domain.User{ID: targetUserID}

	// Expectations
	mockUserRepo.EXPECT().
		GetByID(ctx, reporterID).
		Return(reporter, nil)

	mockUserRepo.EXPECT().
		GetByID(ctx, targetUserID).
		Return(target, nil)

	mockStrikeRepo.EXPECT().
		HasActiveStrikeFromUser(ctx, reporterID.String(), targetUserID.String()).
		Return(true, nil) // Already has active strike

	// Execute
	strike, err := service.CreateStrike(ctx, strikeData)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, errors.ErrDuplicateStrike, err)
	assert.Nil(t, strike)
}

func TestStrikeService_GetStrikesByDateRange_InvalidRange(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStrikeRepo := mocks.NewMockStrikeRepository(ctrl)
	mockUserRepo := mocks.NewMockUserRepository(ctrl)
	mockLogger := mocks.NewMockLogger()

	service := NewStrikeService(mockStrikeRepo, mockUserRepo, mockLogger)

	ctx := context.Background()

	// from is after to - invalid range
	from := time.Now()
	to := from.Add(-24 * time.Hour)

	// Execute
	strikes, err := service.GetStrikesByDateRange(ctx, from, to, 20, 0)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, errors.ErrInvalidDateRange, err)
	assert.Nil(t, strikes)
}

func TestStrikeService_GetStrikesByDateRange_TooLarge(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStrikeRepo := mocks.NewMockStrikeRepository(ctrl)
	mockUserRepo := mocks.NewMockUserRepository(ctrl)
	mockLogger := mocks.NewMockLogger()

	service := NewStrikeService(mockStrikeRepo, mockUserRepo, mockLogger)

	ctx := context.Background()

	// Range is more than 1 year
	from := time.Now().Add(-400 * 24 * time.Hour)
	to := time.Now()

	// Execute
	strikes, err := service.GetStrikesByDateRange(ctx, from, to, 20, 0)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, errors.ErrDateRangeTooLarge, err)
	assert.Nil(t, strikes)
}

func TestStrikeService_GetUserStrikeStats_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStrikeRepo := mocks.NewMockStrikeRepository(ctrl)
	mockUserRepo := mocks.NewMockUserRepository(ctrl)
	mockLogger := mocks.NewMockLogger()

	service := NewStrikeService(mockStrikeRepo, mockUserRepo, mockLogger)

	ctx := context.Background()
	userID := uuid.New()

	user := &domain.User{ID: userID}

	strikes := []*domain.Strike{
		{
			ID:           uuid.New(),
			ReporterID:   uuid.New(),
			TargetUserID: userID,
			Type:         constants.StrikeTypeSpam,
			Status:       constants.StrikeStatusPending,
			CreatedAt:    time.Now().Add(-24 * time.Hour),
		},
		{
			ID:           uuid.New(),
			ReporterID:   uuid.New(),
			TargetUserID: userID,
			Type:         constants.StrikeTypeOffensive,
			Status:       constants.StrikeStatusApproved,
			CreatedAt:    time.Now(),
		},
	}

	// Expectations
	mockUserRepo.EXPECT().
		GetByID(ctx, userID).
		Return(user, nil)

	mockStrikeRepo.EXPECT().
		GetStrikesByUserID(ctx, userID.String(), 1000, 0).
		Return(strikes, nil)

	// Execute
	stats, err := service.GetUserStrikeStats(ctx, userID.String())

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, stats)
	assert.Equal(t, 2, stats.TotalStrikes)
	// Check that status counting works
	assert.Equal(t, 1, stats.StrikeTypes.Pending)
	assert.Equal(t, 1, stats.StrikeTypes.Approved)
	assert.NotNil(t, stats.LastStrikeAt)
}

func TestStrikeService_UpdateStrikeStatus_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStrikeRepo := mocks.NewMockStrikeRepository(ctrl)
	mockUserRepo := mocks.NewMockUserRepository(ctrl)
	mockLogger := mocks.NewMockLogger()

	service := NewStrikeService(mockStrikeRepo, mockUserRepo, mockLogger)

	ctx := context.Background()
	strikeID := "test-strike-id"

	strike := &domain.Strike{
		ID:     uuid.New(),
		Status: constants.StrikeStatusPending,
	}

	// Expectations
	mockStrikeRepo.EXPECT().
		GetStrikeByID(ctx, strikeID).
		Return(strike, nil)

	mockStrikeRepo.EXPECT().
		UpdateStrikeStatus(ctx, strikeID, constants.StrikeStatusApproved).
		Return(nil)

	// Execute
	err := service.UpdateStrikeStatus(ctx, strikeID, constants.StrikeStatusApproved, nil, nil)

	// Assert
	assert.NoError(t, err)
}

func TestStrikeService_UpdateStrikeStatus_SameStatus(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStrikeRepo := mocks.NewMockStrikeRepository(ctrl)
	mockUserRepo := mocks.NewMockUserRepository(ctrl)
	mockLogger := mocks.NewMockLogger()

	service := NewStrikeService(mockStrikeRepo, mockUserRepo, mockLogger)

	ctx := context.Background()
	strikeID := "test-strike-id"

	strike := &domain.Strike{
		ID:     uuid.New(),
		Status: constants.StrikeStatusApproved, // Already approved
	}

	// Expectations
	mockStrikeRepo.EXPECT().
		GetStrikeByID(ctx, strikeID).
		Return(strike, nil)

	// UpdateStrikeStatus should NOT be called because status is the same

	// Execute
	err := service.UpdateStrikeStatus(ctx, strikeID, constants.StrikeStatusApproved, nil, nil)

	// Assert
	assert.NoError(t, err) // Should return nil without error
}

func TestStrikeService_ValidateStrikeCreate_InvalidType(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStrikeRepo := mocks.NewMockStrikeRepository(ctrl)
	mockUserRepo := mocks.NewMockUserRepository(ctrl)
	mockLogger := mocks.NewMockLogger()

	service := NewStrikeService(mockStrikeRepo, mockUserRepo, mockLogger)

	ctx := context.Background()

	strikeData := &dto.StrikeCreateRequest{
		ReporterID:   uuid.New(),
		TargetUserID: uuid.New(),
		Type:         "invalid_type", // Invalid type
		Reason:       "Test",
	}

	// Execute
	err := service.ValidateStrikeCreate(ctx, strikeData)

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid strike type")
}
