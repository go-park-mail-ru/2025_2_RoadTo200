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

var _ interfaces.UserPhotoRepository = (*UserPhotoRepository)(nil)

type UserPhotoRepository struct {
	pool   bdIface.PgxIface
	logger logger.Log
}

func NewUserPhotoRepository(pool bdIface.PgxIface, l logger.Log) *UserPhotoRepository {
	return &UserPhotoRepository{pool: pool, logger: l}
}

func (r *UserPhotoRepository) Create(ctx context.Context, photo *domain.UserPhoto) error {
	r.logger.Trace("UserPhotoRepository.Create")
	conn, err := utils.GetConn(ctx, r.pool)
	if err != nil {
		return err
	}

	query := `
		INSERT INTO user_photo (user_id, photo_url, display_order, is_approved)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at`

	err = conn.QueryRow(ctx, query,
		photo.UserID, photo.PhotoURL, photo.DisplayOrder, photo.IsApproved).
		Scan(&photo.ID, &photo.CreatedAt)

	if err != nil {
		return err
	}
	return nil
}

func (r *UserPhotoRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.UserPhoto, error) {
	r.logger.Trace("UserPhotoRepository.GetByID")
	conn, err := utils.GetConn(ctx, r.pool)
	if err != nil {
		return nil, err
	}

	var photo domain.UserPhoto
	query := `SELECT * FROM user_photo WHERE id = $1`

	err = conn.QueryRow(ctx, query, id).Scan(
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
	conn, err := utils.GetConn(ctx, r.pool)
	if err != nil {
		return nil, err
	}

	query := `SELECT * FROM user_photo WHERE user_id = $1 ORDER BY display_order`

	rows, err := conn.Query(ctx, query, userID)
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
	conn, err := utils.GetConn(ctx, r.pool)
	if err != nil {
		return err
	}

	query := `
		UPDATE user_photo 
		SET photo_url = $1, display_order = $2, is_approved = $3
		WHERE id = $4`

	_, err = conn.Exec(ctx, query,
		photo.PhotoURL, photo.DisplayOrder, photo.IsApproved, photo.ID)

	if err != nil {
		return err
	}
	return nil
}

func (r *UserPhotoRepository) Delete(ctx context.Context, id uuid.UUID) error {
	r.logger.Trace("UserPhotoRepository.Delete")
	conn, err := utils.GetConn(ctx, r.pool)
	if err != nil {
		return err
	}

	query := `DELETE FROM user_photo WHERE id = $1`

	_, err = conn.Exec(ctx, query, id)
	if err != nil {
		return err
	}
	return nil
}

func (r *UserPhotoRepository) UpdateDisplayOrder(ctx context.Context, userID uuid.UUID, photos []domain.UserPhoto) error {
	r.logger.Trace("UserPhotoRepository.UpdateDisplayOrder")
	conn, err := utils.GetConn(ctx, r.pool)
	if err != nil {
		return err
	}

	tx, err := conn.Begin(ctx)
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
