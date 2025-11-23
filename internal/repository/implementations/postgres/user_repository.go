package postgres

import (
	"context"
	"errors"
	"fmt"
	"strings"

	domain "github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/domain/entities"
	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/logger"
	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/repository/interfaces"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

var _ interfaces.UserRepository = (*UserRepository)(nil)

type UserRepository struct {
	pool   interfaces.PgxIface
	logger logger.Log
}

func NewUserRepository(pool interfaces.PgxIface, l logger.Log) *UserRepository {
	return &UserRepository{pool: pool, logger: l}
}

func (r *UserRepository) Create(ctx context.Context, user *domain.User) error {
	r.logger.Trace("FeedService.Create")
	query := `
	INSERT INTO "user" (email, phone, name, password, birth_date, gender, bio, 
					   city, artist, quote, is_verified)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
	RETURNING id, created_at, updated_at, last_active`

	err := r.pool.QueryRow(ctx, query,
		user.Email, user.Phone, user.Name, user.Password, user.BirthDate, user.Gender, user.Bio,
		user.City, user.Artist, user.Quote, user.IsVerified).
		Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt, &user.LastActive)

	if err != nil {
		return err
	}
	return nil
}

func (r *UserRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	r.logger.Trace("FeedService.GetByID")
	var user domain.User
	query := `
        SELECT id, email, phone, name, password, birth_date, gender, bio, 
               city, artist, quote, is_verified, last_active, created_at, updated_at 
        FROM "user" WHERE id = $1`

	err := r.pool.QueryRow(ctx, query, id).Scan(
		&user.ID, &user.Email, &user.Phone, &user.Name, &user.Password,
		&user.BirthDate, &user.Gender, &user.Bio,
		&user.City, &user.Artist, &user.Quote, &user.IsVerified, &user.LastActive,
		&user.CreatedAt, &user.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	r.logger.Trace("FeedService.GetByEmail")
	var user domain.User
	query := `
        SELECT id, email, phone, name, password, birth_date, gender, bio, 
               city, artist, quote, is_verified, last_active, created_at, updated_at 
        FROM "user" WHERE email = $1`

	err := r.pool.QueryRow(ctx, query, email).Scan(
		&user.ID, &user.Email, &user.Phone, &user.Name, &user.Password,
		&user.BirthDate, &user.Gender, &user.Bio,
		&user.City, &user.Artist, &user.Quote, &user.IsVerified, &user.LastActive,
		&user.CreatedAt, &user.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) GetByPhone(ctx context.Context, phone string) (*domain.User, error) {
	r.logger.Trace("FeedService.GetByPhone")
	var user domain.User
	query := `SELECT * FROM "user" WHERE phone = $1`

	err := r.pool.QueryRow(ctx, query, phone).Scan(
		&user.ID, &user.Email, &user.Phone, &user.Name, &user.Password,
		&user.BirthDate, &user.Gender, &user.Bio, &user.City, &user.Artist, &user.Quote,
		&user.IsVerified, &user.LastActive, &user.CreatedAt, &user.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) Update(ctx context.Context, user *domain.User) error {
	r.logger.Trace("FeedService.Update")
	query := `
        UPDATE "user" 
        SET email = $1, phone = $2, name = $3, password = $4, birth_date = $5, 
            gender = $6, bio = $7, city = $8, artist = $9, quote = $10, is_verified = $11,
            updated_at = NOW()
        WHERE id = $12
        RETURNING updated_at`

	err := r.pool.QueryRow(ctx, query,
		user.Email, user.Phone, user.Name, user.Password, user.BirthDate, user.Gender,
		user.Bio, user.City, user.Artist, user.Quote, user.IsVerified,
		user.ID).
		Scan(&user.UpdatedAt)

	if err != nil {
		return err
	}
	return nil
}

func (r *UserRepository) UpdateLastActive(ctx context.Context, userID uuid.UUID) error {
	r.logger.Trace("FeedService.UpdateLastActive")
	query := `UPDATE "user" SET last_active = NOW() WHERE id = $1`

	_, err := r.pool.Exec(ctx, query, userID)
	if err != nil {
		return err
	}
	return nil
}

func (r *UserRepository) Delete(ctx context.Context, id uuid.UUID) error {
	r.logger.Trace("FeedService.Delete")
	query := `DELETE FROM "user" WHERE id = $1`

	_, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		return err
	}
	return nil
}

func (r *UserRepository) GetUsersByIDs(ctx context.Context, ids []uuid.UUID) ([]domain.User, error) {
	r.logger.Trace("FeedService.GetUsersByIDs")
	if len(ids) == 0 {
		return []domain.User{}, nil
	}

	placeholders := make([]string, len(ids))
	args := make([]interface{}, len(ids))
	for i, id := range ids {
		placeholders[i] = fmt.Sprintf("$%d", i+1)
		args[i] = id
	}

	query := fmt.Sprintf(`
        SELECT id, email, phone, name, password, birth_date, gender, bio, 
               city, artist, quote, is_verified, last_active, created_at, updated_at 
        FROM "user" 
        WHERE id IN (%s)
        ORDER BY created_at DESC`, strings.Join(placeholders, ","))

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []domain.User
	for rows.Next() {
		var user domain.User
		err := rows.Scan(
			&user.ID, &user.Email, &user.Phone, &user.Name, &user.Password,
			&user.BirthDate, &user.Gender, &user.Bio,
			&user.City, &user.Artist, &user.Quote, &user.IsVerified, &user.LastActive,
			&user.CreatedAt, &user.UpdatedAt,
		)
		if err != nil {
			r.logger.Errorf("Error while scanning user rows: %v", err)
			return nil, err
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

func (r *UserRepository) GetUsersForFeed(ctx context.Context, userID uuid.UUID, limit, offset int) ([]domain.User, error) {
	r.logger.Tracef("GetUsersForFeed called with userID:", userID, "limit:", limit, "offset:", offset)

	query := `
        SELECT id, email, phone, name, password, birth_date, gender, bio, 
               city, artist, quote, is_verified, last_active, created_at, updated_at
        FROM "user" u
        WHERE u.id != $1
          AND NOT EXISTS (
              SELECT 1 FROM swipe s
              WHERE s.swiper_user_id = $1 AND s.target_user_id = u.id
          )
        ORDER BY u.last_active DESC
        LIMIT $2 OFFSET $3`

	rows, err := r.pool.Query(ctx, query, userID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	r.logger.Debugf("Getting users for feed with userID: %v", userID)
	var users []domain.User
	for rows.Next() {
		r.logger.Trace("Scanning a user row in GetUsersForFeed")
		var user domain.User
		err := rows.Scan(
			&user.ID, &user.Email, &user.Phone, &user.Name, &user.Password,
			&user.BirthDate, &user.Gender, &user.Bio,
			&user.City, &user.Artist, &user.Quote, &user.IsVerified, &user.LastActive,
			&user.CreatedAt, &user.UpdatedAt,
		)
		r.logger.Debugf("Scanned user: %+v\n", user)
		if err != nil {
			r.logger.Errorf("Error while scanning user rows: %v", err)
			return nil, err
		}
		users = append(users, user)
	}
	r.logger.Trace("Finished scanning users in GetUsersForFeed")
	return users, nil
}
