package service

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"

	"path/filepath"
	"strings"
	"time"

	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/domain/constants"
	domain "github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/domain/entities"
	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/domain/errors"
	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/logger"
	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/repository/interfaces"
	"github.com/google/uuid"
)

type ProfileService struct {
	userRepo       interfaces.UserRepository
	userPhotoRepo  interfaces.UserPhotoRepository
	preferenceRepo interfaces.UserPreferenceRepository
	fileStorage    interfaces.FileStorage // Интерфейс для работы с файловым хранилищем
	logger         logger.Log
}

func NewProfileService(
	userRepo interfaces.UserRepository,
	userPhotoRepo interfaces.UserPhotoRepository,
	preferenceRepo interfaces.UserPreferenceRepository,
	fileStorage interfaces.FileStorage,
	l logger.Log,
) *ProfileService {
	return &ProfileService{
		userRepo:       userRepo,
		userPhotoRepo:  userPhotoRepo,
		preferenceRepo: preferenceRepo,
		fileStorage:    fileStorage,
		logger:         l,
	}
}

// GetProfile возвращает полный профиль пользователя
func (s *ProfileService) GetProfile(ctx context.Context, userID uuid.UUID) (*domain.ProfileResponse, error) {
	// Получаем основную информацию пользователя
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.ErrProfileNotFound
	}

	// Получаем предпочтения
	preferences, err := s.preferenceRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Получаем фотографии
	photos, err := s.userPhotoRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Обновляем время последней активности
	go s.userRepo.UpdateLastActive(ctx, userID)

	return &domain.ProfileResponse{
		User:        user,
		Preferences: preferences,
		Photos:      photos,
	}, nil
}

// UpdateProfileInfo обновляет основную информацию профиля
func (s *ProfileService) UpdateProfileInfo(ctx context.Context, userID uuid.UUID, updateData *domain.ProfileUpdateRequest) error {
	// Валидация данных
	if err := s.ValidateProfileUpdate(ctx, updateData); err != nil {
		return err
	}

	// Получаем текущего пользователя
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return err
	}
	if user == nil {
		return errors.ErrProfileNotFound
	}

	// Применяем изменения
	if updateData.Name != "" {
		user.Name = updateData.Name
	}
	if updateData.Gender != "" {
		user.Gender = updateData.Gender
	}

	user.UpdatedAt = time.Now()

	// Сохраняем изменения
	if err := s.userRepo.Update(ctx, user); err != nil {
		return err
	}

	return nil
}

// UpdatePreferences обновляет предпочтения пользователя
func (s *ProfileService) UpdatePreferences(ctx context.Context, userID uuid.UUID, updateData *domain.PreferencesUpdateRequest) error {
	// Валидация
	if err := s.ValidatePreferencesUpdate(ctx, updateData); err != nil {
		return err
	}

	// Получаем текущие предпочтения
	preferences, err := s.preferenceRepo.GetByUserID(ctx, userID)
	if err != nil {
		return err
	}

	// Если предпочтений нет - создаем новые
	if preferences == nil {
		s.logger.Debugf("Create preferences with empty user: %v", userID)
		preferences = &domain.UserPreference{
			UserID:       userID,
			ShowGender:   constants.GenderPrefMale,
			AgeMin:       constants.MinAge,
			AgeMax:       constants.MaxAge,
			MaxDistance:  100,
			GlobalSearch: false,
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		}
	}

	// Применяем изменения
	if updateData.ShowGender != "" {
		preferences.ShowGender = updateData.ShowGender
	}
	if updateData.AgeMin > 0 {
		preferences.AgeMin = updateData.AgeMin
	}
	if updateData.AgeMax > 0 {
		preferences.AgeMax = updateData.AgeMax
	}
	if updateData.MaxDistance > 0 {
		preferences.MaxDistance = updateData.MaxDistance
	}
	preferences.GlobalSearch = updateData.GlobalSearch
	preferences.UpdatedAt = time.Now()

	// Сохраняем
	if preferences.CreatedAt.IsZero() {
		return s.preferenceRepo.Create(ctx, preferences)
	}
	return s.preferenceRepo.Update(ctx, preferences)
}

func (s *ProfileService) UpdateInterests(ctx context.Context, userID uuid.UUID, inter []domain.Interest) error {
	if err := s.validateInterestsUpdate(ctx, inter); err != nil {
		return err
	}
	return s.preferenceRepo.UpdateInterests(ctx, userID, inter)
}

