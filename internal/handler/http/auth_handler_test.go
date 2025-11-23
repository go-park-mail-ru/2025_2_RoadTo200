package handler

// import (
// 	"bytes"
// 	"encoding/json"
// 	"net/http"
// 	"net/http/httptest"
// 	"testing"

// 	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/domain/entities"
// 	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/domain/errors"
// 	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/tests/mocks"
// 	"github.com/google/uuid"
// 	"github.com/stretchr/testify/assert"
// 	"go.uber.org/mock/gomock"
// )

// func TestAuthHandler_Register_Success(t *testing.T) {
// 	ctrl := gomock.NewController(t)
// 	defer ctrl.Finish()

// 	mockLogger := mocks.NewMockLogger()
// 	mockAuthService := mocks.NewMockAuthService(ctrl)
// 	handler := NewAuthHandler(mockAuthService, mockLogger)

// 	email := "test@example.com"
// 	password := "123456"

// 	user := &domain.User{
// 		ID:    uuid.MustParse("b5fafe25-4bd1-466f-9d66-16edab0ebcd3"),
// 		Email: email,
// 	}
// 	session := &domain.Session{
// 		Token: "session-token",
// 	}

// 	mockAuthService.EXPECT().
// 		Register(email, password, password).
// 		Return(user, session, nil)

// 	// ПРОСТОЙ JSON без лишних полей
// 	registerData := `{
// 		"email": "test@example.com",
// 		"password": "123456",
// 		"passwordConfirm": "123456"
// 	}`

// 	req := httptest.NewRequest("POST", "/api/register", bytes.NewReader([]byte(registerData)))
// 	req.Header.Set("Content-Type", "application/json")

// 	rr := httptest.NewRecorder()

// 	handler.Register(rr, req)

// 	t.Logf("Status: %d", rr.Code)
// 	t.Logf("Body: %s", rr.Body.String())

// 	// Пока проверяем только что не 400
// 	if rr.Code == 400 {
// 		t.Fatalf("Expected not 400, got 400. Response: %s", rr.Body.String())
// 	}
// }

// func TestAuthHandler_Register_InvalidJSON(t *testing.T) {
// 	ctrl := gomock.NewController(t)
// 	defer ctrl.Finish()

// 	mockAuthService := mocks.NewMockAuthService(ctrl)
// 	handler := NewAuthHandler(mockAuthService)

// 	// Создаем запрос с невалидным JSON
// 	req := httptest.NewRequest("POST", "/api/register", bytes.NewReader([]byte("invalid json")))
// 	req.Header.Set("Content-Type", "application/json")

// 	rr := httptest.NewRecorder()

// 	// Вызов хендлера
// 	handler.Register(rr, req)

// 	// Проверки
// 	assert.Equal(t, http.StatusBadRequest, rr.Code)

// 	var response map[string]string
// 	json.Unmarshal(rr.Body.Bytes(), &response)
// 	assert.Equal(t, "invalid request body", response["error"])
// }

// func TestAuthHandler_Login_Success(t *testing.T) {
// 	ctrl := gomock.NewController(t)
// 	defer ctrl.Finish()

// 	mockAuthService := mocks.NewMockAuthService(ctrl)
// 	handler := NewAuthHandler(mockAuthService)

// 	email := "test@example.com"
// 	password := "123456"

// 	user := &domain.User{
// 		ID:    uuid.New(),
// 		Email: email,
// 	}
// 	session := &domain.Session{
// 		Token: "session-token",
// 	}

// 	// Настройка ожиданий
// 	mockAuthService.EXPECT().
// 		Login(email, password).
// 		Return(user, session, nil)

// 	// Создаем запрос
// 	loginData := struct {
// 		Email    string `json:"email"`
// 		Password string `json:"password"`
// 	}{
// 		Email:    email,
// 		Password: password,
// 	}
// 	body, _ := json.Marshal(loginData)

// 	req := httptest.NewRequest("POST", "/api/login", bytes.NewReader(body))
// 	req.Header.Set("Content-Type", "application/json")

// 	rr := httptest.NewRecorder()

// 	// Вызов хендлера
// 	handler.Login(rr, req)

// 	// Проверки
// 	assert.Equal(t, http.StatusOK, rr.Code)

