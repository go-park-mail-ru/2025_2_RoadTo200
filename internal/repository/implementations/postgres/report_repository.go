package postgres

import (
	"context"

	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/domain/constants"
	domain "github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/domain/entities"
	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/repository/interfaces"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
)

type ReportRepository struct {
	pool interfaces.PgxIface
}

func NewReportRepository(pool interfaces.PgxIface) *ReportRepository {
	return &ReportRepository{pool}
}

func (r *ReportRepository) Create(report *domain.Report) error {
	query := `
		INSERT INTO reports (user_id, theme, problem, contact)
		VALUES ($1, $2, $3, $4)`

	_, err := r.pool.Exec(
		context.Background(),
		query,
		report.UserID,
		report.Theme,
		report.Problem,
		report.Contact,
	)

	return err
}

func (r *ReportRepository) GetByUser(userID uuid.UUID) ([]domain.Report, error) {
	query := `
		SELECT id, user_id, theme, problem, contact, comment, status, created_at, work_at, closed_at
		FROM reports 
		WHERE user_id = $1
		ORDER BY created_at DESC`

	rows, err := r.pool.Query(context.Background(), query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reports []domain.Report
	for rows.Next() {
		var report domain.Report
		err := rows.Scan(
			&report.ID,
			&report.UserID,
			&report.Theme,
			&report.Problem,
			&report.Contact,
			&report.Comment,
			&report.Status,
			&report.CreatedAt,
			&report.WorkAt,
			&report.ClosedAt,
		)
		if err != nil {
			return nil, err
		}
		reports = append(reports, report)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return reports, nil
}

func (r *ReportRepository) GetById(id uuid.UUID) (*domain.Report, error) {
	query := `
		SELECT id, user_id, theme, problem, contact, comment, status, created_at, work_at, closed_at
		FROM reports 
		WHERE id = $1`

	var report domain.Report
	err := r.pool.QueryRow(context.Background(), query, id).Scan(
		&report.ID,
		&report.UserID,
		&report.Theme,
		&report.Problem,
		&report.Contact,
		&report.Comment,
		&report.Status,
		&report.CreatedAt,
		&report.WorkAt,
		&report.ClosedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &report, nil
}

func (r *ReportRepository) UpdateStatus(id uuid.UUID, status constants.ReportStatus) error {
	query := `
		UPDATE reports 
		SET status = $1
		WHERE id = $2`

	_, err := r.pool.Exec(context.Background(), query, status, id)
	return err
}

func (r *ReportRepository) GetStatistics() (*interfaces.ReportStatistics, error) {
	query := `
        WITH response_times AS (
            SELECT 
                EXTRACT(EPOCH FROM (work_at - created_at)) as response_seconds
            FROM reports 
            WHERE work_at IS NOT NULL 
                AND created_at IS NOT NULL
                AND work_at > created_at
        ),
        avg_response AS (
            SELECT 
                CASE 
                    WHEN COUNT(*) > 0 THEN 
                        FLOOR(AVG(response_seconds) / 3600) || 'h' || 
                        FLOOR((AVG(response_seconds) % 3600) / 60) || 'm'
                    ELSE '0h0m'
                END as avg_response_time
            FROM response_times
        )
        SELECT 
            COUNT(*) as total_tickets,
            COUNT(*) FILTER (WHERE theme = 'technical') as technical,
            COUNT(*) FILTER (WHERE theme = 'feature') as feature,
            COUNT(*) FILTER (WHERE theme = 'question') as question,
            COUNT(*) FILTER (WHERE theme = 'security') as security,
            COUNT(*) FILTER (WHERE theme = 'billing') as billing,
            COUNT(*) FILTER (WHERE theme = 'device') as device,
            COUNT(*) FILTER (WHERE status = 'open') as open,
            COUNT(*) FILTER (WHERE status = 'work') as work,
            COUNT(*) FILTER (WHERE status = 'closed') as closed,
            (SELECT avg_response_time FROM avg_response) as average_response_time
        FROM reports`

	var row interfaces.ReportStatistics
	err := r.pool.QueryRow(context.Background(), query).Scan(
		&row.TotalTickets,
		&row.TicketsByTheme.Technical,
		&row.TicketsByTheme.Feature,
		&row.TicketsByTheme.Question,
		&row.TicketsByTheme.Security,
		&row.TicketsByTheme.Billing,
		&row.TicketsByTheme.Device,
		&row.TicketsByStatus.Open,
		&row.TicketsByStatus.Work,
		&row.TicketsByStatus.Closed,
		&row.AverageResponseTime,
	)
	if err != nil {
		return nil, err
	}
	return &row, nil
}
