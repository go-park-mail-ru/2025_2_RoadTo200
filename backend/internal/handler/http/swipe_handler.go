package handler

import (
	"encoding/json"
	"net/http"

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

// SwipeRequest represents swipe request
type SwipeRequest struct {
	CardID string `json:"card_id" example:"550e8400-e29b-41d4-a716-446655440000"`
	Action string `json:"action" example:"like" enums:"like,dislike,super_like"`
}

// SwipeResponse represents swipe response
type SwipeResponse struct {
	IsMatch bool   `json:"is_match" example:"true"`
	MatchID string `json:"match_id,omitempty" example:"660e8400-e29b-41d4-a716-446655440000"`
	Message string `json:"message,omitempty" example:"It's a match!"`
}

// ProcessSwipe godoc
// @Summary Обработать свайп
// @Description Обрабатывает действие пользователя (лайк, дизлайк, суперлайк) и возвращает результат
// @Tags swipe
// @Accept json
// @Produce json
// @Security SessionToken
// @Param request body SwipeRequest true "Данные свайпа"
// @Success 200 {object} SwipeResponse "Результат свайпа"
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

	var req service.SwipeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteJSONError(w, http.StatusBadRequest, "invalid JSON")
		return
	}

	// Конвертируем в доменную сущность
	swipeRequest := &service.SwipeRequest{
		CardID: req.CardID,
		Action: req.Action,
	}

	response, err := h.swipeService.ProcessSwipe(userID, swipeRequest)
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
