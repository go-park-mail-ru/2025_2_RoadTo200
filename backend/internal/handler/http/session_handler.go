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

func (h *SessionHandler) GetSession(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session_token")
	if err != nil {
		utils.WriteJSON(w, http.StatusOK, map[string]interface{}{
			"authenticated": false,
		})
		return
	}

	// Валидируем сессию
	user, err := h.authService.ValidateSession(cookie.Value)
	if err != nil {
		utils.WriteJSON(w, http.StatusOK, map[string]interface{}{
			"authenticated": false,
		})
		return
	}

	// Формируем успешный ответ
	response := map[string]interface{}{
		"authenticated": true,
		"user": map[string]interface{}{
			"id":    user.ID,
			"email": user.Email,
			"name":  user.Name,
		},
	}

	// ВАЖНО: используем WriteJSON который устанавливает заголовки
	utils.WriteJSON(w, http.StatusOK, response)
}
