package postgres

import (
	"context"
	"errors"
	"time"

	domain "github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/domain/entities"
	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/repository/interfaces"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
)

var _ interfaces.SubscriptionRepository = (*subscriptionRepository)(nil)

type subscriptionRepository struct {
	pool interfaces.PgxIface
}

func NewSubscriptionRepository(pool interfaces.PgxIface) interfaces.SubscriptionRepository {
	return &subscriptionRepository{pool: pool}
}

func (r *subscriptionRepository) Create(ctx context.Context, subscription *domain.Subscription) error {
	query := `
		INSERT INTO subscription (user_id, plan_type, start_date, end_date, is_active)
		VALUES ($1, $2, $3, $4, $5)`

	_, err := r.pool.Exec(ctx, query,
		subscription.UserID, subscription.PlanType, subscription.StartDate, subscription.EndDate, subscription.IsActive)

	return err
}

func (r *subscriptionRepository) GetByUserID(ctx context.Context, userID uuid.UUID) (*domain.Subscription, error) {
	var subscription domain.Subscription
	query := `SELECT * FROM subscription WHERE user_id = $1`

	err := r.pool.QueryRow(ctx, query, userID).Scan(
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

func (r *subscriptionRepository) Update(ctx context.Context, subscription *domain.Subscription) error {
	query := `
		UPDATE subscription 
		SET plan_type = $1, start_date = $2, end_date = $3, is_active = $4
		WHERE user_id = $5`

	_, err := r.pool.Exec(ctx, query,
		subscription.PlanType, subscription.StartDate, subscription.EndDate,
		subscription.IsActive, subscription.UserID)

	if err != nil {
		return err
	}
	return nil
}

func (r *subscriptionRepository) Delete(ctx context.Context, userID uuid.UUID) error {
	query := `DELETE FROM subscription WHERE user_id = $1`

	_, err := r.pool.Exec(ctx, query, userID)
	if err != nil {
		return err
	}
	return nil
}

func (r *subscriptionRepository) GetActiveSubscription(ctx context.Context, userID uuid.UUID) (*domain.Subscription, error) {
	var subscription domain.Subscription
	query := `
		SELECT * FROM subscription 
		WHERE user_id = $1 AND is_active = true AND end_date > $2`

	err := r.pool.QueryRow(ctx, query, userID, time.Now()).Scan(
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
