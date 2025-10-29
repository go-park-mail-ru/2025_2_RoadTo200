package postgres

import (
	"context"
	"errors"

	domain "github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/domain/entities"
	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/repository/interfaces"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type swipeRepository struct {
	pool *pgxpool.Pool
}

func NewSwipeRepository(pool *pgxpool.Pool) interfaces.SwipeRepository {
	return &swipeRepository{pool: pool}
}

func (r *swipeRepository) Create(swipe *domain.Swipe) error {
	query := `
		INSERT INTO swipe (swiper_user_id, target_user_id, swipe_type)
		VALUES ($1, $2, $3)
		RETURNING created_at`

	err := r.pool.QueryRow(context.Background(), query,
		swipe.SwiperUserID, swipe.TargetUserID, swipe.SwipeType).
		Scan(&swipe.CreatedAt)

	if err != nil {
		return err
	}
	return nil
}

func (r *swipeRepository) GetBySwiperAndTarget(swiperID, targetID uuid.UUID) (*domain.Swipe, error) {
	var swipe domain.Swipe
	query := `SELECT * FROM swipe WHERE swiper_user_id = $1 AND target_user_id = $2`

	err := r.pool.QueryRow(context.Background(), query, swiperID, targetID).Scan(
		&swipe.SwiperUserID, &swipe.TargetUserID, &swipe.SwipeType, &swipe.CreatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &swipe, nil
}

func (r *swipeRepository) GetSwipesBySwiper(swiperID uuid.UUID, limit, offset int) ([]domain.Swipe, error) {
	query := `
		SELECT * FROM swipe 
		WHERE swiper_user_id = $1 
		ORDER BY created_at DESC 
		LIMIT $2 OFFSET $3`

	rows, err := r.pool.Query(context.Background(), query, swiperID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var swipes []domain.Swipe
	for rows.Next() {
		var swipe domain.Swipe
		err := rows.Scan(
			&swipe.SwiperUserID, &swipe.TargetUserID, &swipe.SwipeType, &swipe.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		swipes = append(swipes, swipe)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return swipes, nil
}

func (r *swipeRepository) GetSwipesByTarget(targetID uuid.UUID, limit, offset int) ([]domain.Swipe, error) {
	query := `
		SELECT * FROM swipe 
		WHERE target_user_id = $1 
		ORDER BY created_at DESC 
		LIMIT $2 OFFSET $3`

	rows, err := r.pool.Query(context.Background(), query, targetID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var swipes []domain.Swipe
	for rows.Next() {
		var swipe domain.Swipe
		err := rows.Scan(
			&swipe.SwiperUserID, &swipe.TargetUserID, &swipe.SwipeType, &swipe.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		swipes = append(swipes, swipe)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return swipes, nil
}

func (r *swipeRepository) Exists(swiperID, targetID uuid.UUID) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM swipe WHERE swiper_user_id = $1 AND target_user_id = $2)`

	err := r.pool.QueryRow(context.Background(), query, swiperID, targetID).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

func (r *swipeRepository) GetSwipesStats(userID uuid.UUID) (int, int, int, error) {
	var likesCount, dislikesCount, superLikesCount int

	query := `
		SELECT 
			COUNT(CASE WHEN swipe_type = 'like' THEN 1 END) as likes_count,
			COUNT(CASE WHEN swipe_type = 'dislike' THEN 1 END) as dislikes_count,
			COUNT(CASE WHEN swipe_type = 'super_like' THEN 1 END) as super_likes_count
		FROM swipe 
		WHERE swiper_user_id = $1`

	err := r.pool.QueryRow(context.Background(), query, userID).Scan(&likesCount, &dislikesCount, &superLikesCount)
	if err != nil {
		return 0, 0, 0, err
	}

	return likesCount, dislikesCount, superLikesCount, nil
}

func (r *swipeRepository) GetMutualLikes(userID uuid.UUID) ([]uuid.UUID, error) {
	query := `
		SELECT s1.swiper_user_id 
		FROM swipe s1
		INNER JOIN swipe s2 ON s1.swiper_user_id = s2.target_user_id AND s1.target_user_id = s2.swiper_user_id
		WHERE s1.target_user_id = $1 
		AND s1.swipe_type = 'like' 
		AND s2.swipe_type = 'like'
		AND NOT EXISTS (
			SELECT 1 FROM match m 
			WHERE (m.user1_id = s1.swiper_user_id AND m.user2_id = s1.target_user_id)
			OR (m.user1_id = s1.target_user_id AND m.user2_id = s1.swiper_user_id)
		)`

	rows, err := r.pool.Query(context.Background(), query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var userIDs []uuid.UUID
	for rows.Next() {
		var id uuid.UUID
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		userIDs = append(userIDs, id)
	}

	return userIDs, nil
}
