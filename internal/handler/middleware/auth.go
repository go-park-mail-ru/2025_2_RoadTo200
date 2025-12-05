package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	service "github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/auth-service/service/interfaces"
	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/pkg/utils"
	"github.com/google/uuid"
)

type contextKey string

const (
	UserIDKey contextKey = "userID"
)

// Public endpoints that don't require authentication
var publicEndpoints = map[string]bool{
	"/api/register":       true,
	"/api/login":          true,
	"/api/session":        true,
	"/health":             true,
	"/swagger/":           true,
	"/swagger":            true,
	"/swagger/doc.json":   true,
	"/swagger/index.html": true,
}

func AuthMiddleware(authService service.AuthService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Check if this is a public endpoint
			if isPublicEndpoint(r.URL.Path) {
				next.ServeHTTP(w, r)
				return
			}

			var token string

			// Try to get token from header first (for Swagger)
			if authHeader := r.Header.Get("X-Session-Token"); authHeader != "" {
				token = authHeader
			} else {
				// Try to get token from cookie (for browser)
				if cookie, err := r.Cookie("session_token"); err == nil {
					token = cookie.Value
				}
			}

			if token == "" {
				utils.WriteJSONError(w, http.StatusUnauthorized, "unauthorized: no session token")
				return
			}

			// Validate session
			user, err := authService.ValidateSession(context.Background(), token)
			if err != nil {
				utils.WriteJSONError(w, http.StatusUnauthorized, "unauthorized: invalid session")
				return
			}

			// Add userID to context
			ctx := context.WithValue(r.Context(), UserIDKey, user.ID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// isPublicEndpoint checks if the requested path is a public endpoint
func isPublicEndpoint(path string) bool {
	// Exact match
	if publicEndpoints[path] {
		return true
	}

	// Prefix match for Swagger
	if strings.HasPrefix(path, "/swagger/") {
		return true
	}

	return false
}

// GetUserIDFromContext извлекает userID из контекста
func GetUserIDFromContext(ctx context.Context) (uuid.UUID, error) {
	userID, ok := ctx.Value(UserIDKey).(uuid.UUID)
	if !ok {
		return uuid.Nil, fmt.Errorf("userID not found in context")
	}
	return userID, nil
}
