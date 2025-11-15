package interfaces

import (
	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/domain/constants"
	domain "github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/domain/entities"
	"github.com/google/uuid"
)

type ReportRepository interface {
	Create(report *domain.Report) error
	GetByUser(uuid uuid.UUID) ([]domain.Report, error)
	GetById(uuid uuid.UUID) (*domain.Report, error)
	UpdateStatus(id uuid.UUID, status constants.ReportStatus) error
}
