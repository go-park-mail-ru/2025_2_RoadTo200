package interfaces

import (
	domain "github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/domain/entities"
	"github.com/google/uuid"
)

type ScreenRepository interface {
	Create(report *domain.Screen) error
	GetById(id uuid.UUID) (*domain.Screen, error)
	GetByReportId(reportId uuid.UUID) (*domain.Screen, error)
}
