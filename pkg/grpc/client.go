package grpc

import (
	"fmt"

	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/proto/auth"
	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/proto/chat"
	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/proto/core"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func NewAuthClient(addr string) (auth.AuthServiceClient, *grpc.ClientConn, error) {
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, nil, fmt.Errorf("failed to connect to auth service: %w", err)
	}

	return auth.NewAuthServiceClient(conn), conn, nil
}

func NewCoreClient(addr string) (core.CoreServiceClient, *grpc.ClientConn, error) {
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, nil, fmt.Errorf("failed to connect to core service: %w", err)
	}

	return core.NewCoreServiceClient(conn), conn, nil
}

func NewChatClient(addr string) (chat.ChatServiceClient, *grpc.ClientConn, error) {
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, nil, fmt.Errorf("failed to connect to chat service: %w", err)
	}

	return chat.NewChatServiceClient(conn), conn, nil
}
