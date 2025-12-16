package dto

//go:generate easyjson -all -no_std_marshalers swipe_dto.go

// SwipeRequest represents swipe request
type SwipeRequest struct {
	CardID string `json:"card_id" example:"550e8400-e29b-41d4-a716-446655440000"`
	Action string `json:"action" example:"like" enums:"like,dislike,super_like"`
}

// SwipeResponse represents swipe response
type SwipeResponse struct {
	IsMatch bool   `json:"is_match" example:"true"`
	MatchID string `json:"match_id,omitempty" example:"660e8400-e29b-41d4-a716-446655440000"`
	Message string `json:"message,omitempty" example:"It's a match!"`
}
