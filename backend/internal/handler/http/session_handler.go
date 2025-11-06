package handler

import (
	//"encoding/json"
	"net/http"

	//"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/domain/entities"
	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/logger"
	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/service/implementations"
	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/pkg/utils"
)

type SessionHandler struct {
	authService *service.AuthService
	logger      *logger.Logger
}

func NewSessionHandler(authService *service.AuthService) *SessionHandler {
	return &SessionHandler{
		authService: authService,
	}
}

// SessionResponse represents session check response
type SessionResponse struct {
	Authenticated bool                 `json:"authenticated" example:"true"`
	User          *SessionUserResponse `json:"user,omitempty"`
}

// SessionUserResponse represents user data in session response
type SessionUserResponse struct {
	ID    string `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
	Email string `json:"email" example:"user@example.com"`
	Name  string `json:"name" example:"Алексей"`
}

// GetSession godoc
// @Summary Проверка сессии пользователя
// @Description Проверяет валидность текущей сессии и возвращает информацию о пользователе
// @Tags auth
// @Produce json
// @Param X-Session-Token header string false "Session token"
// @Success 200 {object} SessionResponse "Сессия валидна"
// @Success 200 {object} SessionResponse "Сессия не валидна"
// @Router /api/session [get]
func (h *SessionHandler) GetSession(w http.ResponseWriter, r *http.Request) {
	var token string

	// Пробуем получить токен из заголовка
	if authHeader := r.Header.Get("X-Session-Token"); authHeader != "" {
		token = authHeader
	} else {
		// Пробуем получить токен из куки
		if cookie, err := r.Cookie("session_token"); err == nil {
			token = cookie.Value
		}
	}

	if token == "" {
		utils.WriteJSON(w, http.StatusOK, SessionResponse{
			Authenticated: false,
		})
		return
	}

	// Валидируем сессию
	user, err := h.authService.ValidateSession(token)
	if err != nil {
		utils.WriteJSON(w, http.StatusOK, SessionResponse{
			Authenticated: false,
		})
		return
	}

	// Формируем успешный ответ
	response := SessionResponse{
		Authenticated: true,
		User: &SessionUserResponse{
			ID:    user.ID.String(),
			Email: user.Email,
			Name:  user.Name,
		},
	}

	utils.WriteJSON(w, http.StatusOK, response)
}
