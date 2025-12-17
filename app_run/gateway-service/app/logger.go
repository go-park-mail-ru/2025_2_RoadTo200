package app

import (
	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/logger"
)

func (a *App) initLogger() error {
	logg := logger.New(&a.config.Logger)
	a.logger = logg
	logg.Info("✅ Logger initialized")
	return nil
}
