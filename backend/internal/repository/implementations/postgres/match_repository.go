package postgres

import (
	"context"
	"errors"
	"fmt"

	domain "github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/domain/entities"
	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/repository/interfaces"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
)

type matchRepository struct {
	pool interfaces.PgxIface
}

func NewMatchRepository(pool interfaces.PgxIface) interfaces.MatchRepository {
	return &matchRepository{pool: pool}
}

func (r *matchRepository) Create(match *domain.Match) error {
	// Убедимся, что user1_id всегда меньше user2_id для consistency
	user1ID, user2ID := match.User1ID, match.User2ID
	if user1ID.String() > user2ID.String() {
		user1ID, user2ID = user2ID, user1ID
	}

	query := `
		INSERT INTO match (user1_id, user2_id, is_active)
		VALUES ($1, $2, $3)
		RETURNING matched_at`

	err := r.pool.QueryRow(context.Background(), query, user1ID, user2ID, match.IsActive).
		Scan(&match.MatchedAt)

	if err != nil {
		return err
	}

	// Обновим IDs в соответствии с порядком в БД
	match.User1ID, match.User2ID = user1ID, user2ID
	return nil
}

func (r *matchRepository) GetByUsers(user1ID, user2ID uuid.UUID) (*domain.Match, error) {
	// Приводим к consistent порядку
	if user1ID.String() > user2ID.String() {
		user1ID, user2ID = user2ID, user1ID
	}

	var match domain.Match
	query := `SELECT * FROM match WHERE user1_id = $1 AND user2_id = $2`

	err := r.pool.QueryRow(context.Background(), query, user1ID, user2ID).Scan(
		&match.User1ID, &match.User2ID, &match.IsActive, &match.MatchedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &match, nil
}

func (r *matchRepository) GetUserMatches(userID uuid.UUID, limit, offset int) ([]domain.Match, error) {
	query := `
        SELECT user1_id, user2_id, is_active, matched_at 
        FROM match 
        WHERE (user1_id = $1 OR user2_id = $1) 
        AND is_active = true
        ORDER BY matched_at DESC 
        LIMIT $2 OFFSET $3`

	rows, err := r.pool.Query(context.Background(), query, userID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get user matches: %w", err)
	}
	defer rows.Close()

	var matches []domain.Match
	for rows.Next() {
		var match domain.Match
		err := rows.Scan(
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

func (r *matchRepository) UpdateActive(user1ID, user2ID uuid.UUID, isActive bool) error {
	if user1ID.String() > user2ID.String() {
		user1ID, user2ID = user2ID, user1ID
	}

	query := `UPDATE match SET is_active = $1 WHERE user1_id = $2 AND user2_id = $3`

	_, err := r.pool.Exec(context.Background(), query, isActive, user1ID, user2ID)
	if err != nil {
		return err
	}
	return nil
}

func (r *matchRepository) Delete(user1ID, user2ID uuid.UUID) error {
	if user1ID.String() > user2ID.String() {
		user1ID, user2ID = user2ID, user1ID
	}

	query := `DELETE FROM match WHERE user1_id = $1 AND user2_id = $2`

	_, err := r.pool.Exec(context.Background(), query, user1ID, user2ID)
	if err != nil {
		return err
	}
	return nil
}

func (r *matchRepository) CheckMutualLike(user1ID, user2ID uuid.UUID) (bool, error) {
	query := `
		SELECT EXISTS(
			SELECT 1 FROM swipe s1
			INNER JOIN swipe s2 ON s1.swiper_user_id = s2.target_user_id AND s1.target_user_id = s2.swiper_user_id
			WHERE s1.swiper_user_id = $1 AND s1.target_user_id = $2
			AND s2.swiper_user_id = $2 AND s2.target_user_id = $1
			AND s1.swipe_type = 'like' AND s2.swipe_type = 'like'
		)`

	var exists bool
	err := r.pool.QueryRow(context.Background(), query, user1ID, user2ID).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}
