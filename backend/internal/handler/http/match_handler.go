package handler

import (
	"net/http"
	"strconv"

	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/domain/errors"
	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/handler/middleware"
	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/service/interfaces"
	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/pkg/utils"
)

type MatchHandler struct {
	matchService service.MatchService
}

func NewMatchHandler(matchService service.MatchService) *MatchHandler {
	return &MatchHandler{
		matchService: matchService,
	}
}

// GetUserMatches returns user's matches
// @Summary Get user matches
// @Description Get list of user's matches with user details
// @Tags matches
// @Accept json
// @Produce json
// @Security CookieAuth
// @Param limit query int false "Number of matches to return" default(20)
// @Param offset query int false "Offset for pagination" default(0)
// @Success 200 {object} service.MatchesResponse
// @Failure 400 {object} service.ErrorResponse
// @Failure 401 {object} service.ErrorResponse
// @Failure 500 {object} service.ErrorResponse
// @Router /api/matches [get]
func (h *MatchHandler) GetUserMatches(w http.ResponseWriter, r *http.Request) {
	userID, err := middleware.GetUserIDFromContext(r.Context())
	if err != nil {
		utils.WriteJSONError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	// Парсим параметры запроса
	limit, offset := h.parseQueryParams(r)

	// Получаем мэтчи
	matches, err := h.matchService.GetUserMatches(userID, limit, offset)
	if err != nil {
		utils.WriteJSONError(w, http.StatusInternalServerError, "failed to get matches")
		return
	}

	utils.WriteJSON(w, http.StatusOK, matches)
}

// Unmatch removes a match
// @Summary Unmatch user
// @Description Remove match with another user
// @Tags matches
// @Accept json
// @Produce json
// @Security CookieAuth
// @Param request body service.UnmatchRequest true "Unmatch data"
// @Success 200 {object} utils.SuccessResponse
// @Failure 400 {object} service.ErrorResponse
// @Failure 401 {object} service.ErrorResponse
// @Failure 404 {object} service.ErrorResponse
// @Failure 500 {object} service.ErrorResponse
// @Router /api/matches/unmatch [post]
func (h *MatchHandler) Unmatch(w http.ResponseWriter, r *http.Request) {
	userID, err := middleware.GetUserIDFromContext(r.Context())
	if err != nil {
		utils.WriteJSONError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	var req service.UnmatchRequest
	if err := utils.ReadJSON(r, &req); err != nil {
		utils.WriteJSONError(w, http.StatusBadRequest, "invalid JSON")
		return
	}

	if err := h.matchService.Unmatch(userID, req.TargetUserID); err != nil {
		status := http.StatusInternalServerError
		if err == errors.ErrMatchNotFound {
			status = http.StatusNotFound
		}
		utils.WriteJSONError(w, status, err.Error())
		return
	}

	utils.WriteJSON(w, http.StatusOK, utils.SuccessResponse{
		Message: "Successfully unmatched",
	})
}

// parseQueryParams парсит параметры запроса
func (h *MatchHandler) parseQueryParams(r *http.Request) (int, int) {
	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")

	limit := 20 // дефолтное значение
	offset := 0

	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 50 {
			limit = l
		}
	}

	if offsetStr != "" {
		if o, err := strconv.Atoi(offsetStr); err == nil && o >= 0 {
			offset = o
		}
	}

	return limit, offset
}