// 	var response map[string]interface{}
// 	json.Unmarshal(rr.Body.Bytes(), &response)

// 	assert.Equal(t, user.ID.String(), response["id"])
// 	assert.Equal(t, user.Email, response["email"])

// 	// Проверяем что установлена cookie
// 	cookies := rr.Result().Cookies()
// 	var sessionCookie *http.Cookie
// 	for _, cookie := range cookies {
// 		if cookie.Name == "session_token" {
// 			sessionCookie = cookie
// 			break
// 		}
// 	}
// 	assert.NotNil(t, sessionCookie)
// 	assert.Equal(t, session.Token, sessionCookie.Value)
// }

// func TestAuthHandler_Login_InvalidCredentials(t *testing.T) {
// 	ctrl := gomock.NewController(t)
// 	defer ctrl.Finish()

// 	mockAuthService := mocks.NewMockAuthService(ctrl)
// 	handler := NewAuthHandler(mockAuthService)

// 	email := "test@example.com"
// 	password := "wrongpassword"

// 	// Настройка ожиданий
// 	mockAuthService.EXPECT().
// 		Login(email, password).
// 		Return(nil, nil, errors.ErrInvalidCredentials)

// 	// Создаем запрос
// 	loginData := struct {
// 		Email    string `json:"email"`
// 		Password string `json:"password"`
// 	}{
// 		Email:    email,
// 		Password: password,
// 	}
// 	body, _ := json.Marshal(loginData)

// 	req := httptest.NewRequest("POST", "/api/login", bytes.NewReader(body))
// 	req.Header.Set("Content-Type", "application/json")

// 	rr := httptest.NewRecorder()

// 	// Вызов хендлера
// 	handler.Login(rr, req)

// 	// Проверки
// 	assert.Equal(t, http.StatusUnauthorized, rr.Code)

// 	var response map[string]string
// 	json.Unmarshal(rr.Body.Bytes(), &response)
// 	assert.Equal(t, "invalid email or password", response["error"])
// }

// func TestAuthHandler_Logout_Success(t *testing.T) {
// 	ctrl := gomock.NewController(t)
// 	defer ctrl.Finish()

// 	mockAuthService := mocks.NewMockAuthService(ctrl)
// 	handler := NewAuthHandler(mockAuthService)

// 	token := "valid-session-token"

// 	// Настройка ожиданий
// 	mockAuthService.EXPECT().
// 		Logout(token).
// 		Return(nil)

// 	// Создаем запрос с cookie
// 	req := httptest.NewRequest("POST", "/api/logout", nil)
// 	req.AddCookie(&http.Cookie{
// 		Name:  "session_token",
// 		Value: token,
// 	})

// 	rr := httptest.NewRecorder()

// 	// Вызов хендлера
// 	handler.Logout(rr, req)

// 	// Проверки
// 	assert.Equal(t, http.StatusOK, rr.Code)

// 	var response map[string]string
// 	json.Unmarshal(rr.Body.Bytes(), &response)
// 	assert.Equal(t, "logged out", response["message"])

// 	// Проверяем что cookie очищена
// 	cookies := rr.Result().Cookies()
// 	var sessionCookie *http.Cookie
// 	for _, cookie := range cookies {
// 		if cookie.Name == "session_token" {
// 			sessionCookie = cookie
// 			break
// 		}
// 	}
// 	assert.NotNil(t, sessionCookie)
// 	assert.Equal(t, "", sessionCookie.Value)
// 	assert.Equal(t, -1, sessionCookie.MaxAge)
// }

// func TestAuthHandler_Logout_NoSession(t *testing.T) {
// 	ctrl := gomock.NewController(t)
// 	defer ctrl.Finish()

// 	mockAuthService := mocks.NewMockAuthService(ctrl)
// 	handler := NewAuthHandler(mockAuthService)

// 	// Создаем запрос БЕЗ cookie
// 	req := httptest.NewRequest("POST", "/api/logout", nil)
// 	rr := httptest.NewRecorder()

// 	// Вызов хендлера
// 	handler.Logout(rr, req)

// 	// Проверки
// 	assert.Equal(t, http.StatusUnauthorized, rr.Code)

// 	var response map[string]string
// 	json.Unmarshal(rr.Body.Bytes(), &response)
// 	assert.Equal(t, "no session", response["error"])
// }
