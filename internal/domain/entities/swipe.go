package domain

import (
	"time"

	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/domain/constants"
	"github.com/google/uuid"
)

type Swipe struct {
	ID           uuid.UUID           `json:"id" db:"id"`
	SwiperUserID uuid.UUID           `json:"swiper_user_id" db:"swiper_user_id"`
	TargetUserID uuid.UUID           `json:"target_user_id" db:"target_user_id"`
	SwipeType    constants.SwipeType `json:"swipe_type" db:"swipe_type"`
	CreatedAt    time.Time           `json:"created_at" db:"created_at"`
}

type SwipeRequest struct {
	CardID uuid.UUID `json:"card_id"` // target_user_id
	Action string    `json:"action"`  // 'like', 'dislike'
}

type SwipeResponse struct {
	Match   bool   `json:"match,omitempty"`
	Message string `json:"message"`
}
