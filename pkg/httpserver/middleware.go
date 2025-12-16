package httpserver

import (
	"net/http"

	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/pkg/logger"
)

func InitMiddleware(l logger.LogFactory) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			r.WithContext(NewContext(r.Context(), l.Fork(NormalizePath(r.URL.Path))))
			next.ServeHTTP(&Response{rw: w}, r)
		})
	}
}
