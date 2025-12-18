package service

// import (
// 	"context"
// 	"testing"

// 	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/domain/constants"
// 	domain "github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/domain/entities"
// 	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/domain/errors"
// 	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/handler/dto"
// 	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/tests/mocks"
// 	"github.com/google/uuid"
// 	"github.com/stretchr/testify/assert"
// 	"go.uber.org/mock/gomock"
// )

// func TestSwipeService_ProcessSwipe_Like_NoMatch(t *testing.T) {
// 	ctrl := gomock.NewController(t)
// 	defer ctrl.Finish()

// 	mockSwipeRepo := mocks.NewMockSwipeRepository(ctrl)
// 	mockMatchRepo := mocks.NewMockMatchRepository(ctrl)

// 	service := NewSwipeService(mockSwipeRepo, mockMatchRepo)

// 	ctx := context.Background()
// 	swiperID := uuid.New()
// 	targetID := uuid.New()

// 	request := &dto.SwipeRequest{
// 		CardID: targetID.String(),
// 		Action: "like",
// 	}

// 	// Expectations
// 	mockSwipeRepo.EXPECT().
// 		Create(ctx, gomock.Any()).
// 		Return(nil)

// 	mockMatchRepo.EXPECT().
// 		CheckMutualLike(ctx, swiperID, targetID).
// 		Return(false, nil)

// 	// Execute
// 	response, err := service.ProcessSwipe(ctx, swiperID, request)

// 	// Assert
// 	assert.NoError(t, err)
// 	assert.NotNil(t, response)
// 	assert.False(t, response.IsMatch)
// 	assert.Equal(t, "Swipe processed successfully", response.Message)
// }

// func TestSwipeService_ProcessSwipe_Like_WithMatch(t *testing.T) {
// 	ctrl := gomock.NewController(t)
// 	defer ctrl.Finish()

// 	mockSwipeRepo := mocks.NewMockSwipeRepository(ctrl)
// 	mockMatchRepo := mocks.NewMockMatchRepository(ctrl)

// 	service := NewSwipeService(mockSwipeRepo, mockMatchRepo)

// 	ctx := context.Background()
// 	swiperID := uuid.New()
// 	targetID := uuid.New()

// 	request := &dto.SwipeRequest{
// 		CardID: targetID.String(),
// 		Action: "like",
// 	}

// 	// Expectations
// 	mockSwipeRepo.EXPECT().
// 		Create(ctx, gomock.Any()).
// 		DoAndReturn(func(ctx context.Context, swipe *domain.Swipe) error {
// 			assert.Equal(t, swiperID, swipe.SwiperUserID)
// 			assert.Equal(t, targetID, swipe.TargetUserID)
// 			assert.Equal(t, constants.SwipeTypeLike, swipe.SwipeType)
// 			return nil
// 		})

// 	mockMatchRepo.EXPECT().
// 		CheckMutualLike(ctx, swiperID, targetID).
// 		Return(true, nil)

// 	mockMatchRepo.EXPECT().
// 		Create(ctx, gomock.Any()).
// 		DoAndReturn(func(ctx context.Context, match *domain.Match) error {
// 			assert.True(t, match.User1ID == swiperID || match.User1ID == targetID)
// 			assert.True(t, match.User2ID == swiperID || match.User2ID == targetID)
// 			assert.True(t, match.IsActive)
// 			return nil
// 		})

// 	// Execute
// 	response, err := service.ProcessSwipe(ctx, swiperID, request)

// 	// Assert
// 	assert.NoError(t, err)
// 	assert.NotNil(t, response)
// 	assert.True(t, response.IsMatch)
// 	assert.Equal(t, "It's a match!", response.Message)
// }

// func TestSwipeService_ProcessSwipe_Dislike(t *testing.T) {
// 	ctrl := gomock.NewController(t)
// 	defer ctrl.Finish()

// 	mockSwipeRepo := mocks.NewMockSwipeRepository(ctrl)
// 	mockMatchRepo := mocks.NewMockMatchRepository(ctrl)

// 	service := NewSwipeService(mockSwipeRepo, mockMatchRepo)

// 	ctx := context.Background()
// 	swiperID := uuid.New()
// 	targetID := uuid.New()

// 	request := &dto.SwipeRequest{
// 		CardID: targetID.String(),
// 		Action: "dislike",
// 	}

// 	// Expectations
// 	mockSwipeRepo.EXPECT().
// 		Create(ctx, gomock.Any()).
// 		DoAndReturn(func(ctx context.Context, swipe *domain.Swipe) error {
// 			assert.Equal(t, constants.SwipeTypeDislike, swipe.SwipeType)
// 			return nil
// 		})

// 	// Execute
// 	response, err := service.ProcessSwipe(ctx, swiperID, request)

// 	// Assert
// 	assert.NoError(t, err)
// 	assert.NotNil(t, response)
// 	assert.False(t, response.IsMatch)
// }

// func TestSwipeService_ProcessSwipe_InvalidAction(t *testing.T) {
// 	ctrl := gomock.NewController(t)
// 	defer ctrl.Finish()

// 	mockSwipeRepo := mocks.NewMockSwipeRepository(ctrl)
// 	mockMatchRepo := mocks.NewMockMatchRepository(ctrl)

// 	service := NewSwipeService(mockSwipeRepo, mockMatchRepo)

// 	ctx := context.Background()
// 	swiperID := uuid.New()
// 	targetID := uuid.New()

// 	request := &dto.SwipeRequest{
// 		CardID: targetID.String(),
// 		Action: "invalid",
// 	}

// 	// Execute
// 	response, err := service.ProcessSwipe(ctx, swiperID, request)