// UploadPhotos загружает фотографии пользователя
func (s *ProfileService) UploadPhotos(ctx context.Context, userID uuid.UUID, photos []*multipart.FileHeader) ([]domain.UserPhoto, error) {
	// Проверяем лимит фотографий
	existingPhotos, err := s.userPhotoRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	if len(existingPhotos)+len(photos) > constants.MaxPhotosPerUser {
		return nil, errors.ErrPhotoLimitExceeded
	}

	var uploadedPhotos []domain.UserPhoto
	nextOrder := len(existingPhotos)

	for _, photoHeader := range photos {
		// Валидация файла
		if err := s.validatePhotoFile(ctx, photoHeader); err != nil {
			return nil, err
		}

		// Открываем файл
		file, err := photoHeader.Open()
		if err != nil {
			return nil, fmt.Errorf("failed to open photo: %w", err)
		}
		defer file.Close()

		// Читаем содержимое
		fileBytes, err := io.ReadAll(file)
		if err != nil {
			return nil, fmt.Errorf("failed to read photo: %w", err)
		}

		// Генерируем уникальное имя файла
		fileExt := filepath.Ext(photoHeader.Filename)
		fileName := fmt.Sprintf("%s/%s%s", userID.String(), uuid.New().String(), fileExt)

		// Загружаем в файловое хранилище
		photoURL, err := s.fileStorage.Upload(ctx, fileName, fileBytes, photoHeader.Header.Get("Content-Type"))
		if err != nil {
			return nil, fmt.Errorf("failed to upload photo: %w", err)
		}

		// Создаем запись в БД
		userPhoto := domain.UserPhoto{
			ID:           uuid.New(),
			UserID:       userID,
			PhotoURL:     photoURL,
			DisplayOrder: nextOrder,
			IsApproved:   true, // Авто-аппрув для демо
			CreatedAt:    time.Now(),
		}

		if err := s.userPhotoRepo.Create(ctx, &userPhoto); err != nil {
			// Пытаемся удалить загруженный файл при ошибке
			s.fileStorage.Delete(ctx, fileName)
			return nil, err
		}

		uploadedPhotos = append(uploadedPhotos, userPhoto)
		nextOrder++
	}

	return uploadedPhotos, nil
}

// DeletePhoto удаляет фотографию пользователя
func (s *ProfileService) DeletePhoto(ctx context.Context, userID uuid.UUID, photoID uuid.UUID) error {
	// Проверяем что фото принадлежит пользователю
	photo, err := s.userPhotoRepo.GetByID(ctx, photoID)
	if err != nil {
		return err
	}
	if photo == nil {
		return errors.ErrPhotoNotFound
	}
	if photo.UserID != userID {
		return errors.ErrPhotoNotOwned
	}

	// Удаляем файл из хранилища
	if err := s.fileStorage.DeleteByURL(ctx, photo.PhotoURL); err != nil {
		// Логируем ошибку, но продолжаем удаление записи из БД
		s.logger.Errorf("Failed to delete photo file: %v\n", err)
	}

	// Удаляем запись из БД
	if err := s.userPhotoRepo.Delete(ctx, photoID); err != nil {
		return err
	}

	// Обновляем порядок оставшихся фото
	return s.reorderRemainingPhotos(ctx, userID)
}

// SetPrimaryPhoto устанавливает фото как основное (первое в порядке)
func (s *ProfileService) SetPrimaryPhoto(ctx context.Context, userID uuid.UUID, photoID uuid.UUID) error {
	// Проверяем что фото принадлежит пользователю
	photo, err := s.userPhotoRepo.GetByID(ctx, photoID)
	if err != nil {
		return err
	}
	if photo == nil {
		return errors.ErrPhotoNotFound
	}
	if photo.UserID != userID {
		return errors.ErrPhotoNotOwned
	}

	// Получаем все фото пользователя
	allPhotos, err := s.userPhotoRepo.GetByUserID(ctx, userID)
	if err != nil {
		return err
	}

	// Создаем новый порядок: выбранное фото становится первым, остальные сдвигаются
	var newPhotos []domain.UserPhoto
	newPhotos = append(newPhotos, *photo) // Выбранное фото становится первым

	for _, p := range allPhotos {
		if p.ID != photoID {
			newPhotos = append(newPhotos, p)
		}
	}

	// Обновляем порядок отображения
	for i := range newPhotos {
		newPhotos[i].DisplayOrder = i
	}

	return s.userPhotoRepo.UpdateDisplayOrder(ctx, userID, newPhotos)
}

