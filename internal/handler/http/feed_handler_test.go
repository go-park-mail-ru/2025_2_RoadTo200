package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/handler/dto"
	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/handler/middleware"
	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/tests/mocks"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestFeedHandler_GetFeed_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockLogger := mocks.NewMockLogger()
	mockFeedService := mocks.NewMockFeedService(ctrl)
	handler := NewFeedHandler(mockFeedService, mockLogger)

	userID := uuid.New()

	feedUsers := []dto.FeedUser{
		{
			ID:   uuid.New().String(),
			Name: "User 1",
		},
	}

	mockFeedService.EXPECT().
		GetFeed(gomock.Any(), userID, 15, 0).
		Return(feedUsers, nil)

	req := httptest.NewRequest(http.MethodGet, "/api/feed", nil)
	req = req.WithContext(context.WithValue(req.Context(), middleware.UserIDKey, userID))
	rec := httptest.NewRecorder()

	handler.GetFeed(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	var response dto.FeedResponse
	json.Unmarshal(rec.Body.Bytes(), &response)
	assert.Equal(t, 1, response.Total)
	assert.Equal(t, 15, response.Limit)
}

func TestFeedHandler_GetFeed_WithPagination(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockLogger := mocks.NewMockLogger()
	mockFeedService := mocks.NewMockFeedService(ctrl)
	handler := NewFeedHandler(mockFeedService, mockLogger)

	userID := uuid.New()

	mockFeedService.EXPECT().
		GetFeed(gomock.Any(), userID, 10, 20).
		Return([]dto.FeedUser{}, nil)

	req := httptest.NewRequest(http.MethodGet, "/api/feed?limit=10&offset=20", nil)
	req = req.WithContext(context.WithValue(req.Context(), middleware.UserIDKey, userID))
	rec := httptest.NewRecorder()

	handler.GetFeed(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	var response dto.FeedResponse
	json.Unmarshal(rec.Body.Bytes(), &response)
	assert.Equal(t, 10, response.Limit)
	assert.Equal(t, 20, response.Offset)
}

func TestFeedHandler_GetFeed_NoUserID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockLogger := mocks.NewMockLogger()
	mockFeedService := mocks.NewMockFeedService(ctrl)
	handler := NewFeedHandler(mockFeedService, mockLogger)

	req := httptest.NewRequest(http.MethodGet, "/api/feed", nil)
	rec := httptest.NewRecorder()

	handler.GetFeed(rec, req)

	assert.Equal(t, http.StatusUnauthorized, rec.Code)
}
