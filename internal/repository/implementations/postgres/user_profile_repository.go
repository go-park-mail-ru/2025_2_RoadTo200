package postgres

import (
	"context"
	"errors"

	domain "github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/domain/entities"
	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/logger"
	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/repository/interfaces"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
)

var _ interfaces.UserPhotoRepository = (*UserPhotoRepository)(nil)
var _ interfaces.UserPreferenceRepository = (*UserPreferenceRepository)(nil)

type UserPhotoRepository struct {
	pool   interfaces.PgxIface
	logger logger.Log
}

func NewUserPhotoRepository(pool interfaces.PgxIface, l logger.Log) *UserPhotoRepository {
	return &UserPhotoRepository{pool: pool, logger: l}
}

func (r *UserPhotoRepository) Create(ctx context.Context, photo *domain.UserPhoto) error {
	r.logger.Trace("UserPhotoRepository.Create")

	query := `
		INSERT INTO user_photo (user_id, photo_url, display_order, is_approved)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at`

	err := r.pool.QueryRow(ctx, query,
		photo.UserID, photo.PhotoURL, photo.DisplayOrder, photo.IsApproved).
		Scan(&photo.ID, &photo.CreatedAt)

	if err != nil {
		return err
	}
	return nil
}

func (r *UserPhotoRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.UserPhoto, error) {
	r.logger.Trace("UserPhotoRepository.GetByID")

	var photo domain.UserPhoto
	query := `SELECT * FROM user_photo WHERE id = $1`

	err := r.pool.QueryRow(ctx, query, id).Scan(
		&photo.ID, &photo.UserID, &photo.PhotoURL, &photo.DisplayOrder,
		&photo.IsApproved, &photo.CreatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &photo, nil
}

func (r *UserPhotoRepository) GetByUserID(ctx context.Context, userID uuid.UUID) ([]domain.UserPhoto, error) {
	r.logger.Trace("UserPhotoRepository.GetByUserID")

	query := `SELECT * FROM user_photo WHERE user_id = $1 ORDER BY display_order`

	rows, err := r.pool.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var photos []domain.UserPhoto
	for rows.Next() {
		var photo domain.UserPhoto
		err := rows.Scan(
			&photo.ID, &photo.UserID, &photo.PhotoURL, &photo.DisplayOrder,
			&photo.IsApproved, &photo.CreatedAt,
		)
		if err != nil {
			r.logger.Errorf("Error while scanning user photo rows: %v", err)
			continue
		}
		photos = append(photos, photo)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return photos, nil
}

func (r *UserPhotoRepository) Update(ctx context.Context, photo *domain.UserPhoto) error {
	r.logger.Trace("UserPhotoRepository.Update")

	query := `
		UPDATE user_photo 
		SET photo_url = $1, display_order = $2, is_approved = $3
		WHERE id = $4`

	_, err := r.pool.Exec(ctx, query,
		photo.PhotoURL, photo.DisplayOrder, photo.IsApproved, photo.ID)

	if err != nil {
		return err
	}
	return nil
}

func (r *UserPhotoRepository) Delete(ctx context.Context, id uuid.UUID) error {
	r.logger.Trace("UserPhotoRepository.Delete")

	query := `DELETE FROM user_photo WHERE id = $1`

	_, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		return err
	}
	return nil
}

func (r *UserPhotoRepository) UpdateDisplayOrder(ctx context.Context, userID uuid.UUID, photos []domain.UserPhoto) error {
	r.logger.Trace("UserPhotoRepository.UpdateDisplayOrder")

	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	// Delete existing photos for user
	_, err = tx.Exec(ctx, "DELETE FROM user_photo WHERE user_id = $1", userID)
	if err != nil {
		return err
	}

	// Insert new photos
	for _, photo := range photos {
		_, err = tx.Exec(ctx, `
			INSERT INTO user_photo (user_id, photo_url, display_order, is_approved)
			VALUES ($1, $2, $3, $4)`,
			userID, photo.PhotoURL, photo.DisplayOrder, photo.IsApproved)
		if err != nil {
			r.logger.Errorf("Error while updating photo for user %v: %v", userID, err)
			return err
		}
	}

	return tx.Commit(ctx)
}

type UserPreferenceRepository struct {
	pool   interfaces.PgxIface
	logger logger.Log
}

func NewUserPreferenceRepository(pool interfaces.PgxIface, l logger.Log) *UserPreferenceRepository {
	return &UserPreferenceRepository{pool: pool, logger: l}
}

func (r *UserPreferenceRepository) Create(ctx context.Context, preference *domain.UserPreference) error {
	r.logger.Trace("UserPreferenceRepository.Create")

	query := `
		INSERT INTO user_preference (user_id, show_gender, age_min, age_max, max_distance, global_search)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING created_at, updated_at`

	err := r.pool.QueryRow(ctx, query,
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

	var preference domain.UserPreference
	query := `SELECT * FROM user_preference WHERE user_id = $1`

	err := r.pool.QueryRow(ctx, query, userID).Scan(
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

	query := `
		UPDATE user_preference 
		SET show_gender = $1, age_min = $2, age_max = $3, max_distance = $4, global_search = $5, updated_at = NOW()
		WHERE user_id = $6
		RETURNING updated_at`

	err := r.pool.QueryRow(ctx, query,
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

	query := `DELETE FROM user_preference WHERE user_id = $1`

	_, err := r.pool.Exec(ctx, query, userID)
	if err != nil {
		return err
	}
	return nil
}

func (r *UserPreferenceRepository) GetInterests(ctx context.Context, userID uuid.UUID) ([]domain.Interest, error) {
	r.logger.Trace("UserPreferenceRepository.GetInterests")

	query := `
		SELECT i.user_id, i.theme FROM interest i WHERE i.user_id = $1`
	interests := make([]domain.Interest, 0)

	res, err := r.pool.Query(ctx, query, userID)
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

	delQuery := `
		DELETE FROM interest WHERE user_id = $1`
	updQuery := `
		INSERT INTO interest(user_id, theme) VALUES ($1, $2)`

	tx, err := r.pool.Begin(ctx)
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
	defer res.Close()

	_, err = res.Exec()
	if err != nil {
		return err
	}
	return tx.Commit(ctx)
}
