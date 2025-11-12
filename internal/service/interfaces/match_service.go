package service

import (
	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/domain/dto"
	"github.com/google/uuid"
)

type MatchService interface {
	GetUserMatches(userID uuid.UUID, limit, offset int) (*dto.MatchesResponse, error)
	Unmatch(userID, targetUserID uuid.UUID) error
}