// 	// Assert
// 	assert.Error(t, err)
// 	assert.Equal(t, errors.ErrInvalidSwipeAction, err)
// 	assert.Nil(t, response)
// }

// func TestSwipeService_ProcessSwipe_InvalidCardID(t *testing.T) {
// 	ctrl := gomock.NewController(t)
// 	defer ctrl.Finish()

// 	mockSwipeRepo := mocks.NewMockSwipeRepository(ctrl)
// 	mockMatchRepo := mocks.NewMockMatchRepository(ctrl)

// 	service := NewSwipeService(mockSwipeRepo, mockMatchRepo)

// 	ctx := context.Background()
// 	swiperID := uuid.New()

// 	request := &dto.SwipeRequest{
// 		CardID: "invalid-uuid",
// 		Action: "like",
// 	}

// 	// Execute
// 	response, err := service.ProcessSwipe(ctx, swiperID, request)

// 	// Assert
// 	assert.Error(t, err)
// 	assert.Nil(t, response)
// }

// func TestSwipeService_ProcessSwipe_SwipeSelf(t *testing.T) {
// 	ctrl := gomock.NewController(t)
// 	defer ctrl.Finish()

// 	mockSwipeRepo := mocks.NewMockSwipeRepository(ctrl)
// 	mockMatchRepo := mocks.NewMockMatchRepository(ctrl)

// 	service := NewSwipeService(mockSwipeRepo, mockMatchRepo)

// 	ctx := context.Background()
// 	userID := uuid.New()

// 	request := &dto.SwipeRequest{
// 		CardID: userID.String(),
// 		Action: "like",
// 	}

// 	// Execute
// 	response, err := service.ProcessSwipe(ctx, userID, request)

// 	// Assert
// 	assert.Error(t, err)
// 	assert.Equal(t, errors.ErrCannotSwipeSelf, err)
// 	assert.Nil(t, response)
// }

// func TestSwipeService_ProcessSwipe_CreateSwipeError(t *testing.T) {
// 	ctrl := gomock.NewController(t)
// 	defer ctrl.Finish()

// 	mockSwipeRepo := mocks.NewMockSwipeRepository(ctrl)
// 	mockMatchRepo := mocks.NewMockMatchRepository(ctrl)

// 	service := NewSwipeService(mockSwipeRepo, mockMatchRepo)

// 	ctx := context.Background()
// 	swiperID := uuid.New()
// 	targetID := uuid.New()

// 	request := &dto.SwipeRequest{
// 		CardID: targetID.String(),
// 		Action: "like",
// 	}

// 	expectedErr := errors.ErrInternalError

// 	// Expectations
// 	mockSwipeRepo.EXPECT().
// 		Create(ctx, gomock.Any()).
// 		Return(expectedErr)

// 	// Execute
// 	response, err := service.ProcessSwipe(ctx, swiperID, request)

// 	// Assert
// 	assert.Error(t, err)
// 	assert.Equal(t, expectedErr, err)
// 	assert.Nil(t, response)
// }

// func TestSwipeService_ProcessSwipe_CheckMutualLikeError(t *testing.T) {
// 	ctrl := gomock.NewController(t)
// 	defer ctrl.Finish()

// 	mockSwipeRepo := mocks.NewMockSwipeRepository(ctrl)
// 	mockMatchRepo := mocks.NewMockMatchRepository(ctrl)

// 	service := NewSwipeService(mockSwipeRepo, mockMatchRepo)

// 	ctx := context.Background()
// 	swiperID := uuid.New()
// 	targetID := uuid.New()

// 	request := &dto.SwipeRequest{
// 		CardID: targetID.String(),
// 		Action: "like",
// 	}

// 	expectedErr := errors.ErrInternalError

// 	// Expectations
// 	mockSwipeRepo.EXPECT().
// 		Create(ctx, gomock.Any()).
// 		Return(nil)

// 	mockMatchRepo.EXPECT().
// 		CheckMutualLike(ctx, swiperID, targetID).
// 		Return(false, expectedErr)

// 	// Execute
// 	response, err := service.ProcessSwipe(ctx, swiperID, request)

// 	// Assert
// 	assert.Error(t, err)
// 	assert.Equal(t, expectedErr, err)
// 	assert.Nil(t, response)
// }

// func TestSwipeService_ProcessSwipe_CreateMatchError(t *testing.T) {
// 	ctrl := gomock.NewController(t)
// 	defer ctrl.Finish()

// 	mockSwipeRepo := mocks.NewMockSwipeRepository(ctrl)
// 	mockMatchRepo := mocks.NewMockMatchRepository(ctrl)

// 	service := NewSwipeService(mockSwipeRepo, mockMatchRepo)

// 	ctx := context.Background()
// 	swiperID := uuid.New()
// 	targetID := uuid.New()

// 	request := &dto.SwipeRequest{
// 		CardID: targetID.String(),
// 		Action: "like",
// 	}

// 	expectedErr := errors.ErrInternalError

// 	// Expectations
// 	mockSwipeRepo.EXPECT().
// 		Create(ctx, gomock.Any()).
// 		Return(nil)

// 	mockMatchRepo.EXPECT().
// 		CheckMutualLike(ctx, swiperID, targetID).
// 		Return(true, nil)

// 	mockMatchRepo.EXPECT().
// 		Create(ctx, gomock.Any()).
// 		Return(expectedErr)

// 	// Execute
// 	response, err := service.ProcessSwipe(ctx, swiperID, request)

// 	// Assert
// 	assert.Error(t, err)
// 	assert.Equal(t, expectedErr, err)
// 	assert.Nil(t, response)
// }
