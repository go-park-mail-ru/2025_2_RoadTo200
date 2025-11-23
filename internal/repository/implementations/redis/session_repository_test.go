package redis

import (
	"encoding/json"
	"errors"
	"testing"
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/rafaeljusto/redigomock"
	"github.com/stretchr/testify/assert"

	domain "github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/domain/entities"
)

func TestSessionRepository_Set(t *testing.T) {
	tests := []struct {
		name        string
		session     *domain.Session
		setupMock   func(conn *redigomock.Conn)
		expectedErr string
	}{
		{
			name: "Success",
			session: &domain.Session{
				Token:     "test-token",
				UserEmail: "test@example.com",
				ExpiresAt: time.Now().Add(time.Hour),
			},
			setupMock: func(conn *redigomock.Conn) {
				// Используем AnyData для JSON данных, так как время может немного отличаться
				conn.Command("APPEND", "session:test-token", redigomock.NewAnyData()).Expect("OK")
			},
		},
		{
			name: "EmptySession",
			session: &domain.Session{
				Token: "test-token",
			},
			setupMock: func(conn *redigomock.Conn) {
				conn.Command("APPEND", "session:test-token", redigomock.NewAnyData()).Expect("OK")
			},
		},
		{
			name: "RedisError",
			session: &domain.Session{
				Token:     "test-token",
				UserEmail: "test@example.com",
				ExpiresAt: time.Now().Add(time.Hour),
			},
			setupMock: func(conn *redigomock.Conn) {
				conn.Command("APPEND", "session:test-token", redigomock.NewAnyData()).
					ExpectError(errors.New("redis error"))
			},
			expectedErr: "failed to set session in redis",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			conn := redigomock.NewConn()
			tt.setupMock(conn)

			pool := &redis.Pool{
				Dial:         func() (redis.Conn, error) { return conn, nil },
				TestOnBorrow: func(c redis.Conn, t time.Time) error { return nil },
			}
			defer pool.Close()

			repo := NewSessionRepository(pool)

			err := repo.Set(tt.session)

			if tt.expectedErr != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedErr)
			} else {
				assert.NoError(t, err)
			}

			// Проверяем что команда была вызвана
			if tt.expectedErr == "" {
				cmd := conn.Command("APPEND", "session:"+tt.session.Token, redigomock.NewAnyData())
				assert.Equal(t, 1, conn.Stats(cmd), "APPEND command should be called once")
			}
		})
	}
}

