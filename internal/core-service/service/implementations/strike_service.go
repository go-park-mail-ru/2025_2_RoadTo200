package service

import (
	"context"
	"fmt"
	"time"

	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/core-service/domain/constants"
	domain "github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/core-service/domain/entities"
	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/core-service/repository/interfaces"
	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/core-service/service/interfaces"
	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/domain/errors"
	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/handler/dto"
	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/pkg/logger"
	"github.com/google/uuid"
)

var _ service.StrikeService = (*StrikeService)(nil)

type StrikeService struct {
	strikeRepo interfaces.StrikeRepository
	userRepo   interfaces.UserRepository
	logger     logger.Log
}

func NewStrikeService(
	strikeRepo interfaces.StrikeRepository,
	userRepo interfaces.UserRepository,
	l logger.Log,
) *StrikeService {
	return &StrikeService{
		strikeRepo: strikeRepo,
		userRepo:   userRepo,
		logger:     l,
	}
}

// CreateStrike создает новую жалобу
func (s *StrikeService) CreateStrike(ctx context.Context, strikeData *dto.StrikeCreateRequest) (*domain.Strike, error) {
	// Валидация данных
	if err := s.ValidateStrikeCreate(ctx, strikeData); err != nil {
		s.logger.Warnf("Validation error: %s", err)
		return nil, err
	}

	// Проверяем существование пользователей
	reporter, err := s.userRepo.GetByID(ctx, strikeData.ReporterID)
	if err != nil {
		s.logger.Errorf("GetUserByID error: %s", err)
		return nil, err
	}
	if reporter == nil {
		s.logger.Warn("User not found")
		return nil, errors.ErrUserNotFound
	}

	targetUser, err := s.userRepo.GetByID(ctx, strikeData.TargetUserID)
	if err != nil {
		s.logger.Errorf("GetUserByID error: %s", err)
		return nil, err
	}
	if targetUser == nil {
		s.logger.Warn("User not found")
		return nil, errors.ErrUserNotFound
	}

	// Проверяем, что пользователь не жалуется на себя
	if strikeData.ReporterID == strikeData.TargetUserID {
		s.logger.Warn("Strike reporter already exists")
		return nil, errors.ErrSelfStrikeNotAllowed
	}

	// Проверяем, нет ли уже активной жалобы от этого пользователя
	hasActiveStrike, err := s.strikeRepo.HasActiveStrikeFromUser(ctx, strikeData.ReporterID.String(), strikeData.TargetUserID.String())
	if err != nil {
		s.logger.Errorf("HasActiveStrike error: %s", err)
		return nil, err
	}
	if hasActiveStrike {
		s.logger.Warn("Strike reporter already active")
		return nil, errors.ErrDuplicateStrike
	}

	// Создаем жалобу
	strike := &domain.Strike{
		ReporterID:   strikeData.ReporterID,
		TargetUserID: strikeData.TargetUserID,
		Type:         strikeData.Type,
		Reason:       strikeData.Reason,
		Status:       constants.StrikeStatusPending,
		CreatedAt:    time.Now(),
	}

	if err := s.strikeRepo.CreateStrike(ctx, strike); err != nil {
		s.logger.Errorf("CreateStrike error: %s", err)
		return nil, err
	}

	s.logger.Debugf("New strike created: %s by user %s against user %s", strike.ID, strike.ReporterID, strike.TargetUserID)

	return strike, nil
}

// GetStrikeByID возвращает жалобу по ID
func (s *StrikeService) GetStrikeByID(ctx context.Context, strikeID string) (*domain.Strike, error) {
	strike, err := s.strikeRepo.GetStrikeByID(ctx, strikeID)
	if err != nil {
		s.logger.Errorf("GetStrikeByID error: %s", err)
		return nil, err
	}

	return strike, nil
}

