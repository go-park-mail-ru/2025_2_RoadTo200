package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/domain/constants"
	expectation "github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/domain/errors"
	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/handler/dto"
	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/handler/middleware"
	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/logger"
	service "github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/service/interfaces"
	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/pkg/utils"
	"github.com/google/uuid"
)

type StrikeHandler struct {
	strikeService service.StrikeService
	logger        logger.Log
}

func NewStrikeHandler(strikeService service.StrikeService, l logger.Log) *StrikeHandler {
	return &StrikeHandler{
		strikeService: strikeService,
		logger:        l,
	}
}

// getContext функция извлечения контекста и проверки формата тела
func (h *StrikeHandler) getContext(w http.ResponseWriter, r *http.Request) (uuid.UUID, error) {
	userID, err := middleware.GetUserIDFromContext(r.Context())
	if err != nil {
		h.logger.Warnf("GetUserFromContext: %v", err)
		utils.WriteJSONError(w, http.StatusUnauthorized, "unauthorized")
		return uuid.Nil, errors.New("unauthorized")
	}

	return userID, nil
}

// CreateStrike godoc
// @Summary Создать жалобу
// @Description Создает новую жалобу на пользователя
// @Tags strikes
// @Accept json
// @Produce json
// @Security SessionToken
// @Param strike body dto.StrikeCreateRequest true "Данные для создания жалобы"
// @Success 201 {object} dto.StrikeResponse
// @Failure 400 {object} map[string]string "Неверный запрос"
// @Failure 401 {object} map[string]string "Не авторизован"
// @Failure 404 {object} map[string]string "Пользователь не найден"
// @Failure 409 {object} map[string]string "Дублирующая жалоба"
// @Failure 500 {object} map[string]string "Внутренняя ошибка сервера"
// @Router /api/strike [post]
func (h *StrikeHandler) CreateStrike(w http.ResponseWriter, r *http.Request) {
	h.logger.Trace("strikeHandler.CreateStrike")

	userID, err := h.getContext(w, r)
	if err != nil {
		h.logger.Warnf("getContext: %v", err)
		return
	}

	var strikeReq dto.StrikeCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&strikeReq); err != nil {
		h.logger.Warnf("CreateStrike: invalid JSON: %v", err)
		utils.WriteJSONError(w, http.StatusBadRequest, "invalid JSON")
		return
	}

	// Устанавливаем reporter_id из контекста
	strikeReq.ReporterID = userID

	strike, err := h.strikeService.CreateStrike(r.Context(), &strikeReq)
	if err != nil {
		h.logger.Errorf("CreateStrike: %v", err)
		switch {
		case errors.Is(err, expectation.ErrUserNotFound),
			errors.Is(err, expectation.ErrInvalidReporterID),
			errors.Is(err, expectation.ErrInvalidTargetUserID),
			errors.Is(err, expectation.ErrSelfStrikeNotAllowed):
			utils.WriteJSONError(w, http.StatusBadRequest, err.Error())
		case errors.Is(err, expectation.ErrDuplicateStrike):
			utils.WriteJSONError(w, http.StatusConflict, err.Error())
		default:
			utils.WriteJSONError(w, http.StatusInternalServerError, "internal server error")
		}
		return
	}

	strikeResp := dto.ToStrikeResponse(strike)
	utils.WriteJSON(w, http.StatusCreated, strikeResp)
}

