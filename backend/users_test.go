package main

import (
	"testing"
)

func TestCreateAndFindUser(t *testing.T) {
	email := "test@example.com"
	password := "secret123"

	u, _ := createUser(email, password)
	if u.Email != email {
		t.Errorf("expected email %s, got %s", email, u.Email)
	}

	found, ok := findUser(email)
	if !ok {
		t.Errorf("user not found")
	}
	if found.ID != u.ID {
		t.Errorf("expected ID %s, got %s", u.ID, found.ID)
	}
}
