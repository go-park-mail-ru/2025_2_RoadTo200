package domain

import "time"

type Session struct {
	Token     string
	UserEmail string
	ExpiresAt time.Time
}

func NewSession(userEmail string, ttl time.Duration) *Session {
	return &Session{
		UserEmail: userEmail,
		ExpiresAt: time.Now().Add(ttl),
	}
}

func (s *Session) IsExpired() bool {
	return time.Now().After(s.ExpiresAt)
}
