package dto

import "github.com/google/uuid"

// SendMessageRequest - request to send a message
type SendMessageRequest struct {
	MatchID uuid.UUID `json:"match_id" binding:"required"`
	Content string    `json:"content" binding:"required,min=1,max=1000"`
}

// SendMessageResponse - response after sending a message
type SendMessageResponse struct {
	MessageID uuid.UUID `json:"message_id"`
	CreatedAt string    `json:"created_at"`
}