// GetStrike godoc
// @Summary Получить жалобу
// @Description Возвращает жалобу по её ID
// @Tags strikes
// @Produce json
// @Security SessionToken
// @Param id path string true "ID жалобы"
// @Success 200 {object} dto.StrikeResponse
// @Failure 400 {object} map[string]string "Неверный запрос"
// @Failure 401 {object} map[string]string "Не авторизован"
// @Failure 404 {object} map[string]string "Жалоба не найдена"
// @Failure 500 {object} map[string]string "Внутренняя ошибка сервера"
// @Router /api/strike/{id} [get]
func (h *StrikeHandler) GetStrike(w http.ResponseWriter, r *http.Request) {
	h.logger.Trace("strikeHandler.GetStrike")

	_, err := h.getContext(w, r)
	if err != nil {
		h.logger.Warnf("getContext: %v", err)
		return
	}

	strikeID := r.PathValue("id")
	if strikeID == "" {
		utils.WriteJSONError(w, http.StatusBadRequest, "strike ID is required")
		return
	}

	strike, err := h.strikeService.GetStrikeByID(r.Context(), strikeID)
	if err != nil {
		h.logger.Errorf("GetStrike: %v", err)
		if errors.Is(err, expectation.ErrStrikeNotFound) {
			utils.WriteJSONError(w, http.StatusNotFound, "strike not found")
			return
		}
		utils.WriteJSONError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	strikeResp := dto.ToStrikeResponse(strike)
	utils.WriteJSON(w, http.StatusOK, strikeResp)
}

// GetStrikesByUserID godoc
// @Summary Получить жалобы на пользователя
// @Description Возвращает список жалоб на конкретного пользователя
// @Tags strikes
// @Produce json
// @Security SessionToken
// @Param user_id path string true "ID пользователя"
// @Param limit query int false "Лимит" default(20)
// @Param offset query int false "Смещение" default(0)
// @Success 200 {object} dto.StrikesListResponse
// @Failure 400 {object} map[string]string "Неверный запрос"
// @Failure 401 {object} map[string]string "Не авторизован"
// @Failure 404 {object} map[string]string "Пользователь не найден"
// @Failure 500 {object} map[string]string "Внутренняя ошибка сервера"
// @Router /api/strike/user/{user_id} [get]
func (h *StrikeHandler) GetStrikesByUserID(w http.ResponseWriter, r *http.Request) {
	h.logger.Trace("strikeHandler.GetStrikesByUserID")

	_, err := h.getContext(w, r)
	if err != nil {
		h.logger.Warnf("getContext: %v", err)
		return
	}

	userID := r.PathValue("user_id")
	if userID == "" {
		utils.WriteJSONError(w, http.StatusBadRequest, "user ID is required")
		return
	}

	limit, offset, err := h.parsePaginationParams(r)
	if err != nil {
		utils.WriteJSONError(w, http.StatusBadRequest, "invalid pagination parameters")
		return
	}

	strikes, err := h.strikeService.GetStrikesByUserID(r.Context(), userID, limit, offset)
	if err != nil {
		h.logger.Errorf("GetStrikesByUserID: %v", err)
		if errors.Is(err, expectation.ErrUserNotFound) {
			utils.WriteJSONError(w, http.StatusNotFound, "user not found")
			return
		}
		utils.WriteJSONError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	strikesResp := dto.ToStrikesListResponse(strikes)
	utils.WriteJSON(w, http.StatusOK, strikesResp)
}

// GetStrikesByType godoc
// @Summary Получить жалобы по типу
// @Description Возвращает список жалоб по типу нарушения
// @Tags strikes
// @Produce json
// @Security SessionToken
// @Param type path string true "Тип нарушения"
// @Param limit query int false "Лимит" default(20)
// @Param offset query int false "Смещение" default(0)
// @Success 200 {object} dto.StrikesListResponse
// @Failure 400 {object} map[string]string "Неверный запрос"
// @Failure 401 {object} map[string]string "Не авторизован"
// @Failure 500 {object} map[string]string "Внутренняя ошибка сервера"
// @Router /api/strike/type/{type} [get]
func (h *StrikeHandler) GetStrikesByType(w http.ResponseWriter, r *http.Request) {
	h.logger.Trace("strikeHandler.GetStrikesByType")

	_, err := h.getContext(w, r)
	if err != nil {
		h.logger.Warnf("getContext: %v", err)
		return
	}

	strikeType := constants.StrikeType(r.PathValue("type"))
	if strikeType == "" {
		utils.WriteJSONError(w, http.StatusBadRequest, "strike type is required")
		return
	}

	limit, offset, err := h.parsePaginationParams(r)
	if err != nil {
		utils.WriteJSONError(w, http.StatusBadRequest, "invalid pagination parameters")
		return
	}

	strikes, err := h.strikeService.GetStrikesByType(r.Context(), strikeType, limit, offset)
	if err != nil {
		h.logger.Errorf("GetStrikesByType: %v", err)
		utils.WriteJSONError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	strikesResp := dto.ToStrikesListResponse(strikes)
	utils.WriteJSON(w, http.StatusOK, strikesResp)
}

// GetStrikesByDateRange godoc
// @Summary Получить жалобы за период
// @Description Возвращает список жалоб за указанный период. По умолчанию: from=начало времен, to=текущее время + 1 день, limit=20, offset=0
// @Tags strikes
// @Produce json
// @Security SessionToken
// @Param from query string false "Начальная дата (RFC3339). По умолчанию: начало времен"
// @Param to query string false "Конечная дата (RFC3339). По умолчанию: текущее время + 1 день"
// @Param limit query int false "Лимит" default(20)
// @Param offset query int false "Смещение" default(0)
// @Success 200 {object} dto.StrikesListResponse
// @Failure 400 {object} map[string]string "Неверный запрос"
// @Failure 401 {object} map[string]string "Не авторизован"
// @Failure 500 {object} map[string]string "Внутренняя ошибка сервера"
// @Router /api/strike/range [get]
func (h *StrikeHandler) GetStrikesByDateRange(w http.ResponseWriter, r *http.Request) {
	h.logger.Trace("strikeHandler.GetStrikesByDateRange")

	_, err := h.getContext(w, r)
	if err != nil {
		h.logger.Warnf("getContext: %v", err)
		return
	}

	// Параметры по умолчанию
	now := time.Now()
	from := time.Time{}           // Нулевое время = начало времен
	to := now.Add(24 * time.Hour) // Текущее время + 1 день

	// Парсим from, если передан
	fromStr := r.URL.Query().Get("from")
	if fromStr != "" {
		from, err = time.Parse(time.RFC3339, fromStr)
		if err != nil {
			utils.WriteJSONError(w, http.StatusBadRequest, "invalid from date format")
			return
		}
	}

	// Парсим to, если передан
	toStr := r.URL.Query().Get("to")
	if toStr != "" {
		to, err = time.Parse(time.RFC3339, toStr)
		if err != nil {
			utils.WriteJSONError(w, http.StatusBadRequest, "invalid to date format")
			return
		}
	}

	limit, offset, err := h.parsePaginationParams(r)
	if err != nil {
		utils.WriteJSONError(w, http.StatusBadRequest, "invalid pagination parameters")
		return
	}

	strikes, err := h.strikeService.GetStrikesByDateRange(r.Context(), from, to, limit, offset)
	if err != nil {
		h.logger.Errorf("GetStrikesByDateRange: %v", err)
		switch {
		case errors.Is(err, expectation.ErrInvalidDateRange),
			errors.Is(err, expectation.ErrDateRangeTooLarge):
			utils.WriteJSONError(w, http.StatusBadRequest, err.Error())
		default:
			utils.WriteJSONError(w, http.StatusInternalServerError, "internal server error")
		}
		return
	}

	strikesResp := dto.ToStrikesListResponse(strikes)
	utils.WriteJSON(w, http.StatusOK, strikesResp)
}

// UpdateStrikeStatus godoc
// @Summary Обновить статус жалобы
// @Description Обновляет статус жалобы (только для модераторов)
// @Tags strikes
// @Accept json
// @Produce json
// @Security SessionToken
// @Param id path string true "ID жалобы"
// @Param request body dto.StrikeStatusUpdateRequest true "Данные для обновления статуса"
// @Success 200 {object} dto.StrikeResponse
// @Failure 400 {object} map[string]string "Неверный запрос"
// @Failure 401 {object} map[string]string "Не авторизован"
// @Failure 404 {object} map[string]string "Жалоба не найдена"
// @Failure 500 {object} map[string]string "Внутренняя ошибка сервера"
// @Router /api/strike/{id}/status [put]
func (h *StrikeHandler) UpdateStrikeStatus(w http.ResponseWriter, r *http.Request) {
	h.logger.Trace("strikeHandler.UpdateStrikeStatus")

	_, err := h.getContext(w, r)
	if err != nil {
		h.logger.Warnf("getContext: %v", err)
		return
	}

	strikeID := r.PathValue("id")
	if strikeID == "" {
		utils.WriteJSONError(w, http.StatusBadRequest, "strike ID is required")
		return
	}

	var statusReq dto.StrikeStatusUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&statusReq); err != nil {
		h.logger.Warnf("UpdateStrikeStatus: invalid JSON: %v", err)
		utils.WriteJSONError(w, http.StatusBadRequest, "invalid JSON")
		return
	}

	err = h.strikeService.UpdateStrikeStatus(r.Context(), strikeID, statusReq.Status, statusReq.ModeratorID, statusReq.Note)
	if err != nil {
		h.logger.Errorf("UpdateStrikeStatus: %v", err)
		if errors.Is(err, expectation.ErrStrikeNotFound) {
			utils.WriteJSONError(w, http.StatusNotFound, "strike not found")
			return
		}
		utils.WriteJSONError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	// Получаем обновленную жалобу для ответа
	strike, err := h.strikeService.GetStrikeByID(r.Context(), strikeID)
	if err != nil {
		h.logger.Errorf("UpdateStrikeStatus: failed to get updated strike: %v", err)
		utils.WriteJSONError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	strikeResp := dto.ToStrikeResponse(strike)
	utils.WriteJSON(w, http.StatusOK, strikeResp)
}

// DeleteStrike godoc
// @Summary Удалить жалобу
// @Description Удаляет жалобу (жесткое удаление)
// @Tags strikes
// @Produce json
// @Security SessionToken
// @Param id path string true "ID жалобы"
// @Success 204
// @Failure 400 {object} map[string]string "Неверный запрос"
// @Failure 401 {object} map[string]string "Не авторизован"
// @Failure 404 {object} map[string]string "Жалоба не найдена"
// @Failure 500 {object} map[string]string "Внутренняя ошибка сервера"
// @Router /api/strike/{id} [delete]
func (h *StrikeHandler) DeleteStrike(w http.ResponseWriter, r *http.Request) {
	h.logger.Trace("strikeHandler.DeleteStrike")

	_, err := h.getContext(w, r)
	if err != nil {
		h.logger.Warnf("getContext: %v", err)
		return
	}

	strikeID := r.PathValue("id")
	if strikeID == "" {
		utils.WriteJSONError(w, http.StatusBadRequest, "strike ID is required")
		return
	}

	err = h.strikeService.DeleteStrike(r.Context(), strikeID)
	if err != nil {
		h.logger.Errorf("DeleteStrike: %v", err)
		if errors.Is(err, expectation.ErrStrikeNotFound) {
			utils.WriteJSONError(w, http.StatusNotFound, "strike not found")
			return
		}
		utils.WriteJSONError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	utils.WriteJSON(w, http.StatusNoContent, nil)
}

// GetUserStrikeStats godoc
// @Summary Получить статистику жалоб пользователя
// @Description Возвращает статистику по жалобам на пользователя
// @Tags strikes
// @Produce json
// @Security SessionToken
// @Param user_id path string true "ID пользователя"
// @Success 200 {object} dto.StrikeStats
// @Failure 400 {object} map[string]string "Неверный запрос"
// @Failure 401 {object} map[string]string "Не авторизован"
// @Failure 404 {object} map[string]string "Пользователь не найден"
// @Failure 500 {object} map[string]string "Внутренняя ошибка сервера"
// @Router /api/strike/user/{user_id}/stat [get]
func (h *StrikeHandler) GetUserStrikeStats(w http.ResponseWriter, r *http.Request) {
	h.logger.Trace("strikeHandler.GetUserStrikeStats")

	_, err := h.getContext(w, r)
	if err != nil {
		h.logger.Warnf("getContext: %v", err)
		return
	}

	userID := r.PathValue("user_id")
	if userID == "" {
		utils.WriteJSONError(w, http.StatusBadRequest, "user ID is required")
		return
	}

	stats, err := h.strikeService.GetUserStrikeStats(r.Context(), userID)
	if err != nil {
		h.logger.Errorf("GetUserStrikeStats: %v", err)
		if errors.Is(err, expectation.ErrUserNotFound) {
			utils.WriteJSONError(w, http.StatusNotFound, "user not found")
			return
		}
		utils.WriteJSONError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	utils.WriteJSON(w, http.StatusOK, stats)
}

// parsePaginationParams парсит параметры пагинации из запроса
func (h *StrikeHandler) parsePaginationParams(r *http.Request) (int, int, error) {
	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")

	limit := 20 // default
	if limitStr != "" {
		l, err := strconv.Atoi(limitStr)
		if err != nil || l < 1 || l > 100 {
			return 0, 0, errors.New("invalid limit parameter")
		}
		limit = l
	}

	offset := 0 // default
	if offsetStr != "" {
		o, err := strconv.Atoi(offsetStr)
		if err != nil || o < 0 {
			return 0, 0, errors.New("invalid offset parameter")
		}
		offset = o
	}

	return limit, offset, nil
}
