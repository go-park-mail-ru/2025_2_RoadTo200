package service

import (
	"time"

	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/domain/constants"
	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/domain/entities"
	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/domain/errors"
	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/repository/interfaces"
	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/service/interfaces"
	"github.com/google/uuid"
)

type swipeService struct {
	swipeRepo interfaces.SwipeRepository
	matchRepo interfaces.MatchRepository
}

func NewSwipeService(swipeRepo interfaces.SwipeRepository, matchRepo interfaces.MatchRepository) service.SwipeService {
	return &swipeService{
		swipeRepo: swipeRepo,
		matchRepo: matchRepo,
	}
}

func (s *swipeService) ProcessSwipe(swiperID uuid.UUID, request *service.SwipeRequest) (*service.SwipeResponse, error) {
	// Валидация
	if swiperID == request.CardID {
		return nil, errors.ErrCannotSwipeSelf
	}

	// Конвертируем action в SwipeType
	swipeType, err := s.actionToSwipeType(request.Action)
	if err != nil {
		return nil, err
	}

	// Создаем свайп
	swipe := &domain.Swipe{
		SwiperUserID: swiperID,
		TargetUserID: request.CardID,
		SwipeType:    swipeType,
		CreatedAt:    time.Now(),
	}

	// Сохраняем свайп
	if err := s.swipeRepo.Create(swipe); err != nil {
		return nil, err
	}

	response := &service.SwipeResponse{
		Message: "Swipe processed successfully",
	}

	// Если это лайк - проверяем на мэтч
	if swipeType == constants.SwipeTypeLike {
		// Проверяем взаимный лайк
		hasMutualLike, err := s.matchRepo.CheckMutualLike(swiperID, request.CardID)
		if err != nil {
			return nil, err
		}

		if hasMutualLike {
			// Создаем мэтч
			match := &domain.Match{
				User1ID:   swiperID,
				User2ID:   request.CardID,
				IsActive:  true,
				MatchedAt: time.Now(),
			}

			if err := s.matchRepo.Create(match); err != nil {
				return nil, err
			}

			response.Match = true
			response.Message = "It's a match!"
		}
	}

	return response, nil
}

func (s *swipeService) actionToSwipeType(action string) (constants.SwipeType, error) {
	switch action {
	case "like":
		return constants.SwipeTypeLike, nil
	case "dislike":
		return constants.SwipeTypeDislike, nil
	default:
		return "", errors.ErrInvalidSwipeAction
	}
}
