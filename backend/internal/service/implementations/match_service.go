package service

import (
	"fmt"

	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/domain/entities"
	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/domain/errors"
	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/repository/interfaces"
	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/service/interfaces"
	"github.com/google/uuid"
)

type matchService struct {
	matchRepo interfaces.MatchRepository
	userRepo  interfaces.UserRepository
	swipeRepo interfaces.SwipeRepository
}

func NewMatchService(
	matchRepo interfaces.MatchRepository,
	userRepo interfaces.UserRepository,
	swipeRepo interfaces.SwipeRepository,
) service.MatchService {
	return &matchService{
		matchRepo: matchRepo,
		userRepo:  userRepo,
		swipeRepo: swipeRepo,
	}
}

func (s *matchService) GetUserMatches(userID uuid.UUID, limit, offset int) (*service.MatchesResponse, error) {
	// Валидация параметров
	if limit <= 0 || limit > 50 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}

	// Получаем мэтчи пользователя
	matches, err := s.matchRepo.GetUserMatches(userID, limit, offset)
	if err != nil {
		fmt.Printf("Error retrieving matches: %v\n", err)
		return nil, err
	}

	// Собираем ID всех пользователей из мэтчей
	var userIDs []uuid.UUID
	for _, match := range matches {
		if match.User1ID == userID {
			userIDs = append(userIDs, match.User2ID)
		} else {
			userIDs = append(userIDs, match.User1ID)
		}
	}

	// Получаем информацию о пользователях
	users, err := s.userRepo.GetUsersByIDs(userIDs)
	if err != nil {
		return nil, err
	}

	// Создаем мапу для быстрого доступа
	userMap := make(map[uuid.UUID]domain.User)
	for _, user := range users {
		userMap[user.ID] = user
	}

	// Формируем ответ
	var matchResponses []service.MatchResponse
	for _, match := range matches {
		var matchedUser domain.User
		if match.User1ID == userID {
			matchedUser = userMap[match.User2ID]
		} else {
			matchedUser = userMap[match.User1ID]
		}

		matchResponses = append(matchResponses, service.MatchResponse{
			Match: match,
			User:  matchedUser,
		})
	}

	response := &service.MatchesResponse{
		Matches: matchResponses,
		Total:   len(matchResponses),
		Limit:   limit,
		Offset:  offset,
	}

	return response, nil
}

func (s *matchService) Unmatch(userID, targetUserID uuid.UUID) error {
	// Проверяем что мэтч существует
	match, err := s.matchRepo.GetByUsers(userID, targetUserID)
	if err != nil {
		return err
	}
	if match == nil {
		return errors.ErrMatchNotFound
	}

	// Деактивируем мэтч
	return s.matchRepo.UpdateActive(userID, targetUserID, false)
}
