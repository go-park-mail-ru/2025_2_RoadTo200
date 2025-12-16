package domain

//go:generate easyjson -all -no_std_marshalers interest.go

import (
	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/domain/constants"
	"github.com/google/uuid"
)

type Interest struct {
	UserID uuid.UUID              `json:"user_id" db:"user_id"`
	Theme  constants.InterestType `json:"theme" db:"theme"`
}
