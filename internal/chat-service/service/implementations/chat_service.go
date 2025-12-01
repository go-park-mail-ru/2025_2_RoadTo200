package service

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/chat-service/repository/interfaces"
	service "github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/chat-service/service/interfaces"
	interfaces2 "github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/core-service/repository/interfaces"
	domain "github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/domain/entities"
	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/domain/errors"
	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/handler/dto"
	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/pkg/logger"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

var _ service.ChatService = (*ChatService)(nil)

type ChatService struct {
	messageRepo interfaces.MessageRepository
	matchRepo   interfaces2.MatchRepository
	redisClient *redis.Client
	logger      logger.Log
}

func NewChatService(
	messageRepo interfaces.MessageRepository,
	matchRepo interfaces2.MatchRepository,
	redisClient *redis.Client,
	logger logger.Log,
) *ChatService {
	return &ChatService{
		messageRepo: messageRepo,
		matchRepo:   matchRepo,
		redisClient: redisClient,
		logger:      logger,
	}
}

// SendMessage validates match, saves message, and publishes to Redis
func (s *ChatService) SendMessage(ctx context.Context, senderID uuid.UUID, req *dto.SendMessageRequest) (*domain.Message, error) {
	s.logger.Trace("ChatService.SendMessage")

	// 1. Validate that match exists and is active
	match, err := s.matchRepo.GetByID(ctx, req.MatchID)
	if err != nil {
		s.logger.Errorf("Failed to get match: %v", err)
		return nil, errors.ErrInternalError
	}
	if match == nil {
		return nil, errors.ErrMatchNotFound
	}
	if !match.IsActive {
		return nil, errors.ErrMatchNotActive
	}

	// 2. Validate that sender is part of the match
	var receiverID uuid.UUID
	if match.User1ID == senderID {
		receiverID = match.User2ID
	} else if match.User2ID == senderID {
		receiverID = match.User1ID
	} else {
		return nil, errors.ErrNotMatchParticipant
	}

	// 3. Create message
	message := &domain.Message{
		MatchID:    req.MatchID,
		SenderID:   senderID,
		ReceiverID: receiverID,
		Content:    req.Content,
		IsRead:     false,
	}

	// 4. Save to database
	if err := s.messageRepo.Create(ctx, message); err != nil {
		s.logger.Errorf("Failed to create message: %v", err)
		return nil, errors.ErrInternalError
	}

	// 5. Publish to Redis Pub/Sub for real-time delivery
	if err := s.publishMessage(ctx, message); err != nil {
		s.logger.Warnf("Failed to publish message to Redis: %v", err)
		// Don't fail the request, message is already saved
	}

	return message, nil
}

// GetMessages returns message history with validation
func (s *ChatService) GetMessages(ctx context.Context, userID, matchID uuid.UUID, limit, offset int) ([]domain.Message, error) {
	s.logger.Trace("ChatService.GetMessages")

	// 1. Validate that user is part of the match
	match, err := s.matchRepo.GetByID(ctx, matchID)
	if err != nil {
		s.logger.Errorf("Failed to get match: %v", err)
		return nil, errors.ErrInternalError
	}
	if match == nil {
		return nil, errors.ErrMatchNotFound
	}

	if match.User1ID != userID && match.User2ID != userID {
		return nil, errors.ErrNotMatchParticipant
	}

	// 2. Get messages
	messages, err := s.messageRepo.GetByMatchID(ctx, matchID, limit, offset)
	if err != nil {
		s.logger.Errorf("Failed to get messages: %v", err)
		return nil, errors.ErrInternalError
	}

	return messages, nil
}

// MarkAsRead marks messages as read and publishes event
func (s *ChatService) MarkAsRead(ctx context.Context, userID, matchID uuid.UUID) error {
	s.logger.Trace("ChatService.MarkAsRead")

	// 1. Validate that user is part of the match
	match, err := s.matchRepo.GetByID(ctx, matchID)
	if err != nil {
		s.logger.Errorf("Failed to get match: %v", err)
		return errors.ErrInternalError
	}
	if match == nil {
		return errors.ErrMatchNotFound
	}

	if match.User1ID != userID && match.User2ID != userID {
		return errors.ErrNotMatchParticipant
	}

	// 2. Mark messages as read
	if err := s.messageRepo.MarkAsRead(ctx, matchID, userID); err != nil {
		s.logger.Errorf("Failed to mark messages as read: %v", err)
		return errors.ErrInternalError
	}

	// 3. Publish "read" event to Redis
	readEvent := domain.ChatMessage{
		Type:    "read",
		MatchID: matchID,
	}
	if err := s.publishEvent(ctx, matchID, readEvent); err != nil {
		s.logger.Warnf("Failed to publish read event: %v", err)
	}

	return nil
}

// GetConversations returns all conversations for user
func (s *ChatService) GetConversations(ctx context.Context, userID uuid.UUID, searchQuery string) ([]domain.Conversation, error) {
	s.logger.Trace("ChatService.GetConversations")

	conversations, err := s.messageRepo.GetConversations(ctx, userID, searchQuery)
	if err != nil {
		return nil, errors.ErrInternalError
	}

	return conversations, nil
}

// GetUnreadCount returns total unread count
func (s *ChatService) GetUnreadCount(ctx context.Context, userID uuid.UUID) (int, error) {
	s.logger.Trace("ChatService.GetUnreadCount")

	count, err := s.messageRepo.GetUnreadCount(ctx, userID)
	if err != nil {
		s.logger.Errorf("Failed to get unread count: %v", err)
		return 0, errors.ErrInternalError
	}

	return count, nil
}

// publishMessage publishes new message to Redis Pub/Sub
func (s *ChatService) publishMessage(ctx context.Context, message *domain.Message) error {
	chatMsg := domain.ChatMessage{
		Type:      "message",
		MessageID: message.ID,
		MatchID:   message.MatchID,
		SenderID:  message.SenderID,
		Content:   message.Content,
		CreatedAt: message.CreatedAt,
	}

	return s.publishEvent(ctx, message.MatchID, chatMsg)
}

// publishEvent publishes event to Redis channel
func (s *ChatService) publishEvent(ctx context.Context, matchID uuid.UUID, event domain.ChatMessage) error {
	// Get match to find participants
	match, err := s.matchRepo.GetByID(ctx, matchID)
	if err != nil {
		return fmt.Errorf("failed to get match: %w", err)
	}
	if match == nil {
		return fmt.Errorf("match not found")
	}

	data, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal event: %w", err)
	}

	// Publish to both users
	users := []uuid.UUID{match.User1ID, match.User2ID}
	for _, userID := range users {
		channel := fmt.Sprintf("chat:user:%s", userID.String())
		if err := s.redisClient.Publish(ctx, channel, data).Err(); err != nil {
			s.logger.Warnf("Failed to publish to channel %s: %v", channel, err)
		} else {
			s.logger.Debugf("Published event to channel %s: %s", channel, event.Type)
		}
	}

	return nil
}
