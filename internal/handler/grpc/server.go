package grpc

import (
	"fmt"
	"net/http"

	handler "github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/handler/http"
	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/logger"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/grpc"
)

const (
	// MaxMessageSize максимальный размер сообщения gRPC (15MB для поддержки файлов до 10MB)
	MaxMessageSize = 15 * 1024 * 1024
)

func NewGrpcServer(port int) *grpc.Server {
	reg := prometheus.NewRegistry()

	// Регистрируем стандартные метрики
	grpcMetrics := grpc_prometheus.NewServerMetrics()

	// Создаем gRPC сервер с метриками и увеличенным лимитом размера сообщений
	server := grpc.NewServer(
		grpc.StreamInterceptor(grpcMetrics.StreamServerInterceptor()),
		grpc.UnaryInterceptor(grpcMetrics.UnaryServerInterceptor()),
		grpc.MaxRecvMsgSize(MaxMessageSize),
		grpc.MaxSendMsgSize(MaxMessageSize),
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
