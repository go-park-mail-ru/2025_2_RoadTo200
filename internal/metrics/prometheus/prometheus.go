package prometheus

import (
	"strconv"

	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/metrics"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var _ metrics.Metrics = (*PrometheusMetrics)(nil)

type PrometheusMetrics struct {
	serviceName string

	httpRequests   *prometheus.CounterVec
	httpDuration   *prometheus.HistogramVec
	grpcRequests   *prometheus.CounterVec
	grpcDuration   *prometheus.HistogramVec
	grpcCount      *prometheus.CounterVec
	activeSessions prometheus.Gauge
	databasePool   *prometheus.GaugeVec
	errorCount     *prometheus.CounterVec
	cacheHits      *prometheus.CounterVec
	cacheMisses    *prometheus.CounterVec
}

func NewPrometheusMetrics(serviceName string) *PrometheusMetrics {
	return &PrometheusMetrics{
		serviceName: serviceName,

		httpRequests: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "http_requests_total",
				Help: "Total number of HTTP requests",
			},
			[]string{"method", "path", "status", "service"},
		),
		httpDuration: promauto.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "http_request_duration_seconds",
				Help:    "Duration of HTTP requests",
				Buckets: []float64{0.001, 0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1, 2.5, 5, 10},
			},
			[]string{"method", "path", "status", "service"},
		),
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
		activeSessions: promauto.NewGauge(
			prometheus.GaugeOpts{
				Name: "active_sessions_total",
				Help: "Total number of active sessions",
			},
		),
		databasePool: promauto.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "database_connections",
				Help: "Database connection pool statistics",
			},
			[]string{"service", "type"},
		),
		errorCount: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "errors_total",
				Help: "Total number of errors",
			},
			[]string{"service", "operation", "type"},
		),
		cacheHits: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "cache_hits_total",
				Help: "Total number of cache hits",
			},
			[]string{"service", "cache_type"},
		),
		cacheMisses: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "cache_misses_total",
				Help: "Total number of cache misses",
			},
			[]string{"service", "cache_type"},
		),
	}
}

func (p *PrometheusMetrics) IncHTTPRequest(method, path string, statusCode int) {
	p.httpRequests.WithLabelValues(method, path, strconv.Itoa(statusCode), p.serviceName).Inc()
}

func (p *PrometheusMetrics) ObserveHTTPDuration(method, path string, statusCode int, duration float64) {
	p.httpDuration.WithLabelValues(method, path, strconv.Itoa(statusCode), p.serviceName).Observe(duration)
}

func (p *PrometheusMetrics) IncGRPCRequest(method string, statusCode int) {
	p.grpcRequests.WithLabelValues(method, strconv.Itoa(statusCode), p.serviceName).Inc()
}

func (p *PrometheusMetrics) ObserveGRPCDuration(method string, statusCode int, duration float64) {
	p.grpcDuration.WithLabelValues(method, strconv.Itoa(statusCode), p.serviceName).Observe(duration)
}

func (p *PrometheusMetrics) IncError(operation, errorType string) {
	p.grpcCount.WithLabelValues(p.serviceName, operation, errorType).Inc()
}

func (p *PrometheusMetrics) SetActiveSessions(count int) {
	p.activeSessions.Set(float64(count))
}

func (p *PrometheusMetrics) SetDatabasePoolMetrics(service string, poolType string, count int) {
	p.databasePool.WithLabelValues(p.serviceName, poolType).Set(float64(count))
}

func (p *PrometheusMetrics) SetErrorMetrics(model, method, errorType string) {
	p.errorCount.WithLabelValues(p.serviceName, model, method, errorType).Inc()
}

func (p *PrometheusMetrics) IncCacheHit(cacheType string) {
	p.cacheHits.WithLabelValues(p.serviceName, cacheType).Inc()
}

func (p *PrometheusMetrics) IncCacheMiss(cacheType string) {
	p.cacheMisses.WithLabelValues(p.serviceName, cacheType).Inc()
}
