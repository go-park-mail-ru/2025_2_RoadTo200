package domain

import "time"

type User struct {
	ID        string
	Email     string
	Password  string // хэш
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewUser(email, passwordHash string) *User {
	return &User{
		Email:    email,
		Password: passwordHash,
	}
}
