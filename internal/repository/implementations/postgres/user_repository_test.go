package postgres

import (
	"testing"
	"time"

	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/domain/constants"
	domain "github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/domain/entities"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
	"github.com/pashagolub/pgxmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUserRepository_Create(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	repo := NewUserRepository(mock)

	phone := "+1234567890"
	bio := "Test bio"
	lat := 55.7558
	lon := 37.6173

	user := &domain.User{
		Email:      "test@example.com",
		Phone:      &phone,
		Name:       "Test User",
		Password:   "hashed_password",
		BirthDate:  time.Now().Add(-20 * 365 * 24 * time.Hour),
		Gender:     constants.GenderMale,
		Bio:        &bio,
		Latitude:   &lat,
		Longitude:  &lon,
		IsVerified: false,
	}

	rows := mock.NewRows([]string{"id", "created_at", "updated_at", "last_active"}).
		AddRow(uuid.New(), time.Now(), time.Now(), time.Now())

	mock.ExpectQuery("INSERT INTO \"user\"").
		WithArgs(user.Email, user.Phone, user.Name, user.Password, user.BirthDate, user.Gender, user.Bio, user.Latitude, user.Longitude, user.IsVerified).
		WillReturnRows(rows)

	err = repo.Create(user)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserRepository_GetByID(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	repo := NewUserRepository(mock)

	userID := uuid.New()
	phone := "+1234567890"
	bio := "Test bio"
	lat := 55.7558
	lon := 37.6173

	expectedUser := &domain.User{
		ID:         userID,
		Email:      "test@example.com",
		Phone:      &phone,
		Name:       "Test User",
		Password:   "hashed_password",
		BirthDate:  time.Now().Add(-20 * 365 * 24 * time.Hour),
		Gender:     constants.GenderMale,
		Bio:        &bio,
		Latitude:   &lat,
		Longitude:  &lon,
		IsVerified: true,
		LastActive: time.Now(),
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	rows := mock.NewRows([]string{
		"id", "email", "phone", "name", "password", "birth_date", "gender",
		"bio", "latitude", "longitude", "is_verified", "last_active", "created_at", "updated_at",
	}).AddRow(
		expectedUser.ID, expectedUser.Email, expectedUser.Phone, expectedUser.Name,
		expectedUser.Password, expectedUser.BirthDate, expectedUser.Gender, expectedUser.Bio,
		expectedUser.Latitude, expectedUser.Longitude, expectedUser.IsVerified,
		expectedUser.LastActive, expectedUser.CreatedAt, expectedUser.UpdatedAt,
	)

	mock.ExpectQuery("SELECT \\* FROM \"user\" WHERE id = \\$1").
		WithArgs(userID).
		WillReturnRows(rows)

	user, err := repo.GetByID(userID)
	assert.NoError(t, err)
	assert.Equal(t, expectedUser.ID, user.ID)
	assert.Equal(t, expectedUser.Email, user.Email)
	assert.Equal(t, expectedUser.Phone, user.Phone)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserRepository_GetByID_NotFound(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	repo := NewUserRepository(mock)

	userID := uuid.New()

	mock.ExpectQuery("SELECT \\* FROM \"user\" WHERE id = \\$1").
		WithArgs(userID).
		WillReturnError(pgx.ErrNoRows)

	user, err := repo.GetByID(userID)
	assert.NoError(t, err)
	assert.Nil(t, user)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserRepository_Update(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	repo := NewUserRepository(mock)

	phone := "+0987654321"
	bio := "Updated bio"
	lat := 59.9343
	lon := 30.3351

	user := &domain.User{
		ID:         uuid.New(),
		Email:      "updated@example.com",
		Phone:      &phone,
		Name:       "Updated User",
		Password:   "new_hashed_password",
		BirthDate:  time.Now().Add(-25 * 365 * 24 * time.Hour),
		Gender:     constants.GenderFemale,
		Bio:        &bio,
		Latitude:   &lat,
		Longitude:  &lon,
		IsVerified: true,
	}

	rows := mock.NewRows([]string{"updated_at"}).AddRow(time.Now())

	mock.ExpectQuery("UPDATE \"user\"").
		WithArgs(
			user.Email, user.Phone, user.Name, user.Password, user.BirthDate,
			user.Gender, user.Bio, user.Latitude, user.Longitude, user.IsVerified, user.ID,
		).
		WillReturnRows(rows)

	err = repo.Update(user)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserRepository_Delete(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	repo := NewUserRepository(mock)

	userID := uuid.New()

	mock.ExpectExec("DELETE FROM \"user\" WHERE id = \\$1").
		WithArgs(userID).
		WillReturnResult(pgxmock.NewResult("DELETE", 1))

	err = repo.Delete(userID)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserRepository_GetUsersByIDs(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	repo := NewUserRepository(mock)

	userIDs := []uuid.UUID{uuid.New(), uuid.New()}
	phone := "+1234567890"
	bio := "Test bio"
	lat := 55.7558
	lon := 37.6173

	expectedUsers := []domain.User{
		{
			ID:         userIDs[0],
			Email:      "user1@example.com",
			Phone:      &phone,
			Name:       "User 1",
			Password:   "password1",
			BirthDate:  time.Now().Add(-20 * 365 * 24 * time.Hour),
			Gender:     constants.GenderMale,
			Bio:        &bio,
			Latitude:   &lat,
			Longitude:  &lon,
			IsVerified: true,
			LastActive: time.Now(),
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		},
		{
			ID:         userIDs[1],
			Email:      "user2@example.com",
			Phone:      &phone,
			Name:       "User 2",
			Password:   "password2",
			BirthDate:  time.Now().Add(-25 * 365 * 24 * time.Hour),
			Gender:     constants.GenderFemale,
			Bio:        &bio,
			Latitude:   &lat,
			Longitude:  &lon,
			IsVerified: true,
			LastActive: time.Now(),
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		},
	}

	rows := mock.NewRows([]string{
		"id", "email", "phone", "name", "password", "birth_date", "gender",
		"bio", "latitude", "longitude", "is_verified", "last_active", "created_at", "updated_at",
	})
	for _, user := range expectedUsers {
		rows.AddRow(
			user.ID, user.Email, user.Phone, user.Name, user.Password, user.BirthDate,
			user.Gender, user.Bio, user.Latitude, user.Longitude, user.IsVerified,
			user.LastActive, user.CreatedAt, user.UpdatedAt,
		)
	}

	mock.ExpectQuery("SELECT \\* FROM \"user\" WHERE id IN").
		WithArgs(userIDs[0], userIDs[1]).
		WillReturnRows(rows)

	users, err := repo.GetUsersByIDs(userIDs)
	assert.NoError(t, err)
	assert.Len(t, users, 2)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserRepository_GetUsersByIDs_Empty(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	repo := NewUserRepository(mock)

	users, err := repo.GetUsersByIDs([]uuid.UUID{})
	assert.NoError(t, err)
	assert.Empty(t, users)
}
