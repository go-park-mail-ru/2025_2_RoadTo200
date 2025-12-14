package service

import (
	"context"
	"time"

	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/domain/constants"
	domain "github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/domain/entities"
	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/domain/errors"
	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/handler/dto"
	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/repository/interfaces"
	serviceInterfaces "github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/service/interfaces"
	"github.com/google/uuid"
)

type SwipeService struct {
	swipeRepo           interfaces.SwipeRepository
	matchRepo           interfaces.MatchRepository
	userRepo            interfaces.UserRepository
	notificationService serviceInterfaces.NotificationService
}

func NewSwipeService(
	swipeRepo interfaces.SwipeRepository,
	matchRepo interfaces.MatchRepository,
	userRepo interfaces.UserRepository,
	notificationService serviceInterfaces.NotificationService,
) *SwipeService {
	return &SwipeService{
		swipeRepo:           swipeRepo,
		matchRepo:           matchRepo,
		userRepo:            userRepo,
		notificationService: notificationService,
	}
}

func (s *SwipeService) ProcessSwipe(ctx context.Context, swiperID uuid.UUID, request *dto.SwipeRequest) (*dto.SwipeResponse, error) {
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

	// Отправляем уведомления в зависимости от типа свайпа
	if swipeType == constants.SwipeTypeSuperLike {
		// При суперлайке всегда отправляем уведомление получателю
		swiperIDPtr := &swiperID
		if err := s.notificationService.SendNotification(ctx, card, constants.NotificationTypeSuperLike, swiperIDPtr, nil); err != nil {
			// Логируем ошибку, но не прерываем выполнение
			// Уведомление - это дополнительная функция
		}
	} else if swipeType == constants.SwipeTypeLike {
		// При лайке проверяем, есть ли у получателя премиум
		targetUser, err := s.userRepo.GetByID(ctx, card)
		if err == nil && targetUser != nil && targetUser.IsPremium {
			// Если есть премиум - отправляем уведомление о лайке
			swiperIDPtr := &swiperID
			if err := s.notificationService.SendNotification(ctx, card, constants.NotificationTypeLike, swiperIDPtr, nil); err != nil {
				// Логируем ошибку, но не прерываем выполнение
			}
		}

		// Проверяем на мэтч
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

			// Отправляем уведомления о матче обоим пользователям
			matchIDPtr := &match.ID
			swiperIDPtr := &swiperID
			cardIDPtr := &card

			// Уведомление первому пользователю (swiper)
			if err := s.notificationService.SendNotification(ctx, swiperID, constants.NotificationTypeMatch, cardIDPtr, matchIDPtr); err != nil {
				// Логируем ошибку, но не прерываем выполнение
			}

			// Уведомление второму пользователю (target)
			if err := s.notificationService.SendNotification(ctx, card, constants.NotificationTypeMatch, swiperIDPtr, matchIDPtr); err != nil {
				// Логируем ошибку, но не прерываем выполнение
			}

			response.IsMatch = true
			response.Message = "It's a match!"
		}
	}

	return response, nil
}

func (s *SwipeService) actionToSwipeType(ctx context.Context, action string) (constants.SwipeType, error) {
	switch action {
	case "like":
		return constants.SwipeTypeLike, nil
	case "dislike":
		return constants.SwipeTypeDislike, nil
	default:
		return "", errors.ErrInvalidSwipeAction
	}
}
