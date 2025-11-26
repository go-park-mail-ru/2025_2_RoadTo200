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

func TestMatchHandler_GetUserMatches_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockLogger := mocks.NewMockLogger()
	mockMatchService := mocks.NewMockMatchService(ctrl)
	handler := NewMatchHandler(mockMatchService, mockLogger)

	userID := uuid.New()

	matchResponse := &dto.MatchesResponse{
		Matches: []dto.MatchResponse{},
		Total:   1,
	}

	mockMatchService.EXPECT().
		GetUserMatches(gomock.Any(), userID, 20, 0).
		Return(matchResponse, nil)

	req := httptest.NewRequest(http.MethodGet, "/api/match", nil)
	req = req.WithContext(context.WithValue(req.Context(), middleware.UserIDKey, userID))
	rec := httptest.NewRecorder()

	handler.GetUserMatches(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	var response dto.MatchesResponse
	json.Unmarshal(rec.Body.Bytes(), &response)
	assert.Equal(t, 1, response.Total)
}

func TestMatchHandler_GetUserMatches_WithPagination(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockLogger := mocks.NewMockLogger()
	mockMatchService := mocks.NewMockMatchService(ctrl)
	handler := NewMatchHandler(mockMatchService, mockLogger)

	userID := uuid.New()

	mockMatchService.EXPECT().
		GetUserMatches(gomock.Any(), userID, 10, 5).
		Return(&dto.MatchesResponse{Matches: []dto.MatchResponse{}, Total: 0}, nil)

	req := httptest.NewRequest(http.MethodGet, "/api/match?limit=10&offset=5", nil)
	req = req.WithContext(context.WithValue(req.Context(), middleware.UserIDKey, userID))
	rec := httptest.NewRecorder()

	handler.GetUserMatches(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestMatchHandler_Unmatch_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockLogger := mocks.NewMockLogger()
	mockMatchService := mocks.NewMockMatchService(ctrl)
	handler := NewMatchHandler(mockMatchService, mockLogger)

	userID := uuid.New()
	targetUserID := uuid.New()

	unmatchReq := dto.UnmatchRequest{
		TargetUserID: targetUserID.String(),
	}

	mockMatchService.EXPECT().
		Unmatch(gomock.Any(), userID, targetUserID).
		Return(nil)

	body, _ := json.Marshal(unmatchReq)
	req := httptest.NewRequest(http.MethodDelete, "/api/match/unmatch", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req = req.WithContext(context.WithValue(req.Context(), middleware.UserIDKey, userID))
	rec := httptest.NewRecorder()

	handler.Unmatch(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestMatchHandler_Unmatch_InvalidUserID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockLogger := mocks.NewMockLogger()
	mockMatchService := mocks.NewMockMatchService(ctrl)
	handler := NewMatchHandler(mockMatchService, mockLogger)

	userID := uuid.New()

	unmatchReq := dto.UnmatchRequest{
		TargetUserID: "invalid-uuid",
	}

	body, _ := json.Marshal(unmatchReq)
	req := httptest.NewRequest(http.MethodDelete, "/api/match/unmatch", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req = req.WithContext(context.WithValue(req.Context(), middleware.UserIDKey, userID))
	rec := httptest.NewRecorder()

	handler.Unmatch(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestMatchHandler_Unmatch_MatchNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockLogger := mocks.NewMockLogger()
	mockMatchService := mocks.NewMockMatchService(ctrl)
	handler := NewMatchHandler(mockMatchService, mockLogger)

	userID := uuid.New()
	targetUserID := uuid.New()

	unmatchReq := dto.UnmatchRequest{
		TargetUserID: targetUserID.String(),
	}

	mockMatchService.EXPECT().
		Unmatch(gomock.Any(), userID, targetUserID).
		Return(errors.ErrMatchNotFound)

	body, _ := json.Marshal(unmatchReq)
	req := httptest.NewRequest(http.MethodDelete, "/api/match/unmatch", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req = req.WithContext(context.WithValue(req.Context(), middleware.UserIDKey, userID))
	rec := httptest.NewRecorder()

	handler.Unmatch(rec, req)

	assert.Equal(t, http.StatusNotFound, rec.Code)
}
