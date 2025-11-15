package postgres

import (
	"context"

	domain "github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/domain/entities"
	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/repository/interfaces"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
)

type ScreenRepository struct {
	pool interfaces.PgxIface
}

func NewScreenRepository(pool interfaces.PgxIface) *ScreenRepository {
	return &ScreenRepository{pool: pool}
}

func (r *ScreenRepository) Create(screen *domain.Screen) error {
	query := `
		INSERT INTO screen (id, report_id, url)
		VALUES ($1, $2, $3)`

	_, err := r.pool.Exec(
		context.Background(),
		query,
		screen.ID,
		screen.ReportId,
		screen.Url,
	)
	return err
}

func (r *ScreenRepository) GetById(id uuid.UUID) (*domain.Screen, error) {
	query := `
		SELECT id, report_id, url
		FROM screen 
		WHERE id = $1`

	var screen domain.Screen
	err := r.pool.QueryRow(context.Background(), query, id.String()).Scan(
		&screen.ID,
		&screen.ReportId,
		&screen.Url,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &screen, nil
}

func (r *ScreenRepository) GetByReportId(reportId uuid.UUID) (*domain.Screen, error) {
	query := `
		SELECT id, report_id, url
		FROM screen 
		WHERE report_id = $1`

	var screen domain.Screen
	err := r.pool.QueryRow(context.Background(), query, reportId.String()).Scan(
		&screen.ID,
		&screen.ReportId,
		&screen.Url,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &screen, nil
}
