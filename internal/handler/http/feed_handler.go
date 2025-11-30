package handler

import (
	"net/http"
	"strconv"

	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/handler/dto"
	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/handler/middleware"
	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/service/interfaces"
	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/pkg/logger"
	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/pkg/utils"
)

type FeedHandler struct {
	feedService service.FeedService
	logger      logger.Log
}

func NewFeedHandler(feedService service.FeedService, l logger.Log) *FeedHandler {
	return &FeedHandler{
		feedService: feedService,
		logger:      l,
	}
}

// GetFeed godoc
// @Summary Получить ленту пользователей
// @Description Возвращает ленту пользователей для свайпинга с основной информацией и фотографиями
// @Tags feed
// @Produce json
// @Security SessionToken
// @Param limit query int false "Лимит пользователей (максимум 50)" default(15) minimum(1) maximum(50)
// @Param offset query int false "Смещение для пагинации" default(0) minimum(0)
// @Success 200 {object} dto.FeedResponse "Лента пользователей"
// @Failure 401 {object} map[string]string "Не авторизован"
// @Failure 500 {object} map[string]string "Внутренняя ошибка сервера"
// @Router /api/feed [get]
func (h *FeedHandler) GetFeed(w http.ResponseWriter, r *http.Request) {
	h.logger.Trace("feedService.GetFeed")

	userID, err := middleware.GetUserIDFromContext(r.Context())
	if err != nil {
		h.logger.Warnf("GetFeed err: %v\n", err)
		utils.WriteJSONError(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	h.logger.Debugf("UserID from context: %v", userID)

	// Парсим параметры запроса
	limit, offset := h.parseQueryParams(r)

	// Получаем ленту
	users, err := h.feedService.GetFeed(r.Context(), userID, limit, offset)
	if err != nil {
		h.logger.Errorf("GetFeed err: %v\n", err)
		utils.WriteJSONError(w, http.StatusInternalServerError, "failed to get feed")
		return
	}
	h.logger.Debugf("GetFeed users: %v\n", users)

	response := dto.FeedResponse{
		Users:  users,
		Limit:  limit,
		Offset: offset,
		Total:  len(users),
	}

	utils.WriteJSON(w, http.StatusOK, response)
}

// parseQueryParams парсит параметры запроса
func (h *FeedHandler) parseQueryParams(r *http.Request) (int, int) {
	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")

	limit := 15 // дефолтное значение
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
