package converters

import (
	"time"

	domain "github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/chat-service/domain/entities"
	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/handler/dto"
	pb "github.com/go-park-mail-ru/2025_2_RoadTo200/backend/proto/chat"
	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func SendMessageRequestToDTO(req *pb.SendMessageRequest) (*dto.SendMessageRequest, error) {
	matchID, err := uuid.Parse(req.MatchId)
	if err != nil {
		return nil, err
	}
	return &dto.SendMessageRequest{
		MatchID: matchID,
		Content: req.Content,
	}, nil
}

func MessageToProto(msg *domain.Message) *pb.MessageResponse {
	return &pb.MessageResponse{
		Id:         msg.ID.String(),
		MatchId:    msg.MatchID.String(),
		SenderId:   msg.SenderID.String(),
		ReceiverId: msg.ReceiverID.String(),
		Content:    msg.Content,
		IsRead:     msg.IsRead,
		CreatedAt:  timestamppb.New(msg.CreatedAt),
	}
}

func MessagesToProto(messages []domain.Message) *pb.GetMessagesResponse {
	var protoMessages []*pb.MessageResponse
	for _, msg := range messages {
		protoMessages = append(protoMessages, MessageToProto(&msg))
	}
	return &pb.GetMessagesResponse{
		Messages: protoMessages,
	}
}

func ConversationToProto(conv *domain.Conversation) *pb.Conversation {
	lastMessageTime := timestamppb.New(conv.LastMessageTime)
	if conv.LastMessageTime.IsZero() {
		lastMessageTime = timestamppb.New(time.Time{})
	}

	return &pb.Conversation{
		MatchId:         conv.MatchID.String(),
		OtherUserId:     conv.OtherUserID.String(),
		OtherUserName:   conv.OtherUserName,
		OtherUserPhoto:  conv.OtherUserPhoto,
		LastMessage:     conv.LastMessage,
		LastMessageTime: lastMessageTime,
		UnreadCount:     int32(conv.UnreadCount),
	}
}

func ConversationsToProto(conversations []domain.Conversation) *pb.GetConversationsResponse {
	var protoConversations []*pb.Conversation
	for _, conv := range conversations {
		protoConversations = append(protoConversations, ConversationToProto(&conv))
	}
	return &pb.GetConversationsResponse{
		Conversations: protoConversations,
	}
}