func TestSessionRepository_SetWithExpiry(t *testing.T) {
	mockSession := &domain.Session{
		Token:     "test-token",
		UserEmail: "test@example.com",
		ExpiresAt: time.Now().Add(time.Hour),
	}
	expiry := 10 * time.Minute

	tests := []struct {
		name        string
		session     *domain.Session
		expiry      time.Duration
		setupMock   func(conn *redigomock.Conn)
		expectedErr string
	}{
		{
			name:    "Success",
			session: mockSession,
			expiry:  expiry,
			setupMock: func(conn *redigomock.Conn) {
				conn.Command("SETEX", "session:test-token", int(expiry.Seconds()), redigomock.NewAnyData()).
					Expect("OK")
			},
		},
		{
			name:    "EmptySession",
			session: &domain.Session{Token: "test-token"},
			expiry:  expiry,
			setupMock: func(conn *redigomock.Conn) {
				conn.Command("SETEX", "session:test-token", int(expiry.Seconds()), redigomock.NewAnyData()).
					Expect("OK")
			},
		},
		{
			name:    "RedisError",
			session: mockSession,
			expiry:  expiry,
			setupMock: func(conn *redigomock.Conn) {
				conn.Command("SETEX", "session:test-token", int(expiry.Seconds()), redigomock.NewAnyData()).
					ExpectError(errors.New("redis error"))
			},
			expectedErr: "failed to set session with expiry",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			conn := redigomock.NewConn()
			tt.setupMock(conn)

			pool := &redis.Pool{
				Dial:         func() (redis.Conn, error) { return conn, nil },
				TestOnBorrow: func(c redis.Conn, t time.Time) error { return nil },
			}
			defer pool.Close()

			repo := NewSessionRepository(pool)

			err := repo.SetWithExpiry(tt.session, tt.expiry)

			if tt.expectedErr != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedErr)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// Остальные тесты остаются без изменений...
func TestSessionRepository_Get(t *testing.T) {
	mockSession := &domain.Session{
		Token:     "test-token",
		UserEmail: "test@example.com",
		ExpiresAt: time.Now().Add(time.Hour),
	}

	tests := []struct {
		name        string
		sessionID   string
		setupMock   func(conn *redigomock.Conn)
		expected    *domain.Session
		expectedErr string
	}{
		{
			name:      "Success",
			sessionID: "test-token",
			setupMock: func(conn *redigomock.Conn) {
				data, _ := json.Marshal(mockSession)
				conn.Command("GET", "session:test-token").Expect(data)
			},
			expected: mockSession,
		},
		{
			name:      "NotFound",
			sessionID: "non-existent",
			setupMock: func(conn *redigomock.Conn) {
				conn.Command("GET", "session:non-existent").Expect(nil)
			},
			expected: nil,
		},
		{
			name:      "RedisError",
			sessionID: "test-token",
			setupMock: func(conn *redigomock.Conn) {
				conn.Command("GET", "session:test-token").
					ExpectError(errors.New("redis error"))
			},
			expectedErr: "failed to get session from redis",
		},
		{
			name:      "UnmarshalError",
			sessionID: "test-token",
			setupMock: func(conn *redigomock.Conn) {
				conn.Command("GET", "session:test-token").Expect([]byte("invalid json"))
			},
			expectedErr: "failed to unmarshal session",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			conn := redigomock.NewConn()
			tt.setupMock(conn)

			pool := &redis.Pool{
				Dial:         func() (redis.Conn, error) { return conn, nil },
				TestOnBorrow: func(c redis.Conn, t time.Time) error { return nil },
			}
			defer pool.Close()

			repo := NewSessionRepository(pool)

			session, err := repo.Get(tt.sessionID)

			if tt.expectedErr != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedErr)
				assert.Nil(t, session)
			} else {
				assert.NoError(t, err)
				if tt.expected == nil {
					assert.Nil(t, session)
				} else {
					assert.Equal(t, tt.expected.UserEmail, session.UserEmail)
					assert.Equal(t, tt.expected.Token, session.Token)
					assert.WithinDuration(t, tt.expected.ExpiresAt, session.ExpiresAt, time.Second)
				}
			}
		})
	}
}

func TestSessionRepository_Delete(t *testing.T) {
	tests := []struct {
		name        string
		sessionID   string
		setupMock   func(conn *redigomock.Conn)
		expectedErr string
	}{
		{
			name:      "Success",
			sessionID: "test-token",
			setupMock: func(conn *redigomock.Conn) {
				conn.Command("DEL", "session:test-token").Expect(int64(1))
			},
		},
		{
			name:      "RedisError",
			sessionID: "test-token",
			setupMock: func(conn *redigomock.Conn) {
				conn.Command("DEL", "session:test-token").
					ExpectError(errors.New("redis error"))
			},
			expectedErr: "failed to delete session from redis",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			conn := redigomock.NewConn()
			tt.setupMock(conn)

			pool := &redis.Pool{
				Dial:         func() (redis.Conn, error) { return conn, nil },
				TestOnBorrow: func(c redis.Conn, t time.Time) error { return nil },
			}
			defer pool.Close()

			repo := NewSessionRepository(pool)

			err := repo.Delete(tt.sessionID)

			if tt.expectedErr != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedErr)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestSessionRepository_Exists(t *testing.T) {
	tests := []struct {
		name        string
		sessionID   string
		setupMock   func(conn *redigomock.Conn)
		expected    bool
		expectedErr string
	}{
		{
			name:      "Exists",
			sessionID: "test-token",
			setupMock: func(conn *redigomock.Conn) {
				conn.Command("EXISTS", "session:test-token").Expect(int64(1))
			},
			expected: true,
		},
		{
			name:      "NotExists",
			sessionID: "non-existent",
			setupMock: func(conn *redigomock.Conn) {
				conn.Command("EXISTS", "session:non-existent").Expect(int64(0))
			},
			expected: false,
		},
		{
			name:      "RedisError",
			sessionID: "test-token",
			setupMock: func(conn *redigomock.Conn) {
				conn.Command("EXISTS", "session:test-token").
					ExpectError(errors.New("redis error"))
			},
			expectedErr: "failed to check session existence",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			conn := redigomock.NewConn()
			tt.setupMock(conn)

			pool := &redis.Pool{
				Dial:         func() (redis.Conn, error) { return conn, nil },
				TestOnBorrow: func(c redis.Conn, t time.Time) error { return nil },
			}
			defer pool.Close()

			repo := NewSessionRepository(pool)

			exists, err := repo.Exists(tt.sessionID)

			if tt.expectedErr != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedErr)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, exists)
			}
		})
	}
}

