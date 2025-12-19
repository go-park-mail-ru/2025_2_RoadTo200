package handler

import (
	"net/http"
	"strconv"

	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/handler/dto"
	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/handler/middleware"
	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/logger"
	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/service/interfaces"
	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/pkg/utils"
	"github.com/google/uuid"
)

type SupportHandler struct {
	supportService service.SupportService
	logger         logger.Log
}

func NewSupportHandler(supportService service.SupportService, l logger.Log) *SupportHandler {
	return &SupportHandler{
		supportService: supportService,
		logger:         l,
	}
}

// CreateSupportTicket godoc
// @Summary Создать обращение в поддержку
// @Description Создает новое обращение в техническую поддержку
// @Tags support
// @Accept json
// @Produce json
// @Security SessionToken
// @Param request body dto.SupportTicketRequest true "Данные обращения"
// @Success 201 {object} dto.SupportTicketResponse "Обращение создано"
// @Failure 400 {object} map[string]string "Неверный запрос"
// @Failure 401 {object} map[string]string "Не авторизован"
// @Failure 500 {object} map[string]string "Внутренняя ошибка сервера"
// @Router /api/report [post]
func (h *SupportHandler) CreateSupportTicket(w http.ResponseWriter, r *http.Request) {
	userID, err := middleware.GetUserIDFromContext(r.Context())
	if err != nil {
		h.logger.Warnf("get user id from context: %v", err)
		utils.WriteJSONError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	var req dto.SupportTicketRequest
	if err := utils.ReadJSON(r, &req); err != nil {
		h.logger.Warnf("read json body: %v", err)
		utils.WriteJSONError(w, http.StatusBadRequest, "invalid JSON")
		return
	}

	// Валидация обязательных полей
	if req.Category == "" || req.Text == "" || req.Email == "" {
		h.logger.Warnf("invalid request: %v", req)
		utils.WriteJSONError(w, http.StatusBadRequest, "category, text and email are required")
		return
	}

	ticket, err := h.supportService.CreateTicket(userID, &req)
	if err != nil {
		h.logger.Errorf("create ticket: %v", err)
		utils.WriteJSONError(w, http.StatusInternalServerError, "failed to create support ticket")
		return
	}

	utils.WriteJSON(w, http.StatusCreated, ticket)
}

// GetUserSupportTickets godoc
// @Summary Получить мои обращения
// @Description Возвращает список обращений текущего пользователя
// @Tags support
// @Produce json
// @Security SessionToken
// @Param limit query int false "Лимит обращений" default(20) minimum(1) maximum(50)
// @Param offset query int false "Смещение" default(0) minimum(0)
// @Success 200 {object} dto.SupportTicketsListResponse "Список обращений"
// @Failure 400 {object} map[string]string "Неверные параметры запроса"
// @Failure 401 {object} map[string]string "Не авторизован"
// @Failure 500 {object} map[string]string "Внутренняя ошибка сервера"
// @Router /api/report [get]
func (h *SupportHandler) GetUserSupportTickets(w http.ResponseWriter, r *http.Request) {
	userID, err := middleware.GetUserIDFromContext(r.Context())
	if err != nil {
		h.logger.Warnf("get user id from context: %v", err)
		utils.WriteJSONError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	// Парсим параметры пагинации
	limit, offset := h.parseQueryParams(r)

	tickets, err := h.supportService.GetUserTickets(userID, limit, offset)
	if err != nil {
		h.logger.Errorf("get user tickets: %v", err)
		utils.WriteJSONError(w, http.StatusInternalServerError, "failed to get support tickets")
		return
	}

	utils.WriteJSON(w, http.StatusOK, tickets)
}

// GetUserSupportTicket godoc
// @Summary Получить детальную информацию об обращении
// @Description Возвращает детальную информацию об обращении по ID
// @Tags support
// @Produce json
// @Security SessionToken
// @Param ticket_id path string true "ID обращения"
// @Success 200 {object} dto.SupportTicketDetailResponse "Детальная информация об обращении"
// @Failure 400 {object} map[string]string "Неверный ID обращения"
// @Failure 401 {object} map[string]string "Не авторизован"
// @Failure 403 {object} map[string]string "Доступ запрещен"
// @Failure 404 {object} map[string]string "Обращение не найдено"
// @Failure 500 {object} map[string]string "Внутренняя ошибка сервера"
// @Router /api/report/{ticket_id} [get]
func (h *SupportHandler) GetUserSupportTicket(w http.ResponseWriter, r *http.Request) {
	userID, err := middleware.GetUserIDFromContext(r.Context())
	if err != nil {
		h.logger.Warnf("get user id from context: %v", err)
		utils.WriteJSONError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	// Получаем ticket_id из URL
	ticketIDStr := r.URL.Path[len("/api/report/"):]
	if ticketIDStr == "" {
		h.logger.Warnf("ticketIDStr is empty")
		utils.WriteJSONError(w, http.StatusBadRequest, "ticket ID is required")
		return
	}

	ticketID, err := uuid.Parse(ticketIDStr)
	if err != nil {
		h.logger.Warnf("ticketIDStr is invalid: %v", err)
		utils.WriteJSONError(w, http.StatusBadRequest, "invalid ticket ID")
		return
	}

	ticket, err := h.supportService.GetUserTicket(userID, ticketID)
	if err != nil {
		if err.Error() == "access denied" {
			h.logger.Warnf("ticket does not exist: %v", err)
			utils.WriteJSONError(w, http.StatusForbidden, "access denied")
			return
		}
		h.logger.Errorf("get ticket: %v", err)
		utils.WriteJSONError(w, http.StatusInternalServerError, "failed to get ticket")
		return
	}

	utils.WriteJSON(w, http.StatusOK, ticket)
}

// parseQueryParams парсит параметры запроса для пагинации
func (h *SupportHandler) parseQueryParams(r *http.Request) (int, int) {
	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")

	limit := 20 // значение по умолчанию
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
