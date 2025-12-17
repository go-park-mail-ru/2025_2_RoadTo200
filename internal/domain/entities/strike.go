package domain

import (
	"time"

	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/domain/constants"
	"github.com/google/uuid"
)

type Strike struct {
	ID            uuid.UUID              `json:"id" db:"id"`
	ReporterID    uuid.UUID              `json:"reporter_id" db:"reporter_id"`
	TargetUserID  uuid.UUID              `json:"target_user_id" db:"target_user_id"`
	Type          constants.StrikeType   `json:"type" db:"type"`
	Reason        string                 `json:"reason" db:"reason"`
	Status        constants.StrikeStatus `json:"status" db:"status"`
	CreatedAt     time.Time              `json:"created_at" db:"created_at"`
	UpdatedAt     *time.Time             `json:"updated_at" db:"updated_at"`
	ModeratorID   *uuid.UUID             `json:"moderator_id" db:"moderator_id"`
	ModeratorNote *string                `json:"moderator_note" db:"moderator_note"`
}
