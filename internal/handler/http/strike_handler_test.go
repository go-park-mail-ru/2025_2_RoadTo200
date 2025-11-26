package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/domain/constants"
	domain "github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/domain/entities"
	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/domain/errors"
	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/handler/dto"
	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/handler/middleware"
	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/tests/mocks"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestStrikeHandler_CreateStrike_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockLogger := mocks.NewMockLogger()
	mockStrikeService := mocks.NewMockStrikeService(ctrl)
	handler := NewStrikeHandler(mockStrikeService, mockLogger)

	userID := uuid.New()
	targetUserID := uuid.New()

	strikeReq := dto.StrikeCreateRequest{
		TargetUserID: targetUserID,
		Type:         constants.StrikeTypeHarassment,
		Reason:       "Inappropriate behavior",
	}

	strike := &domain.Strike{
		ID:           uuid.New(),
		ReporterID:   userID,
		TargetUserID: targetUserID,
		Type:         constants.StrikeTypeHarassment,
	}

	mockStrikeService.EXPECT().
		CreateStrike(gomock.Any(), gomock.Any()).
		Return(strike, nil)

	body, _ := json.Marshal(strikeReq)
	req := httptest.NewRequest(http.MethodPost, "/api/strike", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req = req.WithContext(context.WithValue(req.Context(), middleware.UserIDKey, userID))
	rec := httptest.NewRecorder()

	handler.CreateStrike(rec, req)

	assert.Equal(t, http.StatusCreated, rec.Code)
}

