package domain

import (
	"time"

	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/domain/constants"
	"github.com/google/uuid"
)

type Report struct {
	ID        uuid.UUID              `json:"id"`
	UserID    uuid.UUID              `json:"user_id"`
	Theme     constants.ReportTheme  `json:"theme"`
	Problem   string                 `json:"problem"`
	Contact   string                 `json:"contact"`
	Comment   *string                `json:"comment"`
	Status    constants.ReportStatus `json:"status"`
	CreatedAt time.Time              `json:"report_at"`
	WorkAt    *time.Time             `json:"work_at"`
	ClosedAt  *time.Time             `json:"close_at"`
}
