package interfaces

import (
	"context"

	domain "github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/chat-service/domain/entities"
	"github.com/google/uuid"
)

type MessageRepository interface {
	// Create new message
	Create(ctx context.Context, message *domain.Message) error

	// GetByMatchID Get message history for a match (with pagination)
	GetByMatchID(ctx context.Context, matchID uuid.UUID, limit, offset int) ([]domain.Message, error)

	// MarkAsRead Mark all messages in match as read by receiver
	MarkAsRead(ctx context.Context, matchID, receiverID uuid.UUID) error

	// GetUnreadCount Get unread message count for user
	GetUnreadCount(ctx context.Context, userID uuid.UUID) (int, error)

	// GetLastMessage Get user last message
	GetLastMessage(ctx context.Context, matchIDs []uuid.UUID) ([]domain.LastMessage, error)

	// GetUnreadCounts Get user unread message count
	GetUnreadCounts(ctx context.Context, userID uuid.UUID, matchIDs []uuid.UUID) ([]domain.UnreadCount, error)

	// GetChats Get user chats
	GetChats(ctx context.Context, userID uuid.UUID) ([]domain.Conversation, error)

	// GetConversations Get all conversations (matches with last message)
	// Deprecated: old base structure
	GetConversations(ctx context.Context, userID uuid.UUID, searchQuery string) ([]domain.Conversation, error)
}
