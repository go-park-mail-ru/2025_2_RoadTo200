package service_test

import (
	"context"
	"errors"
	"testing"

	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/domain/constants"
	domain "github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/domain/entities"
	expectations "github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/domain/errors"
	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/handler/dto"
	service "github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/service/implementations"
	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/tests/mocks"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSwipeService_ProcessSwipe(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSwipeRepo := mocks.NewMockSwipeRepository(ctrl)
	mockMatchRepo := mocks.NewMockMatchRepository(ctrl)

	swipeService := service.NewSwipeService(mockSwipeRepo, mockMatchRepo)

	ctx := context.Background()
	swiperID := uuid.New()
	targetID := uuid.New()

	tests := []struct {
		name           string
		swiperID       uuid.UUID
		request        *dto.SwipeRequest
		setupMocks     func()
		expectedResp   *dto.SwipeResponse
		expectedErr    error
		wantMatchCheck bool
	}{
		{
			name:     "Success like without match",
			swiperID: swiperID,
			request: &dto.SwipeRequest{
				CardID: targetID.String(),
				Action: "like",
			},
			setupMocks: func() {
				mockSwipeRepo.EXPECT().
					Create(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, swipe *domain.Swipe) error {
						assert.Equal(t, swiperID, swipe.SwiperUserID)
						assert.Equal(t, targetID, swipe.TargetUserID)
						assert.Equal(t, constants.SwipeTypeLike, swipe.SwipeType)
						return nil
					})

				mockMatchRepo.EXPECT().
					CheckMutualLike(gomock.Any(), swiperID, targetID).
					Return(false, nil)
			},
			expectedResp: &dto.SwipeResponse{
				Message: "Swipe processed successfully",
				IsMatch: false,
			},
			expectedErr:    nil,
			wantMatchCheck: true,
		},
		{
			name:     "Success like with match",
			swiperID: swiperID,
			request: &dto.SwipeRequest{
				CardID: targetID.String(),
				Action: "like",
			},
			setupMocks: func() {
				mockSwipeRepo.EXPECT().
					Create(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, swipe *domain.Swipe) error {
						assert.Equal(t, swiperID, swipe.SwiperUserID)
						assert.Equal(t, targetID, swipe.TargetUserID)
						assert.Equal(t, constants.SwipeTypeLike, swipe.SwipeType)
						return nil
					})

				mockMatchRepo.EXPECT().
					CheckMutualLike(gomock.Any(), swiperID, targetID).
					Return(true, nil)

				mockMatchRepo.EXPECT().
					Create(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, match *domain.Match) error {
						assert.Equal(t, swiperID, match.User1ID)
						assert.Equal(t, targetID, match.User2ID)
						assert.True(t, match.IsActive)
						return nil
					})
			},
			expectedResp: &dto.SwipeResponse{
				Message: "It's a match!",
				IsMatch: true,
			},
			expectedErr:    nil,
			wantMatchCheck: true,
		},
		{
			name:     "Success dislike",
			swiperID: swiperID,
			request: &dto.SwipeRequest{
				CardID: targetID.String(),
				Action: "dislike",
			},
			setupMocks: func() {
				mockSwipeRepo.EXPECT().
					Create(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, swipe *domain.Swipe) error {
						assert.Equal(t, swiperID, swipe.SwiperUserID)
						assert.Equal(t, targetID, swipe.TargetUserID)
						assert.Equal(t, constants.SwipeTypeDislike, swipe.SwipeType)
						return nil
					})
				// Для дизлайка не должно быть проверки на мэтч
			},
			expectedResp: &dto.SwipeResponse{
				Message: "Swipe processed successfully",
				IsMatch: false,
			},
			expectedErr:    nil,
			wantMatchCheck: false,
		},
		{
			name:     "Invalid card ID",
			swiperID: swiperID,
			request: &dto.SwipeRequest{
				CardID: "invalid-uuid",
				Action: "like",
			},
			setupMocks:   func() {},
			expectedResp: nil,
			expectedErr:  errors.New("invalid UUID format"), // предполагая, что такая ошибка есть, или просто ошибка парсинга
		},
		{
			name:     "Swipe self",
			swiperID: swiperID,
			request: &dto.SwipeRequest{
				CardID: swiperID.String(),
				Action: "like",
			},
			setupMocks:   func() {},
			expectedResp: nil,
			expectedErr:  expectations.ErrCannotSwipeSelf,
		},
		{
			name:     "Invalid action",
			swiperID: swiperID,
			request: &dto.SwipeRequest{
				CardID: targetID.String(),
				Action: "invalid_action",
			},
			setupMocks:   func() {},
			expectedResp: nil,
			expectedErr:  expectations.ErrInvalidSwipeAction,
		},
		{
			name:     "Swipe repo error",
			swiperID: swiperID,
			request: &dto.SwipeRequest{
				CardID: targetID.String(),
				Action: "like",
			},
			setupMocks: func() {
				mockSwipeRepo.EXPECT().
					Create(gomock.Any(), gomock.Any()).
					Return(assert.AnError)
			},
			expectedResp: nil,
			expectedErr:  assert.AnError,
		},
		{
			name:     "Match check error",
			swiperID: swiperID,
			request: &dto.SwipeRequest{
				CardID: targetID.String(),
				Action: "like",
			},
			setupMocks: func() {
				mockSwipeRepo.EXPECT().
					Create(gomock.Any(), gomock.Any()).
					Return(nil)

				mockMatchRepo.EXPECT().
					CheckMutualLike(gomock.Any(), swiperID, targetID).
					Return(false, assert.AnError)
			},
			expectedResp:   nil,
			expectedErr:    assert.AnError,
			wantMatchCheck: true,
		},
		{
			name:     "Match creation error",
			swiperID: swiperID,
			request: &dto.SwipeRequest{
				CardID: targetID.String(),
				Action: "like",
			},
			setupMocks: func() {
				mockSwipeRepo.EXPECT().
					Create(gomock.Any(), gomock.Any()).
					Return(nil)

				mockMatchRepo.EXPECT().
					CheckMutualLike(gomock.Any(), swiperID, targetID).
					Return(true, nil)

				mockMatchRepo.EXPECT().
					Create(gomock.Any(), gomock.Any()).
					Return(assert.AnError)
			},
			expectedResp:   nil,
			expectedErr:    assert.AnError,
			wantMatchCheck: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMocks()

			resp, err := swipeService.ProcessSwipe(ctx, tt.swiperID, tt.request)

			if tt.expectedErr != nil {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedErr.Error())
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.expectedResp, resp)
			}
		})
	}
}

