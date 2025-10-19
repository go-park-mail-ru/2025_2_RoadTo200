package repository

import (
	"sync"

	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/domain/entities"
	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/domain/errors"
)

type inMemoryUserRepository struct {
	users map[string]*domain.User
	mu    sync.RWMutex
}

func NewInMemoryUserRepository() UserRepository {
	return &inMemoryUserRepository{
		users: make(map[string]*domain.User),
	}
}

func (r *inMemoryUserRepository) Create(user *domain.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.users[user.Email]; exists {
		return errors.ErrUserAlreadyExists
	}

	r.users[user.Email] = user
	return nil
}

func (r *inMemoryUserRepository) FindByEmail(email string) (*domain.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	user, exists := r.users[email]
	if !exists {
		return nil, errors.ErrUserNotFound
	}

	return user, nil
}

func (r *inMemoryUserRepository) FindByID(id string) (*domain.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, user := range r.users {
		if user.ID == id {
			return user, nil
		}
	}

	return nil, errors.ErrUserNotFound
}

func (r *inMemoryUserRepository) Exists(email string) (bool, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	_, exists := r.users[email]
	return exists, nil
}
