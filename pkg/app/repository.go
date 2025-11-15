package app

import (
	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/repository/implementations/minio"
	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/repository/implementations/postgres"
	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/repository/implementations/redis"
)

func (a *App) initRepository() {
	// Репозитории
	a.repositories = &Repositories{
		Storage:    minio.NewStorageRepository(a.resources.MinIO, &a.config.MinIO),
		User:       postgres.NewUserRepository(a.resources.Postgres, a.logger),
		Session:    redis.NewSessionRepository(a.resources.Redis),
		Preference: postgres.NewUserPreferenceRepository(a.resources.Postgres, a.logger),
		Photo:      postgres.NewUserPhotoRepository(a.resources.Postgres, a.logger),
		Swipe:      postgres.NewSwipeRepository(a.resources.Postgres),
		Match:      postgres.NewMatchRepository(a.resources.Postgres),
		Report:     postgres.NewReportRepository(a.resources.Postgres),
		Screen:     postgres.NewScreenRepository(a.resources.Postgres),
	}
}
