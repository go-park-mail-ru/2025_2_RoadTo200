package grpc

import (
	"fmt"

	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/proto/auth"
	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/proto/chat"
	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/proto/core"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func initClient(addr, service string) (*grpc.ClientConn, error) {

	// gRPC клиент с метриками
	conn, err := grpc.NewClient(
		addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithStreamInterceptor(grpc_prometheus.StreamClientInterceptor),
		grpc.WithUnaryInterceptor(grpc_prometheus.UnaryClientInterceptor),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to %s service: %w", service, err)
	}
	return conn, nil
}

func NewAuthClient(addr string) (auth.AuthServiceClient, error) {
	conn, err := initClient(addr, "auth")
	if err != nil {
		return nil, err
	}
	return auth.NewAuthServiceClient(conn), nil
}

func NewCoreClient(addr string) (core.CoreServiceClient, error) {
	conn, err := initClient(addr, "core")
	if err != nil {
		return nil, err
	}
	return core.NewCoreServiceClient(conn), nil
}

func NewChatClient(addr string) (chat.ChatServiceClient, error) {
	conn, err := initClient(addr, "chat")
	if err != nil {
		return nil, err
	}
	return chat.NewChatServiceClient(conn), nil
}
