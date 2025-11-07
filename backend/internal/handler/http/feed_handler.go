package handler

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/handler/middleware"
	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/logger"
	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/service/interfaces"
	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/pkg/utils"
	"github.com/google/uuid"
)

type FeedHandler struct {
	feedService service.FeedService
	logger      *logger.Logger
}

func NewFeedHandler(feedService service.FeedService, l *logger.Logger) *FeedHandler {
	return &FeedHandler{
		feedService: feedService,
		logger:      l,
	}
}

// FeedResponse represents feed response
type FeedResponse struct {
	Users  []interface{} `json:"users"`
	Limit  int           `json:"limit"`
	Offset int           `json:"offset"`
	Total  int           `json:"total"`
}

// GetFeed godoc
// @Summary Получить ленту пользователей
// @Description Возвращает ленту пользователей для свайпинга на основе предпочтений и местоположения
// @Tags feed
// @Produce json
// @Security SessionToken
// @Param limit query int false "Лимит пользователей (максимум 50)" default(15) minimum(1) maximum(50)
// @Param offset query int false "Смещение для пагинации" default(0) minimum(0)
// @Success 200 {object} FeedResponse "Лента пользователей"
// @Failure 401 {object} map[string]string "Не авторизован"
// @Failure 500 {object} map[string]string "Внутренняя ошибка сервера"
// @Router /api/feed [get]
func (h *FeedHandler) GetFeed(w http.ResponseWriter, r *http.Request) {
	userID, err := h.getUserIDFromContext(r)
	fmt.Printf("UserID from context: %v, err: %v\n", userID, err)
	if err != nil {
		h.logger.Warnf("GetFeed err: %v\n", err)
		utils.WriteJSONError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	// Парсим параметры запроса
	limit, offset := h.parseQueryParams(r)

	// Получаем ленту
	users, err := h.feedService.GetFeed(userID, limit, offset)
	if err != nil {
		h.logger.Warnf("GetFeed err: %v\n", err)
		utils.WriteJSONError(w, http.StatusInternalServerError, "failed to get feed")
		return
	}
	h.logger.Debugf("GetFeed users: %v\n", users)

	// Преобразуем в интерфейс для ответа
	usersResponse := make([]interface{}, len(users))
	for i, user := range users {
		usersResponse[i] = user
	}

	response := FeedResponse{
		Users:  usersResponse,
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

// getUserIDFromContext извлекает userID из контекста (будет установлен в middleware)
func (h *FeedHandler) getUserIDFromContext(r *http.Request) (uuid.UUID, error) {
	return middleware.GetUserIDFromContext(r.Context())
}
