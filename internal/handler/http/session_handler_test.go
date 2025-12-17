package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	domain "github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/domain/entities"
	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/domain/errors"
	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/handler/dto"
	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/tests/mocks"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestSessionHandler_GetSession_ValidToken_Header(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockLogger := mocks.NewMockLogger()
	mockAuthService := mocks.NewMockAuthService(ctrl)
	handler := NewSessionHandler(mockAuthService, mockLogger)

	token := "valid-token-123"
	userID := uuid.New()

	user := &domain.User{
		ID:    userID,
		Email: "test@example.com",
		Name:  "Test User",
	}

	mockAuthService.EXPECT().
		ValidateSession(gomock.Any(), token).
		Return(user, nil)

	req := httptest.NewRequest(http.MethodGet, "/api/session", nil)
	req.Header.Set("X-Session-Token", token)
	rec := httptest.NewRecorder()

	handler.GetSession(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	var response dto.SessionResponse
	json.Unmarshal(rec.Body.Bytes(), &response)

	assert.True(t, response.Authenticated)
	assert.NotNil(t, response.User)
	assert.Equal(t, user.ID.String(), response.User.ID)
	assert.Equal(t, user.Email, response.User.Email)
	assert.Equal(t, user.Name, response.User.Name)
}

func TestSessionHandler_GetSession_ValidToken_Cookie(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockLogger := mocks.NewMockLogger()
	mockAuthService := mocks.NewMockAuthService(ctrl)
	handler := NewSessionHandler(mockAuthService, mockLogger)

	token := "cookie-token-456"
	user := &domain.User{
		ID:    uuid.New(),
		Email: "test@example.com",
		Name:  "Test User",
	}

	mockAuthService.EXPECT().
		ValidateSession(gomock.Any(), token).
		Return(user, nil)

	req := httptest.NewRequest(http.MethodGet, "/api/session", nil)
	req.AddCookie(&http.Cookie{
		Name:  "session_token",
		Value: token,
	})
	rec := httptest.NewRecorder()

	handler.GetSession(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	var response dto.SessionResponse
	json.Unmarshal(rec.Body.Bytes(), &response)

	assert.True(t, response.Authenticated)
	assert.NotNil(t, response.User)
}

func TestSessionHandler_GetSession_NoToken(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockLogger := mocks.NewMockLogger()
	mockAuthService := mocks.NewMockAuthService(ctrl)
	handler := NewSessionHandler(mockAuthService, mockLogger)

	req := httptest.NewRequest(http.MethodGet, "/api/session", nil)
	rec := httptest.NewRecorder()

	handler.GetSession(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	var response dto.SessionResponse
	json.Unmarshal(rec.Body.Bytes(), &response)

	assert.False(t, response.Authenticated)
	assert.Nil(t, response.User)
}

func TestSessionHandler_GetSession_InvalidToken(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockLogger := mocks.NewMockLogger()
	mockAuthService := mocks.NewMockAuthService(ctrl)
	handler := NewSessionHandler(mockAuthService, mockLogger)

	token := "invalid-token"

	mockAuthService.EXPECT().
		ValidateSession(gomock.Any(), token).
		Return(nil, errors.ErrSessionNotFound)

	req := httptest.NewRequest(http.MethodGet, "/api/session", nil)
	req.Header.Set("X-Session-Token", token)
	rec := httptest.NewRecorder()

	handler.GetSession(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	var response dto.SessionResponse
	json.Unmarshal(rec.Body.Bytes(), &response)

	assert.False(t, response.Authenticated)
	assert.Nil(t, response.User)
}

func TestSessionHandler_GetSession_HeaderTakesPrecedence(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockLogger := mocks.NewMockLogger()
	mockAuthService := mocks.NewMockAuthService(ctrl)
	handler := NewSessionHandler(mockAuthService, mockLogger)

	headerToken := "header-token"
	cookieToken := "cookie-token"

	user := &domain.User{
		ID:    uuid.New(),
		Email: "test@example.com",
		Name:  "Test User",
	}

	// Should use header token, not cookie
	mockAuthService.EXPECT().
		ValidateSession(gomock.Any(), headerToken).
		Return(user, nil)

	req := httptest.NewRequest(http.MethodGet, "/api/session", nil)
	req.Header.Set("X-Session-Token", headerToken)
	req.AddCookie(&http.Cookie{
		Name:  "session_token",
		Value: cookieToken,
	})
	rec := httptest.NewRecorder()

	handler.GetSession(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	var response dto.SessionResponse
	json.Unmarshal(rec.Body.Bytes(), &response)

	assert.True(t, response.Authenticated)
}