// GetStrikesByUserID возвращает все жалобы на конкретного пользователя
func (s *StrikeService) GetStrikesByUserID(ctx context.Context, userID string, limit, offset int) ([]*domain.Strike, error) {
	// Проверяем существование пользователя
	user, err := s.userRepo.GetByID(ctx, uuid.MustParse(userID))
	if err != nil {
		s.logger.Errorf("GetUserByID error: %s", err)
		return nil, err
	}
	if user == nil {
		s.logger.Warn("User not found")
		return nil, errors.ErrUserNotFound
	}

	strikes, err := s.strikeRepo.GetStrikesByUserID(ctx, userID, limit, offset)
	if err != nil {
		s.logger.Errorf("GetStrikesByUserID error: %s", err)
		return nil, err
	}

	return strikes, nil
}

// GetStrikesByType возвращает жалобы по типу нарушения
func (s *StrikeService) GetStrikesByType(ctx context.Context, strikeType constants.StrikeType, limit, offset int) ([]*domain.Strike, error) {
	strikes, err := s.strikeRepo.GetStrikesByType(ctx, strikeType, limit, offset)
	if err != nil {
		s.logger.Errorf("GetStrikesByType error: %s", err)
		return nil, err
	}

	return strikes, nil
}

func (s *StrikeService) GetStrikesByDateRange(ctx context.Context, from, to time.Time, limit, offset int) ([]*domain.Strike, error) {
	// Если from нулевое (не установлено), устанавливаем минимальную дату
	if from.IsZero() {
		from = time.Time{} // Уже нулевое, оставляем как есть
	}

	// Если to нулевое (не установлено), устанавливаем текущее время + 1 день
	if to.IsZero() {
		to = time.Now().Add(24 * time.Hour)
	}

	// Валидация периода
	if from.After(to) {
		s.logger.Warnf("From date %s is after to %s", from, to)
		return nil, errors.ErrInvalidDateRange
	}

	// Ограничиваем максимальный период (например, 1 год)
	maxPeriod := 365 * 24 * time.Hour
	if to.Sub(from) > maxPeriod {
		s.logger.Warnf("Date range too large: from %s to %s (max allowed: %v)", from, to, maxPeriod)
		return nil, errors.ErrDateRangeTooLarge
	}

	// Устанавливаем значения по умолчанию для пагинации, если не установлены
	if limit <= 0 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}

	strikes, err := s.strikeRepo.GetStrikesByDateRange(ctx, from, to, limit, offset)
	if err != nil {
		s.logger.Errorf("GetStrikesByDateRange error: %s", err)
		return nil, err
	}

	return strikes, nil
}

// UpdateStrikeStatus обновляет статус жалобы
func (s *StrikeService) UpdateStrikeStatus(ctx context.Context, strikeID string, status constants.StrikeStatus, moderatorID *uuid.UUID, note *string) error {
	// Валидация статуса
	if err := s.ValidateStrikeStatus(ctx, status); err != nil {
		s.logger.Warnf("Validation error: %s", err)
		return err
	}

	// Получаем текущую жалобу
	strike, err := s.strikeRepo.GetStrikeByID(ctx, strikeID)
	if err != nil {
		s.logger.Errorf("GetStrike error: %s", err)
		return err
	}

	// Проверяем, что статус действительно изменился
	if strike.Status == status {
		s.logger.Warnf("Strike status already set to %s", status)
		return nil
	}

	// Обновляем статус
	if err := s.strikeRepo.UpdateStrikeStatus(ctx, strikeID, status); err != nil {
		s.logger.Errorf("UpdateStrike error: %s", err)
		return err
	}

	// Логируем изменение статуса
	s.logger.Debugf("Strike %s status changed from %s to %s by moderator %v", strikeID, strike.Status, status, moderatorID)

	return nil
}

