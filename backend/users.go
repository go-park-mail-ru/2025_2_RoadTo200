// заглушки для пользователей
package main

import (
	"errors"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"sync"
)

type User struct {
	ID       string
	Email    string
	Password string // хэш
}

var (
	users   = make(map[string]*User) // email -> user
	usersMu sync.Mutex
	userSeq int
)

func createUser(email, password string) (*User, error) {
	usersMu.Lock()
	defer usersMu.Unlock()

	if _, exists := users[email]; exists {
		return nil, errors.New("user already exists")
	}

	userSeq++
	id := fmt.Sprintf("%d", userSeq)

	// хэшируем пароль
	hashed, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	u := &User{
		ID:       id,
		Email:    email,
		Password: string(hashed), // сохраняем хэш
	}
	users[email] = u
	return u, nil
}

func findUser(email string) (*User, bool) {
	usersMu.Lock()
	defer usersMu.Unlock()
	u, ok := users[email]
	return u, ok
}
