package app

import (
	"fmt"

	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/pkg/migration"
	"github.com/golang-migrate/migrate/v4"
)

func (a *App) runMigrations() error {
	if !a.config.Postgres.Migrated {
		a.logger.Infof("Skip migrations")
		return nil
	}
	a.logger.Info("🔄 Running database migrations...")

	// Получаем connection string из конфигурации PostgreSQL
	connStr := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		a.config.Postgres.User,
		a.config.Postgres.Password,
		a.config.Postgres.Host,
		a.config.Postgres.Port,
		a.config.Postgres.Base,
	)

	migrator, err := migration.NewMigrator(connStr, a.config.Postgres.Migrations)
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

	a.logger.Info("✅ Database migrations applied successfully. Version: %d", version)
	return nil
}
