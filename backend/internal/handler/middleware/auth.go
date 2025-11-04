package middleware

import (
	"context"
	"fmt"
	"net/http"

	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/service/implementations"
	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/pkg/utils"
	"github.com/google/uuid"
)

type contextKey string

const (
	UserIDKey contextKey = "userID"
)

func AuthMiddleware(authService *service.AuthService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Получаем session_token из cookie
			cookie, err := r.Cookie("session_token")
			if err != nil {
				utils.WriteJSONError(w, http.StatusUnauthorized, "unauthorized: no session token")
				return
			}

			// Валидируем сессию
			user, err := authService.ValidateSession(cookie.Value)
			if err != nil {
				utils.WriteJSONError(w, http.StatusUnauthorized, "unauthorized: invalid session")
				return
			}

			// Добавляем userID в контекст
			ctx := context.WithValue(r.Context(), UserIDKey, user.ID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// GetUserIDFromContext извлекает userID из контекста
func GetUserIDFromContext(ctx context.Context) (uuid.UUID, error) {
	userID, ok := ctx.Value(UserIDKey).(uuid.UUID)
	if !ok {
		return uuid.Nil, fmt.Errorf("userID not found in context")
	}
	return userID, nil
}
