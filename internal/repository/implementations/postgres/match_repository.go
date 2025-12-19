package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	domain "github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/domain/entities"
	expectation "github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/domain/errors"
	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/repository/interfaces"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

var _ interfaces.MatchRepository = (*MatchRepository)(nil)

type MatchRepository struct {
	pool interfaces.PgxIface
}

func NewMatchRepository(pool interfaces.PgxIface) *MatchRepository {
	return &MatchRepository{pool: pool}
}

func (r *MatchRepository) Create(ctx context.Context, match *domain.Match) error {
	// Убедимся, что user1_id всегда меньше user2_id для consistency
	user1ID, user2ID := match.User1ID, match.User2ID
	if user1ID.String() > user2ID.String() {
		user1ID, user2ID = user2ID, user1ID
	}

	query := `
		INSERT INTO match (user1_id, user2_id, is_active, expires_at)
		VALUES ($1, $2, $3, NOW() + INTERVAL '24 hours')
		RETURNING id, matched_at, expires_at`

	var expiresAt sql.NullTime
	err := r.pool.QueryRow(ctx, query, user1ID, user2ID, match.IsActive).
		Scan(&match.ID, &match.MatchedAt, &expiresAt)
	if err != nil {
		return err
	}
	if expiresAt.Valid {
		match.ExpiresAt = &expiresAt.Time
	}

	if err != nil {
		return err
	}

	// Обновим IDs в соответствии с порядком в БД
	match.User1ID, match.User2ID = user1ID, user2ID
	return nil
}

func (r *MatchRepository) GetByID(ctx context.Context, matchID uuid.UUID) (*domain.Match, error) {
	query := `
		SELECT id, user1_id, user2_id, is_active, matched_at, expires_at
		FROM match
		WHERE id = $1
	`

	var match domain.Match
	err := r.pool.QueryRow(ctx, query, matchID).Scan(
		&match.ID,
		&match.User1ID,
		&match.User2ID,
		&match.IsActive,
		&match.MatchedAt,
		&match.ExpiresAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, expectation.ErrMatchNotFound
		}
		return nil, fmt.Errorf("failed to get match by id: %w", err)
	}

	return &match, nil
}

