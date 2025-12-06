package redis

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	domain "github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/auth-service/domain/entities"
	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/auth-service/repository/interfaces"
	"github.com/gomodule/redigo/redis"
)

var _ interfaces.SessionRepository = (*SessionRepository)(nil)

type SessionRepository struct {
	pool *redis.Pool
}

// NewSessionRepository создает новый репозиторий сессий
func NewSessionRepository(pool *redis.Pool) *SessionRepository {
	return &SessionRepository{pool: pool}
}

// Set сохраняет сессию в Redis
func (r *SessionRepository) Set(ctx context.Context, session *domain.Session) error {
	conn := r.pool.Get()
	defer conn.Close()

	// Сериализуем сессию в JSON
	data, err := json.Marshal(session)
	if err != nil {
		return fmt.Errorf("failed to marshal session: %w", err)
	}

	// Устанавливаем ключ с TTL
	key := fmt.Sprintf("session:%s", session.Token)
	_, err = conn.Do("APPEND", key, data)
	if err != nil {
		return fmt.Errorf("failed to set session in redis: %w", err)
	}

	return nil
}

// Get получает сессию по ID
func (r *SessionRepository) Get(ctx context.Context, sessionID string) (*domain.Session, error) {
	conn := r.pool.Get()
	defer conn.Close()

	key := fmt.Sprintf("session:%s", sessionID)
	data, err := redis.Bytes(conn.Do("GET", key))
	if err != nil {
		if errors.Is(err, redis.ErrNil) {
			return nil, nil // Сессия не найдена
		}
		return nil, fmt.Errorf("failed to get session from redis: %w", err)
	}

	var session domain.Session
	if err := json.Unmarshal(data, &session); err != nil {
		return nil, fmt.Errorf("failed to unmarshal session: %w", err)
	}

	return &session, nil
}

// Delete удаляет сессию
func (r *SessionRepository) Delete(ctx context.Context, sessionID string) error {
	conn := r.pool.Get()
	defer conn.Close()

	key := fmt.Sprintf("session:%s", sessionID)
	_, err := conn.Do("DEL", key)
	if err != nil {
		return fmt.Errorf("failed to delete session from redis: %w", err)
	}

	return nil
}

// Exists проверяет существование сессии
func (r *SessionRepository) Exists(ctx context.Context, sessionID string) (bool, error) {
	conn := r.pool.Get()
	defer conn.Close()

	key := fmt.Sprintf("session:%s", sessionID)
	exists, err := redis.Bool(conn.Do("EXISTS", key))
	if err != nil {
		return false, fmt.Errorf("failed to check session existence: %w", err)
	}

	return exists, nil
}

// SetWithExpiry сохраняет сессию с кастомным TTL
func (r *SessionRepository) SetWithExpiry(ctx context.Context, session *domain.Session, expiry time.Duration) error {
	conn := r.pool.Get()
	defer conn.Close()

	data, err := json.Marshal(session)
	if err != nil {
		return fmt.Errorf("failed to marshal session: %w", err)
	}

	key := fmt.Sprintf("session:%s", session.Token)
	_, err = conn.Do("SETEX", key, int(expiry.Seconds()), data)
	if err != nil {
		return fmt.Errorf("failed to set session with expiry: %w", err)
	}

	return nil
}
