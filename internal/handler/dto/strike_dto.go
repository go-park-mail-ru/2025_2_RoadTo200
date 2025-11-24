package dto

import (
	"time"

	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/domain/constants"
	domain "github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/domain/entities"
	"github.com/google/uuid"
)

type StrikeCreateRequest struct {
	ReporterID   uuid.UUID            `json:"reporter_id"`
	TargetUserID uuid.UUID            `json:"target_user_id"`
	Type         constants.StrikeType `json:"type"`
	Reason       string               `json:"reason"`
}

type StrikeStats struct {
	UserID       string         `json:"user_id"`
	TotalStrikes int            `json:"total_strikes"`
	LastStrikeAt *time.Time     `json:"last_strike_at"`
	StrikeTypes  StrikeTypeStat `json:"strike_types"`
}

type StrikeTypeStat struct {
	Pending  int `json:"pending"`
	Approved int `json:"approved"`
	Rejected int `json:"rejected"`
	Resolved int `json:"resolved"`
}

type StrikeStatusUpdateRequest struct {
	Status      constants.StrikeStatus `json:"status"`
	ModeratorID *uuid.UUID             `json:"moderator_id,omitempty"`
	Note        *string                `json:"note,omitempty"`
}

type StrikeResponse struct {
	ID            uuid.UUID              `json:"id"`
	ReporterID    uuid.UUID              `json:"reporter_id"`
	TargetUserID  uuid.UUID              `json:"target_user_id"`
	Type          constants.StrikeType   `json:"type"`
	Reason        string                 `json:"reason"`
	Status        constants.StrikeStatus `json:"status"`
	CreatedAt     time.Time              `json:"created_at"`
	UpdatedAt     *time.Time             `json:"updated_at,omitempty"`
	ModeratorID   *uuid.UUID             `json:"moderator_id,omitempty"`
	ModeratorNote *string                `json:"moderator_note,omitempty"`
}

type StrikesListResponse struct {
	Strikes []StrikeResponse `json:"strikes"`
	Total   int              `json:"total"`
}

// ToStrikeResponse преобразует доменную сущность в DTO
func ToStrikeResponse(strike *domain.Strike) StrikeResponse {
	return StrikeResponse{
		ID:            strike.ID,
		ReporterID:    strike.ReporterID,
		TargetUserID:  strike.TargetUserID,
		Type:          strike.Type,
		Reason:        strike.Reason,
		Status:        strike.Status,
		CreatedAt:     strike.CreatedAt,
		UpdatedAt:     strike.UpdatedAt,
		ModeratorID:   strike.ModeratorID,
		ModeratorNote: strike.ModeratorNote,
	}
}

// ToStrikesListResponse преобразует список доменных сущностей в DTO
func ToStrikesListResponse(strikes []*domain.Strike) StrikesListResponse {
	strikeResponses := make([]StrikeResponse, 0, len(strikes))
	for _, strike := range strikes {
		strikeResponses = append(strikeResponses, ToStrikeResponse(strike))
	}

	return StrikesListResponse{
		Strikes: strikeResponses,
		Total:   len(strikeResponses),
	}
}
