package service

import (
	"fmt"

	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/domain/entities"
	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/repository/interfaces"
	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/service/interfaces"
	"github.com/google/uuid"
)

type feedService struct {
	userRepo interfaces.UserRepository
}

func NewFeedService(userRepo interfaces.UserRepository) service.FeedService {
	return &feedService{
		userRepo: userRepo,
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
	fmt.Printf("Retrieved %d users for feed\n", len(users))
	if err != nil {
		fmt.Printf("Error retrieving users for feed: %v\n", err)
		return nil, err
	}

	// Обновляем время последней активности текущего пользователя
	go s.userRepo.UpdateLastActive(userID)

	return users, nil
}
