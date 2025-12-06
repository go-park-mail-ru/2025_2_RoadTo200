package service

import (
	"context"
	"time"

	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/core-service/domain/constants"
	domain "github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/core-service/domain/entities"
	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/handler/dto"
	"github.com/google/uuid"
)

type StrikeService interface {
	CreateStrike(ctx context.Context, strikeData *dto.StrikeCreateRequest) (*domain.Strike, error)
	GetStrikeByID(ctx context.Context, strikeID string) (*domain.Strike, error)
	GetStrikesByUserID(ctx context.Context, userID string, limit, offset int) ([]*domain.Strike, error)
	GetStrikesByType(ctx context.Context, strikeType constants.StrikeType, limit, offset int) ([]*domain.Strike, error)
	GetStrikesByDateRange(ctx context.Context, from, to time.Time, limit, offset int) ([]*domain.Strike, error)
	UpdateStrikeStatus(ctx context.Context, strikeID string, status constants.StrikeStatus, moderatorID *uuid.UUID, note *string) error
	DeleteStrike(ctx context.Context, strikeID string) error
	GetUserStrikeStats(ctx context.Context, userID string) (*dto.StrikeStats, error)
}
