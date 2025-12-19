package service

import (
	"context"
	"fmt"
	"time"

	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/domain/constants"
	domain "github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/domain/entities"
	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/handler/dto"
	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/logger"
	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/repository/interfaces"
	"github.com/google/uuid"
)

type FeedService struct {
	userRepo  interfaces.UserRepository
	prefRepo  interfaces.UserPreferenceRepository
	photoRepo interfaces.UserPhotoRepository
	logger    logger.Log
}

func NewFeedService(
	userRepo interfaces.UserRepository,
	prefRepo interfaces.UserPreferenceRepository,
	photoRepo interfaces.UserPhotoRepository,
	l logger.Log,
) *FeedService {
	return &FeedService{
		userRepo:  userRepo,
		prefRepo:  prefRepo,
		photoRepo: photoRepo,
		logger:    l,
	}
}

func (s *FeedService) GetFeed(ctx context.Context, userID uuid.UUID, limit, offset int) ([]dto.FeedUser, error) {
	s.logger.Trace("FeedService.GetFeed")
	// Валидация параметров
	if limit <= 0 || limit > 50 {
		limit = 15 // дефолтное значение
	}
	if offset < 0 {
		offset = 0
	}

	s.logger.Infof("Getting feed for userID: %v with limit: %d and offset: %d", userID, limit, offset)

	// Получаем предпочтения пользователя для отладки
	prefs, err := s.prefRepo.GetByUserID(ctx, userID)
	if err != nil {
		s.logger.Warnf("Failed to get user preferences for feed: %v", err)
	} else if prefs != nil {
		s.logger.Infof("User preferences: show_gender=%s, age_min=%d, age_max=%d", prefs.ShowGender, prefs.AgeMin, prefs.AgeMax)
	} else {
		s.logger.Warnf("User has no preferences set")
	}

	// Получаем пользователей для ленты
	users, err := s.userRepo.GetUsersForFeed(ctx, userID, limit, offset)
	if err != nil {
		s.logger.Errorf("Getting users for feed failed: %v", err)
		return nil, err
	}
	s.logger.Infof("Retrieved %d users for feed", len(users))

	// Преобразуем в формат для ленты
	feedUsers := make([]dto.FeedUser, 0, len(users))
	for _, user := range users {
		res, err := s.prefRepo.GetInterests(ctx, user.ID)
		if err != nil {
			s.logger.Warnf("Getting interests for user failed: %v", err)
			continue
		}
		feedUser, err := s.convertToFeedUser(ctx, user)
		if err != nil {
			s.logger.Warnf("Error converting user %s: %v\n", user.ID, err)
			continue // Пропускаем пользователя с ошибкой
		}
		feedUser.Interests = res
		feedUsers = append(feedUsers, feedUser)
	}

	// Обновляем время последней активности текущего пользователя
	go s.userRepo.UpdateLastActive(ctx, userID)

	return feedUsers, nil
}

// convertToFeedUser преобразует доменного пользователя в формат для ленты
func (s *FeedService) convertToFeedUser(ctx context.Context, user domain.User) (dto.FeedUser, error) {
	s.logger.Trace("convertToFeedUser")

	// Вычисляем возраст
	age := calculateAge(user.BirthDate)

	// Получаем фото пользователя
	images, err := s.getUserPhotos(ctx, user.ID)
	if err != nil {
		return dto.FeedUser{}, err
	}

	return dto.FeedUser{
		ID:          user.ID.String(),
		Name:        user.Name,
		Age:         age,
		Gender:      getGenderString(user.Gender),
		Description: getDescription(user.Bio),
		Images:      images,
		PhotosCount: len(images),
		Artist:      user.Artist,
		Quote:       user.Quote,
		IsPremium:   user.IsPremium,
	}, nil
}

// getGenderString возвращает строковое представление gender
func getGenderString(gender constants.Gender) string {
	if gender != "" {
		return string(gender)
	}
	return "not_specified"
}

// calculateAge вычисляет возраст по дате рождения
func calculateAge(birthDate time.Time) int {
	now := time.Now()
	age := now.Year() - birthDate.Year()

	// Проверяем, был ли уже день рождения в этом году
	if now.YearDay() < birthDate.YearDay() {
		age--
	}

	return age
}

// getDescription возвращает описание или заглушку
func getDescription(bio *string) string {
	if bio != nil && *bio != "" {
		return *bio
	}
	return "Нет описания"
}

// getUserPhotos возвращает фото пользователя из репозитория
func (s *FeedService) getUserPhotos(ctx context.Context, userID uuid.UUID) ([]string, error) {
	photos, err := s.photoRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get photos for user %s: %w", userID, err)
	}

	urls := make([]string, 0, len(photos))
	for _, photo := range photos {
		// Добавляем только approved фото
		if photo.IsApproved {
			urls = append(urls, photo.PhotoURL)
		}
	}

	return urls, nil
}
