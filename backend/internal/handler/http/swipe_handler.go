package handler

import (
	"encoding/json"
	"net/http"

	//"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/domain/entities"
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

// ProcessSwipe handles swipe actions
// @Summary Process swipe
// @Description Process like/dislike swipe and check for matches
// @Tags swipe
// @Accept json
// @Produce json
// @Security CookieAuth
// @Param request body dto.SwipeRequest true "Swipe data"
// @Success 200 {object} entities.SwipeResponse
// @Failure 400 {object} utils.ErrorResponse
// @Failure 401 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
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
