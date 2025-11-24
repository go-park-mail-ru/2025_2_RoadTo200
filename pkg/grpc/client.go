package grpc

import (
	"fmt"

	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/proto/auth"
	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/proto/chat"
	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/proto/core"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/prometheus/client_golang/prometheus"
	"google.golang.org/grpc"
)

func initClient(addr, service string, reg *prometheus.Registry) (*grpc.ClientConn, error) {
	// Метрики для gRPC клиента
	grpcMetrics := grpc_prometheus.NewClientMetrics()
	reg.MustRegister(grpcMetrics)

	// gRPC клиент с метриками
	conn, err := grpc.NewClient(
		addr,
		grpc.WithStreamInterceptor(grpc_prometheus.StreamClientInterceptor),
		grpc.WithUnaryInterceptor(grpc_prometheus.UnaryClientInterceptor),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to %s service: %w", service, err)
	}
	return conn, nil
}

func NewAuthClient(addr string, reg *prometheus.Registry) (auth.AuthServiceClient, *grpc.ClientConn, error) {
	conn, err := initClient(addr, "auth", reg)
	if err != nil {
		return nil, nil, err
	}
	return auth.NewAuthServiceClient(conn), conn, nil
}

func NewCoreClient(addr string, reg *prometheus.Registry) (core.CoreServiceClient, *grpc.ClientConn, error) {
	conn, err := initClient(addr, "auth", reg)
	if err != nil {
		return nil, nil, err
	}
	return core.NewCoreServiceClient(conn), conn, nil
}

func NewChatClient(addr string, reg *prometheus.Registry) (chat.ChatServiceClient, *grpc.ClientConn, error) {
	conn, err := initClient(addr, "auth", reg)
	if err != nil {
		return nil, nil, err
	}
	return chat.NewChatServiceClient(conn), conn, nil
}
