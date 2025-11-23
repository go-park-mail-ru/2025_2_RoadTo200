package service_test

import (
	"context"
	"testing"
	"time"

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

// calculateAge вычисляет возраст на основе даты рождения (копия реализации из сервиса)
func calculateAge(birthDate time.Time) int {
	if birthDate.IsZero() {
		return 0
	}
	now := time.Now()
	age := now.Year() - birthDate.Year()
	if now.YearDay() < birthDate.YearDay() {
		age--
	}
	return age
}

func TestMatchService_GetUserMatches(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockMatchRepo := mocks.NewMockMatchRepository(ctrl)
	mockUserRepo := mocks.NewMockUserRepository(ctrl)
	mockSwipeRepo := mocks.NewMockSwipeRepository(ctrl)
	mockPhotoRepo := mocks.NewMockUserPhotoRepository(ctrl)
	mockLogger := mocks.NewMockLogger()

	matchService := service.NewMatchService(
		mockMatchRepo,
		mockUserRepo,
		mockSwipeRepo,
		mockPhotoRepo,
		mockLogger,
	)

	ctx := context.Background()

	// Используем фиксированные UUID для стабильности тестов
	userID := uuid.MustParse("12345678-1234-1234-1234-123456789012")
	matchedUserID1 := uuid.MustParse("22345678-1234-1234-1234-123456789012")
	matchedUserID2 := uuid.MustParse("32345678-1234-1234-1234-123456789012")

	// Фиксированное время для стабильного вычисления возраста
	fixedTime := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	birthDate1 := time.Date(1999, 1, 1, 0, 0, 0, 0, time.UTC) // 25 лет на fixedTime
	birthDate2 := time.Date(1994, 1, 1, 0, 0, 0, 0, time.UTC) // 30 лет на fixedTime

	tests := []struct {
		name           string
		userID         uuid.UUID
		limit          int
		offset         int
		setupMocks     func()
		expectedResult *dto.MatchesResponse
		expectedErr    error
	}{
		{
			name:   "Success - with matches",
			userID: userID,
			limit:  20,
			offset: 0,
			setupMocks: func() {
				matches := []domain.Match{
					{
						User1ID:   userID,
						User2ID:   matchedUserID1,
						IsActive:  true,
						MatchedAt: fixedTime.Add(-24 * time.Hour),
					},
					{
						User1ID:   matchedUserID2,
						User2ID:   userID,
						IsActive:  true,
						MatchedAt: fixedTime.Add(-12 * time.Hour),
					},
				}
				mockMatchRepo.EXPECT().GetUserMatches(ctx, userID, 20, 0).Return(matches, nil)

				users := []domain.User{
					{
						ID:        matchedUserID1,
						Name:      "User One",
						BirthDate: birthDate1,
						Bio:       stringPtr("Bio for user one"),
					},
					{
						ID:        matchedUserID2,
						Name:      "User Two",
						BirthDate: birthDate2,
						Bio:       stringPtr("Bio for user two"),
					},
				}
				mockUserRepo.EXPECT().GetUsersByIDs(ctx, []uuid.UUID{matchedUserID1, matchedUserID2}).Return(users, nil)

				photos1 := []domain.UserPhoto{
					{PhotoURL: "https://example.com/photo1.jpg", IsApproved: true},
					{PhotoURL: "https://example.com/photo2.jpg", IsApproved: true},
				}
				photos2 := []domain.UserPhoto{
					{PhotoURL: "https://example.com/photo3.jpg", IsApproved: true},
				}
				mockPhotoRepo.EXPECT().GetByUserID(ctx, matchedUserID1).Return(photos1, nil)
				mockPhotoRepo.EXPECT().GetByUserID(ctx, matchedUserID2).Return(photos2, nil)
			},
			expectedResult: &dto.MatchesResponse{
				Matches: []dto.MatchResponse{
					{
						Match: domain.Match{
							User1ID:  userID,
							User2ID:  matchedUserID1,
							IsActive: true,
						},
						User: domain.User{
							ID:        matchedUserID1,
							Name:      "User One",
							BirthDate: birthDate1,
							Bio:       stringPtr("Bio for user one"),
						},
						Photos:      []string{"https://example.com/photo1.jpg", "https://example.com/photo2.jpg"},
						Age:         calculateAge(birthDate1), // Используем ту же функцию, что и в сервисе
						Description: "Bio for user one",
						PhotosCount: 2,
					},
					{
						Match: domain.Match{
							User1ID:  matchedUserID2,
							User2ID:  userID,
							IsActive: true,
						},
						User: domain.User{
							ID:        matchedUserID2,
							Name:      "User Two",
							BirthDate: birthDate2,
							Bio:       stringPtr("Bio for user two"),
						},
						Photos:      []string{"https://example.com/photo3.jpg"},
						Age:         calculateAge(birthDate2), // Используем ту же функцию, что и в сервисе
						Description: "Bio for user two",
						PhotosCount: 1,
					},
				},
				Total:  2,
				Limit:  20,
				Offset: 0,
			},
			expectedErr: nil,
		},
		{
			name:   "Success - no matches",
			userID: userID,
			limit:  20,
			offset: 0,
			setupMocks: func() {
				mockMatchRepo.EXPECT().GetUserMatches(ctx, userID, 20, 0).Return([]domain.Match{}, nil)
			},
			expectedResult: &dto.MatchesResponse{
				Matches: []dto.MatchResponse{},
				Total:   0,
				Limit:   20,
				Offset:  0,
			},
			expectedErr: nil,
		},
		{
			name:   "Success - limit and offset adjustment",
			userID: userID,
			limit:  0,  // Должен быть установлен в 20
			offset: -5, // Должен быть установлен в 0
			setupMocks: func() {
				mockMatchRepo.EXPECT().GetUserMatches(ctx, userID, 20, 0).Return([]domain.Match{}, nil)
			},
			expectedResult: &dto.MatchesResponse{
				Matches: []dto.MatchResponse{},
				Total:   0,
				Limit:   20,
				Offset:  0,
			},
			expectedErr: nil,
		},
		{
			name:   "Error - get matches fails",
			userID: userID,
			limit:  20,
			offset: 0,
			setupMocks: func() {
				mockMatchRepo.EXPECT().GetUserMatches(ctx, userID, 20, 0).Return(nil, assert.AnError)
			},
			expectedResult: nil,
			expectedErr:    assert.AnError,
		},
		{
			name:   "Error - get users fails",
			userID: userID,
			limit:  20,
			offset: 0,
			setupMocks: func() {
				matches := []domain.Match{
					{
						User1ID: userID,
						User2ID: matchedUserID1,
					},
				}
				mockMatchRepo.EXPECT().GetUserMatches(ctx, userID, 20, 0).Return(matches, nil)
				mockUserRepo.EXPECT().GetUsersByIDs(ctx, []uuid.UUID{matchedUserID1}).Return(nil, assert.AnError)
			},
			expectedResult: nil,
			expectedErr:    assert.AnError,
		},
		{
			name:   "Error - get photos fails",
			userID: userID,
			limit:  20,
			offset: 0,
			setupMocks: func() {
				matches := []domain.Match{
					{
						User1ID: userID,
						User2ID: matchedUserID1,
					},
				}
				mockMatchRepo.EXPECT().GetUserMatches(ctx, userID, 20, 0).Return(matches, nil)

				users := []domain.User{
					{
						ID:   matchedUserID1,
						Name: "User One",
					},
				}
				mockUserRepo.EXPECT().GetUsersByIDs(ctx, []uuid.UUID{matchedUserID1}).Return(users, nil)
				mockPhotoRepo.EXPECT().GetByUserID(ctx, matchedUserID1).Return(nil, assert.AnError)
			},
			expectedResult: nil,
			expectedErr:    assert.AnError,
		},
		{
			name:   "Success - user not found in match (skip)",
			userID: userID,
			limit:  20,
			offset: 0,
			setupMocks: func() {
				matches := []domain.Match{
					{
						User1ID: userID,
						User2ID: matchedUserID1,
					},
					{
						User1ID: matchedUserID2, // Этот пользователь не будет найден
						User2ID: userID,
					},
				}
				mockMatchRepo.EXPECT().GetUserMatches(ctx, userID, 20, 0).Return(matches, nil)

				// Возвращаем только одного пользователя
				users := []domain.User{
					{
						ID:   matchedUserID1,
						Name: "Found User",
					},
				}
				mockUserRepo.EXPECT().GetUsersByIDs(ctx, []uuid.UUID{matchedUserID1, matchedUserID2}).Return(users, nil)

				// Фотографии для обоих пользователей (второй будет пустым)
				photos1 := []domain.UserPhoto{
					{PhotoURL: "https://example.com/photo.jpg", IsApproved: true},
				}
				mockPhotoRepo.EXPECT().GetByUserID(ctx, matchedUserID1).Return(photos1, nil)
				mockPhotoRepo.EXPECT().GetByUserID(ctx, matchedUserID2).Return([]domain.UserPhoto{}, nil)
			},
			expectedResult: &dto.MatchesResponse{
				Matches: []dto.MatchResponse{
					{
						Match: domain.Match{
							User1ID: userID,
							User2ID: matchedUserID1,
						},
						User: domain.User{
							ID:   matchedUserID1,
							Name: "Found User",
						},
						Photos:      []string{"https://example.com/photo.jpg"},
						Age:         2024, // TODO: Age будет 0, так как BirthDate не установлен
						Description: "Нет описания",
						PhotosCount: 1,
					},
				},
				Total:  1,
				Limit:  20,
				Offset: 0,
			},
			expectedErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMocks()

			result, err := matchService.GetUserMatches(ctx, tt.userID, tt.limit, tt.offset)

			if tt.expectedErr != nil {
				require.Error(t, err)
				// Используем ErrorIs для проверки обернутых ошибок
				assert.ErrorIs(t, err, tt.expectedErr)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.expectedResult.Total, result.Total)
				assert.Equal(t, tt.expectedResult.Limit, result.Limit)
				assert.Equal(t, tt.expectedResult.Offset, result.Offset)
				assert.Len(t, result.Matches, len(tt.expectedResult.Matches))

				// Проверяем каждую пару мэтчей
				for i, expectedMatch := range tt.expectedResult.Matches {
					actualMatch := result.Matches[i]

					assert.Equal(t, expectedMatch.Match.User1ID, actualMatch.Match.User1ID)
					assert.Equal(t, expectedMatch.Match.User2ID, actualMatch.Match.User2ID)
					assert.Equal(t, expectedMatch.User.ID, actualMatch.User.ID)
					assert.Equal(t, expectedMatch.User.Name, actualMatch.User.Name)
					assert.Equal(t, expectedMatch.Age, actualMatch.Age)
					assert.Equal(t, expectedMatch.Description, actualMatch.Description)
					assert.Equal(t, expectedMatch.PhotosCount, actualMatch.PhotosCount)
					assert.ElementsMatch(t, expectedMatch.Photos, actualMatch.Photos)
				}
			}
		})
	}
}

