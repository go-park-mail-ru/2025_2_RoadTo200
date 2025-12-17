package service

import (
	"context"

	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/handler/dto"
	"github.com/google/uuid"
)

type FeedService interface {
	GetFeed(ctx context.Context, userID uuid.UUID, limit, offset int) ([]dto.FeedUser, error)
}
