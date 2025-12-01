package postgres

import (
	"context"
	"errors"
	"time"

	domain "github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/domain/entities"
	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/repository/interfaces"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

var _ interfaces.SubscriptionRepository = (*SubscriptionRepository)(nil)

type SubscriptionRepository struct {
	pool interfaces.PgxIface
}

func NewSubscriptionRepository(pool interfaces.PgxIface) *SubscriptionRepository {
	return &SubscriptionRepository{pool: pool}
}

func (r *SubscriptionRepository) Create(ctx context.Context, subscription *domain.Subscription) error {
	conn, err := GetConn(ctx, r.pool)
	if err != nil {
		return err
	}

	query := `
		INSERT INTO subscription (user_id, plan_type, start_date, end_date, is_active)
		VALUES ($1, $2, $3, $4, $5)`

	_, err = conn.Exec(ctx, query,
		subscription.UserID, subscription.PlanType, subscription.StartDate, subscription.EndDate, subscription.IsActive)

	return err
}

func (r *SubscriptionRepository) GetByUserID(ctx context.Context, userID uuid.UUID) (*domain.Subscription, error) {
	conn, err := GetConn(ctx, r.pool)
	if err != nil {
		return nil, err
	}

	var subscription domain.Subscription
	query := `SELECT * FROM subscription WHERE user_id = $1`

	err = conn.QueryRow(ctx, query, userID).Scan(
		&subscription.UserID, &subscription.PlanType, &subscription.StartDate,
		&subscription.EndDate, &subscription.IsActive, &subscription.CreatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &subscription, nil
}

func (r *SubscriptionRepository) Update(ctx context.Context, subscription *domain.Subscription) error {
	conn, err := GetConn(ctx, r.pool)
	if err != nil {
		return err
	}

	query := `
		UPDATE subscription 
		SET plan_type = $1, start_date = $2, end_date = $3, is_active = $4
		WHERE user_id = $5`

	_, err = conn.Exec(ctx, query,
		subscription.PlanType, subscription.StartDate, subscription.EndDate,
		subscription.IsActive, subscription.UserID)

	if err != nil {
		return err
	}
	return nil
}

func (r *SubscriptionRepository) Delete(ctx context.Context, userID uuid.UUID) error {
	conn, err := GetConn(ctx, r.pool)
	if err != nil {
		return err
	}

	query := `DELETE FROM subscription WHERE user_id = $1`

	_, err = conn.Exec(ctx, query, userID)
	if err != nil {
		return err
	}
	return nil
}

func (r *SubscriptionRepository) GetActiveSubscription(ctx context.Context, userID uuid.UUID) (*domain.Subscription, error) {
	conn, err := GetConn(ctx, r.pool)
	if err != nil {
		return nil, err
	}

	var subscription domain.Subscription
	query := `
		SELECT * FROM subscription 
		WHERE user_id = $1 AND is_active = true AND end_date > $2`

	err = conn.QueryRow(ctx, query, userID, time.Now()).Scan(
		&subscription.UserID, &subscription.PlanType, &subscription.StartDate,
		&subscription.EndDate, &subscription.IsActive, &subscription.CreatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &subscription, nil
}
