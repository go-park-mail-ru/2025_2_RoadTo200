package redis

import (
	"fmt"
	"strconv"
	"time"

	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/config"
	"github.com/gomodule/redigo/redis"
)

func NewConnection(cfg *config.RedisConfig) (*redis.Pool, error) {
	// Используем значения НАПРЯМУЮ из конфига
	host := cfg.Host         // "localhost"
	port := cfg.Port         // "6377"
	password := cfg.Password // твой пароль
	base := cfg.Base         // "0"

	// Преобразуем порт и базу в числа
	portNum, err := strconv.Atoi(port)
	if err != nil {
		return nil, fmt.Errorf("failed to convert port to int: %w", err)
	}

	baseNum, err := strconv.Atoi(base)
	if err != nil {
		return nil, fmt.Errorf("failed to convert base to int: %w", err)
	}

	address := fmt.Sprintf("%s:%d", host, portNum)

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
			if _, err := c.Do("SELECT", baseNum); err != nil {
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
