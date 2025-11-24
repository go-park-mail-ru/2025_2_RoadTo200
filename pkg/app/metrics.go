package app

import (
	metrics "github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/metrics/prometheus"
	"github.com/prometheus/client_golang/prometheus"
)

func (a *App) registerMetrics() {
	a.metrics = &Metrics{
		HttpMetrics: metrics.NewPrometheusHttpMetrics(a.config.Name),
		GrpcMetrics: prometheus.NewRegistry(),
	}
}
