package service_test

import (
	"bytes"
	"context"
	"fmt"
	"mime/multipart"
	"strings"
	"testing"
	"time"

	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/domain/constants"
	domain "github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/domain/entities"
	expectations "github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/domain/errors"
	service "github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/service/implementations"
	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/tests/mocks"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestProfileService_GetProfile тестирует получение профиля
func TestProfileService_GetProfile(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mocks.NewMockUserRepository(ctrl)
	mockUserPhotoRepo := mocks.NewMockUserPhotoRepository(ctrl)
	mockPreferenceRepo := mocks.NewMockUserPreferenceRepository(ctrl)
	mockFileStorage := mocks.NewMockFileStorage(ctrl)
	mockLogger := mocks.NewMockLogger()

	profileService := service.NewProfileService(
		mockUserRepo,
		mockUserPhotoRepo,
		mockPreferenceRepo,
		mockFileStorage,
		mockLogger,
	)

	ctx := context.Background()
	userID := uuid.New()

	tests := []struct {
		name           string
		userID         uuid.UUID
		setupMocks     func()
		expectedResult *domain.ProfileResponse
		expectedErr    error
	}{
		{
			name:   "Success - full profile",
			userID: userID,
			setupMocks: func() {
				user := &domain.User{
					ID:        userID,
					Name:      "Test User",
					Email:     "test@example.com",
					Gender:    constants.Gender("male"),
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				}
				mockUserRepo.EXPECT().GetByID(ctx, userID).Return(user, nil)

				preferences := &domain.UserPreference{
					UserID:       userID,
					ShowGender:   constants.GenderPreference("female"),
					AgeMin:       18,
					AgeMax:       30,
					MaxDistance:  50,
					GlobalSearch: true,
				}
				mockPreferenceRepo.EXPECT().GetByUserID(ctx, userID).Return(preferences, nil)

				photos := []domain.UserPhoto{
					{
						ID:           uuid.New(),
						UserID:       userID,
						PhotoURL:     "https://example.com/photo1.jpg",
						DisplayOrder: 0,
						IsApproved:   true,
					},
				}
				mockUserPhotoRepo.EXPECT().GetByUserID(ctx, userID).Return(photos, nil)

				// Исправление: разрешаем вызов UpdateLastActive в любом количестве раз
				mockUserRepo.EXPECT().UpdateLastActive(gomock.Any(), gomock.Any()).AnyTimes()
			},
			expectedResult: &domain.ProfileResponse{
				User: &domain.User{
					ID:        userID,
					Name:      "Test User",
					Email:     "test@example.com",
					Gender:    constants.Gender("male"),
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				},
				Preferences: &domain.UserPreference{
					UserID:       userID,
					ShowGender:   constants.GenderPreference("female"),
					AgeMin:       18,
					AgeMax:       30,
					MaxDistance:  50,
					GlobalSearch: true,
				},
				Photos: []domain.UserPhoto{
					{
						ID:           uuid.New(),
						UserID:       userID,
						PhotoURL:     "https://example.com/photo1.jpg",
						DisplayOrder: 0,
						IsApproved:   true,
					},
				},
			},
			expectedErr: nil,
		},
		{
			name:   "Error - user not found",
			userID: userID,
			setupMocks: func() {
				mockUserRepo.EXPECT().GetByID(ctx, userID).Return(nil, expectations.ErrUserNotFound)
				// Исправление: разрешаем вызов UpdateLastActive в любом количестве раз
				mockUserRepo.EXPECT().UpdateLastActive(gomock.Any(), gomock.Any()).AnyTimes()
			},
			expectedResult: nil,
			expectedErr:    expectations.ErrUserNotFound,
		},
		{
			name:   "Error - get preferences fails",
			userID: userID,
			setupMocks: func() {
				user := &domain.User{ID: userID}
				mockUserRepo.EXPECT().GetByID(ctx, userID).Return(user, nil)
				mockPreferenceRepo.EXPECT().GetByUserID(ctx, userID).Return(nil, assert.AnError)
				// Исправление: разрешаем вызов UpdateLastActive в любом количестве раз
				mockUserRepo.EXPECT().UpdateLastActive(gomock.Any(), gomock.Any()).AnyTimes()
			},
			expectedResult: nil,
			expectedErr:    assert.AnError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMocks()

			result, err := profileService.GetProfile(ctx, tt.userID)

			if tt.expectedErr != nil {
				require.Error(t, err)
				assert.Equal(t, tt.expectedErr, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.expectedResult.User.ID, result.User.ID)
				assert.Equal(t, tt.expectedResult.User.Name, result.User.Name)
				assert.NotNil(t, result.Preferences)
				assert.NotNil(t, result.Photos)
			}
		})
	}
}

