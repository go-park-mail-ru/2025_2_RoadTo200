package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/domain/constants"
	domain "github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/domain/entities"
	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/domain/errors"
	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/repository/interfaces"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

var _ interfaces.StrikeRepository = (*StrikeRepository)(nil)

type StrikeRepository struct {
	db interfaces.PgxIface
}

func NewStrikeRepository(db interfaces.PgxIface) *StrikeRepository {
	return &StrikeRepository{
		db: db,
	}
}

// CreateStrike создает новую жалобу
func (r *StrikeRepository) CreateStrike(ctx context.Context, strike *domain.Strike) error {
	query := `
		INSERT INTO strike (id, reporter_id, target_user_id, reason, status)
		VALUES ($1, $2, $3, $4, $5, $6)
	`

	_, err := r.db.Exec(ctx, query,
		uuid.New(),
		strike.ReporterID,
		strike.TargetUserID,
		strike.Reason,
		constants.StrikeStatusPending,
	)

	if err != nil {
		return fmt.Errorf("failed to create strike: %w", err)
	}

	return nil
}

// GetStrikeByID получает жалобу по ID
func (r *StrikeRepository) GetStrikeByID(ctx context.Context, strikeID string) (*domain.Strike, error) {
	query := `
		SELECT id, reporter_id, target_user_id, type, reason, status, created_at, updated_at, moderator_id, moderator_note
		FROM strike 
		WHERE id = $1
	`

	var strike domain.Strike
	err := r.db.QueryRow(ctx, query, strikeID).Scan(
		&strike.ID,
		&strike.ReporterID,
		&strike.TargetUserID,
		&strike.Type,
		&strike.Reason,
		&strike.Status,
		&strike.CreatedAt,
		&strike.UpdatedAt,
		&strike.ModeratorID,
		&strike.ModeratorNote,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, errors.ErrStrikeNotFound
		}
		return nil, fmt.Errorf("failed to get strike by ID: %w", err)
	}

	return &strike, nil
}

// GetStrikesByUserID получает все жалобы на конкретного пользователя
func (r *StrikeRepository) GetStrikesByUserID(ctx context.Context, userID string, limit, offset int) ([]*domain.Strike, error) {
	query := `
		SELECT id, reporter_id, target_user_id, type, reason, status, created_at, updated_at, moderator_id, moderator_note
		FROM strike 
		WHERE target_user_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.Query(ctx, query, userID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get strikes by user ID: %w", err)
	}
	defer rows.Close()

	var strikes []*domain.Strike
	for rows.Next() {
		var strike domain.Strike
		err := rows.Scan(
			&strike.ID,
			&strike.ReporterID,
			&strike.TargetUserID,
			&strike.Type,
			&strike.Reason,
			&strike.Status,
			&strike.CreatedAt,
			&strike.UpdatedAt,
			&strike.ModeratorID,
			&strike.ModeratorNote,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan strike: %w", err)
		}
		strikes = append(strikes, &strike)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating strikes: %w", err)
	}

	return strikes, nil
}

// GetStrikesByType получает жалобы по типу нарушения
func (r *StrikeRepository) GetStrikesByType(ctx context.Context, strikeType constants.StrikeType, limit, offset int) ([]*domain.Strike, error) {
	query := `
		SELECT id, reporter_id, target_user_id, type, reason, status, created_at, updated_at, moderator_id, moderator_note
		FROM strike 
		WHERE type = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.Query(ctx, query, strikeType, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get strikes by type: %w", err)
	}
	defer rows.Close()

	var strikes []*domain.Strike
	for rows.Next() {
		var strike domain.Strike
		err := rows.Scan(
			&strike.ID,
			&strike.ReporterID,
			&strike.TargetUserID,
			&strike.Type,
			&strike.Reason,
			&strike.Status,
			&strike.CreatedAt,
			&strike.UpdatedAt,
			&strike.ModeratorID,
			&strike.ModeratorNote,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan strike: %w", err)
		}
		strikes = append(strikes, &strike)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating strikes: %w", err)
	}

	return strikes, nil
}

// GetStrikesByDateRange получает жалобы за определенный период
func (r *StrikeRepository) GetStrikesByDateRange(ctx context.Context, from, to time.Time, limit, offset int) ([]*domain.Strike, error) {
	query := `
		SELECT id, reporter_id, target_user_id, type, reason, status, created_at, updated_at, moderator_id, moderator_note
		FROM strike 
		WHERE created_at BETWEEN $1 AND $2
		ORDER BY created_at DESC
		LIMIT $3 OFFSET $4
	`

	rows, err := r.db.Query(ctx, query, from, to, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get strikes by date range: %w", err)
	}
	defer rows.Close()

	var strikes []*domain.Strike
	for rows.Next() {
		var strike domain.Strike
		err := rows.Scan(
			&strike.ID,
			&strike.ReporterID,
			&strike.TargetUserID,
			&strike.Type,
			&strike.Reason,
			&strike.Status,
			&strike.CreatedAt,
			&strike.UpdatedAt,
			&strike.ModeratorID,
			&strike.ModeratorNote,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan strike: %w", err)
		}
		strikes = append(strikes, &strike)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating strikes: %w", err)
	}

	return strikes, nil
}

// UpdateStrikeStatus обновляет статус жалобы
func (r *StrikeRepository) UpdateStrikeStatus(ctx context.Context, strikeID string, status constants.StrikeStatus) error {
	query := `
		UPDATE strike 
		SET status = $1, updated_at = $2
		WHERE id = $3
	`

	now := time.Now()
	result, err := r.db.Exec(ctx, query, status, now, strikeID)
	if err != nil {
		return fmt.Errorf("failed to update strike status: %w", err)
	}

	if result.RowsAffected() == 0 {
		return errors.ErrStrikeNotFound
	}

	return nil
}

// DeleteStrike удаляет жалобу
func (r *StrikeRepository) DeleteStrike(ctx context.Context, strikeID string) error {
	query := `DELETE FROM strike WHERE id = $1`

	_, err := r.db.Exec(ctx, query, strikeID)
	if err != nil {
		return fmt.Errorf("failed to delete strike %s: %w", strikeID, err)
	}
	return nil
}

// HasActiveStrikeFromUser проверяет, есть ли уже активная жалоба от этого пользователя на целевого
func (r *StrikeRepository) HasActiveStrikeFromUser(ctx context.Context, reporterID, targetUserID string) (bool, error) {
	query := `
		SELECT COUNT(*) 
		FROM strike 
		WHERE reporter_id = $1 
		AND target_user_id = $2 
		AND status = $3
	`

	var count int
	err := r.db.QueryRow(ctx, query, reporterID, targetUserID, constants.StrikeStatusPending).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("failed to check active strike: %w", err)
	}

	return count > 0, nil
}
