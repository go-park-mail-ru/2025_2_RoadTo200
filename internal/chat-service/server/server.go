package server

import (
	"context"

	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/chat-service/converters"
	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/chat-service/service"
	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/pkg/logger"
	pb "github.com/go-park-mail-ru/2025_2_RoadTo200/backend/proto/chat"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ChatServer struct {
	pb.UnimplementedChatServiceServer
	service *service.ChatService
	logger  logger.Log
}

func NewChatServer(service *service.ChatService, logger logger.Log) *ChatServer {
	return &ChatServer{
		service: service,
		logger:  logger,
	}
}

func (s *ChatServer) SendMessage(ctx context.Context, req *pb.SendMessageRequest) (*pb.MessageResponse, error) {
	senderID, err := uuid.Parse(req.SenderId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid sender_id")
	}

	dtoReq, err := converters.SendMessageRequestToDTO(req)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	message, err := s.service.SendMessage(ctx, senderID, dtoReq)
	if err != nil {
		s.logger.Errorf("Failed to send message: %v", err)
		return nil, status.Error(codes.Internal, "failed to send message")
	}

	return converters.MessageToProto(message), nil
}

func (s *ChatServer) GetMessages(ctx context.Context, req *pb.GetMessagesRequest) (*pb.GetMessagesResponse, error) {
	userID, err := uuid.Parse(req.UserId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid user_id")
	}

	matchID, err := uuid.Parse(req.MatchId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid match_id")
	}

	messages, err := s.service.GetMessages(ctx, userID, matchID, int(req.Limit), int(req.Offset))
	if err != nil {
		s.logger.Errorf("Failed to get messages: %v", err)
		return nil, status.Error(codes.Internal, "failed to get messages")
	}

	return converters.MessagesToProto(messages), nil
}

func (s *ChatServer) MarkAsRead(ctx context.Context, req *pb.MarkAsReadRequest) (*pb.Empty, error) {
	userID, err := uuid.Parse(req.UserId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid user_id")
	}

	matchID, err := uuid.Parse(req.MatchId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid match_id")
	}

	if err := s.service.MarkAsRead(ctx, userID, matchID); err != nil {
		s.logger.Errorf("Failed to mark as read: %v", err)
		return nil, status.Error(codes.Internal, "failed to mark as read")
	}

	return &pb.Empty{}, nil
}

func (s *ChatServer) GetConversations(ctx context.Context, req *pb.GetConversationsRequest) (*pb.GetConversationsResponse, error) {
	userID, err := uuid.Parse(req.UserId)
	if err != nil {
		s.logger.Warnf("Failed to parse user id: %v", err)
		return nil, status.Error(codes.InvalidArgument, "invalid user_id")
	}

	conversations, err := s.service.GetConversations(ctx, userID, req.SearchQuery)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to get conversations")
	}

	return converters.ConversationsToProto(conversations), nil
}

func (s *ChatServer) GetUnreadCount(ctx context.Context, req *pb.GetUnreadCountRequest) (*pb.GetUnreadCountResponse, error) {
	userID, err := uuid.Parse(req.UserId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid user_id")
	}

	count, err := s.service.GetUnreadCount(ctx, userID)
	if err != nil {
		s.logger.Errorf("Failed to get unread count: %v", err)
		return nil, status.Error(codes.Internal, "failed to get unread count")
	}

	return &pb.GetUnreadCountResponse{
		Count: int32(count),
	}, nil
}
