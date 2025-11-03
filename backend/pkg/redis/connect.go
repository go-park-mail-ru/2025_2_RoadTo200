package redis

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/config"
	"github.com/gomodule/redigo/redis"
)

func NewConnection(cfg *config.RedisConfig) (*redis.Pool, error) {
	host := os.Getenv(cfg.Host)
	port, err := strconv.Atoi(os.Getenv(cfg.Port))
	if err != nil {
		return nil, err
	}
	password := os.Getenv(cfg.Password)
	base, err := strconv.Atoi(os.Getenv(cfg.Base))
	if err != nil {
		return nil, err
	}
	address := fmt.Sprintf("%s:%d", host, port)

	pool := &redis.Pool{
		MaxIdle:     cfg.MaxIdle,
		MaxActive:   cfg.MaxActive, // 0 означает нет ограничения
		IdleTimeout: cfg.IdleTimeout,
		Wait:        true, // Ждать свободного соединения если достигнут лимит
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", address)
			if err != nil {
				return nil, fmt.Errorf("failed to dial redis: %w", err)
			}

			// Аутентификация если указан пароль
			if password != "" {
				if _, err := c.Do("AUTH", password); err != nil {
					c.Close()
					return nil, fmt.Errorf("redis auth failed: %w", err)
				}
			}

			// Выбор базы данных
			if _, err := c.Do("SELECT", base); err != nil {
				c.Close()
				return nil, fmt.Errorf("failed to select redis db: %w", err)
			}

			return c, nil
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			if time.Since(t) < time.Minute {
				return nil
			}
			_, err := c.Do("PING")
			return err
		},
	}
	return pool, nil
}