// DeleteStrike удаляет жалобу (жесткое удаление)
func (s *StrikeService) DeleteStrike(ctx context.Context, strikeID string) error {
	// Проверяем существование жалобы
	strike, err := s.strikeRepo.GetStrikeByID(ctx, strikeID)
	if err != nil {
		s.logger.Errorf("GetStrike error: %s", err)
		return err
	}

	// Логируем факт удаления
	s.logger.Debugf("Deleting strike %s (hard delete)", strikeID)

	// Выполняем жесткое удаление
	if err := s.strikeRepo.DeleteStrike(ctx, strikeID); err != nil {
		s.logger.Errorf("DeleteStrike error: %s", err)
		return err
	}

	s.logger.Debugf("Strike %s hard deleted (was in status: %s)", strikeID, strike.Status)
	return nil
}

// GetUserStrikeStats возвращает статистику по жалобам пользователя
func (s *StrikeService) GetUserStrikeStats(ctx context.Context, userID string) (*dto.StrikeStats, error) {
	// Проверяем существование пользователя
	user, err := s.userRepo.GetByID(ctx, uuid.MustParse(userID))
	if err != nil {
		s.logger.Errorf("GetUserStrike error: %s", err)
		return nil, err
	}
	if user == nil {
		s.logger.Warn("User not found")
		return nil, errors.ErrUserNotFound
	}

	// Получаем все жалобы на пользователя
	allStrikes, err := s.strikeRepo.GetStrikesByUserID(ctx, userID, 1000, 0) // Большой лимит для статистики
	if err != nil {
		s.logger.Errorf("GetUserStrike error: %s", err)
		return nil, err
	}

	// Считаем статистику
	stats := &dto.StrikeStats{
		UserID:       userID,
		TotalStrikes: len(allStrikes),
	}

	var lastStrikeAt *time.Time

	for _, strike := range allStrikes {
		// Считаем по типам
		switch strike.Status {
		case constants.StrikeStatusPending:
			stats.StrikeTypes.Pending++
		case constants.StrikeStatusApproved:
			stats.StrikeTypes.Approved++
		case constants.StrikeStatusRejected:
			stats.StrikeTypes.Rejected++
		case constants.StrikeStatusResolved:
			stats.StrikeTypes.Resolved++
		}

		// Находим последнюю жалобу
		if lastStrikeAt == nil || strike.CreatedAt.After(*lastStrikeAt) {
			lastStrikeAt = &strike.CreatedAt
		}
	}

	stats.LastStrikeAt = lastStrikeAt

	return stats, nil
}

// ValidateStrikeCreate валидация данных при создании жалобы
func (s *StrikeService) ValidateStrikeCreate(ctx context.Context, strikeData *dto.StrikeCreateRequest) error {
	if strikeData.ReporterID == uuid.Nil {
		return errors.ErrInvalidReporterID
	}

	if strikeData.TargetUserID == uuid.Nil {
		return errors.ErrInvalidTargetUserID
	}

	// Валидация типа жалобы
	validTypes := map[constants.StrikeType]struct{}{
		constants.StrikeTypeSpam:          {},
		constants.StrikeTypeFakeProfile:   {},
		constants.StrikeTypeOffensive:     {},
		constants.StrikeTypeHarassment:    {},
		constants.StrikeTypeInappropriate: {},
		constants.StrikeTypeUnderage:      {},
		constants.StrikeTypeCopyright:     {},
		constants.StrikeTypeOther:         {},
	}

	if _, ok := validTypes[strikeData.Type]; !ok {
		return fmt.Errorf("invalid strike type: %s", strikeData.Type)
	}

	// Валидация длины причины
	if len(strikeData.Reason) > constants.MaxStrikeReasonLength {
		return fmt.Errorf("reason too long, max %d characters", constants.MaxStrikeReasonLength)
	}

	return nil
}

// ValidateStrikeStatus валидация статуса жалобы
func (s *StrikeService) ValidateStrikeStatus(ctx context.Context, status constants.StrikeStatus) error {
	validStatuses := map[constants.StrikeStatus]struct{}{
		constants.StrikeStatusPending:  {},
		constants.StrikeStatusApproved: {},
		constants.StrikeStatusRejected: {},
		constants.StrikeStatusResolved: {},
	}

	if _, ok := validStatuses[status]; !ok {
		return fmt.Errorf("invalid strike status: %s", status)
	}

	return nil
}
