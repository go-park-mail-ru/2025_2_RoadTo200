package service

import (
	"context"
	"time"

	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/domain/constants"
	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/domain/entities"
	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/domain/errors"
	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/handler/dto"
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

func (s *swipeService) ProcessSwipe(ctx context.Context, swiperID uuid.UUID, request *dto.SwipeRequest) (*dto.SwipeResponse, error) {
	card, err := uuid.Parse(request.CardID)
	if err != nil {
		return nil, err
	}
	// Валидация
	if swiperID == card {
		return nil, errors.ErrCannotSwipeSelf
	}

	// Конвертируем action в SwipeType
	swipeType, err := s.actionToSwipeType(ctx, request.Action)
	if err != nil {
		return nil, err
	}

	// Создаем свайп
	swipe := &domain.Swipe{
		SwiperUserID: swiperID,
		TargetUserID: card,
		SwipeType:    swipeType,
		CreatedAt:    time.Now(),
	}

	// Сохраняем свайп
	if err := s.swipeRepo.Create(ctx, swipe); err != nil {
		return nil, err
	}

	response := &dto.SwipeResponse{
		Message: "Swipe processed successfully",
	}

	// Если это лайк - проверяем на мэтч
	if swipeType == constants.SwipeTypeLike {
		// Проверяем взаимный лайк
		hasMutualLike, err := s.matchRepo.CheckMutualLike(ctx, swiperID, card)
		if err != nil {
			return nil, err
		}

		if hasMutualLike {
			// Создаем мэтч
			match := &domain.Match{
				User1ID:   swiperID,
				User2ID:   card,
				IsActive:  true,
				MatchedAt: time.Now(),
			}

			if err := s.matchRepo.Create(ctx, match); err != nil {
				return nil, err
			}

			response.IsMatch = true
			response.Message = "It's a match!"
		}
	}

	return response, nil
}

func (s *swipeService) actionToSwipeType(ctx context.Context, action string) (constants.SwipeType, error) {
	switch action {
	case "like":
		return constants.SwipeTypeLike, nil
	case "dislike":
		return constants.SwipeTypeDislike, nil
	default:
		return "", errors.ErrInvalidSwipeAction
	}
}
