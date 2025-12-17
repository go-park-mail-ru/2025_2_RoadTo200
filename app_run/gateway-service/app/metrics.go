package app

import (
	"github.com/prometheus/client_golang/prometheus"
)

func (a *App) registerMetrics() {
	a.metrics = prometheus.NewRegistry()
}
