package domain

import (
	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/domain/constants"
	"github.com/google/uuid"
)

type Interest struct {
	User_id uuid.UUID              `json:"user_id" db:"user_id"`
	Theme   constants.InterestType `json:"theme" db:"theme"`
}
