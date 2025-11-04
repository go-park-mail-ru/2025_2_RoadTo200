package service

import (
	"github.com/google/uuid"
)

type SwipeService interface {
	ProcessSwipe(swiperID uuid.UUID, request *SwipeRequest) (*SwipeResponse, error)
}

type SwipeRequest struct {
	CardID uuid.UUID `json:"card_id"` // target_user_id
	Action string    `json:"action"`  // 'like', 'dislike'
}

type SwipeResponse struct {
	Match   bool   `json:"match,omitempty"`
	Message string `json:"message"`
}
