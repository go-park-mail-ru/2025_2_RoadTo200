package grpc

import (
	"fmt"
	"net/http"

	handler "github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/handler/http"
	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/logger"
	pb "github.com/go-park-mail-ru/2025_2_RoadTo200/backend/proto/core"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/grpc"
)

func NewGrpcServer(coreServiceServer pb.CoreServiceServer, port int) *grpc.Server {
	reg := prometheus.NewRegistry()

	// Регистрируем стандартные метрики
	grpcMetrics := grpc_prometheus.NewServerMetrics()
	reg.MustRegister(grpcMetrics)

	// Создаем gRPC сервер с метриками
	server := grpc.NewServer(
		grpc.StreamInterceptor(grpc_prometheus.StreamServerInterceptor),
		grpc.UnaryInterceptor(grpc_prometheus.UnaryServerInterceptor),
	)

	grpcServer := grpc.NewServer()

	pb.RegisterCoreServiceServer(grpcServer, coreServiceServer)
	// Инициализируем метрики
	grpcMetrics.InitializeMetrics(server)

	http.Handle("/health", http.HandlerFunc(handler.HealthHandler))
	http.Handle("/metrics", promhttp.HandlerFor(reg, promhttp.HandlerOpts{}))
	go func() {
		err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil) // порт для метрик
		if err != nil {
			logger.Fatal("Fatal metric endpoint %s", err)
		}
	}()

	return grpcServer
}
