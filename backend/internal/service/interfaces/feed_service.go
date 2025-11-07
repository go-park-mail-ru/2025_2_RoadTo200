package service

import (
	//"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/domain/entities"
	"github.com/google/uuid"
)

type FeedService interface {
	GetFeed(userID uuid.UUID, limit, offset int) ([]FeedUser, error)
}

// FeedUser represents user data for feed
type FeedUser struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Age         int      `json:"age"`
	Gender      string   `json:"gender"`
	Description string   `json:"description"`
	Images      []string `json:"images"`
	PhotosCount int      `json:"photosCount"`
}

type FeedRequest struct {
	Limit  int `json:"limit,omitempty" form:"limit,omitempty"`
	Offset int `json:"offset,omitempty" form:"offset,omitempty"`
}

type FeedResponse struct {
	Users  []FeedUser `json:"users"`
	Total  int        `json:"total"`
	Limit  int        `json:"limit"`
	Offset int        `json:"offset"`
}
