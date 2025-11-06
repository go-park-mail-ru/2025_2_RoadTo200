package middleware

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/config"
)

func CORSMiddleware(cfg *config.CORSConfig) func(http.Handler) http.Handler {
	for _, el := range cfg.AllowedOrigins {
		if el == "*" {
			cfg.AllowedOrigins = []string{"*"}
			break
		}
	}
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			origin := r.Header.Get("Origin")
			if len(cfg.AllowedOrigins) > 0 && cfg.AllowedOrigins[0] == "*" {
				origin = "*"
			}
			for _, or := range cfg.AllowedOrigins {
				if or == origin {
					w.Header().Set("Access-Control-Allow-Origin", origin)
					break
				}
			}

			if cfg.AllowCredentials {
				w.Header().Set("Access-Control-Allow-Credentials", "true")
			}
			if len(cfg.AllowedMethods) > 0 {
				w.Header().Set("Access-Control-Allow-Methods", strings.Join(cfg.AllowedMethods, ", "))
			}

			if len(cfg.AllowedHeaders) > 0 {
				w.Header().Set("Access-Control-Allow-Headers", strings.Join(cfg.AllowedHeaders, ", "))
			}

			if len(cfg.ExposeHeaders) > 0 {
				w.Header().Set("Access-Control-Expose-Headers", strings.Join(cfg.ExposeHeaders, ", "))
			}
			if cfg.MaxAgeSeconds > 0 {
				w.Header().Set("Access-Control-Max-Age", strconv.Itoa(cfg.MaxAgeSeconds))
			}

			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusOK)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
