package main

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func cleanup() {
	usersMu.Lock()
	users = make(map[string]*User)
	userSeq = 0
	usersMu.Unlock()

	sessionsMu.Lock()
	sessions = make(map[string]*Session)
	sessionsMu.Unlock()
}

func TestRegisterHandler(t *testing.T) {
	cleanup()
	t.Cleanup(cleanup)

	payload := []byte(`{"email":"alice@test.com","password":"123456","passwordConfirm":"123456"}`)
	req := httptest.NewRequest(http.MethodPost, "/api/register", bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	registerHandler(w, req)

	res := w.Result()
	if res.StatusCode != http.StatusCreated {
		t.Errorf("expected 201 Created, got %d", res.StatusCode)
	}
}

func TestLoginHandler(t *testing.T) {
	cleanup()
	t.Cleanup(cleanup)

	// сначала создаём пользователя
	_, err := createUser("bob@test.com", "mypassword")
	if err != nil {
		t.Fatalf("failed to create user: %v", err)
	}

	payload := []byte(`{"email":"bob@test.com","password":"mypassword"}`)
	req := httptest.NewRequest(http.MethodPost, "/api/login", bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	loginHandler(w, req)

	res := w.Result()
	if res.StatusCode != http.StatusCreated {
		t.Errorf("expected 201 Created, got %d", res.StatusCode)
	}

	// проверяем cookie
	cookies := res.Cookies()
	if len(cookies) == 0 || cookies[0].Name != "session_token" {
		t.Errorf("expected session_token cookie, got %v", cookies)
	}
}

func TestSessionHandler(t *testing.T) {
	cleanup()
	t.Cleanup(cleanup)

	// Case 1: No session cookie
	req := httptest.NewRequest(http.MethodGet, "/api/session", nil)
	w := httptest.NewRecorder()
	sessionHandler(w, req)
	res := w.Result()
	if res.StatusCode != http.StatusUnauthorized {
		t.Errorf("expected 401 Unauthorized, got %d", res.StatusCode)
	}

	// Case 2: Valid session
	user, _ := createUser("test@test.com", "password")
	session := createSession(user.Email, time.Hour)

	req = httptest.NewRequest(http.MethodGet, "/api/session", nil)
	req.AddCookie(&http.Cookie{Name: "session_token", Value: session.Token})
	w = httptest.NewRecorder()
	sessionHandler(w, req)
	res = w.Result()
	if res.StatusCode != http.StatusOK {
		t.Errorf("expected 200 OK, got %d", res.StatusCode)
	}
}

func TestLogoutHandler(t *testing.T) {
	cleanup()
	t.Cleanup(cleanup)

	user, _ := createUser("test@test.com", "password")
	session := createSession(user.Email, time.Hour)

	req := httptest.NewRequest(http.MethodPost, "/api/logout", nil)
	req.AddCookie(&http.Cookie{Name: "session_token", Value: session.Token})

	w := httptest.NewRecorder()
	logoutHandler(w, req)

	res := w.Result()
	if res.StatusCode != http.StatusOK {
		t.Errorf("expected 200 OK, got %d", res.StatusCode)
	}

	// Check if session is deleted
	_, ok := getSession(session.Token)
	if ok {
		t.Errorf("session should be deleted")
	}

	// Check for cookie expiration
	cookies := res.Cookies()
	if len(cookies) == 0 || cookies[0].MaxAge != -1 {
		t.Errorf("expected cookie to be expired, got %v", cookies)
	}
}

func TestFeedHandler(t *testing.T) {
	cleanup()
	t.Cleanup(cleanup)

	// Case 1: Unauthorized
	req := httptest.NewRequest(http.MethodGet, "/api/feed", nil)
	w := httptest.NewRecorder()
	feedHandler(w, req)
	res := w.Result()
	if res.StatusCode != http.StatusUnauthorized {
		t.Errorf("expected 401 Unauthorized, got %d", res.StatusCode)
	}

	// Case 2: Authorized
	user, _ := createUser("test@test.com", "password")
	session := createSession(user.Email, time.Hour)
	req = httptest.NewRequest(http.MethodGet, "/api/feed", nil)
	req.AddCookie(&http.Cookie{Name: "session_token", Value: session.Token})
	w = httptest.NewRecorder()
	feedHandler(w, req)
	res = w.Result()
	if res.StatusCode != http.StatusOK {
		t.Errorf("expected 200 OK, got %d", res.StatusCode)
	}
}

func TestSwipeHandler(t *testing.T) {
	cleanup()
	t.Cleanup(cleanup)

	// Case 1: Unauthorized
	req := httptest.NewRequest(http.MethodPost, "/api/swipe", nil)
	w := httptest.NewRecorder()
	swipeHandler(w, req)
	res := w.Result()
	if res.StatusCode != http.StatusUnauthorized {
		t.Errorf("expected 401 Unauthorized, got %d", res.StatusCode)
	}

	// Case 2: Authorized
	user, _ := createUser("test@test.com", "password")
	session := createSession(user.Email, time.Hour)
	payload := []byte(`{"target_id":"2","direction":"right"}`)
	req = httptest.NewRequest(http.MethodPost, "/api/swipe", bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")
	req.AddCookie(&http.Cookie{Name: "session_token", Value: session.Token})
	w = httptest.NewRecorder()
	swipeHandler(w, req)
	res = w.Result()
	if res.StatusCode != http.StatusOK {
		t.Errorf("expected 200 OK, got %d", res.StatusCode)
	}
}
