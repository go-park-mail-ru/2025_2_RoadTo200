package prometheus

import (
	"strconv"

	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/metrics"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var _ metrics.GrpcMetrics = (*PrometheusGrpcMetrics)(nil)

type PrometheusGrpcMetrics struct {
	serviceName string

	grpcRequests *prometheus.CounterVec
	grpcDuration *prometheus.HistogramVec
	grpcCount    *prometheus.CounterVec
}

func NewPrometheusGrpcMetrics(serviceName string) *PrometheusGrpcMetrics {
	return &PrometheusGrpcMetrics{
		serviceName: serviceName,

		grpcRequests: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "grpc_requests_total",
				Help: "Total number of gRPC requests",
			},
			[]string{"method", "status", "service"},
		),
		grpcDuration: promauto.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "grpc_request_duration_seconds",
				Help:    "Duration of gRPC requests",
				Buckets: []float64{0.001, 0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1, 2.5, 5, 10},
			},
			[]string{"method", "status", "service"},
		),
		grpcCount: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "errors_total",
				Help: "Total number of errors",
			},
			[]string{"service", "operation", "type"},
		),
	}
}

func (p *PrometheusGrpcMetrics) IncGRPCRequest(method string, statusCode int) {
	p.grpcRequests.WithLabelValues(method, strconv.Itoa(statusCode), p.serviceName).Inc()
}

func (p *PrometheusGrpcMetrics) ObserveGRPCDuration(method string, statusCode int, duration float64) {
	p.grpcDuration.WithLabelValues(method, strconv.Itoa(statusCode), p.serviceName).Observe(duration)
}

func (p *PrometheusGrpcMetrics) IncError(operation, errorType string) {
	p.grpcCount.WithLabelValues(p.serviceName, operation, errorType).Inc()
}
