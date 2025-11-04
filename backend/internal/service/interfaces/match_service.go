package service

import (
	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/domain/entities"
	"github.com/google/uuid"
)

type MatchResponse struct {
	Match domain.Match `json:"match"`
	User  domain.User  `json:"user"`
}

type MatchesResponse struct {
	Matches []MatchResponse `json:"matches"`
	Total   int             `json:"total"`
	Limit   int             `json:"limit"`
	Offset  int             `json:"offset"`
}

type MatchService interface {
	GetUserMatches(userID uuid.UUID, limit, offset int) (*MatchesResponse, error)
	Unmatch(userID, targetUserID uuid.UUID) error
}

type UnmatchRequest struct {
	TargetUserID uuid.UUID `json:"target_user_id"`
}
