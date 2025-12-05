package implementations

import (
	"context"
	"errors"
	"fmt"

	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/chat-service/repository/interfaces"
	domain "github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/domain/entities"
	expectation "github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/domain/errors"
	utils "github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/repository/implementations/postgres"
	bdIface "github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/repository/interfaces"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

var _ interfaces.MatchRepository = (*MatchRepository)(nil)

type MatchRepository struct {
	pool bdIface.PgxIface
}

func NewMatchRepository(pool bdIface.PgxIface) *MatchRepository {
	return &MatchRepository{pool: pool}
}

func (r *MatchRepository) Create(ctx context.Context, match *domain.Match) error {
	conn, err := utils.GetConn(ctx, r.pool)
	if err != nil {
		return err
	}

	// Убедимся, что user1_id всегда меньше user2_id для consistency
	user1ID, user2ID := match.User1ID, match.User2ID
	if user1ID.String() > user2ID.String() {
		user1ID, user2ID = user2ID, user1ID
	}

	query := `
		INSERT INTO match (user1_id, user2_id, is_active)
		VALUES ($1, $2, $3)
		RETURNING matched_at`

	err = conn.QueryRow(ctx, query, user1ID, user2ID, match.IsActive).
		Scan(&match.MatchedAt)

	if err != nil {
		return err
	}

	// Обновим IDs в соответствии с порядком в БД
	match.User1ID, match.User2ID = user1ID, user2ID
	return nil
}

func (r *MatchRepository) GetByID(ctx context.Context, matchID uuid.UUID) (*domain.Match, error) {
	conn, err := utils.GetConn(ctx, r.pool)
	if err != nil {
		return nil, err
	}

	query := `
		SELECT id, user1_id, user2_id, is_active, matched_at
		FROM match
		WHERE id = $1
	`

	var match domain.Match
	err = conn.QueryRow(ctx, query, matchID).Scan(
		&match.ID,
		&match.User1ID,
		&match.User2ID,
		&match.IsActive,
		&match.MatchedAt,
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
	conn, err := utils.GetConn(ctx, r.pool)
	if err != nil {
		return nil, err
	}

	// Приводим к consistent порядку
	if user1ID.String() > user2ID.String() {
		user1ID, user2ID = user2ID, user1ID
	}

	var match domain.Match
	query := `SELECT id, user1_id, user2_id, is_active, matched_at FROM match WHERE user1_id = $1 AND user2_id = $2`

	err = conn.QueryRow(ctx, query, user1ID, user2ID).Scan(
		&match.ID, &match.User1ID, &match.User2ID, &match.IsActive, &match.MatchedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &match, nil
}

func (r *MatchRepository) GetUserMatches(ctx context.Context, userID uuid.UUID, limit, offset int) ([]domain.Match, error) {
	conn, err := utils.GetConn(ctx, r.pool)
	if err != nil {
		return nil, err
	}

	query := `
        SELECT id, user1_id, user2_id, is_active, matched_at 
        FROM match 
        WHERE (user1_id = $1 OR user2_id = $1) 
        AND is_active = true
        ORDER BY matched_at DESC 
        LIMIT $2 OFFSET $3`

	rows, err := conn.Query(ctx, query, userID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get user matches: %w", err)
	}
	defer rows.Close()

	var matches []domain.Match
	for rows.Next() {
		var match domain.Match
		err := rows.Scan(
			&match.ID,
			&match.User1ID,
			&match.User2ID,
			&match.IsActive,
			&match.MatchedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan match: %w", err)
		}
		matches = append(matches, match)
	}

	return matches, nil
}

func (r *MatchRepository) UpdateActive(ctx context.Context, user1ID, user2ID uuid.UUID, isActive bool) error {
	conn, err := utils.GetConn(ctx, r.pool)
	if err != nil {
		return err
	}

	if user1ID.String() > user2ID.String() {
		user1ID, user2ID = user2ID, user1ID
	}

	query := `UPDATE match SET is_active = $1 WHERE user1_id = $2 AND user2_id = $3`

	_, err = conn.Exec(ctx, query, isActive, user1ID, user2ID)
	if err != nil {
		return err
	}
	return nil
}

func (r *MatchRepository) Delete(ctx context.Context, user1ID, user2ID uuid.UUID) error {
	conn, err := utils.GetConn(ctx, r.pool)
	if err != nil {
		return err
	}

	if user1ID.String() > user2ID.String() {
		user1ID, user2ID = user2ID, user1ID
	}

	query := `DELETE FROM match WHERE user1_id = $1 AND user2_id = $2`

	_, err = conn.Exec(ctx, query, user1ID, user2ID)
	if err != nil {
		return err
	}
	return nil
}

func (r *MatchRepository) CheckMutualLike(ctx context.Context, user1ID, user2ID uuid.UUID) (bool, error) {
	conn, err := utils.GetConn(ctx, r.pool)
	if err != nil {
		return false, err
	}

	query := `
		SELECT EXISTS(
			SELECT 1 FROM swipe s1
			INNER JOIN swipe s2 ON s1.swiper_user_id = s2.target_user_id AND s1.target_user_id = s2.swiper_user_id
			WHERE s1.swiper_user_id = $1 AND s1.target_user_id = $2
			AND s2.swiper_user_id = $2 AND s2.target_user_id = $1
			AND s1.swipe_type = 'like' AND s2.swipe_type = 'like'
		)`

	var exists bool
	err = conn.QueryRow(ctx, query, user1ID, user2ID).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}
