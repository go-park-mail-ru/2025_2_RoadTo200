package service

import (
	"context"

	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/handler/dto"
	"github.com/google/uuid"
)

type SwipeService interface {
	ProcessSwipe(ctx context.Context, swiperID uuid.UUID, request *dto.SwipeRequest) (*dto.SwipeResponse, error)
}
