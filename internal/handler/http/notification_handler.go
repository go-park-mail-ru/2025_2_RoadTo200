package handler

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/domain/errors"
	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/handler/dto"
	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/handler/middleware"
	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/logger"
	service "github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/service/interfaces"
	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/pkg/utils"
	"github.com/google/uuid"
)

type NotificationHandler struct {
	notificationService service.NotificationService
	logger              logger.Log
}

func NewNotificationHandler(notificationService service.NotificationService, l logger.Log) *NotificationHandler {
	return &NotificationHandler{
		notificationService: notificationService,
		logger:              l,
	}
}

// GetNotifications godoc
// @Summary Получить список уведомлений
// @Description Возвращает список уведомлений текущего пользователя с пагинацией
// @Tags notifications
// @Produce json
// @Security SessionToken
// @Param limit query int false "Количество уведомлений (по умолчанию 20)"
// @Param offset query int false "Смещение (по умолчанию 0)"
// @Success 200 {object} dto.NotificationsResponse "Список уведомлений"
// @Failure 401 {object} map[string]string "Не авторизован"
// @Failure 500 {object} map[string]string "Внутренняя ошибка сервера"
// @Router /api/notifications [get]
func (h *NotificationHandler) GetNotifications(w http.ResponseWriter, r *http.Request) {
	h.logger.Trace("NotificationHandler.GetNotifications")

	// Получаем userID из контекста (установлен middleware)
	userID, err := middleware.GetUserIDFromContext(r.Context())
	if err != nil {
		h.logger.Warnf("GetUserIDFromContext: %v", err)
		utils.WriteJSONError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	// Парсим параметры пагинации
	limit := 20
	offset := 0

	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		if parsedLimit, err := strconv.Atoi(limitStr); err == nil && parsedLimit > 0 && parsedLimit <= 100 {
			limit = parsedLimit
		}
	}

	if offsetStr := r.URL.Query().Get("offset"); offsetStr != "" {
		if parsedOffset, err := strconv.Atoi(offsetStr); err == nil && parsedOffset >= 0 {
			offset = parsedOffset
		}
	}

	// Получаем уведомления
	notifications, err := h.notificationService.GetNotifications(r.Context(), userID, limit, offset)
	if err != nil {
		h.logger.Errorf("Failed to get notifications: %v", err)
		utils.WriteJSONError(w, http.StatusInternalServerError, errors.ErrInternalError.Error())
		return
	}

	utils.WriteJSON(w, http.StatusOK, dto.NotificationsResponse{
		Notifications: notifications,
		Total:         len(notifications),
		Limit:         limit,
		Offset:        offset,
	})
}

// MarkAsRead godoc
// @Summary Пометить уведомление как прочитанное
// @Description Помечает указанное уведомление как прочитанное
// @Tags notifications
// @Accept json
// @Produce json
// @Security SessionToken
// @Param notification_id path string true "ID уведомления"
// @Success 200 {object} map[string]string "Уведомление помечено как прочитанное"
// @Failure 400 {object} map[string]string "Неверный ID уведомления"
// @Failure 401 {object} map[string]string "Не авторизован"
// @Failure 404 {object} map[string]string "Уведомление не найдено"
// @Failure 500 {object} map[string]string "Внутренняя ошибка сервера"
// @Router /api/notifications/{notification_id}/read [put]
func (h *NotificationHandler) MarkAsRead(w http.ResponseWriter, r *http.Request) {
	h.logger.Trace("NotificationHandler.MarkAsRead")

	// Получаем userID из контекста
	userID, err := middleware.GetUserIDFromContext(r.Context())
	if err != nil {
		h.logger.Warnf("GetUserIDFromContext: %v", err)
		utils.WriteJSONError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	// Получаем notification_id из пути
	// Извлекаем ID из URL: /api/notifications/{id}/read
	path := r.URL.Path
	parts := strings.Split(path, "/")
	var notificationIDStr string
	for i, part := range parts {
		if part == "notifications" && i+1 < len(parts) {
			notificationIDStr = parts[i+1]
			break
		}
	}

	if notificationIDStr == "" {
		utils.WriteJSONError(w, http.StatusBadRequest, "notification_id is required")
		return
	}

	notificationID, err := uuid.Parse(notificationIDStr)
	if err != nil {
		utils.WriteJSONError(w, http.StatusBadRequest, "invalid notification_id")
		return
	}

	// Помечаем как прочитанное
	if err := h.notificationService.MarkAsRead(r.Context(), userID, notificationID); err != nil {
		h.logger.Errorf("Failed to mark notification as read: %v", err)
		if err.Error() == "notification not found or access denied" {
			utils.WriteJSONError(w, http.StatusNotFound, "notification not found")
			return
		}
		utils.WriteJSONError(w, http.StatusInternalServerError, errors.ErrInternalError.Error())
		return
	}

	utils.WriteJSON(w, http.StatusOK, map[string]string{
		"message": "Notification marked as read",
	})
}
