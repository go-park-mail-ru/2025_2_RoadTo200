package service

import (
	"context"

	domain "github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/chat-service/domain/entities"
	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/handler/dto"
	"github.com/google/uuid"
)

type ChatService interface {
	// SendMessage validates match, saves message to DB, publishes to Redis
	SendMessage(ctx context.Context, senderID uuid.UUID, req *dto.SendMessageRequest) (*domain.Message, error)

	// GetMessages returns message history for a match (with pagination)
	GetMessages(ctx context.Context, userID, matchID uuid.UUID, limit, offset int) ([]domain.Message, error)

	// MarkAsRead marks all messages in match as read by user
	MarkAsRead(ctx context.Context, userID, matchID uuid.UUID) error

	// GetConversations returns list of all chats with last message
	GetConversations(ctx context.Context, userID uuid.UUID, searchQuery string) ([]domain.Conversation, error)

	// GetUnreadCount returns total unread message count for user
	GetUnreadCount(ctx context.Context, userID uuid.UUID) (int, error)
}
