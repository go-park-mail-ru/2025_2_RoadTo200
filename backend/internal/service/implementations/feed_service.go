package service

import (
	"fmt"
	"time"

	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/domain/entities"
	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/logger"
	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/repository/interfaces"
	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/service/interfaces"
	"github.com/google/uuid"
)

type feedService struct {
	userRepo  interfaces.UserRepository
	photoRepo interfaces.UserPhotoRepository
	logger    *logger.Logger
}

func NewFeedService(
	userRepo interfaces.UserRepository,
	photoRepo interfaces.UserPhotoRepository,
	l *logger.Logger,
) service.FeedService {
	return &feedService{
		userRepo:  userRepo,
		photoRepo: photoRepo,
		logger:    l,
	}
}

func (s *feedService) GetFeed(userID uuid.UUID, limit, offset int) ([]service.FeedUser, error) {
	// Валидация параметров
	if limit <= 0 || limit > 50 {
		limit = 15 // дефолтное значение
	}
	if offset < 0 {
		offset = 0
	}

	s.logger.Debugf("Getting feed for userID: %v with limit: %d and offset: %d\n", userID, limit, offset)

	// Получаем пользователей для ленты
	users, err := s.userRepo.GetUsersForFeed(userID, limit, offset)
	if err != nil {
		s.logger.Warnf("Getting users for feed failed: %v", err)
		return nil, err
	}
	s.logger.Warnf("Retrieved %d users for feed\n", len(users))

	// Преобразуем в формат для ленты
	feedUsers := make([]service.FeedUser, 0, len(users))
	for _, user := range users {
		feedUser, err := s.convertToFeedUser(user)
		if err != nil {
			s.logger.Warnf("Error converting user %s: %v\n", user.ID, err)
			continue // Пропускаем пользователя с ошибкой
		}
		feedUsers = append(feedUsers, feedUser)
	}

	// Обновляем время последней активности текущего пользователя
	go s.userRepo.UpdateLastActive(userID)

	return feedUsers, nil
}

// convertToFeedUser преобразует доменного пользователя в формат для ленты
func (s *feedService) convertToFeedUser(user domain.User) (service.FeedUser, error) {
	// Вычисляем возраст
	age := calculateAge(user.BirthDate)

	// Получаем фото пользователя
	images, err := s.getUserPhotos(user.ID)
	if err != nil {
		return service.FeedUser{}, err
	}

	return service.FeedUser{
		ID:          user.ID.String(),
		Name:        user.Name,
		Age:         age,
		Gender:      string(user.Gender),
		Description: getDescription(user.Bio),
		Images:      images,
		PhotosCount: len(images),
	}, nil
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
func (s *feedService) getUserPhotos(userID uuid.UUID) ([]string, error) {
	photos, err := s.photoRepo.GetByUserID(userID)
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