// TestProfileService_UpdateProfileInfo тестирует обновление информации профиля
func TestProfileService_UpdateProfileInfo(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mocks.NewMockUserRepository(ctrl)
	mockUserPhotoRepo := mocks.NewMockUserPhotoRepository(ctrl)
	mockPreferenceRepo := mocks.NewMockUserPreferenceRepository(ctrl)
	mockFileStorage := mocks.NewMockFileStorage(ctrl)
	mockLogger := mocks.NewMockLogger()

	profileService := service.NewProfileService(
		mockUserRepo,
		mockUserPhotoRepo,
		mockPreferenceRepo,
		mockFileStorage,
		mockLogger,
	)

	ctx := context.Background()
	userID := uuid.New()

	tests := []struct {
		name        string
		userID      uuid.UUID
		updateData  *domain.ProfileUpdateRequest
		setupMocks  func()
		expectedErr error
	}{
		{
			name:   "Success - update name and gender",
			userID: userID,
			updateData: &domain.ProfileUpdateRequest{
				Name:   "Updated Name",
				Gender: constants.Gender("female"),
			},
			setupMocks: func() {
				user := &domain.User{
					ID:        userID,
					Name:      "Original Name",
					Gender:    constants.Gender("male"),
					UpdatedAt: time.Now().Add(-time.Hour),
				}
				mockUserRepo.EXPECT().GetByID(ctx, userID).Return(user, nil)
				mockUserRepo.EXPECT().Update(ctx, gomock.Any()).DoAndReturn(func(ctx context.Context, user *domain.User) error {
					assert.Equal(t, "Updated Name", user.Name)
					assert.Equal(t, constants.Gender("female"), user.Gender)
					return nil
				})
			},
			expectedErr: nil,
		},
		{
			name:   "Error - validation fails - name too long",
			userID: userID,
			updateData: &domain.ProfileUpdateRequest{
				Name: strings.Repeat("a", constants.MaxNameLength+1), // Too long name
			},
			setupMocks:  func() {},
			expectedErr: fmt.Errorf("name too long, max %d characters", constants.MaxNameLength),
		},
		{
			name:   "Error - user not found",
			userID: userID,
			updateData: &domain.ProfileUpdateRequest{
				Name: "New Name",
			},
			setupMocks: func() {
				mockUserRepo.EXPECT().GetByID(ctx, userID).Return(nil, expectations.ErrUserNotFound)
			},
			expectedErr: expectations.ErrUserNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMocks()

			err := profileService.UpdateProfileInfo(ctx, tt.userID, tt.updateData)

			if tt.expectedErr != nil {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedErr.Error())
			} else {
				require.NoError(t, err)
			}
		})
	}
}

