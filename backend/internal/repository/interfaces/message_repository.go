package interfaces

import (
	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/domain/dto"
	domain "github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/domain/entities"
	"github.com/google/uuid"
)

type MessageRepository interface {
	Create(message *domain.Message) error
	GetByID(id uuid.UUID) (*domain.Message, error)
	GetByMatchID(matchID uuid.UUID, limit, offset int) ([]domain.Message, error)
	UpdateStatus(id uuid.UUID, status int) error
	Delete(id uuid.UUID) error
	MarkMessagesAsRead(matchID, userID uuid.UUID) error
	GetUnreadCount(userID uuid.UUID) (int, error)
	GetConversations(userID uuid.UUID) ([]dto.Conversation, error)
}
