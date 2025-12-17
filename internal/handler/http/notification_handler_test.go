package handler

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/domain/constants"
	domain "github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/domain/entities"
	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/handler/middleware"
	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/tests/mocks"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestNotificationHandler_GetNotifications(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	notificationService := mocks.NewMockNotificationService(ctrl)
	logger := mocks.NewMockLogger()
	handler := NewNotificationHandler(notificationService, logger)

	userID := uuid.New()
	notificationID := uuid.New()

	fromUserID := uuid.New()
	fixedTime := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)
	notifications := []domain.Notification{
		{
			ID:         notificationID,
			UserID:     userID,
			Type:       constants.NotificationTypeLike,
			FromUserID: &fromUserID,
			IsRead:     false,
			CreatedAt:  fixedTime,
		},
	}

	tests := []struct {
		name           string
		queryParams    string
		userID         uuid.UUID
		setupMocks     func()
		expectedStatus int
		expectedBody   string
	}{
		{
			name:        "successful get notifications with default params",
			queryParams: "",
			userID:      userID,
			setupMocks: func() {
				notificationService.EXPECT().
					GetNotifications(gomock.Any(), userID, 20, 0).
					Return(notifications, nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody: `{
				"notifications": [
					{
						"id": "` + notificationID.String() + `",
						"user_id": "` + userID.String() + `",
						"type": "like",
						"from_user_id": "` + fromUserID.String() + `",
						"is_read": false,
						"created_at": "2024-01-01T12:00:00Z"
					}
				],
				"total": 1,
				"limit": 20,
				"offset": 0
			}`,
		},
		{
			name:        "successful get notifications with custom params",
			queryParams: "?limit=10&offset=5",
			userID:      userID,
			setupMocks: func() {
				notificationService.EXPECT().
					GetNotifications(gomock.Any(), userID, 10, 5).
					Return([]domain.Notification{}, nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody: `{
				"notifications": [],
				"total": 0,
				"limit": 10,
				"offset": 5
			}`,
		},
		{
			name:        "invalid limit parameter (too high)",
			queryParams: "?limit=200",
			userID:      userID,
			setupMocks: func() {
				notificationService.EXPECT().
					GetNotifications(gomock.Any(), userID, 20, 0). // Should use default limit
					Return([]domain.Notification{}, nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody: `{
				"notifications": [],
				"total": 0,
				"limit": 20,
				"offset": 0
			}`,
		},
		{
			name:        "invalid offset parameter (negative)",
			queryParams: "?offset=-5",
			userID:      userID,
			setupMocks: func() {
				notificationService.EXPECT().
					GetNotifications(gomock.Any(), userID, 20, 0). // Should use default offset
					Return([]domain.Notification{}, nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody: `{
				"notifications": [],
				"total": 0,
				"limit": 20,
				"offset": 0
			}`,
		},
		{
			name:        "service error",
			queryParams: "",
			userID:      userID,
			setupMocks: func() {
				notificationService.EXPECT().
					GetNotifications(gomock.Any(), userID, 20, 0).
					Return(nil, assert.AnError)
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   `{"error":"internal error"}`,
		},
		{
			name:        "unauthorized - no user ID in context",
			queryParams: "",
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

			req := httptest.NewRequest(http.MethodGet, "/api/notifications"+tt.queryParams, nil)

			// Set user ID in context if provided
			if tt.userID != uuid.Nil {
				ctx := context.WithValue(req.Context(), middleware.UserIDKey, tt.userID)
				req = req.WithContext(ctx)
			}

			w := httptest.NewRecorder()
			handler.GetNotifications(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			// For JSON responses, compare as JSON
			if tt.expectedStatus == http.StatusOK {
				assert.JSONEq(t, tt.expectedBody, w.Body.String())
			} else {
				assert.Equal(t, tt.expectedBody, w.Body.String())
			}
		})
	}
}

func TestNotificationHandler_MarkAsRead(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	notificationService := mocks.NewMockNotificationService(ctrl)
	logger := mocks.NewMockLogger()
	handler := NewNotificationHandler(notificationService, logger)

	userID := uuid.New()
	notificationID := uuid.New()

	tests := []struct {
		name           string
		urlPath        string
		userID         uuid.UUID
		setupMocks     func()
		expectedStatus int
		expectedBody   string
	}{
		{
			name:    "successful mark as read",
			urlPath: "/api/notifications/" + notificationID.String() + "/read",
			userID:  userID,
			setupMocks: func() {
				notificationService.EXPECT().
					MarkAsRead(gomock.Any(), userID, notificationID).
					Return(nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody:   `{"message":"Notification marked as read"}`,
		},
		{
			name:    "notification not found",
			urlPath: "/api/notifications/" + notificationID.String() + "/read",
			userID:  userID,
			setupMocks: func() {
				notificationService.EXPECT().
					MarkAsRead(gomock.Any(), userID, notificationID).
					Return(fmt.Errorf("notification not found or access denied"))
			},
			expectedStatus: http.StatusNotFound,
			expectedBody:   `{"error":"notification not found"}`,
		},
		{
			name:    "invalid notification ID format",
			urlPath: "/api/notifications/invalid-uuid/read",
			userID:  userID,
			setupMocks: func() {
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"error":"invalid notification_id"}`,
		},
		{
			name:    "missing notification ID in path",
			urlPath: "/api/notifications//read",
			userID:  userID,
			setupMocks: func() {
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"error":"notification_id is required"}`,
		},
		{
			name:    "service error",
			urlPath: "/api/notifications/" + notificationID.String() + "/read",
			userID:  userID,
			setupMocks: func() {
				notificationService.EXPECT().
					MarkAsRead(gomock.Any(), userID, notificationID).
					Return(fmt.Errorf("some other error")) // Service returns generic error
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   `{"error":"internal error"}`,
		},
		{
			name:    "unauthorized - no user ID in context",
			urlPath: "/api/notifications/" + notificationID.String() + "/read",
			userID:  uuid.Nil,
			setupMocks: func() {
			},
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   `{"error":"unauthorized"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMocks()

			req := httptest.NewRequest(http.MethodPut, tt.urlPath, nil)

			// Set user ID in context if provided
			if tt.userID != uuid.Nil {
				ctx := context.WithValue(req.Context(), middleware.UserIDKey, tt.userID)
				req = req.WithContext(ctx)
			}

			w := httptest.NewRecorder()
			handler.MarkAsRead(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			// For JSON responses, compare as JSON
			if tt.expectedStatus == http.StatusOK {
				assert.JSONEq(t, tt.expectedBody, w.Body.String())
			} else {
				assert.Equal(t, tt.expectedBody, w.Body.String())
			}
		})
	}
}
