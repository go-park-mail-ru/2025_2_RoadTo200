package handler

import (
	"net/http"
	"strconv"

	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/domain/errors"
	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/handler/dto"
	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/handler/middleware"
	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/service/interfaces"
	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/pkg/logger"
	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/pkg/utils"
	"github.com/google/uuid"
)

type MatchHandler struct {
	matchService service.MatchService
	logger       logger.Log
}

func NewMatchHandler(matchService service.MatchService, l logger.Log) *MatchHandler {
	return &MatchHandler{
		matchService: matchService,
		logger:       l,
	}
}

// GetUserMatches godoc
// @Summary Получить мэтчи пользователя
// @Description Возвращает список мэтчей текущего пользователя
// @Tags matches
// @Produce json
// @Security SessionToken
// @Param limit query int false "Лимит мэтчей (максимум 50)" default(20) minimum(1) maximum(50)
// @Param offset query int false "Смещение для пагинации" default(0) minimum(0)
// @Success 200 {object} dto.MatchResponse "Список мэтчей"
// @Failure 401 {object} map[string]string "Не авторизован"
// @Failure 500 {object} map[string]string "Внутренняя ошибка сервера"
// @Router /api/match [get]
func (h *MatchHandler) GetUserMatches(w http.ResponseWriter, r *http.Request) {
	h.logger.Tracef("matchHandler.GetUserMatches")

	userID, err := middleware.GetUserIDFromContext(r.Context())
	if err != nil {
		h.logger.Warnf("GetUserMatches parse context err: %v", err)
		utils.WriteJSONError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	// Парсим параметры запроса
	limit, offset := h.parseQueryParams(r)

	// Получаем мэтчи
	matches, err := h.matchService.GetUserMatches(r.Context(), userID, limit, offset)
	if err != nil {
		h.logger.Errorf("GetUserMatches err: %v", err)
		utils.WriteJSONError(w, http.StatusInternalServerError, "failed to get matches")
		return
	}

	utils.WriteJSON(w, http.StatusOK, matches)
}

// Unmatch godoc
// @Summary Удалить мэтч
// @Description Удаляет мэтч с указанным пользователем
// @Tags matches
// @Accept json
// @Produce json
// @Security SessionToken
// @Param request body dto.UnmatchRequest true "Данные для удаления мэтча"
// @Success 200 {object} dto.SuccessResponse "Мэтч успешно удален"
// @Failure 400 {object} map[string]string "Неверный запрос"
// @Failure 401 {object} map[string]string "Не авторизован"
// @Failure 404 {object} map[string]string "Мэтч не найден"
// @Failure 500 {object} map[string]string "Внутренняя ошибка сервера"
// @Router /api/match/unmatch [delete]
func (h *MatchHandler) Unmatch(w http.ResponseWriter, r *http.Request) {
	h.logger.Tracef("matchHandler.Unmatch")

	userID, err := middleware.GetUserIDFromContext(r.Context())
	if err != nil {
		h.logger.Warnf("Unmatch get context err: %v", err)
		utils.WriteJSONError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	var req dto.UnmatchRequest
	if err := utils.ReadJSON(r, &req); err != nil {
		h.logger.Warnf("Unmatch parse body err: %v", err)
		utils.WriteJSONError(w, http.StatusBadRequest, "invalid JSON")
		return
	}

	// Преобразуем строку в UUID
	targetUserID, err := uuid.Parse(req.TargetUserID)
	if err != nil {
		h.logger.Warnf("Unmatch parse uuid err: %v", err)
		utils.WriteJSONError(w, http.StatusBadRequest, "invalid user ID format")
		return
	}

	if err := h.matchService.Unmatch(r.Context(), userID, targetUserID); err != nil {
		h.logger.Errorf("Unmatch err: %v", err)
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