func (r *MatchRepository) GetByUsers(ctx context.Context, user1ID, user2ID uuid.UUID) (*domain.Match, error) {
	// Приводим к consistent порядку
	if user1ID.String() > user2ID.String() {
		user1ID, user2ID = user2ID, user1ID
	}

	var match domain.Match
	query := `SELECT id, user1_id, user2_id, is_active, matched_at, expires_at FROM match WHERE user1_id = $1 AND user2_id = $2`

	var expiresAt sql.NullTime
	err := r.pool.QueryRow(ctx, query, user1ID, user2ID).Scan(
		&match.ID, &match.User1ID, &match.User2ID, &match.IsActive, &match.MatchedAt, &expiresAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	if expiresAt.Valid {
		match.ExpiresAt = &expiresAt.Time
	}

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &match, nil
}

func (r *MatchRepository) GetUserMatches(ctx context.Context, userID uuid.UUID, limit, offset int) ([]domain.Match, error) {
	query := `
        SELECT id, user1_id, user2_id, is_active, matched_at, expires_at
        FROM match 
        WHERE (user1_id = $1 OR user2_id = $1) 
        AND is_active = true
        ORDER BY matched_at DESC 
        LIMIT $2 OFFSET $3`

	rows, err := r.pool.Query(ctx, query, userID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get user matches: %w", err)
	}
	defer rows.Close()

	var matches []domain.Match
	for rows.Next() {
		var match domain.Match
		var expiresAt sql.NullTime
		err := rows.Scan(
			&match.ID,
			&match.User1ID,
			&match.User2ID,
			&match.IsActive,
			&match.MatchedAt,
			&expiresAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan match: %w", err)
		}
		if expiresAt.Valid {
			match.ExpiresAt = &expiresAt.Time
		}
		matches = append(matches, match)
	}

	return matches, nil
}

func (r *MatchRepository) UpdateActive(ctx context.Context, user1ID, user2ID uuid.UUID, isActive bool) error {
	if user1ID.String() > user2ID.String() {
		user1ID, user2ID = user2ID, user1ID
	}

	query := `UPDATE match SET is_active = $1 WHERE user1_id = $2 AND user2_id = $3`

	_, err := r.pool.Exec(ctx, query, isActive, user1ID, user2ID)
	if err != nil {
		return err
	}
	return nil
}

func (r *MatchRepository) Delete(ctx context.Context, user1ID, user2ID uuid.UUID) error {
	if user1ID.String() > user2ID.String() {
		user1ID, user2ID = user2ID, user1ID
	}

	query := `DELETE FROM match WHERE user1_id = $1 AND user2_id = $2`

	_, err := r.pool.Exec(ctx, query, user1ID, user2ID)
	if err != nil {
		return err
	}
	return nil
}

func (r *MatchRepository) CheckMutualLike(ctx context.Context, user1ID, user2ID uuid.UUID) (bool, error) {
	query := `
		SELECT EXISTS(
			SELECT 1 FROM swipe s1
			INNER JOIN swipe s2 ON s1.swiper_user_id = s2.target_user_id AND s1.target_user_id = s2.swiper_user_id
			WHERE s1.swiper_user_id = $1 AND s1.target_user_id = $2
			AND s2.swiper_user_id = $2 AND s2.target_user_id = $1
			AND s1.swipe_type IN ('like', 'super_like')
			AND s2.swipe_type IN ('like', 'super_like')
		)`

	var exists bool
	err := r.pool.QueryRow(ctx, query, user1ID, user2ID).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

// DeactivateExpiredMatches деактивирует мэтчи, у которых истекло 24-часовое окно
func (r *MatchRepository) DeactivateExpiredMatches(ctx context.Context) ([]domain.Match, error) {
	query := `
		UPDATE match 
		SET is_active = FALSE 
		WHERE expires_at IS NOT NULL 
		  AND expires_at < NOW() 
		  AND is_active = TRUE
		RETURNING id, user1_id, user2_id, is_active, matched_at, expires_at`

	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to deactivate expired matches: %w", err)
	}
	defer rows.Close()

	var deactivatedMatches []domain.Match
	for rows.Next() {
		var match domain.Match
		var expiresAt sql.NullTime
		if err := rows.Scan(&match.ID, &match.User1ID, &match.User2ID, &match.IsActive, &match.MatchedAt, &expiresAt); err != nil {
			return deactivatedMatches, fmt.Errorf("failed to scan match: %w", err)
		}
		if expiresAt.Valid {
			match.ExpiresAt = &expiresAt.Time
		}
		deactivatedMatches = append(deactivatedMatches, match)
	}

	return deactivatedMatches, nil
}

// SetExpiresAtNull устанавливает expires_at в NULL для матча (когда начинается чат)
func (r *MatchRepository) SetExpiresAtNull(ctx context.Context, user1ID, user2ID uuid.UUID) error {
	fmt.Printf("[MatchRepository.SetExpiresAtNull] Called with user1ID=%s, user2ID=%s\n", user1ID, user2ID)

	// Приводим к consistent порядку
	originalUser1ID, originalUser2ID := user1ID, user2ID
	if user1ID.String() > user2ID.String() {
		user1ID, user2ID = user2ID, user1ID
		fmt.Printf("[MatchRepository.SetExpiresAtNull] Reordered: user1ID=%s, user2ID=%s\n", user1ID, user2ID)
	} else {
		fmt.Printf("[MatchRepository.SetExpiresAtNull] No reordering needed: user1ID=%s, user2ID=%s\n", user1ID, user2ID)
	}

	query := `
		UPDATE match
		SET expires_at = NULL
		WHERE user1_id = $1 AND user2_id = $2 AND expires_at IS NOT NULL`

	fmt.Printf("[MatchRepository.SetExpiresAtNull] Executing query: UPDATE match SET expires_at = NULL WHERE user1_id = %s AND user2_id = %s AND expires_at IS NOT NULL\n", user1ID, user2ID)

	result, err := r.pool.Exec(ctx, query, user1ID, user2ID)
	if err != nil {
		fmt.Printf("[MatchRepository.SetExpiresAtNull] ERROR: Query execution failed: %v\n", err)
		return fmt.Errorf("failed to set expires_at to null: %w", err)
	}

	rowsAffected := result.RowsAffected()
	fmt.Printf("[MatchRepository.SetExpiresAtNull] Query executed successfully. Rows affected: %d\n", rowsAffected)

	if rowsAffected == 0 {
		fmt.Printf("[MatchRepository.SetExpiresAtNull] WARNING: No rows affected. Possible reasons:\n")
		fmt.Printf("  - Match with user1_id=%s and user2_id=%s not found\n", user1ID, user2ID)
		fmt.Printf("  - Match already has expires_at = NULL\n")
		fmt.Printf("  - Original IDs were: user1ID=%s, user2ID=%s\n", originalUser1ID, originalUser2ID)
	} else {
		fmt.Printf("[MatchRepository.SetExpiresAtNull] SUCCESS: Match expires_at set to NULL (%d row(s) updated)\n", rowsAffected)
	}

	return nil
}
