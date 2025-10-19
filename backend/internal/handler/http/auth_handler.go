package handler

import (
	"net/http"

	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/domain/errors"
	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/service/interfaces"
	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/pkg/utils"
)

type AuthHandler struct {
	authService *service.AuthService
}

func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

// Register godoc
// @Summary Регистрация пользователя
// @Description Создает нового пользователя и устанавливает сессию
// @Tags auth
// @Accept json
// @Produce json
// @Param request body RegisterRequest true "Данные для регистрации"
// @Success 201 {object} RegisterResponse
// @Failure 400 {object} ErrorResponse
// @Failure 409 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/register [post]
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email           string `json:"email"`
		Password        string `json:"password"`
		PasswordConfirm string `json:"passwordConfirm"`
	}

	if err := utils.ReadJSON(r, &req); err != nil {
		utils.WriteJSONError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	user, session, err := h.authService.Register(req.Email, req.Password, req.PasswordConfirm)
	if err != nil {
		status := http.StatusBadRequest
		if err == errors.ErrInternalError {
			status = http.StatusInternalServerError
		}
		utils.WriteJSONError(w, status, err.Error())
		return
	}

	utils.SetSessionCookie(w, session.Token, 3600)

	utils.WriteJSON(w, http.StatusCreated, map[string]interface{}{
		"id":    user.ID,
		"email": user.Email,
	})
}

// Login godoc
// @Summary Авторизация пользователя
// @Description Выполняет вход пользователя и устанавливает сессию
// @Tags auth
// @Accept json
// @Produce json
// @Param request body LoginRequest true "Данные для входа"
// @Success 200 {object} LoginResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/login [post]
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := utils.ReadJSON(r, &req); err != nil {
		utils.WriteJSONError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	user, session, err := h.authService.Login(req.Email, req.Password)
	if err != nil {
		utils.WriteJSONError(w, http.StatusUnauthorized, "invalid email or password")
		return
	}

	utils.SetSessionCookie(w, session.Token, 3600)

	utils.WriteJSON(w, http.StatusOK, map[string]interface{}{
		"id":    user.ID,
		"email": user.Email,
	})
}

// Logout godoc
// @Summary Выход пользователя
// @Description Завершает сессию пользователя
// @Tags auth
// @Accept json
// @Produce json
// @Success 200 {object} LogoutResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/logout [post]
func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session_token")
	if err != nil {
		utils.WriteJSONError(w, http.StatusUnauthorized, "no session")
		return
	}

	if err := h.authService.Logout(cookie.Value); err != nil {
		utils.WriteJSONError(w, http.StatusInternalServerError, "internal error")
		return
	}

	utils.ClearSessionCookie(w)
	utils.WriteJSON(w, http.StatusOK, map[string]string{"message": "logged out"})
}
