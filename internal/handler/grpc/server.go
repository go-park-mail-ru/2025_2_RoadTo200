package grpc

import (
	"fmt"
	"net/http"

	handler "github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/handler/http"
	grpcserver "github.com/go-park-mail-ru/2025_2_RoadTo200/backend/pkg/grpc"
	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/pkg/logger"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/grpc"
)

func NewGrpcServer(port int, l logger.Log) *grpc.Server {
	reg := prometheus.NewRegistry()

	// Регистрируем стандартные метрики
	grpcMetrics := grpc_prometheus.NewServerMetrics()

	// Создаем gRPC сервер с метриками
	server := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			grpcserver.ContextInjectorInterceptor(l),
			grpcMetrics.UnaryServerInterceptor(),
		),
		grpc.ChainStreamInterceptor(
			grpcserver.StreamContextInjectorInterceptor(l),
			grpcMetrics.StreamServerInterceptor(),
		),
	)

	// Инициализируем метрики
	grpcMetrics.InitializeMetrics(server)
	reg.MustRegister(grpcMetrics)

	http.Handle("/api/health", http.HandlerFunc(handler.HealthHandler))
	http.Handle("/metrics", promhttp.HandlerFor(reg, promhttp.HandlerOpts{}))
	go func() {
		err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil) // порт для метрик
		if err != nil {
			logger.Fatal("Fatal metric endpoint %s", err)
		}
	}()

	return server
}
