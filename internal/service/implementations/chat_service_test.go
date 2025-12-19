package service

import (
	"context"
	"testing"
	"time"

	domain "github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/domain/entities"
	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/domain/errors"
	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/handler/dto"
	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/tests/mocks"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestChatService_SendMessage_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockMessageRepo := mocks.NewMockMessageRepository(ctrl)
	mockMatchRepo := mocks.NewMockMatchRepository(ctrl)
	mockLogger := mocks.NewMockLogger()

	// Pass nil for notificationService and redisClient as we're not testing pub/sub
	service := NewChatService(mockMessageRepo, mockMatchRepo, nil, nil, mockLogger)

	ctx := context.Background()
	senderID := uuid.New()
	receiverID := uuid.New()
	matchID := uuid.New()

	match := &domain.Match{
		ID:        matchID,
		User1ID:   senderID,
		User2ID:   receiverID,
		IsActive:  true,
		MatchedAt: time.Now(),
	}

	request := &dto.SendMessageRequest{
		MatchID: matchID,
		Content: "Hello!",
	}

	// Expectations
	mockMatchRepo.EXPECT().
		GetByID(ctx, matchID).
		Return(match, nil)

	mockMessageRepo.EXPECT().
		Create(ctx, gomock.Any()).
		DoAndReturn(func(ctx context.Context, msg *domain.Message) error {
			assert.Equal(t, matchID, msg.MatchID)
			assert.Equal(t, senderID, msg.SenderID)
			assert.Equal(t, receiverID, msg.ReceiverID)
			assert.Equal(t, "Hello!", msg.Content)
			assert.False(t, msg.IsRead)
			return nil
		})

	// Execute
	message, err := service.SendMessage(ctx, senderID, request)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, message)
	assert.Equal(t, "Hello!", message.Content)
}

func TestChatService_SendMessage_MatchNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockMessageRepo := mocks.NewMockMessageRepository(ctrl)
	mockMatchRepo := mocks.NewMockMatchRepository(ctrl)
	mockLogger := mocks.NewMockLogger()

	service := NewChatService(mockMessageRepo, mockMatchRepo, nil, nil, mockLogger)

	ctx := context.Background()
	senderID := uuid.New()
	matchID := uuid.New()

	request := &dto.SendMessageRequest{
		MatchID: matchID,
		Content: "Hello!",
	}

	// Expectations
	mockMatchRepo.EXPECT().
		GetByID(ctx, matchID).
		Return(nil, nil)

	// Execute
	message, err := service.SendMessage(ctx, senderID, request)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, errors.ErrMatchNotFound, err)
	assert.Nil(t, message)
}

func TestChatService_SendMessage_MatchNotActive(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockMessageRepo := mocks.NewMockMessageRepository(ctrl)
	mockMatchRepo := mocks.NewMockMatchRepository(ctrl)
	mockLogger := mocks.NewMockLogger()

	service := NewChatService(mockMessageRepo, mockMatchRepo, nil, nil, mockLogger)

	ctx := context.Background()
	senderID := uuid.New()
	receiverID := uuid.New()
	matchID := uuid.New()

	match := &domain.Match{
		ID:        matchID,
		User1ID:   senderID,
		User2ID:   receiverID,
		IsActive:  false, // Not active
		MatchedAt: time.Now(),
	}

	request := &dto.SendMessageRequest{
		MatchID: matchID,
		Content: "Hello!",
	}

	// Expectations
	mockMatchRepo.EXPECT().
		GetByID(ctx, matchID).
		Return(match, nil)

	// Execute
	message, err := service.SendMessage(ctx, senderID, request)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, errors.ErrMatchNotActive, err)
	assert.Nil(t, message)
}

func TestChatService_SendMessage_NotParticipant(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockMessageRepo := mocks.NewMockMessageRepository(ctrl)
	mockMatchRepo := mocks.NewMockMatchRepository(ctrl)
	mockLogger := mocks.NewMockLogger()

	service := NewChatService(mockMessageRepo, mockMatchRepo, nil, nil, mockLogger)

	ctx := context.Background()
	senderID := uuid.New() // Not part of match
	user1ID := uuid.New()
	user2ID := uuid.New()
	matchID := uuid.New()

	match := &domain.Match{
		ID:        matchID,
		User1ID:   user1ID,
		User2ID:   user2ID,
		IsActive:  true,
		MatchedAt: time.Now(),
	}

	request := &dto.SendMessageRequest{
		MatchID: matchID,
		Content: "Hello!",
	}

	// Expectations
	mockMatchRepo.EXPECT().
		GetByID(ctx, matchID).
		Return(match, nil)

	// Execute
	message, err := service.SendMessage(ctx, senderID, request)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, errors.ErrNotMatchParticipant, err)
	assert.Nil(t, message)
}

func TestChatService_GetMessages_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockMessageRepo := mocks.NewMockMessageRepository(ctrl)
	mockMatchRepo := mocks.NewMockMatchRepository(ctrl)
	mockLogger := mocks.NewMockLogger()

	service := NewChatService(mockMessageRepo, mockMatchRepo, nil, nil, mockLogger)

	ctx := context.Background()
	userID := uuid.New()
	otherUserID := uuid.New()
	matchID := uuid.New()

	match := &domain.Match{
		ID:        matchID,
		User1ID:   userID,
		User2ID:   otherUserID,
		IsActive:  true,
		MatchedAt: time.Now(),
	}

	messages := []domain.Message{
		{
			ID:       uuid.New(),
			MatchID:  matchID,
			SenderID: userID,
			Content:  "Hello!",
		},
	}

	// Expectations
	mockMatchRepo.EXPECT().
		GetByID(ctx, matchID).
		Return(match, nil)

	mockMessageRepo.EXPECT().
		GetByMatchID(ctx, matchID, 20, 0).
		Return(messages, nil)

	// Execute
	result, err := service.GetMessages(ctx, userID, matchID, 20, 0)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 1, len(result))
}

