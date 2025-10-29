package domain

import (
	"time"

	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/domain/constants"
	"github.com/google/uuid"
)

type Subscription struct {
	UserID    uuid.UUID          `json:"user_id" db:"user_id"`
	PlanType  constants.PlanType `json:"plan_type" db:"plan_type"`
	StartDate time.Time          `json:"start_date" db:"start_date"`
	EndDate   time.Time          `json:"end_date" db:"end_date"`
	IsActive  bool               `json:"is_active" db:"is_active"`
	CreatedAt time.Time          `json:"created_at" db:"created_at"`
}
