package dto

//go:generate easyjson -all -no_std_marshalers feed_dto.go

import domain "github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/domain/entities"

// FeedUser represents user data for feed
type FeedUser struct {
	ID          string            `json:"id"`
	Name        string            `json:"name"`
	Age         int               `json:"age"`
	Gender      string            `json:"gender"`
	Description string            `json:"description"`
	Images      []string          `json:"images"`
	PhotosCount int               `json:"photos_count"`
	Artist      *string           `db:"artist" json:"artist,omitempty"`
	Quote       *string           `db:"quote" json:"quote,omitempty"`
	IsPremium   bool              `json:"is_premium"`
	Interests   []domain.Interest `json:"interests"`
}

// FeedResponse represents feed response
type FeedResponse struct {
	Users  []FeedUser `json:"users"`
	Total  int        `json:"total"`
	Limit  int        `json:"limit"`
	Offset int        `json:"offset"`
}
