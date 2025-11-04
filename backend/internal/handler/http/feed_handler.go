package handler

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/handler/middleware"
	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/service/interfaces"
	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/pkg/utils"
	"github.com/google/uuid"
)

type FeedHandler struct {
	feedService service.FeedService
}

func NewFeedHandler(feedService service.FeedService) *FeedHandler {
	return &FeedHandler{
		feedService: feedService,
	}
}

// GetFeed returns feed of users
// @Summary Get user feed
// @Description Get feed of potential matches with pagination
// @Tags feed
// @Accept json
// @Produce json
// @Security CookieAuth
// @Param limit query int false "Number of users to return" default(15)
// @Param offset query int false "Offset for pagination" default(0)
// @Success 200 {object} service.FeedResponse
// @Failure 400 {object} service.ErrorResponse
// @Failure 401 {object} service.ErrorResponse
// @Failure 500 {object} service.ErrorResponse
// @Router /api/feed [get]
func (h *FeedHandler) GetFeed(w http.ResponseWriter, r *http.Request) {
	userID, err := h.getUserIDFromContext(r)
	fmt.Printf("UserID from context: %v, err: %v\n", userID, err)
	if err != nil {
		utils.WriteJSONError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	// Парсим параметры запроса
	limit, offset := h.parseQueryParams(r)

	// Получаем ленту
	users, err := h.feedService.GetFeed(userID, limit, offset)
	if err != nil {
		utils.WriteJSONError(w, http.StatusInternalServerError, "failed to get feed")
		return
	}

	// Преобразуем в интерфейс для ответа
	usersResponse := make([]interface{}, len(users))
	for i, user := range users {
		usersResponse[i] = user
	}

	response := service.FeedResponse{
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
