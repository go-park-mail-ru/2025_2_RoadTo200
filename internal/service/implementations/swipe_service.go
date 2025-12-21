package service

import (
	"context"
	"fmt"
	"time"

	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/domain/constants"
	domain "github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/domain/entities"
	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/domain/errors"
	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/handler/dto"
	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/logger"
	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/repository/interfaces"
	serviceInterfaces "github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/service/interfaces"
	"github.com/google/uuid"
)

type SwipeService struct {
	swipeRepo           interfaces.SwipeRepository
	matchRepo           interfaces.MatchRepository
	userRepo            interfaces.UserRepository
	notificationService serviceInterfaces.NotificationService
	logger              logger.Log
}

func NewSwipeService(
	swipeRepo interfaces.SwipeRepository,
	matchRepo interfaces.MatchRepository,
	userRepo interfaces.UserRepository,
	notificationService serviceInterfaces.NotificationService,
	logger logger.Log,
) *SwipeService {
	return &SwipeService{
		swipeRepo:           swipeRepo,
		matchRepo:           matchRepo,
		userRepo:            userRepo,
		notificationService: notificationService,
		logger:              logger,
	}
}

func (s *SwipeService) ProcessSwipe(ctx context.Context, swiperID uuid.UUID, request *dto.SwipeRequest) (*dto.SwipeResponse, error) {
	s.logger.Infof("ProcessSwipe called: swiperID=%s, cardID=%s, action=%s", swiperID, request.CardID, request.Action)

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

	s.logger.Infof("SwipeType determined: %s", swipeType)

	// Создаем свайп
	swipe := &domain.Swipe{
		SwiperUserID: swiperID,
		TargetUserID: card,
		SwipeType:    swipeType,
		CreatedAt:    time.Now(),
	}

	// Сохраняем свайп
	if err := s.swipeRepo.Create(ctx, swipe); err != nil {
		s.logger.Errorf("Failed to create swipe: %v", err)
		return nil, err
	}
	s.logger.Infof("Swipe created successfully: swiper=%s, target=%s, type=%s", swiperID, card, swipeType)

	response := &dto.SwipeResponse{
		Message: "Swipe processed successfully",
	}

	// Отправляем уведомления в зависимости от типа свайпа
	if swipeType == constants.SwipeTypeSuperLike {
		// Проверяем, есть ли у пользователя суперлайки
		swiper, err := s.userRepo.GetByID(ctx, swiperID)
		if err != nil {
			s.logger.Errorf("Failed to get user for super_like check: %v", err)
			return nil, err
		}
		if swiper == nil {
			return nil, fmt.Errorf("user not found")
		}

		// Проверяем, есть ли суперлайки
		if swiper.SuperLikesCount <= 0 {
			return nil, fmt.Errorf("no super likes available")
		}

		// Уменьшаем количество суперлайков
		swiper.SuperLikesCount--
		if err := s.userRepo.Update(ctx, swiper); err != nil {
			s.logger.Errorf("Failed to update super_likes_count: %v", err)
			return nil, fmt.Errorf("failed to update super likes count: %w", err)
		}
		s.logger.Infof("Super likes count decreased: user=%s, remaining=%d", swiperID, swiper.SuperLikesCount)

		// Проверяем на мэтч (суперлайк тоже может создать матч, если второй пользователь лайкнул или суперлайкнул)
		s.logger.Debugf("Checking for mutual like (super_like): swiper=%s, target=%s", swiperID, card)
		hasMutualLike, err := s.matchRepo.CheckMutualLike(ctx, swiperID, card)
		if err != nil {
			s.logger.Errorf("Error checking mutual like for super_like: %v", err)
			return nil, err
		}

		s.logger.Debugf("Mutual like check result (super_like): hasMutualLike=%v", hasMutualLike)

		if hasMutualLike {
			// Проверяем, не существует ли уже мэтч между этими пользователями
			existingMatch, err := s.matchRepo.GetByUsers(ctx, swiperID, card)
			if err != nil {
				s.logger.Errorf("Failed to check existing match for super_like: %v", err)
				return nil, err
			}

			var match *domain.Match
			if existingMatch != nil {
				// Мэтч уже существует, используем его
				s.logger.Infof("Match already exists for super_like: user1=%s, user2=%s, matchID=%s", swiperID, card, existingMatch.ID)
				match = existingMatch
			} else {
				// Создаем новый мэтч
				match = &domain.Match{
				User1ID:   swiperID,
				User2ID:   card,
				IsActive:  true,
				MatchedAt: time.Now(),
			}

			if err := s.matchRepo.Create(ctx, match); err != nil {
				s.logger.Errorf("Failed to create match for super_like: %v", err)
				return nil, err
			}

			s.logger.Infof("Match created from super_like: user1=%s, user2=%s, matchID=%s", swiperID, card, match.ID)
			}

			// Удаляем все старые уведомления (лайк, суперлайк, мэтч) между этими пользователями
			if s.notificationService != nil {
				if err := s.notificationService.DeleteUserNotifications(ctx, swiperID, card); err != nil {
					s.logger.Warnf("Failed to delete old user notifications for super_like: %v", err)
					// Не прерываем выполнение, продолжаем отправку новых уведомлений
				}
			}

			// Отправляем уведомления о матче обоим пользователям
			matchIDPtr := &match.ID
			swiperIDPtr := &swiperID
			cardIDPtr := &card

			if s.notificationService == nil {
				s.logger.Error("CRITICAL: NotificationService is nil! Cannot send match notifications")
			} else {
				// Уведомление первому пользователю (swiper)
				s.logger.Infof("Sending match notification to user1 (super_like): swiperID=%s, fromUserID=%s, matchID=%s", swiperID, card, match.ID)
				if err := s.notificationService.SendNotification(ctx, swiperID, constants.NotificationTypeMatch, cardIDPtr, matchIDPtr); err != nil {
					s.logger.Errorf("Failed to send match notification to user1 (super_like): %v", err)
				} else {
					s.logger.Infof("Match notification sent successfully to user1 (super_like): %s", swiperID)
				}

				// Уведомление второму пользователю (target)
				s.logger.Infof("Sending match notification to user2 (super_like): targetID=%s, fromUserID=%s, matchID=%s", card, swiperID, match.ID)
				if err := s.notificationService.SendNotification(ctx, card, constants.NotificationTypeMatch, swiperIDPtr, matchIDPtr); err != nil {
					s.logger.Errorf("Failed to send match notification to user2 (super_like): %v", err)
				} else {
					s.logger.Infof("Match notification sent successfully to user2 (super_like): %s", card)
				}
			}

			response.IsMatch = true
			response.Message = "It's a match!"
		} else {
			// При суперлайке всегда отправляем уведомление получателю (если нет матча)
			swiperIDPtr := &swiperID
			if err := s.notificationService.SendNotification(ctx, card, constants.NotificationTypeSuperLike, swiperIDPtr, nil); err != nil {
				s.logger.Errorf("Failed to send super_like notification: %v", err)
			} else {
				s.logger.Infof("Super_like notification sent: from=%s, to=%s", swiperID, card)
			}
		}
	} else if swipeType == constants.SwipeTypeLike {
		// Сначала проверяем на мэтч (важно делать это до отправки уведомления о лайке)
		s.logger.Debugf("Checking for mutual like: swiper=%s, target=%s", swiperID, card)
		hasMutualLike, err := s.matchRepo.CheckMutualLike(ctx, swiperID, card)
		if err != nil {
			s.logger.Errorf("Error checking mutual like: %v", err)
			return nil, err
		}

		s.logger.Debugf("Mutual like check result: hasMutualLike=%v", hasMutualLike)

		if hasMutualLike {
			// Проверяем, не существует ли уже мэтч между этими пользователями
			existingMatch, err := s.matchRepo.GetByUsers(ctx, swiperID, card)
			if err != nil {
				s.logger.Errorf("Failed to check existing match: %v", err)
				return nil, err
			}

			var match *domain.Match
			if existingMatch != nil {
				// Мэтч уже существует, используем его
				s.logger.Infof("Match already exists: user1=%s, user2=%s, matchID=%s", swiperID, card, existingMatch.ID)
				match = existingMatch
			} else {
				// Создаем новый мэтч
				match = &domain.Match{
				User1ID:   swiperID,
				User2ID:   card,
				IsActive:  true,
				MatchedAt: time.Now(),
			}

			if err := s.matchRepo.Create(ctx, match); err != nil {
				return nil, err
			}

			s.logger.Infof("Match created: user1=%s, user2=%s, matchID=%s", swiperID, card, match.ID)
			}

			// Удаляем все старые уведомления (лайк, суперлайк, мэтч) между этими пользователями
			if s.notificationService != nil {
				if err := s.notificationService.DeleteUserNotifications(ctx, swiperID, card); err != nil {
					s.logger.Warnf("Failed to delete old user notifications: %v", err)
					// Не прерываем выполнение, продолжаем отправку новых уведомлений
				}
			}

			// Отправляем уведомления о матче обоим пользователям
			// НЕ отправляем уведомление о лайке, только о матче
			matchIDPtr := &match.ID
			swiperIDPtr := &swiperID
			cardIDPtr := &card

			// Проверяем, что NotificationService инициализирован
			if s.notificationService == nil {
				s.logger.Error("CRITICAL: NotificationService is nil! Cannot send match notifications")
			} else {
				s.logger.Infof("NotificationService is available, sending notifications...")

				// Уведомление первому пользователю (swiper)
				s.logger.Infof("=== Sending match notification to user1: swiperID=%s, fromUserID=%s, matchID=%s ===", swiperID, card, match.ID)
				err1 := s.notificationService.SendNotification(ctx, swiperID, constants.NotificationTypeMatch, cardIDPtr, matchIDPtr)
				if err1 != nil {
					s.logger.Errorf("ERROR: Failed to send match notification to user1 (swiper): %v", err1)
				} else {
					s.logger.Infof("SUCCESS: Match notification sent successfully to user1: %s", swiperID)
				}

				// Уведомление второму пользователю (target)
				s.logger.Infof("=== Sending match notification to user2: targetID=%s, fromUserID=%s, matchID=%s ===", card, swiperID, match.ID)
				err2 := s.notificationService.SendNotification(ctx, card, constants.NotificationTypeMatch, swiperIDPtr, matchIDPtr)
				if err2 != nil {
					s.logger.Errorf("ERROR: Failed to send match notification to user2 (target): %v", err2)
				} else {
					s.logger.Infof("SUCCESS: Match notification sent successfully to user2: %s", card)
				}
			}

			response.IsMatch = true
			response.Message = "It's a match!"
		} else {
			// Если нет матча, проверяем премиум и отправляем уведомление о лайке
			targetUser, err := s.userRepo.GetByID(ctx, card)
			if err != nil {
				s.logger.Warnf("Failed to get target user for like notification: %v", err)
			} else if targetUser != nil && targetUser.IsPremium {
				// Если есть премиум - отправляем уведомление о лайке
				swiperIDPtr := &swiperID
				if err := s.notificationService.SendNotification(ctx, card, constants.NotificationTypeLike, swiperIDPtr, nil); err != nil {
					s.logger.Errorf("Failed to send like notification: %v", err)
				} else {
					s.logger.Infof("Like notification sent: from=%s, to=%s (premium user)", swiperID, card)
				}
			} else {
				s.logger.Debugf("Like notification skipped: user %s is not premium", card)
			}
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
	case "super_like":
		return constants.SwipeTypeSuperLike, nil
	default:
		return "", errors.ErrInvalidSwipeAction
	}
}
