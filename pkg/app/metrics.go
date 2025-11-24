package app

import "github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/metrics/prometheus"

func (a *App) RegisterMetrics() *Metrics {
	return &Metrics{
		HttpMetrics: prometheus.NewPrometheusHttpMetrics(a.config.Name),
		GrpcMetrics: prometheus.NewPrometheusGrpcMetrics(a.config.Name),
	}
}
