package repository

import (
	"sync"
	//"time"

	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/domain/entities"
	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/domain/errors"
)

type inMemorySessionRepository struct {
	sessions map[string]*domain.Session
	mu       sync.RWMutex
}

func NewInMemorySessionRepository() SessionRepository {
	return &inMemorySessionRepository{
		sessions: make(map[string]*domain.Session),
	}
}

func (r *inMemorySessionRepository) Create(session *domain.Session) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.sessions[session.Token] = session
	return nil
}

func (r *inMemorySessionRepository) FindByToken(token string) (*domain.Session, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	session, exists := r.sessions[token]
	if !exists {
		return nil, errors.ErrSessionNotFound
	}

	if session.IsExpired() {
		go r.Delete(token)
		return nil, errors.ErrSessionExpired
	}

	return session, nil
}

func (r *inMemorySessionRepository) Delete(token string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	delete(r.sessions, token)
	return nil
}

func (r *inMemorySessionRepository) DeleteExpired() error {
	r.mu.Lock()
	defer r.mu.Unlock()

	for token, session := range r.sessions {
		if session.IsExpired() {
			delete(r.sessions, token)
		}
	}

	return nil
}
