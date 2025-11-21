package service

import (
	"context"

	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/handler/dto"
	"github.com/google/uuid"
)

type MatchService interface {
	GetUserMatches(ctx context.Context, userID uuid.UUID, limit, offset int) (*dto.MatchesResponse, error)
	Unmatch(ctx context.Context, userID, targetUserID uuid.UUID) error
}
