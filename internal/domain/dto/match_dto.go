package dto

import domain "github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/domain/entities"

// UnmatchRequest represents unmatch request
type UnmatchRequest struct {
	TargetUserID string `json:"target_user_id" example:"550e8400-e29b-41d4-a716-446655440000"`
}

type MatchResponse struct {
	Match       domain.Match `json:"match"`
	User        domain.User  `json:"user"`
	Photos      []string     `json:"photos"`       // Добавляем фотографии
	Age         int          `json:"age"`          // Добавляем возраст
	Description string       `json:"description"`  // Добавляем описание
	PhotosCount int          `json:"photos_count"` // Добавляем количество фото
}

type MatchesResponse struct {
	Matches []MatchResponse `json:"matches"`
	Total   int             `json:"total"`
	Limit   int             `json:"limit"`
	Offset  int             `json:"offset"`
}
