package service

import (
	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/domain/entities"
	"github.com/google/uuid"
)

type MatchService interface {
	GetUserMatches(userID uuid.UUID, limit, offset int) (*MatchesResponse, error)
	Unmatch(userID, targetUserID uuid.UUID) error
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

type UnmatchRequest struct {
	TargetUserID string `json:"target_user_id"`
}
