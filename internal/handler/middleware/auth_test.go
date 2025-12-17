package middleware

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	domain "github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/domain/entities"
	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/domain/errors"
	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/tests/mocks"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestAuthMiddleware_PublicEndpoint(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAuthService := mocks.NewMockAuthService(ctrl)

	middleware := AuthMiddleware(mockAuthService)

	handler := middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}))

	tests := []struct {
		name string
		path string
	}{
		{"Register endpoint", "/api/register"},
		{"Login endpoint", "/api/login"},
		{"Session endpoint", "/api/session"},
		{"Health endpoint", "/health"},
		{"Swagger index", "/swagger/index.html"},
		{"Swagger doc", "/swagger/doc.json"},
		{"Swagger path", "/swagger/something.html"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, tt.path, nil)
			rec := httptest.NewRecorder()

			handler.ServeHTTP(rec, req)

			assert.Equal(t, http.StatusOK, rec.Code)
			assert.Equal(t, "OK", rec.Body.String())
		})
	}
}

func TestAuthMiddleware_ValidToken_FromHeader(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAuthService := mocks.NewMockAuthService(ctrl)

	userID := uuid.New()
	user := &domain.User{
		ID:   userID,
		Name: "Test User",
	}

	token := "valid-token-123"

	mockAuthService.EXPECT().
		ValidateSession(gomock.Any(), token).
		Return(user, nil)

	middleware := AuthMiddleware(mockAuthService)

	var capturedUserID uuid.UUID
	handler := middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		uid, err := GetUserIDFromContext(r.Context())
		assert.NoError(t, err)
		capturedUserID = uid
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodGet, "/api/protected", nil)
	req.Header.Set("X-Session-Token", token)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, userID, capturedUserID)
}

func TestAuthMiddleware_ValidToken_FromCookie(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAuthService := mocks.NewMockAuthService(ctrl)

	userID := uuid.New()
	user := &domain.User{
		ID:   userID,
		Name: "Test User",
	}

	token := "valid-cookie-token"

	mockAuthService.EXPECT().
		ValidateSession(gomock.Any(), token).
		Return(user, nil)

	middleware := AuthMiddleware(mockAuthService)

	var capturedUserID uuid.UUID
	handler := middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		uid, err := GetUserIDFromContext(r.Context())
		assert.NoError(t, err)
		capturedUserID = uid
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodGet, "/api/protected", nil)
	req.AddCookie(&http.Cookie{
		Name:  "session_token",
		Value: token,
	})
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, userID, capturedUserID)
}

func TestAuthMiddleware_NoToken(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAuthService := mocks.NewMockAuthService(ctrl)

	middleware := AuthMiddleware(mockAuthService)

	handler := middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("Handler should not be called")
	}))

	req := httptest.NewRequest(http.MethodGet, "/api/protected", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusUnauthorized, rec.Code)
}

func TestAuthMiddleware_InvalidToken(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAuthService := mocks.NewMockAuthService(ctrl)

	token := "invalid-token"

	mockAuthService.EXPECT().
		ValidateSession(gomock.Any(), token).
		Return(nil, errors.ErrSessionNotFound)

	middleware := AuthMiddleware(mockAuthService)

	handler := middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("Handler should not be called")
	}))

	req := httptest.NewRequest(http.MethodGet, "/api/protected", nil)
	req.Header.Set("X-Session-Token", token)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusUnauthorized, rec.Code)
}

func TestAuthMiddleware_HeaderTakesPrecedenceOverCookie(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAuthService := mocks.NewMockAuthService(ctrl)

	userID := uuid.New()
	user := &domain.User{ID: userID}

	headerToken := "header-token"
	cookieToken := "cookie-token"

	// Should use header token, not cookie
	mockAuthService.EXPECT().
		ValidateSession(gomock.Any(), headerToken).
		Return(user, nil)

	middleware := AuthMiddleware(mockAuthService)

	handler := middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodGet, "/api/protected", nil)
	req.Header.Set("X-Session-Token", headerToken)
	req.AddCookie(&http.Cookie{
		Name:  "session_token",
		Value: cookieToken,
	})
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestGetUserIDFromContext_Success(t *testing.T) {
	userID := uuid.New()
	ctx := context.WithValue(context.Background(), UserIDKey, userID)

	extractedID, err := GetUserIDFromContext(ctx)

	assert.NoError(t, err)
	assert.Equal(t, userID, extractedID)
}

func TestGetUserIDFromContext_NotFound(t *testing.T) {
	ctx := context.Background()

	extractedID, err := GetUserIDFromContext(ctx)

	assert.Error(t, err)
	assert.Equal(t, uuid.Nil, extractedID)
	assert.Contains(t, err.Error(), "userID not found in context")
}

func TestIsPublicEndpoint(t *testing.T) {
	tests := []struct {
		name     string
		path     string
		expected bool
	}{
		{"Register", "/api/register", true},
		{"Login", "/api/login", true},
		{"Health", "/health", true},
		{"Swagger prefix", "/swagger/something", true},
		{"Protected endpoint", "/api/users", false},
		{"Protected feed", "/api/feed", false},
		{"Root", "/", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isPublicEndpoint(tt.path)
			assert.Equal(t, tt.expected, result)
		})
	}
}
