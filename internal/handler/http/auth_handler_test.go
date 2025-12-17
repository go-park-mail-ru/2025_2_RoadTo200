package handler

import (
	"bytes"
	"encoding/json"
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

func TestAuthHandler_Register_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockLogger := mocks.NewMockLogger()
	mockAuthService := mocks.NewMockAuthService(ctrl)
	handler := NewAuthHandler(mockAuthService, mockLogger)

	email := "test@example.com"
	password := "123456"

	user := &domain.User{
		ID:    uuid.MustParse("b5fafe25-4bd1-466f-9d66-16edab0ebcd3"),
		Email: email,
	}
	session := &domain.Session{
		Token: "session-token",
	}

	mockAuthService.EXPECT().
		Register(gomock.Any(), email, password, password).
		Return(user, session, nil)

	registerData := `{
		"email": "test@example.com",
		"password": "123456",
		"passwordConfirm": "123456"
	}`

	req := httptest.NewRequest("POST", "/api/register", bytes.NewReader([]byte(registerData)))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()

	handler.Register(rr, req)

	assert.Equal(t, http.StatusCreated, rr.Code)
}

func TestAuthHandler_Register_InvalidJSON(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockLogger := mocks.NewMockLogger()
	mockAuthService := mocks.NewMockAuthService(ctrl)
	handler := NewAuthHandler(mockAuthService, mockLogger)

	req := httptest.NewRequest("POST", "/api/register", bytes.NewReader([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()

	handler.Register(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)

	var response map[string]string
	json.Unmarshal(rr.Body.Bytes(), &response)
	assert.Equal(t, "invalid request body", response["error"])
}

func TestAuthHandler_Login_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockLogger := mocks.NewMockLogger()
	mockAuthService := mocks.NewMockAuthService(ctrl)
	handler := NewAuthHandler(mockAuthService, mockLogger)

	email := "test@example.com"
	password := "123456"

	user := &domain.User{
		ID:    uuid.New(),
		Email: email,
	}
	session := &domain.Session{
		Token: "session-token",
	}

	mockAuthService.EXPECT().
		Login(gomock.Any(), email, password).
		Return(user, session, nil)

	loginData := struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}{
		Email:    email,
		Password: password,
	}
	body, _ := json.Marshal(loginData)

	req := httptest.NewRequest("POST", "/api/login", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()

	handler.Login(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var response map[string]interface{}
	json.Unmarshal(rr.Body.Bytes(), &response)

	assert.Equal(t, user.ID.String(), response["id"])
	assert.Equal(t, user.Email, response["email"])

	// Проверяем что установлена cookie
	cookies := rr.Result().Cookies()
	var sessionCookie *http.Cookie
	for _, cookie := range cookies {
		if cookie.Name == "session_token" {
			sessionCookie = cookie
			break
		}
	}
	assert.NotNil(t, sessionCookie)
	assert.Equal(t, session.Token, sessionCookie.Value)
}

func TestAuthHandler_Login_InvalidCredentials(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockLogger := mocks.NewMockLogger()
	mockAuthService := mocks.NewMockAuthService(ctrl)
	handler := NewAuthHandler(mockAuthService, mockLogger)

	email := "test@example.com"
	password := "wrongpassword"

	mockAuthService.EXPECT().
		Login(gomock.Any(), email, password).
		Return(nil, nil, errors.ErrInvalidCredentials)

	loginData := struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}{
		Email:    email,
		Password: password,
	}
	body, _ := json.Marshal(loginData)

	req := httptest.NewRequest("POST", "/api/login", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()

	handler.Login(rr, req)

	assert.Equal(t, http.StatusUnauthorized, rr.Code)

	var response map[string]string
	json.Unmarshal(rr.Body.Bytes(), &response)
	assert.Contains(t, response["error"], "invalid")
}

func TestAuthHandler_Logout_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockLogger := mocks.NewMockLogger()
	mockAuthService := mocks.NewMockAuthService(ctrl)
	handler := NewAuthHandler(mockAuthService, mockLogger)

	token := "valid-session-token"

	mockAuthService.EXPECT().
		Logout(gomock.Any(), token).
		Return(nil)

	req := httptest.NewRequest("POST", "/api/logout", nil)
	req.AddCookie(&http.Cookie{
		Name:  "session_token",
		Value: token,
	})

	rr := httptest.NewRecorder()

	handler.Logout(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var response map[string]string
	json.Unmarshal(rr.Body.Bytes(), &response)
	assert.Equal(t, "logged out", response["message"])

	// Проверяем что cookie очищена
	cookies := rr.Result().Cookies()
	var sessionCookie *http.Cookie
	for _, cookie := range cookies {
		if cookie.Name == "session_token" {
			sessionCookie = cookie
			break
		}
	}
	assert.NotNil(t, sessionCookie)
	assert.Equal(t, "", sessionCookie.Value)
	assert.Equal(t, -1, sessionCookie.MaxAge)
}

func TestAuthHandler_Logout_NoSession(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockLogger := mocks.NewMockLogger()
	mockAuthService := mocks.NewMockAuthService(ctrl)
	handler := NewAuthHandler(mockAuthService, mockLogger)

	req := httptest.NewRequest("POST", "/api/logout", nil)
	rr := httptest.NewRecorder()

	handler.Logout(rr, req)

	assert.Equal(t, http.StatusUnauthorized, rr.Code)

	var response map[string]string
	json.Unmarshal(rr.Body.Bytes(), &response)
	assert.Equal(t, "no session", response["error"])
}
