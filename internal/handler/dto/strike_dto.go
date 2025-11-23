package dto

import (
	"time"

	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/domain/constants"
	"github.com/google/uuid"
)

type StrikeCreateRequest struct {
	ReporterID   uuid.UUID            `json:"reporter_id"`
	TargetUserID uuid.UUID            `json:"target_user_id"`
	Type         constants.StrikeType `json:"type"`
	Reason       string               `json:"reason"`
}

type StrikeStats struct {
	UserID        string                       `json:"user_id"`
	TotalStrikes  int                          `json:"total_strikes"`
	ActiveStrikes int                          `json:"active_strikes"`
	LastStrikeAt  *time.Time                   `json:"last_strike_at"`
	StrikeTypes   map[constants.StrikeType]int `json:"strike_types"`
}