func TestStrikeHandler_CreateStrike_InvalidJSON(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockLogger := mocks.NewMockLogger()
	mockStrikeService := mocks.NewMockStrikeService(ctrl)
	handler := NewStrikeHandler(mockStrikeService, mockLogger)

	userID := uuid.New()

	req := httptest.NewRequest(http.MethodPost, "/api/strike", bytes.NewReader([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")
	req = req.WithContext(context.WithValue(req.Context(), middleware.UserIDKey, userID))
	rec := httptest.NewRecorder()

	handler.CreateStrike(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestStrikeHandler_CreateStrike_SelfStrike(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockLogger := mocks.NewMockLogger()
	mockStrikeService := mocks.NewMockStrikeService(ctrl)
	handler := NewStrikeHandler(mockStrikeService, mockLogger)

	userID := uuid.New()

	strikeReq := dto.StrikeCreateRequest{
		TargetUserID: uuid.New(),
		Type:         constants.StrikeTypeHarassment,
		Reason:       "Test",
	}

	mockStrikeService.EXPECT().
		CreateStrike(gomock.Any(), gomock.Any()).
		Return(nil, errors.ErrSelfStrikeNotAllowed)

	body, _ := json.Marshal(strikeReq)
	req := httptest.NewRequest(http.MethodPost, "/api/strike", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req = req.WithContext(context.WithValue(req.Context(), middleware.UserIDKey, userID))
	rec := httptest.NewRecorder()

	handler.CreateStrike(rec, req)

	assert.Equal(t, http.StatusForbidden, rec.Code)
}

func TestStrikeHandler_CreateStrike_DuplicateStrike(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockLogger := mocks.NewMockLogger()
	mockStrikeService := mocks.NewMockStrikeService(ctrl)
	handler := NewStrikeHandler(mockStrikeService, mockLogger)

	userID := uuid.New()

	strikeReq := dto.StrikeCreateRequest{
		TargetUserID: uuid.New(),
		Type:         constants.StrikeTypeHarassment,
		Reason:       "Test",
	}

	mockStrikeService.EXPECT().
		CreateStrike(gomock.Any(), gomock.Any()).
		Return(nil, errors.ErrDuplicateStrike)

	body, _ := json.Marshal(strikeReq)
	req := httptest.NewRequest(http.MethodPost, "/api/strike", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req = req.WithContext(context.WithValue(req.Context(), middleware.UserIDKey, userID))
	rec := httptest.NewRecorder()

	handler.CreateStrike(rec, req)

	assert.Equal(t, http.StatusConflict, rec.Code)
}

func TestStrikeHandler_GetStrikesByUserID_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockLogger := mocks.NewMockLogger()
	mockStrikeService := mocks.NewMockStrikeService(ctrl)
	handler := NewStrikeHandler(mockStrikeService, mockLogger)

	userID := uuid.New()
	targetUserID := uuid.New().String()

	strikes := []*domain.Strike{
		{ID: uuid.New(), TargetUserID: uuid.MustParse(targetUserID)},
	}

	mockStrikeService.EXPECT().
		GetStrikesByUserID(gomock.Any(), targetUserID, 20, 0).
		Return(strikes, nil)

	req := httptest.NewRequest(http.MethodGet, "/api/strike/user/"+targetUserID, nil)
	req.SetPathValue("user_id", targetUserID)
	req = req.WithContext(context.WithValue(req.Context(), middleware.UserIDKey, userID))
	rec := httptest.NewRecorder()

	handler.GetStrikesByUserID(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestStrikeHandler_GetStrikesByUserID_UserNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockLogger := mocks.NewMockLogger()
	mockStrikeService := mocks.NewMockStrikeService(ctrl)
	handler := NewStrikeHandler(mockStrikeService, mockLogger)

	userID := uuid.New()
	targetUserID := uuid.New().String()

	mockStrikeService.EXPECT().
		GetStrikesByUserID(gomock.Any(), targetUserID, 20, 0).
		Return(nil, errors.ErrUserNotFound)

	req := httptest.NewRequest(http.MethodGet, "/api/strike/user/"+targetUserID, nil)
	req.SetPathValue("user_id", targetUserID)
	req = req.WithContext(context.WithValue(req.Context(), middleware.UserIDKey, userID))
	rec := httptest.NewRecorder()

	handler.GetStrikesByUserID(rec, req)

	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestStrikeHandler_GetStrikesByType_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockLogger := mocks.NewMockLogger()
	mockStrikeService := mocks.NewMockStrikeService(ctrl)
	handler := NewStrikeHandler(mockStrikeService, mockLogger)

	userID := uuid.New()
	strikeType := constants.StrikeTypeHarassment

	strikes := []*domain.Strike{
		{ID: uuid.New(), Type: strikeType},
	}

	mockStrikeService.EXPECT().
		GetStrikesByType(gomock.Any(), gomock.Any(), 20, 0).
		Return(strikes, nil)

	req := httptest.NewRequest(http.MethodGet, "/api/strike/type/"+string(strikeType), nil)
	req.SetPathValue("type", string(strikeType))
	req = req.WithContext(context.WithValue(req.Context(), middleware.UserIDKey, userID))
	rec := httptest.NewRecorder()

	handler.GetStrikesByType(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestStrikeHandler_DeleteStrike_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockLogger := mocks.NewMockLogger()
	mockStrikeService := mocks.NewMockStrikeService(ctrl)
	handler := NewStrikeHandler(mockStrikeService, mockLogger)

	userID := uuid.New()
	strikeID := uuid.New().String()

	mockStrikeService.EXPECT().
		DeleteStrike(gomock.Any(), strikeID).
		Return(nil)

	req := httptest.NewRequest(http.MethodDelete, "/api/strike/"+strikeID, nil)
	req.SetPathValue("id", strikeID)
	req = req.WithContext(context.WithValue(req.Context(), middleware.UserIDKey, userID))
	rec := httptest.NewRecorder()

	handler.DeleteStrike(rec, req)

	assert.Equal(t, http.StatusNoContent, rec.Code)
}

func TestStrikeHandler_DeleteStrike_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockLogger := mocks.NewMockLogger()
	mockStrikeService := mocks.NewMockStrikeService(ctrl)
	handler := NewStrikeHandler(mockStrikeService, mockLogger)

	userID := uuid.New()
	strikeID := uuid.New().String()

	mockStrikeService.EXPECT().
		DeleteStrike(gomock.Any(), strikeID).
		Return(errors.ErrStrikeNotFound)

	req := httptest.NewRequest(http.MethodDelete, "/api/strike/"+strikeID, nil)
	req.SetPathValue("id", strikeID)
	req = req.WithContext(context.WithValue(req.Context(), middleware.UserIDKey, userID))
	rec := httptest.NewRecorder()

	handler.DeleteStrike(rec, req)

	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestStrikeHandler_GetUserStrikeStats_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockLogger := mocks.NewMockLogger()
	mockStrikeService := mocks.NewMockStrikeService(ctrl)
	handler := NewStrikeHandler(mockStrikeService, mockLogger)

	userID := uuid.New()
	targetUserID := uuid.New().String()

	stats := &dto.StrikeStats{
		TotalStrikes: 5,
		StrikeTypes: dto.StrikeTypeStat{
			Pending: 2,
		},
	}

	mockStrikeService.EXPECT().
		GetUserStrikeStats(gomock.Any(), targetUserID).
		Return(stats, nil)

	req := httptest.NewRequest(http.MethodGet, "/api/strike/user/"+targetUserID+"/stat", nil)
	req.SetPathValue("user_id", targetUserID)
	req = req.WithContext(context.WithValue(req.Context(), middleware.UserIDKey, userID))
	rec := httptest.NewRecorder()

	handler.GetUserStrikeStats(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	var response dto.StrikeStats
	json.Unmarshal(rec.Body.Bytes(), &response)
	assert.Equal(t, 5, response.TotalStrikes)
}
