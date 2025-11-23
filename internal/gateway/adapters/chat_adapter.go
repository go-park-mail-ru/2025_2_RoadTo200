package adapters

import (
	"context"
	"time"

	domain "github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/domain/entities"
	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/handler/dto"
	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/logger"
	pb "github.com/go-park-mail-ru/2025_2_RoadTo200/backend/proto/chat"
	"github.com/google/uuid"
)

type ChatServiceAdapter struct {
	client pb.ChatServiceClient
	logger logger.Log
}

func NewChatServiceAdapter(client pb.ChatServiceClient, logger logger.Log) *ChatServiceAdapter {
	return &ChatServiceAdapter{
		client: client,
		logger: logger,
	}
}

func (a *ChatServiceAdapter) SendMessage(ctx context.Context, senderID uuid.UUID, req *dto.SendMessageRequest) (*domain.Message, error) {
	pbReq := &pb.SendMessageRequest{
		SenderId: senderID.String(),
		MatchId:  req.MatchID.String(),
		Content:  req.Content,
	}

	resp, err := a.client.SendMessage(ctx, pbReq)
	if err != nil {
		return nil, err
	}

	return protoToMessage(resp)
}

func (a *ChatServiceAdapter) GetMessages(ctx context.Context, userID, matchID uuid.UUID, limit, offset int) ([]domain.Message, error) {
	req := &pb.GetMessagesRequest{
		UserId:  userID.String(),
		MatchId: matchID.String(),
		Limit:   int32(limit),
		Offset:  int32(offset),
	}

	resp, err := a.client.GetMessages(ctx, req)
	if err != nil {
		return nil, err
	}

	var messages []domain.Message
	for _, msg := range resp.Messages {
		m, err := protoToMessage(msg)
		if err != nil {
			continue
		}
		messages = append(messages, *m)
	}

	return messages, nil
}

func (a *ChatServiceAdapter) MarkAsRead(ctx context.Context, userID, matchID uuid.UUID) error {
	req := &pb.MarkAsReadRequest{
		UserId:  userID.String(),
		MatchId: matchID.String(),
	}

	_, err := a.client.MarkAsRead(ctx, req)
	return err
}

func (a *ChatServiceAdapter) GetConversations(ctx context.Context, userID uuid.UUID) ([]domain.Conversation, error) {
	req := &pb.GetConversationsRequest{
		UserId: userID.String(),
	}

	resp, err := a.client.GetConversations(ctx, req)
	if err != nil {
		return nil, err
	}

	var conversations []domain.Conversation
	for _, conv := range resp.Conversations {
		c, err := protoToConversation(conv)
		if err != nil {
			continue
		}
		conversations = append(conversations, *c)
	}

	return conversations, nil
}

func (a *ChatServiceAdapter) GetUnreadCount(ctx context.Context, userID uuid.UUID) (int, error) {
	req := &pb.GetUnreadCountRequest{
		UserId: userID.String(),
	}

	resp, err := a.client.GetUnreadCount(ctx, req)
	if err != nil {
		return 0, err
	}

	return int(resp.Count), nil
}

// Helpers

func protoToMessage(pbMsg *pb.MessageResponse) (*domain.Message, error) {
	id, _ := uuid.Parse(pbMsg.Id)
	matchID, _ := uuid.Parse(pbMsg.MatchId)
	senderID, _ := uuid.Parse(pbMsg.SenderId)
	receiverID, _ := uuid.Parse(pbMsg.ReceiverId)

	return &domain.Message{
		ID:         id,
		MatchID:    matchID,
		SenderID:   senderID,
		ReceiverID: receiverID,
		Content:    pbMsg.Content,
		IsRead:     pbMsg.IsRead,
		CreatedAt:  pbMsg.CreatedAt.AsTime(),
	}, nil
}

func protoToConversation(pbConv *pb.Conversation) (*domain.Conversation, error) {
	matchID, _ := uuid.Parse(pbConv.MatchId)
	otherUserID, _ := uuid.Parse(pbConv.OtherUserId)

	lastMessageTime := time.Time{}
	if pbConv.LastMessageTime != nil {
		lastMessageTime = pbConv.LastMessageTime.AsTime()
	}

	return &domain.Conversation{
		MatchID:         matchID,
		OtherUserID:     otherUserID,
		OtherUserName:   pbConv.OtherUserName,
		OtherUserPhoto:  pbConv.OtherUserPhoto,
		LastMessage:     pbConv.LastMessage,
		LastMessageTime: lastMessageTime,
		UnreadCount:     int(pbConv.UnreadCount),
	}, nil
}