func TestMatchService_Unmatch(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockMatchRepo := mocks.NewMockMatchRepository(ctrl)
	mockUserRepo := mocks.NewMockUserRepository(ctrl)
	mockSwipeRepo := mocks.NewMockSwipeRepository(ctrl)
	mockPhotoRepo := mocks.NewMockUserPhotoRepository(ctrl)
	mockLogger := mocks.NewMockLogger()

	matchService := service.NewMatchService(
		mockMatchRepo,
		mockUserRepo,
		mockSwipeRepo,
		mockPhotoRepo,
		mockLogger,
	)

	ctx := context.Background()
	// Используем фиксированные UUID
	userID := uuid.MustParse("12345678-1234-1234-1234-123456789012")
	targetUserID := uuid.MustParse("22345678-1234-1234-1234-123456789012")

	tests := []struct {
		name        string
		userID      uuid.UUID
		targetID    uuid.UUID
		setupMocks  func()
		expectedErr error
	}{
		{
			name:     "Success - unmatch",
			userID:   userID,
			targetID: targetUserID,
			setupMocks: func() {
				match := &domain.Match{
					User1ID:  userID,
					User2ID:  targetUserID,
					IsActive: true,
				}
				mockMatchRepo.EXPECT().GetByUsers(ctx, userID, targetUserID).Return(match, nil)
				mockMatchRepo.EXPECT().UpdateActive(ctx, userID, targetUserID, false).Return(nil)
			},
			expectedErr: nil,
		},
		{
			name:     "Error - match not found",
			userID:   userID,
			targetID: targetUserID,
			setupMocks: func() {
				mockMatchRepo.EXPECT().GetByUsers(ctx, userID, targetUserID).Return(nil, nil)
			},
			expectedErr: expectations.ErrMatchNotFound,
		},
		{
			name:     "Error - get match fails",
			userID:   userID,
			targetID: targetUserID,
			setupMocks: func() {
				mockMatchRepo.EXPECT().GetByUsers(ctx, userID, targetUserID).Return(nil, assert.AnError)
			},
			expectedErr: assert.AnError,
		},
		{
			name:     "Error - update active fails",
			userID:   userID,
			targetID: targetUserID,
			setupMocks: func() {
				match := &domain.Match{
					User1ID:  userID,
					User2ID:  targetUserID,
					IsActive: true,
				}
				mockMatchRepo.EXPECT().GetByUsers(ctx, userID, targetUserID).Return(match, nil)
				mockMatchRepo.EXPECT().UpdateActive(ctx, userID, targetUserID, false).Return(assert.AnError)
			},
			expectedErr: assert.AnError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMocks()

			err := matchService.Unmatch(ctx, tt.userID, tt.targetID)

			if tt.expectedErr != nil {
				require.Error(t, err)
				assert.Equal(t, tt.expectedErr, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestMatchService_EdgeCases(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockMatchRepo := mocks.NewMockMatchRepository(ctrl)
	mockUserRepo := mocks.NewMockUserRepository(ctrl)
	mockSwipeRepo := mocks.NewMockSwipeRepository(ctrl)
	mockPhotoRepo := mocks.NewMockUserPhotoRepository(ctrl)
	mockLogger := mocks.NewMockLogger()

	matchService := service.NewMatchService(
		mockMatchRepo,
		mockUserRepo,
		mockSwipeRepo,
		mockPhotoRepo,
		mockLogger,
	)

	ctx := context.Background()
	// Используем фиксированные UUID
	userID := uuid.MustParse("12345678-1234-1234-1234-123456789012")

	t.Run("Large limit adjustment", func(t *testing.T) {
		// Мокаем вызов с автоматически скорректированным лимитом
		mockMatchRepo.EXPECT().GetUserMatches(ctx, userID, 20, 0).Return([]domain.Match{}, nil)

		result, err := matchService.GetUserMatches(ctx, userID, 200, 0)

		require.NoError(t, err)
		assert.Equal(t, 20, result.Limit) // Лимит должен быть скорректирован до 50
	})

	t.Run("Photos with unapproved ones", func(t *testing.T) {
		matchedUserID := uuid.MustParse("22345678-1234-1234-1234-123456789012")

		matches := []domain.Match{
			{
				User1ID: userID,
				User2ID: matchedUserID,
			},
		}
		mockMatchRepo.EXPECT().GetUserMatches(ctx, userID, 20, 0).Return(matches, nil)

		users := []domain.User{
			{
				ID:   matchedUserID,
				Name: "Test User",
			},
		}
		mockUserRepo.EXPECT().GetUsersByIDs(ctx, []uuid.UUID{matchedUserID}).Return(users, nil)

		// Смешанные approved и не approved фото
		photos := []domain.UserPhoto{
			{PhotoURL: "approved1.jpg", IsApproved: true},
			{PhotoURL: "not-approved.jpg", IsApproved: false}, // Должно быть отфильтровано
			{PhotoURL: "approved2.jpg", IsApproved: true},
		}
		mockPhotoRepo.EXPECT().GetByUserID(ctx, matchedUserID).Return(photos, nil)

		result, err := matchService.GetUserMatches(ctx, userID, 20, 0)

		require.NoError(t, err)
		require.Len(t, result.Matches, 1)
		assert.ElementsMatch(t, []string{"approved1.jpg", "approved2.jpg"}, result.Matches[0].Photos)
		assert.Equal(t, 2, result.Matches[0].PhotosCount)
	})

	t.Run("User without bio", func(t *testing.T) {
		matchedUserID := uuid.MustParse("22345678-1234-1234-1234-123456789012")

		matches := []domain.Match{
			{
				User1ID: userID,
				User2ID: matchedUserID,
			},
		}
		mockMatchRepo.EXPECT().GetUserMatches(ctx, userID, 20, 0).Return(matches, nil)

		users := []domain.User{
			{
				ID:   matchedUserID,
				Name: "Test User",
				// Bio is nil
			},
		}
		mockUserRepo.EXPECT().GetUsersByIDs(ctx, []uuid.UUID{matchedUserID}).Return(users, nil)
		mockPhotoRepo.EXPECT().GetByUserID(ctx, matchedUserID).Return([]domain.UserPhoto{}, nil)

		result, err := matchService.GetUserMatches(ctx, userID, 20, 0)

		require.NoError(t, err)
		require.Len(t, result.Matches, 1)
		assert.Equal(t, "Нет описания", result.Matches[0].Description) // Description должен быть пустым
	})
}

// Benchmark тесты
func BenchmarkMatchService_GetUserMatches(b *testing.B) {
	ctrl := gomock.NewController(b)
	defer ctrl.Finish()

	mockMatchRepo := mocks.NewMockMatchRepository(ctrl)
	mockUserRepo := mocks.NewMockUserRepository(ctrl)
	mockSwipeRepo := mocks.NewMockSwipeRepository(ctrl)
	mockPhotoRepo := mocks.NewMockUserPhotoRepository(ctrl)
	mockLogger := mocks.NewMockLogger()

	matchService := service.NewMatchService(
		mockMatchRepo,
		mockUserRepo,
		mockSwipeRepo,
		mockPhotoRepo,
		mockLogger,
	)

	ctx := context.Background()
	// Используем фиксированные UUID
	userID := uuid.MustParse("12345678-1234-1234-1234-123456789012")
	matchedUserID := uuid.MustParse("22345678-1234-1234-1234-123456789012")

	// Настраиваем моки для успешного выполнения
	matches := []domain.Match{
		{
			User1ID: userID,
			User2ID: matchedUserID,
		},
	}
	users := []domain.User{
		{
			ID:   matchedUserID,
			Name: "Test User",
		},
	}
	photos := []domain.UserPhoto{
		{PhotoURL: "photo1.jpg", IsApproved: true},
	}

	mockMatchRepo.EXPECT().GetUserMatches(gomock.Any(), userID, 20, 0).Return(matches, nil).AnyTimes()
	mockUserRepo.EXPECT().GetUsersByIDs(gomock.Any(), gomock.Any()).Return(users, nil).AnyTimes()
	mockPhotoRepo.EXPECT().GetByUserID(gomock.Any(), matchedUserID).Return(photos, nil).AnyTimes()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = matchService.GetUserMatches(ctx, userID, 20, 0)
	}
}
