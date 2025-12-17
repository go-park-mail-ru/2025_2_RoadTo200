package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/domain/errors"
	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/handler/dto"
	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/handler/middleware"
	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/tests/mocks"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestSwipeHandler_ProcessSwipe_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockLogger := mocks.NewMockLogger()
	mockSwipeService := mocks.NewMockSwipeService(ctrl)
	handler := NewSwipeHandler(mockSwipeService, mockLogger)

	userID := uuid.New()
	cardID := uuid.New()

	swipeReq := dto.SwipeRequest{
		CardID: cardID.String(),
		Action: "like",
	}

	swipeResp := &dto.SwipeResponse{
		IsMatch: true,
		Message: "It's a match!",
	}

	mockSwipeService.EXPECT().
		ProcessSwipe(gomock.Any(), userID, gomock.Any()).
		Return(swipeResp, nil)

	body, _ := json.Marshal(swipeReq)
	req := httptest.NewRequest(http.MethodPost, "/api/swipe", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req = req.WithContext(context.WithValue(req.Context(), middleware.UserIDKey, userID))
	rec := httptest.NewRecorder()

	handler.ProcessSwipe(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	var response dto.SwipeResponse
	json.Unmarshal(rec.Body.Bytes(), &response)
	assert.True(t, response.IsMatch)
}

func TestSwipeHandler_ProcessSwipe_InvalidJSON(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockLogger := mocks.NewMockLogger()
	mockSwipeService := mocks.NewMockSwipeService(ctrl)
	handler := NewSwipeHandler(mockSwipeService, mockLogger)

	userID := uuid.New()

	req := httptest.NewRequest(http.MethodPost, "/api/swipe", bytes.NewReader([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")
	req = req.WithContext(context.WithValue(req.Context(), middleware.UserIDKey, userID))
	rec := httptest.NewRecorder()

	handler.ProcessSwipe(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestSwipeHandler_ProcessSwipe_CannotSwipeSelf(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockLogger := mocks.NewMockLogger()
	mockSwipeService := mocks.NewMockSwipeService(ctrl)
	handler := NewSwipeHandler(mockSwipeService, mockLogger)

	userID := uuid.New()
	swipeReq := dto.SwipeRequest{
		CardID: uuid.New().String(),
		Action: "like",
	}

	mockSwipeService.EXPECT().
		ProcessSwipe(gomock.Any(), userID, gomock.Any()).
		Return(nil, errors.ErrCannotSwipeSelf)

	body, _ := json.Marshal(swipeReq)
	req := httptest.NewRequest(http.MethodPost, "/api/swipe", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req = req.WithContext(context.WithValue(req.Context(), middleware.UserIDKey, userID))
	rec := httptest.NewRecorder()

	handler.ProcessSwipe(rec, req)

	assert.Equal(t, http.StatusForbidden, rec.Code)
}

func TestSwipeHandler_ProcessSwipe_InvalidAction(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockLogger := mocks.NewMockLogger()
	mockSwipeService := mocks.NewMockSwipeService(ctrl)
	handler := NewSwipeHandler(mockSwipeService, mockLogger)

	userID := uuid.New()
	swipeReq := dto.SwipeRequest{
		CardID: uuid.New().String(),
		Action: "invalid",
	}

	mockSwipeService.EXPECT().
		ProcessSwipe(gomock.Any(), userID, gomock.Any()).
		Return(nil, errors.ErrInvalidSwipeAction)

	body, _ := json.Marshal(swipeReq)
	req := httptest.NewRequest(http.MethodPost, "/api/swipe", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req = req.WithContext(context.WithValue(req.Context(), middleware.UserIDKey, userID))
	rec := httptest.NewRecorder()

	handler.ProcessSwipe(rec, req)

	assert.Equal(t, http.StatusUnprocessableEntity, rec.Code)
}
