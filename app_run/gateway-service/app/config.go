package app

import (
	"fmt"

	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/config"
	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/logger"
)

func (a *App) initConfig() error {
	logger.Println("🚀 Starting server...")

	cfg, err := config.NewConfig()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	a.config = cfg
	logger.Printf("✅ Config loaded")
	return nil
}
