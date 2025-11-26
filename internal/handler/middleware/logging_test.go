package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/tests/mocks"
	"github.com/stretchr/testify/assert"
)

func TestLogMiddleware_LogsRequest(t *testing.T) {
	mockLogger := mocks.NewMockLogger()

	middleware := LogMiddleware(mockLogger)
	handler := middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}))

	req := httptest.NewRequest(http.MethodGet, "/api/test?param=value", nil)
	req.Header.Set("User-Agent", "TestAgent")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "OK", rec.Body.String())
}

func TestLogMiddleware_LogsPostRequest(t *testing.T) {
	mockLogger := mocks.NewMockLogger()

	middleware := LogMiddleware(mockLogger)
	handler := middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusCreated)
	}))

	req := httptest.NewRequest(http.MethodPost, "/api/users", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusCreated, rec.Code)
}

func TestLogMiddleware_CallsNextHandler(t *testing.T) {
	mockLogger := mocks.NewMockLogger()

	middleware := LogMiddleware(mockLogger)

	handlerCalled := false
	handler := middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handlerCalled = true
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodGet, "/api/test", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.True(t, handlerCalled, "Next handler should be called")
	assert.Equal(t, http.StatusOK, rec.Code)
}
