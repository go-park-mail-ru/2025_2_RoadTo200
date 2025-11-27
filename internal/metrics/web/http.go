package web

import (
	"bufio"
	"fmt"
	"net"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/metrics"
)

var _ http.Hijacker = (*responseWriter)(nil)

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
			ww := responseWriter{w, 0}
			next.ServeHTTP(&ww, r)
			dur := time.Since(start).Seconds()

			path := h.normalizePath(r.URL.Path)
			h.metrics.IncHTTPRequest(r.Method, path, ww.statusCode)
			h.metrics.ObserveHTTPDuration(r.Method, path, ww.statusCode, float64(dur))
		})
	}
}

// normalizePath нормализует путь URL, заменяя числовые ID и UUID на плейсхолдеры
func (c *HttpMetricCollector) normalizePath(path string) string {
	path, _ = strings.CutSuffix(path, "?")

	// Заменяем UUID на :id
	uuidRegex := regexp.MustCompile(`[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{12}`)
	path = uuidRegex.ReplaceAllString(path, ":id")

	// Заменяем числовые ID на :id
	numericRegex := regexp.MustCompile(`/\d+`)
	path = numericRegex.ReplaceAllString(path, "/:id")

	// Заменяем email-like строки
	emailRegex := regexp.MustCompile(`/[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}`)
	path = emailRegex.ReplaceAllString(path, "/:email")

	// Заменяем хеши и токены
	hashRegex := regexp.MustCompile(`/[a-fA-F0-9]{32,}`)
	path = hashRegex.ReplaceAllString(path, "/:hash")

	path, _ = strings.CutSuffix(path, "/")

	return path
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
