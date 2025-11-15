package app

import (
	"fmt"
	"os"
	"strconv"

	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/pkg/migration"
	"github.com/golang-migrate/migrate/v4"
)

func (a *App) runMigrations() error {
	if !a.config.Postgres.Migrated {
		a.logger.Infof("Skip migrations")
		return nil
	}
	a.logger.Info("🔄 Running database migrations...")
	host := os.Getenv(a.config.Postgres.Host)         // ❌ Это ищет env переменную с именем "localhost"
	sport := os.Getenv(a.config.Postgres.Port)        // ❌ Это ищет env переменную с именем "5435"
	user := os.Getenv(a.config.Postgres.User)         // ❌ Это ищет env переменную с именем "postgres"
	password := os.Getenv(a.config.Postgres.Password) // ❌ Это ищет env переменную с именем "password"
	base := os.Getenv(a.config.Postgres.Base)         // ❌ Это ищет env переменную с именем "dating_app"

	// Преобразуем порт в число
	port, err := strconv.Atoi(sport)
	if err != nil {
		return fmt.Errorf("failed to convert port to int: %w", err)
	}
	// Получаем connection string из конфигурации PostgreSQL
	dbURL := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s",
		user, password, host, port, base,
		"disable",
	)

	migrator, err := migration.NewMigrator(dbURL, a.config.Postgres.Migrations)
	if err != nil {
		return fmt.Errorf("failed to create migrator: %w", err)
	}
	defer migrator.Close()

	if err := migrator.Up(); err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	version, dirty, err := migrator.Version()
	if err != nil && err != migrate.ErrNilVersion {
		return fmt.Errorf("failed to get migration version: %w", err)
	}

	if dirty {
		a.logger.Warn("Database is in dirty state, consider fixing manually")
	}

	a.logger.Infof("✅ Database migrations applied successfully. Version: %d", version)
	return nil
}
