package handler

import (
	"net/http"

	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/domain/errors"
	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/handler/dto"
	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/service/interfaces"
	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/pkg/logger"
	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/pkg/utils"
)

type AuthHandler struct {
	authService service.AuthService
	logger      logger.Log
}

func NewAuthHandler(authService service.AuthService, l logger.Log) *AuthHandler {
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
// @Param request body dto.RegisterRequest true "Данные для регистрации"
// @Success 201 {object} dto.RegisterResponse "Пользователь создан"
// @Failure 400 {object} map[string]string "Неверный формат запроса"
// @Failure 400 {object} map[string]string "Ошибка валидации данных"
// @Failure 500 {object} map[string]string "Внутренняя ошибка сервера"
// @Router /api/register [post]
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	h.logger.Trace("authHandler.Register")
	var req dto.RegisterRequest

	if err := utils.ReadJSON(r, &req); err != nil {
		h.logger.Warnf("handler.Register: %v", err)
		utils.WriteJSONError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	user, session, err := h.authService.Register(r.Context(), req.Email, req.Password, req.PasswordConfirm)
	if err != nil {
		h.logger.Errorf("handler.Register: %v", err)
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
// @Param request body dto.LoginRequest true "Данные для входа"
// @Success 200 {object} dto.LoginResponse "Успешный вход"
// @Failure 400 {object} map[string]string "Неверный формат запроса"
// @Failure 401 {object} map[string]string "Неверный email или пароль"
// @Router /api/login [post]
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	h.logger.Trace("authHandler.Login")
	var req dto.LoginRequest

	if err := utils.ReadJSON(r, &req); err != nil {
		h.logger.Warnf("handler.Login Error body: %v", err)
		utils.WriteJSONError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	user, session, err := h.authService.Login(r.Context(), req.Email, req.Password)
	if err != nil {
		h.logger.Errorf("handler.Login: %v", err)
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
// @Security dto.SessionToken
// @Success 200 {object} map[string]string "Успешный выход"
// @Failure 401 {object} map[string]string "Сессия не найдена"
// @Failure 500 {object} map[string]string "Внутренняя ошибка сервера"
// @Router /api/logout [post]
func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	h.logger.Trace("authHandler.Logout")

	cookie, err := r.Cookie("session_token")
	if err != nil {
		h.logger.Warnf("handler.Logout no session: %v", err)
		utils.WriteJSONError(w, http.StatusUnauthorized, "no session")
		return
	}

	if err := h.authService.Logout(r.Context(), cookie.Value); err != nil {
		h.logger.Errorf("handler.Logout: %v", err)
		utils.WriteJSONError(w, http.StatusInternalServerError, "internal error")
		return
	}

	utils.ClearSessionCookie(w)
	utils.WriteJSON(w, http.StatusOK, map[string]string{"message": "logged out"})
}
