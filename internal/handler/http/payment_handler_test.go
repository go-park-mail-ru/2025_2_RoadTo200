package handler

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/handler/middleware"
	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/tests/mocks"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestPaymentHandler_CreatePayment(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	paymentService := mocks.NewMockPaymentService(ctrl)
	logger := mocks.NewMockLogger()
	handler := NewPaymentHandler(paymentService, logger)

	userID := uuid.New()

	tests := []struct {
		name           string
		requestBody    string
		userID         uuid.UUID
		setupMocks     func()
		expectedStatus int
		expectedBody   string
	}{
		{
			name:        "successful payment creation",
			requestBody: `{"planId": "week", "amount": 100.0}`,
			userID:      userID,
			setupMocks: func() {
				paymentService.EXPECT().
					CreatePaymentLink(gomock.Any(), userID, "week", 100.0).
					Return("https://payment.link", nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody:   `{"payment_url":"https://payment.link"}`,
		},
		{
			name:        "invalid amount (zero)",
			requestBody: `{"planId": "week", "amount": 0}`,
			userID:      userID,
			setupMocks: func() {
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"error":"amount must be greater than 0"}`,
		},
		{
			name:        "invalid amount (negative)",
			requestBody: `{"planId": "week", "amount": -10}`,
			userID:      userID,
			setupMocks: func() {
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"error":"amount must be greater than 0"}`,
		},
		{
			name:        "empty planId",
			requestBody: `{"planId": "", "amount": 100}`,
			userID:      userID,
			setupMocks: func() {
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"error":"planId is required"}`,
		},
		{
			name:        "user already has active subscription",
			requestBody: `{"planId": "month", "amount": 200}`,
			userID:      userID,
			setupMocks: func() {
				paymentService.EXPECT().
					CreatePaymentLink(gomock.Any(), userID, "month", 200.0).
					Return("", fmt.Errorf("user already has active subscription"))
			},
			expectedStatus: http.StatusConflict,
			expectedBody:   `{"error":"user already has active subscription"}`,
		},
		{
			name:        "invalid plan ID",
			requestBody: `{"planId": "invalid", "amount": 100}`,
			userID:      userID,
			setupMocks: func() {
				paymentService.EXPECT().
					CreatePaymentLink(gomock.Any(), userID, "invalid", 100.0).
					Return("", fmt.Errorf("invalid plan ID: invalid"))
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"error":"invalid plan ID"}`,
		},
		{
			name:        "invalid JSON",
			requestBody: `{"planId": "week", "amount": "invalid"}`,
			userID:      userID,
			setupMocks: func() {
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"error":"invalid JSON"}`,
		},
		{
			name:        "unauthorized - no user ID in context",
			requestBody: `{"planId": "week", "amount": 100}`,
			userID:      uuid.Nil,
			setupMocks: func() {
			},
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   `{"error":"unauthorized"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMocks()

			req := httptest.NewRequest(http.MethodPost, "/api/payment/create",
				strings.NewReader(tt.requestBody))
			req.Header.Set("Content-Type", "application/json")

			// Set user ID in context if provided
			if tt.userID != uuid.Nil {
				ctx := context.WithValue(req.Context(), middleware.UserIDKey, tt.userID)
				req = req.WithContext(ctx)
			}

			w := httptest.NewRecorder()
			handler.CreatePayment(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			assert.JSONEq(t, tt.expectedBody, w.Body.String())
		})
	}
}

func TestPaymentHandler_HandleWebhook(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	paymentService := mocks.NewMockPaymentService(ctrl)
	logger := mocks.NewMockLogger()
	handler := NewPaymentHandler(paymentService, logger)

	tests := []struct {
		name           string
		formData       url.Values
		setupMocks     func()
		expectedStatus int
		expectedBody   string
	}{
		{
			name: "successful webhook processing",
			formData: url.Values{
				"notification_type": {"p2p-incoming"},
				"operation_id":      {"123456"},
				"amount":            {"100.00"},
				"currency":          {"RUB"},
				"datetime":          {"2024-01-01T12:00:00Z"},
				"sender":            {"sender@example.com"},
				"codepro":           {"false"},
				"label":             {"user-id:week:timestamp"},
				"sha1_hash":         {"hash"},
			},
			setupMocks: func() {
				paymentService.EXPECT().
					ProcessWebhook(gomock.Any(), gomock.Any()).
					Return(nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody:   "OK",
		},
		{
			name:     "invalid form data",
			formData: url.Values{},
			setupMocks: func() {
				paymentService.EXPECT().
					ProcessWebhook(gomock.Any(), gomock.Any()).
					Return(fmt.Errorf("some error"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   "failed to process webhook\n",
		},
		{
			name: "webhook processing error",
			formData: url.Values{
				"notification_type": {"p2p-incoming"},
				"operation_id":      {"123456"},
				"amount":            {"100.00"},
				"currency":          {"RUB"},
				"datetime":          {"2024-01-01T12:00:00Z"},
				"sender":            {"sender@example.com"},
				"codepro":           {"false"},
				"label":             {"user-id:week:timestamp"},
				"sha1_hash":         {"hash"},
			},
			setupMocks: func() {
				paymentService.EXPECT().
					ProcessWebhook(gomock.Any(), gomock.Any()).
					Return(assert.AnError)
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   "failed to process webhook\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMocks()

			req := httptest.NewRequest(http.MethodPost, "/notificate_premium",
				strings.NewReader(tt.formData.Encode()))
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
			req.RemoteAddr = "127.0.0.1:12345"

			w := httptest.NewRecorder()
			handler.HandleWebhook(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			assert.Equal(t, tt.expectedBody, w.Body.String())
		})
	}
}
