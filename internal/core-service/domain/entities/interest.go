package domain

import (
	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/core-service/domain/constants"
	"github.com/google/uuid"
)

type Interest struct {
	UserID uuid.UUID              `json:"user_id" db:"user_id"`
	Theme  constants.InterestType `json:"theme" db:"theme"`
}
