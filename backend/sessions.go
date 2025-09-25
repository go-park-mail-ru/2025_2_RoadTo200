// заглушки для сессий
package main

import (
	"net/http"
	"sync"
	"time"

	"github.com/google/uuid"
)

type Session struct {
	Token     string
	UserEmail string
	ExpiresAt time.Time
}

var (
	sessions   = make(map[string]*Session)
	sessionsMu sync.Mutex
)

func createSession(userEmail string, ttl time.Duration) *Session {
	token := uuid.NewString()
	s := &Session{
		Token:     token,
		UserEmail: userEmail,
		ExpiresAt: time.Now().Add(ttl),
	}
	sessionsMu.Lock()
	sessions[token] = s
	sessionsMu.Unlock()
	return s
}

func getSession(token string) (*Session, bool) {
	sessionsMu.Lock()
	defer sessionsMu.Unlock()
	s, ok := sessions[token]
	if !ok {
		return nil, false
	}
	if time.Now().After(s.ExpiresAt) {
		delete(sessions, token)
		return nil, false
	}
	return s, true
}

func deleteSession(token string) {
	sessionsMu.Lock()
	delete(sessions, token)
	sessionsMu.Unlock()
}

func setSessionCookie(w http.ResponseWriter, token string, ttl time.Duration) {
	cookie := &http.Cookie{
		Name:     "session_token",
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		MaxAge:   int(ttl.Seconds()),
	}
	http.SetCookie(w, cookie)
}

func clearSessionCookie(w http.ResponseWriter) {
	cookie := &http.Cookie{
		Name:     "session_token",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		MaxAge:   -1, // удалить
	}
	http.SetCookie(w, cookie)
}
