package prometheus

import (
	"strconv"

	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/metrics"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var _ metrics.HttpMetrics = (*PrometheusHttpMetrics)(nil)

type PrometheusHttpMetrics struct {
	serviceName string

	httpRequests *prometheus.CounterVec
	httpDuration *prometheus.HistogramVec
}

func NewPrometheusHttpMetrics(serviceName string) *PrometheusHttpMetrics {
	return &PrometheusHttpMetrics{
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
	}
}

func (p *PrometheusHttpMetrics) IncHTTPRequest(method, path string, statusCode int) {
	p.httpRequests.WithLabelValues(method, path, strconv.Itoa(statusCode), p.serviceName).Inc()
}

func (p *PrometheusHttpMetrics) ObserveHTTPDuration(method, path string, statusCode int, duration float64) {
	p.httpDuration.WithLabelValues(method, path, strconv.Itoa(statusCode), p.serviceName).Observe(duration)
}
