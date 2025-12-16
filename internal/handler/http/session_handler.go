package handler

import (
	//"encoding/json"
	"net/http"

	service "github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/auth-service/service/interfaces"
	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/handler/dto"
	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/pkg/logger"
	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/pkg/utils"
)

type SessionHandler struct {
	authService service.AuthService
	logger      logger.Log
}

func NewSessionHandler(authService service.AuthService, l logger.Log) *SessionHandler {
	return &SessionHandler{
		authService: authService,
		logger:      l,
	}
}

// GetSession godoc
// @Summary Проверка сессии пользователя
// @Description Проверяет валидность текущей сессии и возвращает информацию о пользователе
// @Tags auth
// @Produce json
// @Param X-Session-Token header string false "Session token"
// @Success 200 {object} dto.SessionResponse "Сессия валидна"
// @Success 200 {object} dto.SessionResponse "Сессия не валидна"
// @Router /api/session [get]
func (h *SessionHandler) GetSession(w http.ResponseWriter, r *http.Request) {
	h.logger.Tracef("sessionHandler.GetSession")
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
		h.logger.Warnf("token is empty")
		utils.WriteJSON(w, http.StatusOK, dto.SessionResponse{
			Authenticated: false,
		})
		return
	}

	// Валидируем сессию
	user, err := h.authService.ValidateSession(r.Context(), token)
	if err != nil {
		h.logger.Errorf("SessionHandler.ValidateSession: %v", err)
		utils.WriteJSON(w, http.StatusOK, dto.SessionResponse{
			Authenticated: false,
		})
		return
	}
	h.logger.Debugf("SessionHandler.ValidateSession: %v", user)
	// Формируем успешный ответ
	response := dto.SessionResponse{
		Authenticated: true,
		User: &dto.SessionUserResponse{
			ID:    user.ID.String(),
			Email: user.Email,
			Name:  user.Name,
		},
	}

	utils.WriteJSON(w, http.StatusOK, response)
}