// TestProfileService_UpdatePreferences тестирует обновление предпочтений
func TestProfileService_UpdatePreferences(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mocks.NewMockUserRepository(ctrl)
	mockUserPhotoRepo := mocks.NewMockUserPhotoRepository(ctrl)
	mockPreferenceRepo := mocks.NewMockUserPreferenceRepository(ctrl)
	mockFileStorage := mocks.NewMockFileStorage(ctrl)
	mockLogger := mocks.NewMockLogger()

	profileService := service.NewProfileService(
		mockUserRepo,
		mockUserPhotoRepo,
		mockPreferenceRepo,
		mockFileStorage,
		mockLogger,
	)

	ctx := context.Background()
	userID := uuid.New()

	tests := []struct {
		name        string
		userID      uuid.UUID
		updateData  *domain.PreferencesUpdateRequest
		setupMocks  func()
		expectedErr error
	}{
		{
			name:   "Success - update existing preferences",
			userID: userID,
			updateData: &domain.PreferencesUpdateRequest{
				ShowGender:   constants.GenderPreference("male"),
				AgeMin:       20,
				AgeMax:       35,
				MaxDistance:  25,
				GlobalSearch: true,
			},
			setupMocks: func() {
				existingPrefs := &domain.UserPreference{
					UserID:       userID,
					ShowGender:   constants.GenderPreference("female"),
					AgeMin:       18,
					AgeMax:       30,
					MaxDistance:  50,
					GlobalSearch: false,
					CreatedAt:    time.Now().Add(-time.Hour),
				}
				mockPreferenceRepo.EXPECT().GetByUserID(ctx, userID).Return(existingPrefs, nil)
				mockPreferenceRepo.EXPECT().Update(ctx, gomock.Any()).DoAndReturn(func(ctx context.Context, prefs *domain.UserPreference) error {
					assert.Equal(t, constants.Gender("male"), prefs.ShowGender)
					assert.Equal(t, 20, prefs.AgeMin)
					assert.Equal(t, 35, prefs.AgeMax)
					assert.Equal(t, 25, prefs.MaxDistance)
					assert.True(t, prefs.GlobalSearch)
					return nil
				})
			},
			expectedErr: nil,
		},
		{
			name:   "Success - create new preferences",
			userID: userID,
			updateData: &domain.PreferencesUpdateRequest{
				ShowGender:   constants.GenderPreference("female"),
				AgeMin:       22,
				AgeMax:       40,
				MaxDistance:  30,
				GlobalSearch: false,
			},
			setupMocks: func() {
				mockPreferenceRepo.EXPECT().GetByUserID(ctx, userID).Return(nil, nil)
				mockPreferenceRepo.EXPECT().Create(ctx, gomock.Any()).DoAndReturn(func(ctx context.Context, prefs *domain.UserPreference) error {
					assert.Equal(t, userID, prefs.UserID)
					assert.Equal(t, constants.Gender("female"), prefs.ShowGender)
					assert.Equal(t, 22, prefs.AgeMin)
					assert.Equal(t, 40, prefs.AgeMax)
					assert.Equal(t, 30, prefs.MaxDistance)
					assert.False(t, prefs.GlobalSearch)
					return nil
				})
				// Исправление: разрешаем вызов UpdateLastActive в любом количестве раз
				mockUserRepo.EXPECT().UpdateLastActive(gomock.Any(), gomock.Any()).AnyTimes()
			},
			expectedErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMocks()

			err := profileService.UpdatePreferences(ctx, tt.userID, tt.updateData)

			if tt.expectedErr != nil {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

// mockMultipartFile реализует интерфейс multipart.File для тестирования
type mockMultipartFile struct {
	reader *bytes.Reader
}

func (m *mockMultipartFile) Read(p []byte) (n int, err error) {
	return m.reader.Read(p)
}

func (m *mockMultipartFile) ReadAt(p []byte, off int64) (n int, err error) {
	return m.reader.ReadAt(p, off)
}

func (m *mockMultipartFile) Seek(offset int64, whence int) (int64, error) {
	return m.reader.Seek(offset, whence)
}

func (m *mockMultipartFile) Close() error {
	return nil
}

// mockFileHeader реализует интерфейс multipart.FileHeader для тестирования
type mockFileHeader struct {
	filename    string
	size        int64
	contentType string
	content     []byte
}

func (m *mockFileHeader) Open() (multipart.File, error) {
	return &mockMultipartFile{reader: bytes.NewReader(m.content)}, nil
}

func (m *mockFileHeader) Filename() string {
	return m.filename
}

func (m *mockFileHeader) Size() int64 {
	return m.size
}

func (m *mockFileHeader) Header() map[string][]string {
	return map[string][]string{
		"Content-Type": {m.contentType},
	}
}

// TestProfileService_UploadPhotos тестирует загрузку фотографий
func TestProfileService_UploadPhotos(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mocks.NewMockUserRepository(ctrl)
	mockUserPhotoRepo := mocks.NewMockUserPhotoRepository(ctrl)
	mockPreferenceRepo := mocks.NewMockUserPreferenceRepository(ctrl)
	mockFileStorage := mocks.NewMockFileStorage(ctrl)
	mockLogger := mocks.NewMockLogger()

	profileService := service.NewProfileService(
		mockUserRepo,
		mockUserPhotoRepo,
		mockPreferenceRepo,
		mockFileStorage,
		mockLogger,
	)

	ctx := context.Background()
	userID := uuid.New()

	// Создаем тестовые файлы с правильной реализацией
	createTestFile := func(filename string, size int, contentType string) *mockFileHeader {
		content := make([]byte, size)
		for i := range content {
			content[i] = 'A' // Заполняем тестовыми данными
		}
		return &mockFileHeader{
			filename:    filename,
			size:        int64(size),
			contentType: contentType,
			content:     content,
		}
	}

	tests := []struct {
		name           string
		userID         uuid.UUID
		photos         []*mockFileHeader
		setupMocks     func()
		expectedPhotos int
		expectedErr    error
	}{
		{
			name:   "Success - upload single photo",
			userID: userID,
			photos: []*mockFileHeader{
				createTestFile("test.jpg", 1024, "image/jpeg"),
			},
			setupMocks: func() {
				// Нет существующих фото
				mockUserPhotoRepo.EXPECT().GetByUserID(ctx, userID).Return([]domain.UserPhoto{}, nil)

				// Мокаем вызовы fileStorage
				mockFileStorage.EXPECT().Upload(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return("https://storage.example.com/photo.jpg", nil)

				mockUserPhotoRepo.EXPECT().Create(ctx, gomock.Any()).Return(nil)
			},
			expectedPhotos: 1,
			expectedErr:    nil,
		},
		{
			name:   "Error - photo limit exceeded",
			userID: userID,
			photos: []*mockFileHeader{
				createTestFile("test1.jpg", 1024, "image/jpeg"),
				createTestFile("test2.jpg", 1024, "image/jpeg"),
			},
			setupMocks: func() {
				// Уже есть максимальное количество фото
				existingPhotos := make([]domain.UserPhoto, constants.MaxPhotosPerUser)
				mockUserPhotoRepo.EXPECT().GetByUserID(ctx, userID).Return(existingPhotos, nil)
			},
			expectedPhotos: 0,
			expectedErr:    expectations.ErrPhotoLimitExceeded,
		},
		{
			name:   "Error - file too large",
			userID: userID,
			photos: []*mockFileHeader{
				createTestFile("large.jpg", int(constants.MaxPhotoSize)+1, "image/jpeg"),
			},
			setupMocks: func() {
				mockUserPhotoRepo.EXPECT().GetByUserID(ctx, userID).Return([]domain.UserPhoto{}, nil)
			},
			expectedPhotos: 0,
			expectedErr:    fmt.Errorf("file too large, max size is %dMB", constants.MaxPhotoSize/(1024*1024)),
		},
		{
			name:   "Error - invalid file type",
			userID: userID,
			photos: []*mockFileHeader{
				createTestFile("test.pdf", 1024, "application/pdf"),
			},
			setupMocks: func() {
				mockUserPhotoRepo.EXPECT().GetByUserID(ctx, userID).Return([]domain.UserPhoto{}, nil)
			},
			expectedPhotos: 0,
			expectedErr:    fmt.Errorf("invalid file type: %s, allowed: %s", "application/pdf", constants.AllowedMimeTypes),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMocks()

			// Конвертируем mockFileHeader в multipart.FileHeader для совместимости
			var fileHeaders []*multipart.FileHeader
			for _, photo := range tt.photos {
				fileHeaders = append(fileHeaders, photo)
			}

			result, err := profileService.UploadPhotos(ctx, tt.userID, fileHeaders)

			if tt.expectedErr != nil {
				require.Error(t, err)
				if tt.expectedErr != expectations.ErrPhotoLimitExceeded {
					assert.Contains(t, err.Error(), tt.expectedErr.Error())
				} else {
					assert.Equal(t, tt.expectedErr, err)
				}
			} else {
				require.NoError(t, err)
				assert.Len(t, result, tt.expectedPhotos)
				for _, photo := range result {
					assert.Equal(t, userID, photo.UserID)
					assert.True(t, photo.IsApproved)
				}
			}
		})
	}
}

// TestProfileService_DeletePhoto тестирует удаление фотографии
func TestProfileService_DeletePhoto(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mocks.NewMockUserRepository(ctrl)
	mockUserPhotoRepo := mocks.NewMockUserPhotoRepository(ctrl)
	mockPreferenceRepo := mocks.NewMockUserPreferenceRepository(ctrl)
	mockFileStorage := mocks.NewMockFileStorage(ctrl)
	mockLogger := mocks.NewMockLogger()

	profileService := service.NewProfileService(
		mockUserRepo,
		mockUserPhotoRepo,
		mockPreferenceRepo,
		mockFileStorage,
		mockLogger,
	)

	ctx := context.Background()
	userID := uuid.New()
	photoID := uuid.New()

	tests := []struct {
		name        string
		userID      uuid.UUID
		photoID     uuid.UUID
		setupMocks  func()
		expectedErr error
	}{
		{
			name:    "Success - delete photo",
			userID:  userID,
			photoID: photoID,
			setupMocks: func() {
				photo := &domain.UserPhoto{
					ID:       photoID,
					UserID:   userID,
					PhotoURL: "https://storage.example.com/photo1.jpg",
				}
				mockUserPhotoRepo.EXPECT().GetByID(ctx, photoID).Return(photo, nil)
				mockFileStorage.EXPECT().DeleteByURL(ctx, photo.PhotoURL).Return(nil)
				mockUserPhotoRepo.EXPECT().Delete(ctx, photoID).Return(nil)
				mockUserPhotoRepo.EXPECT().GetByUserID(ctx, userID).Return([]domain.UserPhoto{}, nil)
				mockUserPhotoRepo.EXPECT().UpdateDisplayOrder(ctx, userID, []domain.UserPhoto{}).Return(nil)
			},
			expectedErr: nil,
		},
		{
			name:    "Error - photo not found",
			userID:  userID,
			photoID: photoID,
			setupMocks: func() {
				mockUserPhotoRepo.EXPECT().GetByID(ctx, photoID).Return(nil, nil)
			},
			expectedErr: expectations.ErrPhotoNotFound,
		},
		{
			name:    "Error - photo not owned by user",
			userID:  userID,
			photoID: photoID,
			setupMocks: func() {
				photo := &domain.UserPhoto{
					ID:     photoID,
					UserID: uuid.New(), // Different user
				}
				mockUserPhotoRepo.EXPECT().GetByID(ctx, photoID).Return(photo, nil)
			},
			expectedErr: expectations.ErrPhotoNotOwned,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMocks()

			err := profileService.DeletePhoto(ctx, tt.userID, tt.photoID)

			if tt.expectedErr != nil {
				require.Error(t, err)
				assert.Equal(t, tt.expectedErr, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

// TestProfileService_Validation тестирует методы валидации
func TestProfileService_Validation(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mocks.NewMockUserRepository(ctrl)
	mockUserPhotoRepo := mocks.NewMockUserPhotoRepository(ctrl)
	mockPreferenceRepo := mocks.NewMockUserPreferenceRepository(ctrl)
	mockFileStorage := mocks.NewMockFileStorage(ctrl)
	mockLogger := mocks.NewMockLogger()

	profileService := service.NewProfileService(
		mockUserRepo,
		mockUserPhotoRepo,
		mockPreferenceRepo,
		mockFileStorage,
		mockLogger,
	)

	ctx := context.Background()

	t.Run("ValidateProfileUpdate", func(t *testing.T) {
		tests := []struct {
			name        string
			updateData  *domain.ProfileUpdateRequest
			expectedErr error
		}{
			{
				name: "Valid data",
				updateData: &domain.ProfileUpdateRequest{
					Name: "Valid Name",
					Bio:  stringPtr("Valid bio"),
				},
				expectedErr: nil,
			},
			{
				name: "Name too long",
				updateData: &domain.ProfileUpdateRequest{
					Name: strings.Repeat("a", constants.MaxNameLength+1),
				},
				expectedErr: fmt.Errorf("name too long, max %d characters", constants.MaxNameLength),
			},
			{
				name: "Bio too long",
				updateData: &domain.ProfileUpdateRequest{
					Bio: stringPtr(strings.Repeat("a", constants.MaxBioLength+1)),
				},
				expectedErr: fmt.Errorf("bio too long, max %d characters", constants.MaxBioLength),
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				err := profileService.ValidateProfileUpdate(ctx, tt.updateData)

				if tt.expectedErr != nil {
					require.Error(t, err)
					assert.Contains(t, err.Error(), tt.expectedErr.Error())
				} else {
					require.NoError(t, err)
				}
			})
		}
	})

	t.Run("ValidatePreferencesUpdate", func(t *testing.T) {
		tests := []struct {
			name        string
			updateData  *domain.PreferencesUpdateRequest
			expectedErr error
		}{
			{
				name: "Valid data",
				updateData: &domain.PreferencesUpdateRequest{
					AgeMin:      20,
					AgeMax:      30,
					MaxDistance: 25,
				},
				expectedErr: nil,
			},
			{
				name: "Age min greater than max",
				updateData: &domain.PreferencesUpdateRequest{
					AgeMin: 30,
					AgeMax: 20,
				},
				expectedErr: fmt.Errorf("minimum age must be less than maximum age"),
			},
			{
				name: "Distance too small",
				updateData: &domain.PreferencesUpdateRequest{
					MaxDistance: constants.MinDistance - 1,
				},
				expectedErr: fmt.Errorf("distance must be at least %d km", constants.MinDistance),
			},
			{
				name: "Distance too large",
				updateData: &domain.PreferencesUpdateRequest{
					MaxDistance: constants.MaxDistance + 1,
				},
				expectedErr: fmt.Errorf("distance cannot exceed %d km", constants.MaxDistance),
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				err := profileService.ValidatePreferencesUpdate(ctx, tt.updateData)

				if tt.expectedErr != nil {
					require.Error(t, err)
					assert.Contains(t, err.Error(), tt.expectedErr.Error())
				} else {
					require.NoError(t, err)
				}
			})
		}
	})
}

// Вспомогательная функция для создания указателя на строку
func stringPtr(s string) *string {
	return &s
}
