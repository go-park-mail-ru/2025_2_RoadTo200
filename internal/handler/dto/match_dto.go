package dto

import (
	auth "github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/auth-service/domain/entities"
	chat "github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/chat-service/domain/entities"
)

// UnmatchRequest represents unmatch request
type UnmatchRequest struct {
	TargetUserID string `json:"target_user_id" example:"550e8400-e29b-41d4-a716-446655440000"`
}

type MatchResponse struct {
	Match       chat.Match `json:"match"`
	User        auth.User  `json:"user"`
	Photos      []string   `json:"photos"`       // Добавляем фотографии
	Age         int        `json:"age"`          // Добавляем возраст
	Description string     `json:"description"`  // Добавляем описание
	PhotosCount int        `json:"photos_count"` // Добавляем количество фото
}

type MatchesResponse struct {
	Matches []MatchResponse `json:"matches"`
	Total   int             `json:"total"`
	Limit   int             `json:"limit"`
	Offset  int             `json:"offset"`
}
