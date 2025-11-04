package service

import (
	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/domain/entities"
	"github.com/google/uuid"
)

type FeedService interface {
	GetFeed(userID uuid.UUID, limit, offset int) ([]domain.User, error)
}

type FeedRequest struct {
	Limit  int `json:"limit,omitempty" form:"limit,omitempty"`
	Offset int `json:"offset,omitempty" form:"offset,omitempty"`
}

type FeedResponse struct {
	Users  []interface{} `json:"users"`
	Total  int           `json:"total"`
	Limit  int           `json:"limit"`
	Offset int           `json:"offset"`
}
