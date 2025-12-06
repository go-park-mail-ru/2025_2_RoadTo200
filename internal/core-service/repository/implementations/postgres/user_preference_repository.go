package postgres

import (
	"context"
	"errors"

	domain "github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/core-service/domain/entities"
	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/core-service/repository/interfaces"
	utils "github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/repository/implementations/postgres"
	bdIface "github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/repository/interfaces"
	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/pkg/logger"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

var _ interfaces.UserPreferenceRepository = (*UserPreferenceRepository)(nil)

type UserPreferenceRepository struct {
	pool   bdIface.PgxIface
	logger logger.Log
}

func NewUserPreferenceRepository(pool bdIface.PgxIface, l logger.Log) *UserPreferenceRepository {
	return &UserPreferenceRepository{pool: pool, logger: l}
}

func (r *UserPreferenceRepository) Create(ctx context.Context, preference *domain.UserPreference) error {
	r.logger.Trace("UserPreferenceRepository.Create")
	conn, err := utils.GetConn(ctx, r.pool)
	if err != nil {
		return err
	}

	query := `
		INSERT INTO user_preference (user_id, show_gender, age_min, age_max, max_distance, global_search)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING created_at, updated_at`

	err = conn.QueryRow(ctx, query,
		preference.UserID, preference.ShowGender, preference.AgeMin, preference.AgeMax,
		preference.MaxDistance, preference.GlobalSearch).
		Scan(&preference.CreatedAt, &preference.UpdatedAt)

	if err != nil {
		return err
	}
	return nil
}

func (r *UserPreferenceRepository) GetByUserID(ctx context.Context, userID uuid.UUID) (*domain.UserPreference, error) {
	r.logger.Trace("UserPreferenceRepository.GetByUserID")
	conn, err := utils.GetConn(ctx, r.pool)
	if err != nil {
		return nil, err
	}

	var preference domain.UserPreference
	query := `SELECT * FROM user_preference WHERE user_id = $1`

	err = conn.QueryRow(ctx, query, userID).Scan(
		&preference.UserID, &preference.ShowGender, &preference.AgeMin, &preference.AgeMax,
		&preference.MaxDistance, &preference.GlobalSearch, &preference.CreatedAt, &preference.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &preference, nil
}

func (r *UserPreferenceRepository) Update(ctx context.Context, preference *domain.UserPreference) error {
	r.logger.Trace("UserPreferenceRepository.Update")
	conn, err := utils.GetConn(ctx, r.pool)
	if err != nil {
		return err
	}

	query := `
		UPDATE user_preference 
		SET show_gender = $1, age_min = $2, age_max = $3, max_distance = $4, global_search = $5, updated_at = NOW()
		WHERE user_id = $6
		RETURNING updated_at`

	err = conn.QueryRow(ctx, query,
		preference.ShowGender, preference.AgeMin, preference.AgeMax, preference.MaxDistance,
		preference.GlobalSearch, preference.UserID).
		Scan(&preference.UpdatedAt)

	if err != nil {
		return err
	}
	return nil
}

func (r *UserPreferenceRepository) Delete(ctx context.Context, userID uuid.UUID) error {
	r.logger.Trace("UserPreferenceRepository.Delete")
	conn, err := utils.GetConn(ctx, r.pool)
	if err != nil {
		return err
	}

	query := `DELETE FROM user_preference WHERE user_id = $1`

	_, err = conn.Exec(ctx, query, userID)
	if err != nil {
		return err
	}
	return nil
}

func (r *UserPreferenceRepository) GetInterests(ctx context.Context, userID uuid.UUID) ([]domain.Interest, error) {
	r.logger.Trace("UserPreferenceRepository.GetInterests")
	conn, err := utils.GetConn(ctx, r.pool)
	if err != nil {
		return nil, err
	}

	query := `
		SELECT i.user_id, i.theme FROM interest i WHERE i.user_id = $1`
	interests := make([]domain.Interest, 0)

	res, err := conn.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer res.Close()

	for res.Next() {
		var el domain.Interest
		err = res.Scan(&el.UserID, &el.Theme)
		if err != nil {
			r.logger.Errorf("Error while getting interests: %v", err)
			continue
		}
		interests = append(interests, el)
	}
	return interests, nil
}

func (r *UserPreferenceRepository) UpdateInterests(ctx context.Context, userID uuid.UUID, inter []domain.Interest) error {
	r.logger.Trace("UserPreferenceRepository.UpdateInterests")
	conn, err := utils.GetConn(ctx, r.pool)
	if err != nil {
		return err
	}

	delQuery := `
		DELETE FROM interest WHERE user_id = $1`
	updQuery := `
		INSERT INTO interest(user_id, theme) VALUES ($1, $2)`

	tx, err := conn.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	_, err = tx.Exec(ctx, delQuery, userID)
	if err != nil {
		return err
	}
	batch := pgx.Batch{}
	for _, el := range inter {
		batch.Queue(updQuery, el.UserID, el.Theme)
	}
	res := tx.SendBatch(ctx, &batch)
	err = res.Close()
	if err != nil {
		return err
	}
	return tx.Commit(ctx)
}
