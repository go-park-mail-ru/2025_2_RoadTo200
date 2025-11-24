package prometheus

//import (
//	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/metrics"
//	"github.com/prometheus/client_golang/prometheus"
//	"github.com/prometheus/client_golang/prometheus/promauto"
//)
//
//var _ metrics.BusinessMetrics = (*PrometheusBisinessMetrics)(nil)
//
//type PrometheusBisinessMetrics struct {
//	serviceName string
//
//	activeSessions prometheus.Gauge
//	databasePool   *prometheus.GaugeVec
//}
//
//func NewPrometheusBisinessMetrics(serviceName string) *PrometheusBisinessMetrics {
//	return &PrometheusBisinessMetrics{
//		serviceName: serviceName,
//
//		activeSessions: promauto.NewGauge(
//			prometheus.GaugeOpts{
//				Name: "active_sessions_total",
//				Help: "Total number of active sessions",
//			},
//		),
//		databasePool: promauto.NewGaugeVec(
//			prometheus.GaugeOpts{
//				Name: "database_connections",
//				Help: "Database connection pool statistics",
//			},
//			[]string{"service", "type"},
//		),
//	}
//}
//
//func (p *PrometheusBisinessMetrics) SetActiveSessions(count int) {
//	p.activeSessions.Set(float64(count))
//}
//
//func (p *PrometheusBisinessMetrics) SetDatabasePoolMetrics(service string, poolType string, count int) {
//	p.databasePool.WithLabelValues(p.serviceName, poolType).Set(float64(count))
//}
