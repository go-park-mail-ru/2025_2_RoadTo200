package postgres

import (
	"context"
	"errors"

	domain "github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/core-service/domain/entities"
	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/core-service/repository/interfaces"
	bdIface "github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/repository/interfaces"
	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/utils"
	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/pkg/logger"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

var _ interfaces.SwipeRepository = (*SwipeRepository)(nil)

type SwipeRepository struct {
	pool   bdIface.PgxIface
	logger logger.Log
}

func NewSwipeRepository(pool bdIface.PgxIface) *SwipeRepository {
	return &SwipeRepository{pool: pool}
}

func (r *SwipeRepository) Create(ctx context.Context, swipe *domain.Swipe) error {
	conn, err := utils.GetConn(ctx, r.pool)
	if err != nil {
		return err
	}

	query := `
		INSERT INTO swipe (swiper_user_id, target_user_id, swipe_type)
		VALUES ($1, $2, $3)
		RETURNING created_at`

	err = conn.QueryRow(ctx, query,
		swipe.SwiperUserID, swipe.TargetUserID, swipe.SwipeType).
		Scan(&swipe.CreatedAt)

	if err != nil {
		return err
	}
	return nil
}

func (r *SwipeRepository) GetBySwiperAndTarget(ctx context.Context, swiperID, targetID uuid.UUID) (*domain.Swipe, error) {
	conn, err := utils.GetConn(ctx, r.pool)
	if err != nil {
		return nil, err
	}

	var swipe domain.Swipe
	query := `SELECT * FROM swipe WHERE swiper_user_id = $1 AND target_user_id = $2`

	err = conn.QueryRow(ctx, query, swiperID, targetID).Scan(
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

func (r *SwipeRepository) GetSwipesBySwiper(ctx context.Context, swiperID uuid.UUID, limit, offset int) ([]domain.Swipe, error) {
	conn, err := utils.GetConn(ctx, r.pool)
	if err != nil {
		return nil, err
	}

	query := `
		SELECT * FROM swipe 
		WHERE swiper_user_id = $1 
		ORDER BY created_at DESC 
		LIMIT $2 OFFSET $3`

	rows, err := conn.Query(ctx, query, swiperID, limit, offset)
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

func (r *SwipeRepository) GetSwipesByTarget(ctx context.Context, targetID uuid.UUID, limit, offset int) ([]domain.Swipe, error) {
	conn, err := utils.GetConn(ctx, r.pool)
	if err != nil {
		return nil, err
	}

	query := `
		SELECT * FROM swipe 
		WHERE target_user_id = $1 
		ORDER BY created_at DESC 
		LIMIT $2 OFFSET $3`

	rows, err := conn.Query(ctx, query, targetID, limit, offset)
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

func (r *SwipeRepository) Exists(ctx context.Context, swiperID, targetID uuid.UUID) (bool, error) {
	conn, err := utils.GetConn(ctx, r.pool)
	if err != nil {
		return false, err
	}

	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM swipe WHERE swiper_user_id = $1 AND target_user_id = $2)`

	err = conn.QueryRow(ctx, query, swiperID, targetID).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

func (r *SwipeRepository) GetSwipesStats(ctx context.Context, userID uuid.UUID) (int, int, int, error) {
	conn, err := utils.GetConn(ctx, r.pool)
	if err != nil {
		return 0, 0, 0, err
	}

	var likesCount, dislikesCount, superLikesCount int

	query := `
		SELECT 
			COUNT(CASE WHEN swipe_type = 'like' THEN 1 END) as likes_count,
			COUNT(CASE WHEN swipe_type = 'dislike' THEN 1 END) as dislikes_count,
			COUNT(CASE WHEN swipe_type = 'super_like' THEN 1 END) as super_likes_count
		FROM swipe 
		WHERE swiper_user_id = $1`

	err = conn.QueryRow(ctx, query, userID).Scan(&likesCount, &dislikesCount, &superLikesCount)
	if err != nil {
		return 0, 0, 0, err
	}

	return likesCount, dislikesCount, superLikesCount, nil
}

func (r *SwipeRepository) GetMutualLikes(ctx context.Context, userID uuid.UUID) ([]uuid.UUID, error) {
	conn, err := utils.GetConn(ctx, r.pool)
	if err != nil {
		return nil, err
	}

	query := `
		SELECT s1.swiper_user_id FROM swipe s1
		INNER JOIN swipe s2 ON s1.swiper_user_id = s2.target_user_id AND s1.target_user_id = s2.swiper_user_id
					AND s1.swipe_type = 'like' AND s2.swipe_type = 'like'
		WHERE s1.target_user_id = $1
	`

	rows, err := conn.Query(ctx, query, userID)
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

func (r *SwipeRepository) GetUsersForFeed(ctx context.Context, userID uuid.UUID, limit, offset int) ([]uuid.UUID, error) {
	r.logger.Tracef("GetUsersForFeed called with userID:", userID, "limit:", limit, "offset:", offset)
	conn, err := utils.GetConn(ctx, r.pool)
	if err != nil {
		r.logger.Errorf("get DB context error: %s", err.Error())
		return nil, err
	}

	query := `
			SELECT i1.user_id FROM interest i
			INNER JOIN interest i1 ON i1.theme = i.theme AND i1.user_id <> i.user_id
			LEFT JOIN swipe s ON i.user_id = s.swiper_user_id AND i1.user_id = s.target_user_id
			WHERE s.id IS NULL
			LIMIT $2 OFFSET $3`

	rows, err := conn.Query(ctx, query, userID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	r.logger.Debugf("Getting users for feed with userID: %v", userID)
	var users []uuid.UUID
	for rows.Next() {
		r.logger.Trace("Scanning a user row in GetUsersForFeed")
		var user uuid.UUID
		err := rows.Scan(&user)
		r.logger.Debugf("Scanned user: %s", user)
		if err != nil {
			r.logger.Errorf("Error while scanning user rows: %s", err)
			return nil, err
		}
		users = append(users, user)
	}
	r.logger.Trace("Finished scanning users in GetUsersForFeed")
	return users, nil
}
