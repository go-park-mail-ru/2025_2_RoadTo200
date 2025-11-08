package handler

import (
	"net/http"

	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/domain/dto"
	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/domain/errors"
	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/logger"
	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/service/implementations"
	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/pkg/utils"
)

type AuthHandler struct {
	authService *service.AuthService
	logger      logger.Log
}

func NewAuthHandler(authService *service.AuthService, l logger.Log) *AuthHandler {
	return &AuthHandler{
		authService: authService,
		logger:      l,
	}
}

// Register godoc
// @Summary Регистрация нового пользователя
// @Description Создает нового пользователя и устанавливает сессию
// @Tags auth
// @Accept json
// @Produce json
// @Param request body RegisterRequest true "Данные для регистрации"
// @Success 201 {object} RegisterResponse "Пользователь создан"
// @Failure 400 {object} map[string]string "Неверный формат запроса"
// @Failure 400 {object} map[string]string "Ошибка валидации данных"
// @Failure 500 {object} map[string]string "Внутренняя ошибка сервера"
// @Router /api/register [post]
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req dto.RegisterRequest

	if err := utils.ReadJSON(r, &req); err != nil {
		h.logger.Warnf("handler.Register: %v", err)
		utils.WriteJSONError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	user, session, err := h.authService.Register(req.Email, req.Password, req.PasswordConfirm)
	if err != nil {
		h.logger.Warnf("handler.Register: %v", err)
		status := http.StatusBadRequest
		if err == errors.ErrInternalError {
			status = http.StatusInternalServerError
		}
		utils.WriteJSONError(w, status, err.Error())
		return
	}

	utils.SetSessionCookie(w, session.Token, 3600)

	utils.WriteJSON(w, http.StatusCreated, dto.RegisterResponse{
		ID:    user.ID.String(),
		Email: user.Email,
	})
}

// Login godoc
// @Summary Аутентификация пользователя
// @Description Вход пользователя в систему и установка сессии
// @Tags auth
// @Accept json
// @Produce json
// @Param request body LoginRequest true "Данные для входа"
// @Success 200 {object} LoginResponse "Успешный вход"
// @Failure 400 {object} map[string]string "Неверный формат запроса"
// @Failure 401 {object} map[string]string "Неверный email или пароль"
// @Router /api/login [post]
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req dto.LoginRequest

	if err := utils.ReadJSON(r, &req); err != nil {
		h.logger.Warnf("handler.Login: %v", err)
		utils.WriteJSONError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	user, session, err := h.authService.Login(req.Email, req.Password)
	if err != nil {
		h.logger.Warnf("handler.Login: %v", err)
		utils.WriteJSONError(w, http.StatusUnauthorized, "invalid email or password")
		return
	}

	utils.SetSessionCookie(w, session.Token, 3600)

	utils.WriteJSON(w, http.StatusOK, dto.LoginResponse{
		ID:    user.ID.String(),
		Email: user.Email,
	})
}

// Logout godoc
// @Summary Выход пользователя
// @Description Завершает сессию пользователя
// @Tags auth
// @Produce json
// @Security SessionToken
// @Success 200 {object} map[string]string "Успешный выход"
// @Failure 401 {object} map[string]string "Сессия не найдена"
// @Failure 500 {object} map[string]string "Внутренняя ошибка сервера"
// @Router /api/logout [post]
func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session_token")
	if err != nil {
		h.logger.Warnf("handler.Logout: %v", err)
		utils.WriteJSONError(w, http.StatusUnauthorized, "no session")
		return
	}

	if err := h.authService.Logout(cookie.Value); err != nil {
		h.logger.Warnf("handler.Logout: %v", err)
		utils.WriteJSONError(w, http.StatusInternalServerError, "internal error")
		return
	}

	utils.ClearSessionCookie(w)
	utils.WriteJSON(w, http.StatusOK, map[string]string{"message": "logged out"})
}
