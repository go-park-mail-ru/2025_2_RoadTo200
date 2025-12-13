package handler

import (
	"encoding/json"
	"net/http"

	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/handler/dto"
	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/handler/middleware"
	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/logger"
	service "github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/service/interfaces"
	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/pkg/utils"
)

type PaymentHandler struct {
	paymentService service.PaymentService
	logger         logger.Log
}

func NewPaymentHandler(paymentService service.PaymentService, l logger.Log) *PaymentHandler {
	return &PaymentHandler{
		paymentService: paymentService,
		logger:         l,
	}
}

// CreatePayment godoc
// @Summary Создать ссылку на оплату премиум-подписки
// @Description Создает ссылку для оплаты выбранного плана премиум-подписки
// @Tags payment
// @Accept json
// @Produce json
// @Security SessionToken
// @Param request body dto.CreatePaymentRequest true "Данные для создания платежа"
// @Success 200 {object} dto.CreatePaymentResponse "Ссылка на оплату"
// @Failure 400 {object} map[string]string "Неверный запрос"
// @Failure 401 {object} map[string]string "Не авторизован"
// @Failure 409 {object} map[string]string "У пользователя уже есть активная подписка"
// @Failure 500 {object} map[string]string "Внутренняя ошибка сервера"
// @Router /api/payment/create [post]
func (h *PaymentHandler) CreatePayment(w http.ResponseWriter, r *http.Request) {
	h.logger.Trace("PaymentHandler.CreatePayment")

	userID, err := middleware.GetUserIDFromContext(r.Context())
	if err != nil {
		h.logger.Warnf("GetUserIDFromContext: %v", err)
		utils.WriteJSONError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	var req dto.CreatePaymentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Warnf("Decode request body: %v", err)
		utils.WriteJSONError(w, http.StatusBadRequest, "invalid JSON")
		return
	}

	// Валидация
	if req.Amount <= 0 {
		utils.WriteJSONError(w, http.StatusBadRequest, "amount must be greater than 0")
		return
	}
	if req.PlanID == "" {
		utils.WriteJSONError(w, http.StatusBadRequest, "planId is required")
		return
	}

	paymentURL, err := h.paymentService.CreatePaymentLink(r.Context(), userID, req.PlanID, req.Amount)
	if err != nil {
		h.logger.Errorf("CreatePaymentLink error: %v", err)
		if err.Error() == "user already has active subscription" {
			utils.WriteJSONError(w, http.StatusConflict, "user already has active subscription")
			return
		}
		if err.Error() == "invalid plan ID: "+req.PlanID {
			utils.WriteJSONError(w, http.StatusBadRequest, "invalid plan ID")
			return
		}
		utils.WriteJSONError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	response := dto.CreatePaymentResponse{
		PaymentURL: paymentURL,
	}

	utils.WriteJSON(w, http.StatusOK, response)
}

// HandleWebhook godoc
// @Summary Обработать webhook от ЮMoney
// @Description Обрабатывает уведомление о поступлении платежа от ЮMoney
// @Tags payment
// @Accept application/x-www-form-urlencoded
// @Produce text/plain
// @Success 200 {string} string "OK"
// @Failure 400 {string} string "Неверный запрос"
// @Failure 500 {string} string "Внутренняя ошибка сервера"
// @Router /notificate_premium [post]
func (h *PaymentHandler) HandleWebhook(w http.ResponseWriter, r *http.Request) {
	h.logger.Infof("PaymentHandler.HandleWebhook: received webhook from %s", r.RemoteAddr)
	h.logger.Tracef("PaymentHandler.HandleWebhook: headers: %+v", r.Header)

	// Парсим form-data из webhook
	if err := r.ParseForm(); err != nil {
		h.logger.Errorf("ParseForm error: %v", err)
		http.Error(w, "invalid form data", http.StatusBadRequest)
		return
	}

	// Преобразуем form values в map[string]string
	webhookData := make(map[string]string)
	for key, values := range r.PostForm {
		if len(values) > 0 {
			webhookData[key] = values[0]
		}
	}

	h.logger.Infof("Received webhook data: %+v", webhookData)

	// Обрабатываем webhook
	if err := h.paymentService.ProcessWebhook(r.Context(), webhookData); err != nil {
		h.logger.Errorf("ProcessWebhook error: %v", err)
		http.Error(w, "failed to process webhook", http.StatusInternalServerError)
		return
	}

	h.logger.Infof("Webhook processed successfully")

	// Возвращаем 200 OK как требуется по документации ЮMoney
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}