func TestChatService_GetMessages_NotParticipant(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockMessageRepo := mocks.NewMockMessageRepository(ctrl)
	mockMatchRepo := mocks.NewMockMatchRepository(ctrl)
	mockLogger := mocks.NewMockLogger()

	service := NewChatService(mockMessageRepo, mockMatchRepo, nil, nil, mockLogger)

	ctx := context.Background()
	userID := uuid.New() // Not part of match
	user1ID := uuid.New()
	user2ID := uuid.New()
	matchID := uuid.New()

	match := &domain.Match{
		ID:        matchID,
		User1ID:   user1ID,
		User2ID:   user2ID,
		IsActive:  true,
		MatchedAt: time.Now(),
	}

	// Expectations
	mockMatchRepo.EXPECT().
		GetByID(ctx, matchID).
		Return(match, nil)

	// Execute
	result, err := service.GetMessages(ctx, userID, matchID, 20, 0)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, errors.ErrNotMatchParticipant, err)
	assert.Nil(t, result)
}

func TestChatService_MarkAsRead_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockMessageRepo := mocks.NewMockMessageRepository(ctrl)
	mockMatchRepo := mocks.NewMockMatchRepository(ctrl)
	mockLogger := mocks.NewMockLogger()

	service := NewChatService(mockMessageRepo, mockMatchRepo, nil, nil, mockLogger)

	ctx := context.Background()
	userID := uuid.New()
	otherUserID := uuid.New()
	matchID := uuid.New()

	match := &domain.Match{
		ID:        matchID,
		User1ID:   userID,
		User2ID:   otherUserID,
		IsActive:  true,
		MatchedAt: time.Now(),
	}

	// Expectations
	mockMatchRepo.EXPECT().
		GetByID(ctx, matchID).
		Return(match, nil)

	mockMessageRepo.EXPECT().
		MarkAsRead(ctx, matchID, userID).
		Return(nil)

	// Execute
	err := service.MarkAsRead(ctx, userID, matchID)

	// Assert
	assert.NoError(t, err)
}

func TestChatService_MarkAsRead_NotParticipant(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockMessageRepo := mocks.NewMockMessageRepository(ctrl)
	mockMatchRepo := mocks.NewMockMatchRepository(ctrl)
	mockLogger := mocks.NewMockLogger()

	service := NewChatService(mockMessageRepo, mockMatchRepo, nil, nil, mockLogger)

	ctx := context.Background()
	userID := uuid.New() // Not part of match
	user1ID := uuid.New()
	user2ID := uuid.New()
	matchID := uuid.New()

	match := &domain.Match{
		ID:        matchID,
		User1ID:   user1ID,
		User2ID:   user2ID,
		IsActive:  true,
		MatchedAt: time.Now(),
	}

	// Expectations
	mockMatchRepo.EXPECT().
		GetByID(ctx, matchID).
		Return(match, nil)

	// Execute
	err := service.MarkAsRead(ctx, userID, matchID)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, errors.ErrNotMatchParticipant, err)
}

func TestChatService_GetConversations_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockMessageRepo := mocks.NewMockMessageRepository(ctrl)
	mockMatchRepo := mocks.NewMockMatchRepository(ctrl)
	mockLogger := mocks.NewMockLogger()

	service := NewChatService(mockMessageRepo, mockMatchRepo, nil, nil, mockLogger)

	ctx := context.Background()
	userID := uuid.New()

	conversations := []domain.Conversation{
		{
			MatchID:       uuid.New(),
			OtherUserID:   uuid.New(),
			OtherUserName: "Test User",
		},
	}

	// Expectations
	mockMessageRepo.EXPECT().
		GetConversations(ctx, userID, "").
		Return(conversations, nil)

	// Execute
	result, err := service.GetConversations(ctx, userID, "")

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 1, len(result))
}

func TestChatService_GetUnreadCount_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockMessageRepo := mocks.NewMockMessageRepository(ctrl)
	mockMatchRepo := mocks.NewMockMatchRepository(ctrl)
	mockLogger := mocks.NewMockLogger()

	service := NewChatService(mockMessageRepo, mockMatchRepo, nil, nil, mockLogger)

	ctx := context.Background()
	userID := uuid.New()

	// Expectations
	mockMessageRepo.EXPECT().
		GetUnreadCount(ctx, userID).
		Return(5, nil)

	// Execute
	count, err := service.GetUnreadCount(ctx, userID)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, 5, count)
}

func TestChatService_GetConversations_WithSearch(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockMessageRepo := mocks.NewMockMessageRepository(ctrl)
	mockMatchRepo := mocks.NewMockMatchRepository(ctrl)
	mockLogger := mocks.NewMockLogger()

	service := NewChatService(mockMessageRepo, mockMatchRepo, nil, nil, mockLogger)

	ctx := context.Background()
	userID := uuid.New()
	searchQuery := "test"

	conversations := []domain.Conversation{
		{
			MatchID:       uuid.New(),
			OtherUserID:   uuid.New(),
			OtherUserName: "Test User",
		},
	}

	// Expectations
	mockMessageRepo.EXPECT().
		GetConversations(ctx, userID, searchQuery).
		Return(conversations, nil)

	// Execute
	result, err := service.GetConversations(ctx, userID, searchQuery)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 1, len(result))
	assert.Equal(t, "Test User", result[0].OtherUserName)
}
