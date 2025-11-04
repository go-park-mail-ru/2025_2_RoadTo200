package handler

import (
	"net/http"

	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/service/implementations"
	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/pkg/utils"
)

type SessionHandler struct {
	authService *service.AuthService
}

func NewSessionHandler(authService *service.AuthService) *SessionHandler {
	return &SessionHandler{
		authService: authService,
	}
}

// GetSession проверяет валидность сессии пользователя
// @Summary Проверка сессии
// @Description Проверяет валидность текущей сессии и возвращает статус аутентификации
// @Tags auth
// @Produce json
// @Success 200 {object} object "Статус сессии" example:{"authenticated":true,"user":{"id":"123","email":"user@example.com"}}
// @Success 200 {object} object "Сессия не валидна" example:{"authenticated":false}
// @Router /session [get]
func (h *SessionHandler) GetSession(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session_token")
	if err != nil {
		utils.WriteJSON(w, http.StatusOK, map[string]interface{}{
			"authenticated": false,
		})
		return
	}

	user, err := h.authService.ValidateSession(cookie.Value)
	if err != nil {
		utils.WriteJSON(w, http.StatusOK, map[string]interface{}{
			"authenticated": false,
		})
		return
	}

	utils.WriteJSON(w, http.StatusOK, map[string]interface{}{
		"authenticated": true,
		"user": map[string]interface{}{
			"id":    user.ID,
			"email": user.Email,
		},
	})
}
