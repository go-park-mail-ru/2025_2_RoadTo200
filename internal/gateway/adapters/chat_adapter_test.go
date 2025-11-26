package adapters

import (
	"context"
	"testing"

	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/handler/dto"
	pb "github.com/go-park-mail-ru/2025_2_RoadTo200/backend/proto/chat"
	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/tests/mocks"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestChatServiceAdapter_SendMessage(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mocks.NewMockChatServiceClient(ctrl)
	mockLogger := mocks.NewMockLogger()
	adapter := NewChatServiceAdapter(mockClient, mockLogger)

	senderID := uuid.New()
	matchID := uuid.New()
	content := "Hello"
	messageID := uuid.New()

	req := &dto.SendMessageRequest{
		MatchID: matchID,
		Content: content,
	}

	resp := &pb.MessageResponse{
		Id:         messageID.String(),
		MatchId:    matchID.String(),
		SenderId:   senderID.String(),
		ReceiverId: uuid.New().String(),
		Content:    content,
		IsRead:     false,
		CreatedAt:  timestamppb.Now(),
	}

	mockClient.EXPECT().
		SendMessage(gomock.Any(), &pb.SendMessageRequest{
			SenderId: senderID.String(),
			MatchId:  matchID.String(),
			Content:  content,
		}).
		Return(resp, nil)

	msg, err := adapter.SendMessage(context.Background(), senderID, req)
	assert.NoError(t, err)
	assert.NotNil(t, msg)
	assert.Equal(t, messageID, msg.ID)
	assert.Equal(t, content, msg.Content)
}

func TestChatServiceAdapter_GetMessages(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mocks.NewMockChatServiceClient(ctrl)
	mockLogger := mocks.NewMockLogger()
	adapter := NewChatServiceAdapter(mockClient, mockLogger)

	userID := uuid.New()
	matchID := uuid.New()
	limit := 10
	offset := 0

	resp := &pb.GetMessagesResponse{
		Messages: []*pb.MessageResponse{
			{
				Id:        uuid.New().String(),
				MatchId:   matchID.String(),
				SenderId:  userID.String(),
				Content:   "Hi",
				CreatedAt: timestamppb.Now(),
			},
		},
	}

	mockClient.EXPECT().
		GetMessages(gomock.Any(), &pb.GetMessagesRequest{
			UserId:  userID.String(),
			MatchId: matchID.String(),
			Limit:   int32(limit),
			Offset:  int32(offset),
		}).
		Return(resp, nil)

	msgs, err := adapter.GetMessages(context.Background(), userID, matchID, limit, offset)
	assert.NoError(t, err)
	assert.Len(t, msgs, 1)
	assert.Equal(t, "Hi", msgs[0].Content)
}

func TestChatServiceAdapter_MarkAsRead(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mocks.NewMockChatServiceClient(ctrl)
	mockLogger := mocks.NewMockLogger()
	adapter := NewChatServiceAdapter(mockClient, mockLogger)

	userID := uuid.New()
	matchID := uuid.New()

	mockClient.EXPECT().
		MarkAsRead(gomock.Any(), &pb.MarkAsReadRequest{
			UserId:  userID.String(),
			MatchId: matchID.String(),
		}).
		Return(&pb.Empty{}, nil)

	err := adapter.MarkAsRead(context.Background(), userID, matchID)
	assert.NoError(t, err)
}

func TestChatServiceAdapter_GetConversations(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mocks.NewMockChatServiceClient(ctrl)
	mockLogger := mocks.NewMockLogger()
	adapter := NewChatServiceAdapter(mockClient, mockLogger)

	userID := uuid.New()
	searchQuery := "test"

	resp := &pb.GetConversationsResponse{
		Conversations: []*pb.Conversation{
			{
				MatchId:       uuid.New().String(),
				OtherUserId:   uuid.New().String(),
				OtherUserName: "User",
				LastMessage:   "Hi",
				UnreadCount:   2,
			},
		},
	}

	mockClient.EXPECT().
		GetConversations(gomock.Any(), &pb.GetConversationsRequest{
			UserId:      userID.String(),
			SearchQuery: searchQuery,
		}).
		Return(resp, nil)

	convs, err := adapter.GetConversations(context.Background(), userID, searchQuery)
	assert.NoError(t, err)
	assert.Len(t, convs, 1)
	assert.Equal(t, 2, convs[0].UnreadCount)
}

func TestChatServiceAdapter_GetUnreadCount(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mocks.NewMockChatServiceClient(ctrl)
	mockLogger := mocks.NewMockLogger()
	adapter := NewChatServiceAdapter(mockClient, mockLogger)

	userID := uuid.New()

	mockClient.EXPECT().
		GetUnreadCount(gomock.Any(), &pb.GetUnreadCountRequest{
			UserId: userID.String(),
		}).
		Return(&pb.GetUnreadCountResponse{Count: 5}, nil)

	count, err := adapter.GetUnreadCount(context.Background(), userID)
	assert.NoError(t, err)
	assert.Equal(t, 5, count)
}
