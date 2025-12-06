package interfaces

import (
	"context"
	"time"

	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/core-service/domain/constants"
	domain "github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/core-service/domain/entities"
)

type StrikeRepository interface {
	// Создать новую жалобу
	CreateStrike(ctx context.Context, strike *domain.Strike) error

	// Получить жалобу по ID
	GetStrikeByID(ctx context.Context, strikeID string) (*domain.Strike, error)

	// Получить все жалобы на конкретного пользователя
	GetStrikesByUserID(ctx context.Context, userID string, limit, offset int) ([]*domain.Strike, error)

	// Получить жалобы по типу нарушения
	GetStrikesByType(ctx context.Context, strikeType constants.StrikeType, limit, offset int) ([]*domain.Strike, error)

	// Получить жалобы за определенный период
	GetStrikesByDateRange(ctx context.Context, from, to time.Time, limit, offset int) ([]*domain.Strike, error)

	// Обновить статус жалобы (например, при модерации)
	UpdateStrikeStatus(ctx context.Context, strikeID string, status constants.StrikeStatus) error

	// Удалить жалобу (мягкое удаление)
	DeleteStrike(ctx context.Context, strikeID string) error

	// Получить статистику по жалобам для пользователя
	//GetUserStrikeStats(ctx context.context, userID string) (*StrikeStats, error)

	// Проверить, есть ли уже активная жалоба от этого пользователя на целевого
	HasActiveStrikeFromUser(ctx context.Context, reporterID, targetUserID string) (bool, error)

	// Получить количество жалоб по типам (для аналитики)
	//GetStrikeCountByType(ctx context.context) (map[constants.StrikeType]int, error)
}
