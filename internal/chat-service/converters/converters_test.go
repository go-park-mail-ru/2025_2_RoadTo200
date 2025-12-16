package converters

import (
	"testing"
	"time"

	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/chat-service/domain/entities"
	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/handler/dto"
	pb "github.com/go-park-mail-ru/2025_2_RoadTo200/backend/proto/chat"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSendMessageRequestToDTO(t *testing.T) {
	tests := []struct {
		name        string
		request     *pb.SendMessageRequest
		expectError bool
		validate    func(t *testing.T, result *dto.SendMessageRequest)
	}{
		{
			name: "valid_request",
			request: &pb.SendMessageRequest{
				MatchId: "550e8400-e29b-41d4-a716-446655440000",
				Content: "Hello, world!",
			},
			expectError: false,
			validate: func(t *testing.T, result *dto.SendMessageRequest) {
				assert.Equal(t, "Hello, world!", result.Content)
				assert.NotEqual(t, uuid.Nil, result.MatchID)
			},
		},
		{
			name: "invalid_uuid",
			request: &pb.SendMessageRequest{
				MatchId: "not-a-uuid",
				Content: "Hello",
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := SendMessageRequestToDTO(tt.request)

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				require.NoError(t, err)
				require.NotNil(t, result)
				if tt.validate != nil {
					tt.validate(t, result)
				}
			}
		})
	}
}

func TestMessageToProto(t *testing.T) {
	msgID := uuid.New()
	matchID := uuid.New()
	senderID := uuid.New()
	receiverID := uuid.New()
	createdAt := time.Now()

	message := &domain.Message{
		ID:         msgID,
		MatchID:    matchID,
		SenderID:   senderID,
		ReceiverID: receiverID,
		Content:    "Test message",
		IsRead:     false,
		CreatedAt:  createdAt,
	}

	result := MessageToProto(message)

	assert.NotNil(t, result)
	assert.Equal(t, msgID.String(), result.Id)
	assert.Equal(t, matchID.String(), result.MatchId)
	assert.Equal(t, senderID.String(), result.SenderId)
	assert.Equal(t, receiverID.String(), result.ReceiverId)
	assert.Equal(t, "Test message", result.Content)
	assert.False(t, result.IsRead)
	assert.NotNil(t, result.CreatedAt)
}

func TestMessagesToProto(t *testing.T) {
	tests := []struct {
		name     string
		messages []domain.Message
		expected int
	}{
		{
			name:     "empty_list",
			messages: []domain.Message{},
			expected: 0,
		},
		{
			name: "single_message",
			messages: []domain.Message{
				{
					ID:         uuid.New(),
					MatchID:    uuid.New(),
					SenderID:   uuid.New(),
					ReceiverID: uuid.New(),
					Content:    "Message 1",
					CreatedAt:  time.Now(),
				},
			},
			expected: 1,
		},
		{
			name: "multiple_messages",
			messages: []domain.Message{
				{
					ID:         uuid.New(),
					MatchID:    uuid.New(),
					SenderID:   uuid.New(),
					ReceiverID: uuid.New(),
					Content:    "Message 1",
					CreatedAt:  time.Now(),
				},
				{
					ID:         uuid.New(),
					MatchID:    uuid.New(),
					SenderID:   uuid.New(),
					ReceiverID: uuid.New(),
					Content:    "Message 2",
					CreatedAt:  time.Now(),
				},
			},
			expected: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := MessagesToProto(tt.messages)

			assert.NotNil(t, result)
			assert.Equal(t, tt.expected, len(result.Messages))

			for i, msg := range result.Messages {
				assert.Equal(t, tt.messages[i].Content, msg.Content)
			}
		})
	}
}

func TestConversationToProto(t *testing.T) {
	tests := []struct {
		name string
		conv *domain.Conversation
	}{
		{
			name: "full_conversation",
			conv: &domain.Conversation{
				MatchID:         uuid.New(),
				OtherUserID:     uuid.New(),
				OtherUserName:   "John Doe",
				OtherUserPhoto:  "https://example.com/photo.jpg",
				LastMessage:     "Hello!",
				LastMessageTime: time.Now(),
				UnreadCount:     5,
			},
		},
		{
			name: "conversation_with_zero_time",
			conv: &domain.Conversation{
				MatchID:         uuid.New(),
				OtherUserID:     uuid.New(),
				OtherUserName:   "Jane Doe",
				OtherUserPhoto:  "",
				LastMessage:     "",
				LastMessageTime: time.Time{},
				UnreadCount:     0,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ConversationToProto(tt.conv)

			assert.NotNil(t, result)
			assert.Equal(t, tt.conv.MatchID.String(), result.MatchId)
			assert.Equal(t, tt.conv.OtherUserID.String(), result.OtherUserId)
			assert.Equal(t, tt.conv.OtherUserName, result.OtherUserName)
			assert.Equal(t, tt.conv.OtherUserPhoto, result.OtherUserPhoto)
			assert.Equal(t, tt.conv.LastMessage, result.LastMessage)
			assert.Equal(t, int32(tt.conv.UnreadCount), result.UnreadCount)
			assert.NotNil(t, result.LastMessageTime)
		})
	}
}

func TestConversationsToProto(t *testing.T) {
	conversations := []domain.Conversation{
		{
			MatchID:         uuid.New(),
			OtherUserID:     uuid.New(),
			OtherUserName:   "User 1",
			LastMessage:     "Hi",
			LastMessageTime: time.Now(),
			UnreadCount:     1,
		},
		{
			MatchID:         uuid.New(),
			OtherUserID:     uuid.New(),
			OtherUserName:   "User 2",
			LastMessage:     "Hello",
			LastMessageTime: time.Now(),
			UnreadCount:     0,
		},
	}

	result := ConversationsToProto(conversations)

	assert.NotNil(t, result)
	assert.Equal(t, 2, len(result.Conversations))
	assert.Equal(t, "User 1", result.Conversations[0].OtherUserName)
	assert.Equal(t, "User 2", result.Conversations[1].OtherUserName)
}
