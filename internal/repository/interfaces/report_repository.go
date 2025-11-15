package interfaces

import (
	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/domain/constants"
	domain "github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/domain/entities"
	"github.com/google/uuid"
)

// Альтернативный вариант с строгой типизацией категорий и статусов
type ReportStatistics struct {
	TotalTickets        int64       `json:"total_tickets"`
	TicketsByTheme      ThemeStats  `json:"tickets_by_theme"`
	TicketsByStatus     StatusStats `json:"tickets_by_status"`
	AverageResponseTime string      `json:"average_response_time"`
}

type ThemeStats struct {
	Technical int64 `json:"technical"`
	Feature   int64 `json:"feature"`
	Question  int64 `json:"question"`
	Security  int64 `json:"security"`
	Billing   int64 `json:"billing"`
	Device    int64 `json:"device"`
}

type StatusStats struct {
	Active int64 `json:"open"`
	Work   int64 `json:"work"`
	Closed int64 `json:"closed"`
}

type ReportRepository interface {
	Create(report *domain.Report) error
	GetByUser(uuid uuid.UUID) ([]domain.Report, error)
	GetById(uuid uuid.UUID) (*domain.Report, error)
	UpdateStatus(id uuid.UUID, status constants.ReportStatus) error
	GetStatistics() (*ReportStatistics, error)
}
