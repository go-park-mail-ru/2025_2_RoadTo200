package interfaces

import (
	"context"

	domain "github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/domain/entities"
	"github.com/google/uuid"
)

type MessageRepository interface {
	// Create new message
	Create(ctx context.Context, message *domain.Message) error

	// Get message history for a match (with pagination)
	GetByMatchID(ctx context.Context, matchID uuid.UUID, limit, offset int) ([]domain.Message, error)

	// Mark all messages in match as read by receiver
	MarkAsRead(ctx context.Context, matchID, receiverID uuid.UUID) error

	// Get unread message count for user
	GetUnreadCount(ctx context.Context, userID uuid.UUID) (int, error)

	// Get all conversations (matches with last message)
	GetConversations(ctx context.Context, userID uuid.UUID, searchQuery string) ([]domain.Conversation, error)

	// HasMessages checks if match has at least one message
	HasMessages(ctx context.Context, matchID uuid.UUID) (bool, error)
}
