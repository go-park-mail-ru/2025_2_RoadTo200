package postgres

import (
	"context"

	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/domain/constants"
	domain "github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/domain/entities"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type ReportRepository struct {
	pool *pgxpool.Pool
}

func NewReportRepository(pool *pgxpool.Pool) *ReportRepository {
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
