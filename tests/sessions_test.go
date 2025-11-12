package main

// import (
// 	"testing"
// 	"time"
// )

// func TestCreateAndGetSession(t *testing.T) {
// 	email := "user@example.com"
// 	session := createSession(email, time.Hour)

// 	got, ok := getSession(session.Token)
// 	if !ok {
// 		t.Errorf("session not found")
// 	}
// 	if got.UserEmail != email {
// 		t.Errorf("expected email %s, got %s", email, got.UserEmail)
// 	}

// 	// тестируем истечение срока
// 	got.ExpiresAt = time.Now().Add(-time.Minute)
// 	if _, ok := getSession(session.Token); ok {
// 		t.Errorf("session should be expired")
// 	}
// }
