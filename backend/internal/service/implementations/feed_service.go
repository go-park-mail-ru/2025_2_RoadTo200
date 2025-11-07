package service

import (
	"fmt"

	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/domain/entities"
	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/logger"
	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/repository/interfaces"
	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/service/interfaces"
	"github.com/google/uuid"
)

type feedService struct {
	userRepo interfaces.UserRepository
	logger   *logger.Logger
}

func NewFeedService(userRepo interfaces.UserRepository, l *logger.Logger) service.FeedService {
	return &feedService{
		userRepo: userRepo,
		logger:   l,
	}
}

func (s *feedService) GetFeed(userID uuid.UUID, limit, offset int) ([]domain.User, error) {
	// Валидация параметров
	if limit <= 0 || limit > 50 {
		limit = 15 // дефолтное значение
	}
	if offset < 0 {
		offset = 0
	}

	fmt.Printf("Getting feed for userID: %v with limit: %d and offset: %d\n", userID, limit, offset)

	// Получаем пользователей для ленты
	users, err := s.userRepo.GetUsersForFeed(userID, limit, offset)
	if err != nil {
		s.logger.Warnf("Getting users for feed failed: %v", err)
		return nil, err
	}
	s.logger.Warnf("Retrieved %d users for feed\n", len(users))

	// Обновляем время последней активности текущего пользователя
	go s.userRepo.UpdateLastActive(userID)

	return users, nil
}