// Добавляем недостающие тесты

func TestSessionRepository_ConnectionErrors(t *testing.T) {
	t.Run("PoolGetError", func(t *testing.T) {
		pool := &redis.Pool{
			Dial: func() (redis.Conn, error) {
				return nil, errors.New("connection failed")
			},
		}
		defer pool.Close()

		repo := NewSessionRepository(pool)

		// Test all methods with connection error
		err := repo.Set(&domain.Session{Token: "test"})
		assert.Error(t, err)

		_, err = repo.Get("test")
		assert.Error(t, err)

		err = repo.Delete("test")
		assert.Error(t, err)

		_, err = repo.Exists("test")
		assert.Error(t, err)

		err = repo.SetWithExpiry(&domain.Session{Token: "test"}, time.Minute)
		assert.Error(t, err)
	})
}

func TestSessionRepository_KeyFormat(t *testing.T) {
	conn := redigomock.NewConn()

	session := &domain.Session{
		Token:     "unique-token",
		UserEmail: "user@example.com",
		ExpiresAt: time.Now().Add(time.Hour),
	}

	conn.Command("APPEND", "session:unique-token", redigomock.NewAnyData()).Expect("OK")

	pool := &redis.Pool{
		Dial:         func() (redis.Conn, error) { return conn, nil },
		TestOnBorrow: func(c redis.Conn, t time.Time) error { return nil },
	}
	defer pool.Close()

	repo := NewSessionRepository(pool)

	err := repo.Set(session)
	assert.NoError(t, err)

	// Проверяем, что команда была вызвана с правильным ключом
	cmd := conn.Command("APPEND", "session:unique-token", redigomock.NewAnyData())
	assert.Equal(t, 1, conn.Stats(cmd))
}

func TestSession_IsExpired(t *testing.T) {
	tests := []struct {
		name     string
		session  *domain.Session
		expected bool
	}{
		{
			name: "NotExpired",
			session: &domain.Session{
				ExpiresAt: time.Now().Add(time.Hour),
			},
			expected: false,
		},
		{
			name: "Expired",
			session: &domain.Session{
				ExpiresAt: time.Now().Add(-time.Hour),
			},
			expected: true,
		},
		{
			name: "JustExpired",
			session: &domain.Session{
				ExpiresAt: time.Now().Add(-time.Second),
			},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.session.IsExpired())
		})
	}
}

func TestSessionRepository_FullCycle(t *testing.T) {
	conn := redigomock.NewConn()

	// Создаем сессию через конструктор
	session := domain.NewSession("test@example.com", time.Hour)

	// Mock для Set
	conn.Command("APPEND", "session:"+session.Token, redigomock.NewAnyData()).Expect("OK")

	// Mock для Get
	data, _ := json.Marshal(session)
	conn.Command("GET", "session:"+session.Token).Expect(data)

	// Mock для Exists
	conn.Command("EXISTS", "session:"+session.Token).Expect(int64(1))

	// Mock для Delete
	conn.Command("DEL", "session:"+session.Token).Expect(int64(1))

	pool := &redis.Pool{
		Dial:         func() (redis.Conn, error) { return conn, nil },
		TestOnBorrow: func(c redis.Conn, t time.Time) error { return nil },
	}
	defer pool.Close()

	repo := NewSessionRepository(pool)

	// Set
	err := repo.Set(session)
	assert.NoError(t, err)

	// Exists
	exists, err := repo.Exists(session.Token)
	assert.NoError(t, err)
	assert.True(t, exists)

	// Get
	retrieved, err := repo.Get(session.Token)
	assert.NoError(t, err)
	assert.Equal(t, session.UserEmail, retrieved.UserEmail)
	assert.Equal(t, session.Token, retrieved.Token)

	// Проверяем, что сессия не истекла
	assert.False(t, retrieved.IsExpired())

	// Delete
	err = repo.Delete(session.Token)
	assert.NoError(t, err)
}
