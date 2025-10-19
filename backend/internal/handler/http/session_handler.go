package handler

import (
	"net/http"

	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/service/interfaces"
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

// GetSession godoc
// @Summary Проверка сессии
// @Description Проверяет валидность текущей сессии пользователя
// @Tags auth
// @Accept json
// @Produce json
// @Success 200 {object} SessionResponse
// @Router /api/session [get]
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