// ReorderPhotos изменяет порядок фотографий
func (s *ProfileService) ReorderPhotos(ctx context.Context, userID uuid.UUID, photoIDs []uuid.UUID) error {
	// Получаем все фото пользователя для проверки принадлежности
	allPhotos, err := s.userPhotoRepo.GetByUserID(ctx, userID)
	if err != nil {
		return err
	}

	// Проверяем что все photoIDs принадлежат пользователю
	userPhotoMap := make(map[uuid.UUID]domain.UserPhoto)
	for _, photo := range allPhotos {
		userPhotoMap[photo.ID] = photo
	}

	var photosToUpdate []domain.UserPhoto
	for order, photoID := range photoIDs {
		photo, exists := userPhotoMap[photoID]
		if !exists {
			return errors.ErrPhotoNotOwned
		}
		photo.DisplayOrder = order
		photosToUpdate = append(photosToUpdate, photo)
	}

	return s.userPhotoRepo.UpdateDisplayOrder(ctx, userID, photosToUpdate)
}

// ValidateProfileUpdate валидация данных профиля
func (s *ProfileService) ValidateProfileUpdate(ctx context.Context, updateData *domain.ProfileUpdateRequest) error {
	if updateData.Name != "" && len(updateData.Name) > constants.MaxNameLength {
		return fmt.Errorf("name too long, max %d characters", constants.MaxNameLength)
	}

	if updateData.Bio != nil && len(*updateData.Bio) > constants.MaxBioLength {
		return fmt.Errorf("bio too long, max %d characters", constants.MaxBioLength)
	}

	if updateData.BirthDate != nil {
		age := time.Since(*updateData.BirthDate).Hours() / 24 / 365
		if age < float64(constants.MinAge) {
			return fmt.Errorf("user must be at least %d years old", constants.MinAge)
		}
	}

	return nil
}

// ValidatePreferencesUpdate валидация предпочтений
func (s *ProfileService) ValidatePreferencesUpdate(ctx context.Context, updateData *domain.PreferencesUpdateRequest) error {
	if updateData.AgeMin > 0 && updateData.AgeMax > 0 {
		if updateData.AgeMin < constants.MinAge {
			return fmt.Errorf("minimum age must be at least %d", constants.MinAge)
		}
		if updateData.AgeMax > constants.MaxAge {
			return fmt.Errorf("maximum age cannot exceed %d", constants.MaxAge)
		}
		if updateData.AgeMin >= updateData.AgeMax {
			return fmt.Errorf("minimum age must be less than maximum age")
		}
	}

	if updateData.MaxDistance > 0 {
		if updateData.MaxDistance < constants.MinDistance {
			return fmt.Errorf("distance must be at least %d km", constants.MinDistance)
		}
		if updateData.MaxDistance > constants.MaxDistance {
			return fmt.Errorf("distance cannot exceed %d km", constants.MaxDistance)
		}
	}

	return nil
}

// ValidateInterestsUpdat валидация интересов
func (s *ProfileService) validateInterestsUpdate(ctx context.Context, updateData []domain.Interest) error {
	validTypes := map[constants.InterestType]struct{}{
		constants.InterestTypeWorkout:    {},
		constants.InterestTypeFun:        {},
		constants.InterestTypeParty:      {},
		constants.InterestTypeChill:      {},
		constants.InterestTypeLove:       {},
		constants.InterestTypeRelax:      {},
		constants.InterestTypeYoga:       {},
		constants.InterestTypeFriendship: {},
		constants.InterestTypeCulture:    {},
		constants.InterestTypeCinema:     {},
	}
	for _, el := range updateData {
		// Проверяем, что тип интереса допустим
		if _, ok := validTypes[el.Theme]; !ok {
			return fmt.Errorf("invalid interest type: %s", el.Theme)
		}
	}
	return nil
}

// validatePhotoFile валидация загружаемого фото
func (s *ProfileService) validatePhotoFile(ctx context.Context, photo *multipart.FileHeader) error {
	// Проверка размера файла
	if photo.Size > constants.MaxPhotoSize {
		return fmt.Errorf("file too large, max size is %dMB", constants.MaxPhotoSize/(1024*1024))
	}

	// Проверка MIME типа
	mimeType := photo.Header.Get("Content-Type")
	allowedTypes := strings.Split(constants.AllowedMimeTypes, ",")
	valid := false
	for _, allowedType := range allowedTypes {
		if mimeType == strings.TrimSpace(allowedType) {
			valid = true
			break
		}
	}
	if !valid {
		return fmt.Errorf("invalid file type: %s, allowed: %s", mimeType, constants.AllowedMimeTypes)
	}

	return nil
}

// reorderRemainingPhotos обновляет порядок фото после удаления
func (s *ProfileService) reorderRemainingPhotos(ctx context.Context, userID uuid.UUID) error {
	photos, err := s.userPhotoRepo.GetByUserID(ctx, userID)
	if err != nil {
		return err
	}

	for i := range photos {
		photos[i].DisplayOrder = i
	}

	return s.userPhotoRepo.UpdateDisplayOrder(ctx, userID, photos)
}