func TestSwipeService_EdgeCases(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSwipeRepo := mocks.NewMockSwipeRepository(ctrl)
	mockMatchRepo := mocks.NewMockMatchRepository(ctrl)

	swipeService := service.NewSwipeService(mockSwipeRepo, mockMatchRepo)

	ctx := context.Background()

	t.Run("Same user ID validation", func(t *testing.T) {
		userID := uuid.New()
		request := &dto.SwipeRequest{
			CardID: userID.String(),
			Action: "like",
		}

		resp, err := swipeService.ProcessSwipe(ctx, userID, request)

		assert.Nil(t, resp)
		assert.Error(t, err)
		assert.Equal(t, expectations.ErrCannotSwipeSelf, err)
	})

	t.Run("Malformed UUID in card ID", func(t *testing.T) {
		userID := uuid.New()
		request := &dto.SwipeRequest{
			CardID: "not-a-valid-uuid",
			Action: "like",
		}

		resp, err := swipeService.ProcessSwipe(ctx, userID, request)

		assert.Nil(t, resp)
		assert.Error(t, err)
	})

	t.Run("Case sensitivity in actions", func(t *testing.T) {
		userID := uuid.New()
		targetID := uuid.New()

		// Uppercase should fail
		request := &dto.SwipeRequest{
			CardID: targetID.String(),
			Action: "LIKE", // uppercase
		}

		resp, err := swipeService.ProcessSwipe(ctx, userID, request)

		assert.Nil(t, resp)
		assert.Error(t, err)
		assert.Equal(t, expectations.ErrInvalidSwipeAction, err)
	})
}

// Benchmark тесты для проверки производительности
func BenchmarkSwipeService_ProcessSwipe(b *testing.B) {
	ctrl := gomock.NewController(b)
	defer ctrl.Finish()

	mockSwipeRepo := mocks.NewMockSwipeRepository(ctrl)
	mockMatchRepo := mocks.NewMockMatchRepository(ctrl)

	swipeService := service.NewSwipeService(mockSwipeRepo, mockMatchRepo)

	ctx := context.Background()
	swiperID := uuid.New()
	targetID := uuid.New()

	// Настраиваем моки для успешного лайка без мэтча
	mockSwipeRepo.EXPECT().
		Create(gomock.Any(), gomock.Any()).
		Return(nil).
		AnyTimes()

	mockMatchRepo.EXPECT().
		CheckMutualLike(gomock.Any(), swiperID, targetID).
		Return(false, nil).
		AnyTimes()

	request := &dto.SwipeRequest{
		CardID: targetID.String(),
		Action: "like",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = swipeService.ProcessSwipe(ctx, swiperID, request)
	}
}
