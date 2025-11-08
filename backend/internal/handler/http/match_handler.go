package handler

import (
	"net/http"
	"strconv"

	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/domain/dto"
	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/domain/errors"
	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/handler/middleware"
	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/service/interfaces"
	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/pkg/utils"
	"github.com/google/uuid"
)

type MatchHandler struct {
	matchService service.MatchService
}

func NewMatchHandler(matchService service.MatchService) *MatchHandler {
	return &MatchHandler{
		matchService: matchService,
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
// @Success 200 {object} interface{} "Список мэтчей"
// @Failure 401 {object} map[string]string "Не авторизован"
// @Failure 500 {object} map[string]string "Внутренняя ошибка сервера"
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

// Unmatch godoc
// @Summary Удалить мэтч
// @Description Удаляет мэтч с указанным пользователем
// @Tags matches
// @Accept json
// @Produce json
// @Security SessionToken
// @Param request body UnmatchRequest true "Данные для удаления мэтча"
// @Success 200 {object} SuccessResponse "Мэтч успешно удален"
// @Failure 400 {object} map[string]string "Неверный запрос"
// @Failure 401 {object} map[string]string "Не авторизован"
// @Failure 404 {object} map[string]string "Мэтч не найден"
// @Failure 500 {object} map[string]string "Внутренняя ошибка сервера"
// @Router /api/matches/unmatch [post]
func (h *MatchHandler) Unmatch(w http.ResponseWriter, r *http.Request) {
	userID, err := middleware.GetUserIDFromContext(r.Context())
	if err != nil {
		utils.WriteJSONError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	var req dto.UnmatchRequest
	if err := utils.ReadJSON(r, &req); err != nil {
		utils.WriteJSONError(w, http.StatusBadRequest, "invalid JSON")
		return
	}

	// Преобразуем строку в UUID
	targetUserID, err := uuid.Parse(req.TargetUserID)
	if err != nil {
		utils.WriteJSONError(w, http.StatusBadRequest, "invalid user ID format")
		return
	}

	if err := h.matchService.Unmatch(userID, targetUserID); err != nil {
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
