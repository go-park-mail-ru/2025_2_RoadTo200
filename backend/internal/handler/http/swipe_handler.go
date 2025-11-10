package handler

import (
	"encoding/json"
	"net/http"

	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/domain/dto"
	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/domain/errors"
	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/handler/middleware"
	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/service/interfaces"
	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/pkg/utils"
)

type SwipeHandler struct {
	swipeService service.SwipeService
}

func NewSwipeHandler(swipeService service.SwipeService) *SwipeHandler {
	return &SwipeHandler{
		swipeService: swipeService,
	}
}

// ProcessSwipe godoc
// @Summary Обработать свайп
// @Description Обрабатывает действие пользователя (лайк, дизлайк, суперлайк) и возвращает результат
// @Tags swipe
// @Accept json
// @Produce json
// @Security SessionToken
// @Param request body dto.SwipeRequest true "Данные свайпа"
// @Success 200 {object} dto.SwipeResponse "Результат свайпа"
// @Failure 400 {object} map[string]string "Неверный запрос"
// @Failure 401 {object} map[string]string "Не авторизован"
// @Failure 403 {object} map[string]string "Нельзя свайпнуть себя"
// @Failure 422 {object} map[string]string "Неверное действие свайпа"
// @Router /api/swipe [post]
func (h *SwipeHandler) ProcessSwipe(w http.ResponseWriter, r *http.Request) {
	userID, err := middleware.GetUserIDFromContext(r.Context())
	if err != nil {
		utils.WriteJSONError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	var req dto.SwipeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteJSONError(w, http.StatusBadRequest, "invalid JSON")
		return
	}

	response, err := h.swipeService.ProcessSwipe(userID, &req)
	if err != nil {
		status := http.StatusBadRequest
		switch err {
		case errors.ErrCannotSwipeSelf:
			status = http.StatusForbidden
		case errors.ErrInvalidSwipeAction:
			status = http.StatusUnprocessableEntity
		}
		utils.WriteJSONError(w, status, err.Error())
		return
	}

	utils.WriteJSON(w, http.StatusOK, response)
}
