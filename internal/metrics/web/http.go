package web

import (
	"bufio"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/metrics"
	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/pkg/httpserver"
)

type HttpMetricCollector struct {
	metrics metrics.HttpMetrics
}

func NewHttpMetricCollector(metric metrics.HttpMetrics) *HttpMetricCollector {
	return &HttpMetricCollector{
		metrics: metric,
	}
}

func (h *HttpMetricCollector) Middleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			next.ServeHTTP(w, r)
			dur := time.Since(start).Seconds()
			ww := w.(*httpserver.Response)

			path := httpserver.NormalizePath(r.URL.Path)
			h.metrics.IncHTTPRequest(r.Method, path, ww.StatusCode)
			h.metrics.ObserveHTTPDuration(r.Method, path, ww.StatusCode, dur)
		})
	}
}

// CustomWriter
type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func (rw *responseWriter) Write(b []byte) (int, error) {
	return rw.ResponseWriter.Write(b)
}

func (rw *responseWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	// Проверяем, реализует ли исходный ResponseWriter интерфейс Hijacker
	if hijacker, ok := rw.ResponseWriter.(http.Hijacker); ok {
		return hijacker.Hijack()
	}
	return nil, nil, fmt.Errorf("ResponseWriter does not implement http.Hijacker")
}
