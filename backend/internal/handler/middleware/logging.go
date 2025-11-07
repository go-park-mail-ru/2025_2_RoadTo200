package middleware

import (
	"net/http"

	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/logger"
)

func LogMiddleware(l *logger.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			l.Infof("%s %s", r.Method, r.RequestURI)
			l.Debugf("Headers: %v", r.Header)
			l.Debugf("Body: %v", r.Body)

			next.ServeHTTP(w, r)
		})
	}
}
