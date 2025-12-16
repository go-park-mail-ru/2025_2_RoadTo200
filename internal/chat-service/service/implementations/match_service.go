package service

import (
	"context"
	"fmt"

	_ "github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/chat-service/domain/entities"
	interfaces2 "github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/chat-service/repository/interfaces"
	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/chat-service/service/interfaces"
	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/core-service/repository/interfaces"
	service2 "github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/core-service/service/interfaces"
	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/domain/errors"
	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/handler/dto"
	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/pkg/logger"
	"github.com/google/uuid"
)

var _ service.MatchService = (*MatchService)(nil)

type MatchService struct {
	matchRepo interfaces2.MatchRepository
	swipeRepo interfaces.SwipeRepository
	photoRepo interfaces.UserPhotoRepository // Добавляем репозиторий фотографий
	logger    logger.Log
}

func NewMatchService(
	matchRepo interfaces2.MatchRepository,
	swipeRepo interfaces.SwipeRepository,
	photoRepo interfaces.UserPhotoRepository, // Добавляем параметр
	l logger.Log,
) *MatchService {
	return &MatchService{
		matchRepo: matchRepo,
		swipeRepo: swipeRepo,
		photoRepo: photoRepo, // Инициализируем
		logger:    l,
	}
}

func (s *MatchService) GetUserMatches(ctx context.Context, userID uuid.UUID, limit, offset int) (*dto.MatchesResponse, error) {
	// Валидация параметров
	if limit <= 0 || limit > 50 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}

	// Получаем мэтчи пользователя
	matches, err := s.matchRepo.GetUserMatches(ctx, userID, limit, offset)
	if err != nil {
		return nil, err
	}

	// Если нет мэтчей, возвращаем пустой ответ
	if len(matches) == 0 {
		s.logger.Debugf("matchService.GetUserMatches: no matches")
		return &dto.MatchesResponse{
			Matches: []dto.MatchResponse{},
			Total:   0,
			Limit:   limit,
			Offset:  offset,
		}, nil
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
	//users, err := s.userRepo.GetUsersByIDs(ctx, userIDs)
	//if err != nil {
	//	return nil, err
	//}

	// Получаем фотографии для всех пользователей
	userPhotos, err := s.getUsersPhotos(ctx, userIDs)
	if err != nil {
		return nil, err
	}

	// Создаем мапу для быстрого доступа
	//userMap := make(map[uuid.UUID]domain.User)
	//for _, user := range users {
	//	userMap[user.ID] = user
	//}

	// Формируем ответ
	var matchResponses []dto.MatchResponse
	for _, match := range matches {
		var ok bool
		var matchedUserID uuid.UUID

		if match.User1ID == userID {
			matchedUserID = match.User2ID
		} else {
			matchedUserID = match.User1ID
		}

		if !ok {
			s.logger.Warnf("User not found for match: %v\n", match)
			continue
		}

		// Получаем фотографии для этого пользователя
		photos := userPhotos[matchedUserID]

		//// Вычисляем возраст
		//age := service2.calculateAge(matchedUser.BirthDate)
		//
		//// Получаем описание
		//description := service2.getDescription(matchedUser.Bio)

		matchResponses = append(matchResponses, dto.MatchResponse{
			Match:       match,
			User:        matchedUserID,
			Photos:      photos,      // Добавляем фотографии
			PhotosCount: len(photos), // Добавляем количество фото
		})
	}

	response := &dto.MatchesResponse{
		Matches: matchResponses,
		Total:   len(matchResponses),
		Limit:   limit,
		Offset:  offset,
	}

	s.logger.Debugf("Successfully formed response with %d matches\n", len(matchResponses))
	return response, nil
}

// getUsersPhotos возвращает фотографии для списка пользователей
func (s *MatchService) getUsersPhotos(ctx context.Context, userIDs []uuid.UUID) (map[uuid.UUID][]string, error) {
	userPhotos := make(map[uuid.UUID][]string)

	for _, userID := range userIDs {
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
		userPhotos[userID] = urls
	}

	return userPhotos, nil
}

func (s *MatchService) Unmatch(ctx context.Context, userID, targetUserID uuid.UUID) error {
	// Проверяем что мэтч существует
	match, err := s.matchRepo.GetByUsers(ctx, userID, targetUserID)
	if err != nil {
		return err
	}
	if match == nil {
		return errors.ErrMatchNotFound
	}

	// Деактивируем мэтч
	return s.matchRepo.UpdateActive(ctx, userID, targetUserID, false)
}
