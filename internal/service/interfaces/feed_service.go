package service

import (
	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/handler/dto"
	"github.com/google/uuid"
)

type FeedService interface {
	GetFeed(userID uuid.UUID, limit, offset int) ([]dto.FeedUser, error)
}
