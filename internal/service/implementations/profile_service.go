package service

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"regexp"

	"path/filepath"
	"strings"
	"time"

	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/domain/constants"
	domain "github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/domain/entities"
	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/domain/errors"
	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/logger"
	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/repository/interfaces"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type ProfileService struct {
	userRepo       interfaces.UserRepository
	userPhotoRepo  interfaces.UserPhotoRepository
	preferenceRepo interfaces.UserPreferenceRepository
	swipeRepo      interfaces.SwipeRepository
	matchRepo      interfaces.MatchRepository
	fileStorage    interfaces.FileStorage // Интерфейс для работы с файловым хранилищем
	logger         logger.Log
}

func NewProfileService(
	userRepo interfaces.UserRepository,
	userPhotoRepo interfaces.UserPhotoRepository,
	preferenceRepo interfaces.UserPreferenceRepository,
	swipeRepo interfaces.SwipeRepository,
	matchRepo interfaces.MatchRepository,
	fileStorage interfaces.FileStorage,
	l logger.Log,
) *ProfileService {
	return &ProfileService{
		userRepo:       userRepo,
		userPhotoRepo:  userPhotoRepo,
		preferenceRepo: preferenceRepo,
		swipeRepo:      swipeRepo,
		matchRepo:      matchRepo,
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

	// Получаем интересы
	interests, err := s.preferenceRepo.GetInterests(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Обновляем время последней активности
	go s.userRepo.UpdateLastActive(ctx, userID)

	return &domain.ProfileResponse{
		User:        user,
		Preferences: preferences,
		Photos:      photos,
		Interests:   interests,
	}, nil
}

// GetProfileWithRelations возвращает профиль пользователя с информацией о отношениях с viewer
func (s *ProfileService) GetProfileWithRelations(ctx context.Context, viewerID, targetID uuid.UUID) (*domain.ProfileResponse, error) {
	s.logger.Infof("GetProfileWithRelations called: viewerID=%s, targetID=%s", viewerID, targetID)

	// Получаем базовый профиль
	profile, err := s.GetProfile(ctx, targetID)
	if err != nil {
		return nil, err
	}

	// Проверяем, лайкнул ли target (пользователь из URL) viewer (пользователь из контекста)
	// swiper_user_id = targetID (тот, чей профиль смотрим)
	// target_user_id = viewerID (тот, кто смотрит)
	s.logger.Infof("Checking swipe: swiper_user_id=%s (targetID), target_user_id=%s (viewerID)", targetID, viewerID)
	swipe, err := s.swipeRepo.GetBySwiperAndTarget(ctx, targetID, viewerID)
	if err != nil {
		s.logger.Warnf("GetBySwiperAndTarget error: %v (swiper=%s, target=%s)", err, targetID, viewerID)
	} else {
		if swipe != nil {
			s.logger.Infof("Swipe found: swiper=%s, target=%s, type=%s", swipe.SwiperUserID, swipe.TargetUserID, swipe.SwipeType)
		} else {
			s.logger.Infof("No swipe found: swiper=%s, target=%s", targetID, viewerID)
		}
	}

	isLiked := false
	if err == nil && swipe != nil {
		// Проверяем, что это лайк или суперлайк
		isLiked = swipe.SwipeType == constants.SwipeTypeLike || swipe.SwipeType == constants.SwipeTypeSuperLike
		s.logger.Infof("isLiked calculated: %v (swipe_type=%s)", isLiked, swipe.SwipeType)
	} else {
		s.logger.Infof("isLiked set to false (err=%v, swipe=%v)", err, swipe != nil)
	}

	// Проверяем, есть ли матч между пользователями
	s.logger.Infof("Checking match: user1=%s, user2=%s", viewerID, targetID)
	match, err := s.matchRepo.GetByUsers(ctx, viewerID, targetID)
	if err != nil {
		s.logger.Warnf("GetByUsers error: %v (user1=%s, user2=%s)", err, viewerID, targetID)
	} else {
		if match != nil {
			s.logger.Infof("Match found: id=%s, user1=%s, user2=%s, is_active=%v", match.ID, match.User1ID, match.User2ID, match.IsActive)
		} else {
			s.logger.Infof("No match found: user1=%s, user2=%s", viewerID, targetID)
		}
	}

	isMatched := false
	if err == nil && match != nil && match.IsActive {
		isMatched = true
		s.logger.Infof("isMatched calculated: %v", isMatched)
	} else {
		s.logger.Infof("isMatched set to false (err=%v, match=%v, is_active=%v)", err, match != nil, match != nil && match.IsActive)
	}

	// Создаем расширенный ответ
	response := profile

	// Добавляем информацию об отношениях только если viewer != target
	if viewerID != targetID {
		// Используем указатели для опциональных полей
		response.IsLiked = &isLiked
		response.IsMatched = &isMatched
	}

	return response, nil
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
	// Проверка и обновление email с проверкой уникальности
	if updateData.Email != nil && *updateData.Email != "" {
		// Проверяем, изменился ли email
		if user.Email != *updateData.Email {
			// Проверяем, не занят ли новый email другим пользователем
			existingUser, err := s.userRepo.GetByEmail(ctx, *updateData.Email)
			if err != nil {
				return err
			}
			if existingUser != nil && existingUser.ID != userID {
				return errors.ErrUserAlreadyExists
			}
			user.Email = *updateData.Email
		}
	}
	if updateData.Name != "" {
		user.Name = updateData.Name
	}
	if updateData.Gender != "" {
		user.Gender = updateData.Gender
	}
	if updateData.Phone != nil {
		user.Phone = updateData.Phone
	}
	if updateData.BirthDate != nil {
		user.BirthDate = *updateData.BirthDate
	}
	if updateData.Bio != nil {
		user.Bio = updateData.Bio
	}
	if updateData.City != nil {
		user.City = updateData.City
	}
	if updateData.Artist != nil {
		user.Artist = updateData.Artist
	}
	if updateData.Quote != nil {
		user.Quote = updateData.Quote
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
	isNew := false
	if preferences == nil {
		s.logger.Debugf("Create preferences with empty user: %v", userID)
		isNew = true
		preferences = &domain.UserPreference{
			UserID:     userID,
			ShowGender: constants.GenderPrefMale,
			AgeMin:     constants.MinAge,
			AgeMax:     constants.MaxAge,
			// MaxDistance и GlobalSearch оставляем пустыми (не используются)
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
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
	preferences.UpdatedAt = time.Now()

	// Сохраняем
	if isNew {
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

// UploadPhotos загружает фотографии пользователя (Multipart)
func (s *ProfileService) UploadPhotos(ctx context.Context, userID uuid.UUID, photos []*multipart.FileHeader) ([]domain.UserPhoto, error) {
	s.logger.Infof("UploadPhotos called: userID=%s, photos_count=%d", userID, len(photos))

	// Проверяем лимит фотографий
	existingPhotos, err := s.userPhotoRepo.GetByUserID(ctx, userID)
	if err != nil {
		s.logger.Errorf("Failed to get existing photos: %v", err)
		return nil, err
	}
	s.logger.Infof("Existing photos count: %d, max allowed: %d", len(existingPhotos), constants.MaxPhotosPerUser)

	if len(existingPhotos)+len(photos) > constants.MaxPhotosPerUser {
		s.logger.Warnf("Photo limit exceeded: existing=%d, new=%d, max=%d",
			len(existingPhotos), len(photos), constants.MaxPhotosPerUser)
		return nil, errors.ErrPhotoLimitExceeded
	}

	var uploadedPhotos []domain.UserPhoto
	nextOrder := len(existingPhotos)

	for i, photoHeader := range photos {
		s.logger.Infof("Processing photo %d/%d: filename=%s, size=%d bytes (%.2f MB)",
			i+1, len(photos), photoHeader.Filename, photoHeader.Size, float64(photoHeader.Size)/(1024*1024))

		// Валидация файла
		s.logger.Debugf("Validating photo %d: filename=%s", i+1, photoHeader.Filename)
		if err := s.validatePhotoFile(ctx, photoHeader); err != nil {
			s.logger.Errorf("Photo validation failed for %s: %v", photoHeader.Filename, err)
			return nil, err
		}
		s.logger.Debugf("Photo %d validation passed", i+1)

		// Открываем файл
		s.logger.Debugf("Opening photo file %d: filename=%s", i+1, photoHeader.Filename)
		file, err := photoHeader.Open()
		if err != nil {
			s.logger.Errorf("Failed to open photo file %s: %v", photoHeader.Filename, err)
			return nil, fmt.Errorf("failed to open photo: %w", err)
		}

		// Читаем содержимое
		s.logger.Debugf("Reading photo file %d: filename=%s, size=%d bytes", i+1, photoHeader.Filename, photoHeader.Size)
		fileBytes, err := io.ReadAll(file)
		file.Close() // Close immediately after reading
		if err != nil {
			s.logger.Errorf("Failed to read photo file %s: %v", photoHeader.Filename, err)
			return nil, fmt.Errorf("failed to read photo: %w", err)
		}
		s.logger.Infof("Photo %d read successfully: filename=%s, bytes_read=%d (%.2f MB)",
			i+1, photoHeader.Filename, len(fileBytes), float64(len(fileBytes))/(1024*1024))

		// Загружаем
		contentType := photoHeader.Header.Get("Content-Type")
		ext := filepath.Ext(photoHeader.Filename)
		s.logger.Debugf("Uploading photo %d to storage: filename=%s, contentType=%s, ext=%s, order=%d",
			i+1, photoHeader.Filename, contentType, ext, nextOrder)

		userPhoto, err := s.uploadPhotoFromBytes(ctx, userID, fileBytes, contentType, ext, nextOrder)
		if err != nil {
			s.logger.Errorf("Failed to upload photo %d (%s): %v", i+1, photoHeader.Filename, err)
			return nil, err
		}

		s.logger.Infof("Photo %d uploaded successfully: id=%s, url=%s", i+1, userPhoto.ID, userPhoto.PhotoURL)
		uploadedPhotos = append(uploadedPhotos, *userPhoto)
		nextOrder++
	}

	s.logger.Infof("UploadPhotos completed: userID=%s, uploaded_count=%d", userID, len(uploadedPhotos))
	return uploadedPhotos, nil
}

// UploadPhoto загружает одну фотографию (bytes) - для gRPC
func (s *ProfileService) UploadPhoto(ctx context.Context, userID uuid.UUID, content []byte, contentType string) (*domain.UserPhoto, error) {
	// Проверяем лимит фотографий
	existingPhotos, err := s.userPhotoRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	if len(existingPhotos) >= constants.MaxPhotosPerUser {
		return nil, errors.ErrPhotoLimitExceeded
	}

	// Определяем расширение по Content-Type
	ext := ".jpg" // default
	if contentType == "image/png" {
		ext = ".png"
	} else if contentType == "image/jpeg" {
		ext = ".jpg"
	} else if contentType == "image/webp" {
		ext = ".webp"
	}

	// Валидация размера и типа (упрощенная для bytes)
	if len(content) > constants.MaxPhotoSize {
		return nil, fmt.Errorf("file too large")
	}
	// TODO: Validate mime type from bytes if needed, or trust contentType

	return s.uploadPhotoFromBytes(ctx, userID, content, contentType, ext, len(existingPhotos))
}

func (s *ProfileService) uploadPhotoFromBytes(ctx context.Context, userID uuid.UUID, content []byte, contentType, ext string, order int) (*domain.UserPhoto, error) {
	s.logger.Debugf("uploadPhotoFromBytes: userID=%s, content_size=%d bytes (%.2f MB), contentType=%s, ext=%s, order=%d",
		userID, len(content), float64(len(content))/(1024*1024), contentType, ext, order)

	// Генерируем уникальное имя файла
	fileName := fmt.Sprintf("%s/%s%s", userID.String(), uuid.New().String(), ext)
	s.logger.Debugf("Generated file name: %s", fileName)

	// Загружаем в файловое хранилище
	s.logger.Debugf("Uploading to storage: fileName=%s, size=%d bytes", fileName, len(content))
	photoURL, err := s.fileStorage.Upload(ctx, fileName, content, contentType)
	if err != nil {
		s.logger.Errorf("Storage upload failed: fileName=%s, error=%v", fileName, err)
		return nil, fmt.Errorf("failed to upload photo: %w", err)
	}
	s.logger.Infof("Storage upload successful: fileName=%s, photoURL=%s", fileName, photoURL)

	// Создаем запись в БД
	photoID := uuid.New()
	userPhoto := domain.UserPhoto{
		ID:           photoID,
		UserID:       userID,
		PhotoURL:     photoURL,
		DisplayOrder: order,
		IsApproved:   true, // Авто-аппрув для демо
		CreatedAt:    time.Now(),
	}
	s.logger.Debugf("Creating DB record: photoID=%s, userID=%s, photoURL=%s, order=%d",
		photoID, userID, photoURL, order)

	if err := s.userPhotoRepo.Create(ctx, &userPhoto); err != nil {
		s.logger.Errorf("DB record creation failed: photoID=%s, error=%v. Attempting to delete uploaded file", photoID, err)
		// Пытаемся удалить загруженный файл при ошибке
		if delErr := s.fileStorage.Delete(ctx, fileName); delErr != nil {
			s.logger.Errorf("Failed to delete uploaded file after DB error: fileName=%s, error=%v", fileName, delErr)
		} else {
			s.logger.Infof("Deleted uploaded file after DB error: fileName=%s", fileName)
		}
		return nil, err
	}

	s.logger.Infof("Photo uploaded and saved successfully: photoID=%s, userID=%s, photoURL=%s",
		photoID, userID, photoURL)
	return &userPhoto, nil
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
	// Валидация email формата
	if updateData.Email != nil && *updateData.Email != "" {
		emailRegex := regexp.MustCompile(`^[A-Z0-9._%+-]+@[A-Z0-9.-]+\.[A-Z]{2,}$`)
		if !emailRegex.MatchString(strings.ToUpper(*updateData.Email)) {
			return fmt.Errorf("invalid email format")
		}
	}

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
	// Валидация пола
	if updateData.ShowGender != "" {
		validGenders := map[constants.GenderPreference]struct{}{
			constants.GenderPrefBoth:   {},
			constants.GenderPrefMale:   {},
			constants.GenderPrefFemale: {},
		}
		if _, ok := validGenders[updateData.ShowGender]; !ok {
			return fmt.Errorf("invalid gender preference: %s (allowed: both, male, female)", updateData.ShowGender)
		}
	}

	// Валидация возраста
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
	s.logger.Debugf("validatePhotoFile: filename=%s, size=%d bytes (%.2f MB), max_allowed=%d bytes (%.2f MB)",
		photo.Filename, photo.Size, float64(photo.Size)/(1024*1024), constants.MaxPhotoSize, float64(constants.MaxPhotoSize)/(1024*1024))

	// Проверка размера файла
	if photo.Size > constants.MaxPhotoSize {
		s.logger.Warnf("File size validation failed: filename=%s, size=%d bytes (%.2f MB), max=%d bytes (%.2f MB)",
			photo.Filename, photo.Size, float64(photo.Size)/(1024*1024), constants.MaxPhotoSize, float64(constants.MaxPhotoSize)/(1024*1024))
		return fmt.Errorf("file too large, max size is %dMB", constants.MaxPhotoSize/(1024*1024))
	}
	s.logger.Debugf("File size validation passed: filename=%s, size=%d bytes", photo.Filename, photo.Size)

	// Проверка MIME типа
	mimeType := photo.Header.Get("Content-Type")
	s.logger.Debugf("MIME type check: filename=%s, mimeType=%s, allowedTypes=%s",
		photo.Filename, mimeType, constants.AllowedMimeTypes)

	allowedTypes := strings.Split(constants.AllowedMimeTypes, ",")
	valid := false
	for _, allowedType := range allowedTypes {
		trimmed := strings.TrimSpace(allowedType)
		if mimeType == trimmed {
			valid = true
			s.logger.Debugf("MIME type validation passed: filename=%s, mimeType=%s matches allowed=%s",
				photo.Filename, mimeType, trimmed)
			break
		}
	}
	if !valid {
		s.logger.Warnf("MIME type validation failed: filename=%s, mimeType=%s, allowedTypes=%s",
			photo.Filename, mimeType, constants.AllowedMimeTypes)
		return fmt.Errorf("invalid file type: %s, allowed: %s", mimeType, constants.AllowedMimeTypes)
	}

	s.logger.Debugf("Photo validation completed successfully: filename=%s", photo.Filename)
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

// ChangePassword изменяет пароль пользователя
func (s *ProfileService) ChangePassword(ctx context.Context, userID uuid.UUID, oldPassword, newPassword, newPasswordConfirm string) error {
	s.logger.Infof("ChangePassword called: userID=%s", userID)

	// Проверяем что новый пароль и подтверждение совпадают
	if newPassword != newPasswordConfirm {
		return errors.ErrPasswordsDontMatch
	}

	// Проверяем минимальную длину нового пароля
	if len(newPassword) < 6 {
		return errors.ErrPasswordTooShort
	}

	// Получаем текущего пользователя
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		s.logger.Errorf("Failed to get user: %v", err)
		return err
	}
	if user == nil {
		return errors.ErrProfileNotFound
	}

	// Проверяем старый пароль
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(oldPassword)); err != nil {
		s.logger.Warnf("Old password verification failed for userID=%s", userID)
		return fmt.Errorf("old password is incorrect")
	}

	// Хешируем новый пароль
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		s.logger.Errorf("Failed to hash new password: %v", err)
		return fmt.Errorf("failed to hash password: %w", err)
	}

	// Обновляем пароль
	user.Password = string(hashedPassword)
	user.UpdatedAt = time.Now()

	if err := s.userRepo.Update(ctx, user); err != nil {
		s.logger.Errorf("Failed to update user password: %v", err)
		return err
	}

	s.logger.Infof("Password changed successfully for userID=%s", userID)
	return nil
}

// DeleteAccount удаляет аккаунт пользователя
func (s *ProfileService) DeleteAccount(ctx context.Context, userID uuid.UUID) error {
	s.logger.Infof("DeleteAccount called: userID=%s", userID)

	// Проверяем что пользователь существует
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		s.logger.Errorf("Failed to get user: %v", err)
		return err
	}
	if user == nil {
		return errors.ErrProfileNotFound
	}

	// Удаляем фотографии из хранилища
	photos, err := s.userPhotoRepo.GetByUserID(ctx, userID)
	if err != nil {
		s.logger.Warnf("Failed to get user photos: %v", err)
		// Продолжаем удаление даже если не удалось получить фото
	} else {
		for _, photo := range photos {
			if err := s.fileStorage.DeleteByURL(ctx, photo.PhotoURL); err != nil {
				s.logger.Errorf("Failed to delete photo file: %v", err)
				// Продолжаем даже если не удалось удалить файл
			}
		}
	}

	// Удаляем пользователя (каскадное удаление в БД должно удалить связанные данные)
	if err := s.userRepo.Delete(ctx, userID); err != nil {
		s.logger.Errorf("Failed to delete user: %v", err)
		return err
	}

	s.logger.Infof("Account deleted successfully for userID=%s", userID)
	return nil
}
