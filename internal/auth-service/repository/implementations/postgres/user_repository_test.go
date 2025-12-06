package postgres

import (
	"context"
	"testing"
	"time"

	domain "github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/auth-service/domain/entities"
	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/domain/constants"
	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/tests/mocks"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/pashagolub/pgxmock/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUserRepository_Create(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	log := mocks.NewMockLogger()
	repo := NewUserRepository(mock, log)

	phone := "+1234567890"
	bio := "Test bio"
	city := "Moscow"
	artist := "Test Artist"
	quote := "Test Quote"

	user := &domain.User{
		Email:      "test@example.com",
		Phone:      &phone,
		Name:       "Test User",
		Password:   "hashed_password",
		BirthDate:  time.Now().Add(-20 * 365 * 24 * time.Hour),
		Gender:     constants.GenderMale,
		Bio:        &bio,
		City:       &city,
		Artist:     &artist,
		Quote:      &quote,
		IsVerified: false,
	}

	rows := mock.NewRows([]string{"id", "created_at", "updated_at", "last_active"}).
		AddRow(uuid.New(), time.Now(), time.Now(), time.Now())

	mock.ExpectQuery("INSERT INTO \"user\"").
		WithArgs(
			user.Email, user.Phone, user.Name, user.Password, user.BirthDate,
			user.Gender, user.Bio, user.City, user.Artist, user.Quote, user.IsVerified,
		).
		WillReturnRows(rows)

	err = repo.Create(context.Background(), user)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
	assert.True(t, log.WasCalled("Trace"))
}

func TestUserRepository_Create_Error(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	log := mocks.NewMockLogger()
	repo := NewUserRepository(mock, log)

	phone := "+1234567890"
	bio := "Test bio"
	city := "Moscow"
	artist := "Test Artist"
	quote := "Test Quote"

	user := &domain.User{
		Email:      "test@example.com",
		Phone:      &phone,
		Name:       "Test User",
		Password:   "hashed_password",
		BirthDate:  time.Now().Add(-20 * 365 * 24 * time.Hour),
		Gender:     constants.GenderMale,
		Bio:        &bio,
		City:       &city,
		Artist:     &artist,
		Quote:      &quote,
		IsVerified: false,
	}

	mock.ExpectQuery("INSERT INTO \"user\"").
		WithArgs(
			user.Email, user.Phone, user.Name, user.Password, user.BirthDate,
			user.Gender, user.Bio, user.City, user.Artist, user.Quote, user.IsVerified,
		).
		WillReturnError(pgx.ErrTxClosed)

	err = repo.Create(context.Background(), user)
	assert.Error(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserRepository_GetByID(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	log := mocks.NewMockLogger()
	repo := NewUserRepository(mock, log)

	userID := uuid.New()
	phone := "+1234567890"
	bio := "Test bio"
	city := "Moscow"
	artist := "Test Artist"
	quote := "Test Quote"

	expectedUser := &domain.User{
		ID:         userID,
		Email:      "test@example.com",
		Phone:      &phone,
		Name:       "Test User",
		Password:   "hashed_password",
		BirthDate:  time.Now().Add(-20 * 365 * 24 * time.Hour),
		Gender:     constants.GenderMale,
		Bio:        &bio,
		City:       &city,
		Artist:     &artist,
		Quote:      &quote,
		IsVerified: true,
		LastActive: time.Now(),
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	rows := mock.NewRows([]string{
		"id", "email", "phone", "name", "password", "birth_date", "gender",
		"bio", "city", "artist", "quote", "is_verified", "last_active", "created_at", "updated_at",
	}).AddRow(
		expectedUser.ID, expectedUser.Email, expectedUser.Phone, expectedUser.Name,
		expectedUser.Password, expectedUser.BirthDate, expectedUser.Gender, expectedUser.Bio,
		expectedUser.City, expectedUser.Artist, expectedUser.Quote, expectedUser.IsVerified,
		expectedUser.LastActive, expectedUser.CreatedAt, expectedUser.UpdatedAt,
	)

	mock.ExpectQuery("SELECT id, email, phone, name, password, birth_date, gender, bio, city, artist, quote, is_verified, last_active, created_at, updated_at FROM \"user\" WHERE id = \\$1").
		WithArgs(userID).
		WillReturnRows(rows)

	user, err := repo.GetByID(context.Background(), userID)
	assert.NoError(t, err)
	assert.Equal(t, expectedUser.ID, user.ID)
	assert.Equal(t, expectedUser.Email, user.Email)
	assert.Equal(t, expectedUser.Phone, user.Phone)
	assert.Equal(t, expectedUser.City, user.City)
	assert.Equal(t, expectedUser.Artist, user.Artist)
	assert.Equal(t, expectedUser.Quote, user.Quote)
	assert.NoError(t, mock.ExpectationsWereMet())
	assert.True(t, log.WasCalled("Trace"))
}

func TestUserRepository_GetByID_NotFound(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	log := mocks.NewMockLogger()
	repo := NewUserRepository(mock, log)

	userID := uuid.New()

	mock.ExpectQuery("SELECT id, email, phone, name, password, birth_date, gender, bio, city, artist, quote, is_verified, last_active, created_at, updated_at FROM \"user\" WHERE id = \\$1").
		WithArgs(userID).
		WillReturnError(pgx.ErrNoRows)

	user, err := repo.GetByID(context.Background(), userID)
	assert.NoError(t, err)
	assert.Nil(t, user)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserRepository_GetByID_Error(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	log := mocks.NewMockLogger()
	repo := NewUserRepository(mock, log)

	userID := uuid.New()

	mock.ExpectQuery("SELECT id, email, phone, name, password, birth_date, gender, bio, city, artist, quote, is_verified, last_active, created_at, updated_at FROM \"user\" WHERE id = \\$1").
		WithArgs(userID).
		WillReturnError(pgx.ErrTxClosed)

	user, err := repo.GetByID(context.Background(), userID)
	assert.Error(t, err)
	assert.Nil(t, user)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserRepository_GetByEmail(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	log := mocks.NewMockLogger()
	repo := NewUserRepository(mock, log)

	userID := uuid.New()
	phone := "+1234567890"
	bio := "Test bio"
	city := "Moscow"
	artist := "Test Artist"
	quote := "Test Quote"
	email := "test@example.com"

	expectedUser := &domain.User{
		ID:         userID,
		Email:      email,
		Phone:      &phone,
		Name:       "Test User",
		Password:   "hashed_password",
		BirthDate:  time.Now().Add(-20 * 365 * 24 * time.Hour),
		Gender:     constants.GenderMale,
		Bio:        &bio,
		City:       &city,
		Artist:     &artist,
		Quote:      &quote,
		IsVerified: true,
		LastActive: time.Now(),
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	rows := mock.NewRows([]string{
		"id", "email", "phone", "name", "password", "birth_date", "gender",
		"bio", "city", "artist", "quote", "is_verified", "last_active", "created_at", "updated_at",
	}).AddRow(
		expectedUser.ID, expectedUser.Email, expectedUser.Phone, expectedUser.Name,
		expectedUser.Password, expectedUser.BirthDate, expectedUser.Gender, expectedUser.Bio,
		expectedUser.City, expectedUser.Artist, expectedUser.Quote, expectedUser.IsVerified,
		expectedUser.LastActive, expectedUser.CreatedAt, expectedUser.UpdatedAt,
	)

	mock.ExpectQuery("SELECT id, email, phone, name, password, birth_date, gender, bio, city, artist, quote, is_verified, last_active, created_at, updated_at FROM \"user\" WHERE email = \\$1").
		WithArgs(email).
		WillReturnRows(rows)

	user, err := repo.GetByEmail(context.Background(), email)
	assert.NoError(t, err)
	assert.Equal(t, expectedUser.ID, user.ID)
	assert.Equal(t, expectedUser.Email, user.Email)
	assert.NoError(t, mock.ExpectationsWereMet())
	assert.True(t, log.WasCalled("Trace"))
}

func TestUserRepository_GetByEmail_NotFound(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	log := mocks.NewMockLogger()
	repo := NewUserRepository(mock, log)

	email := "test@example.com"

	mock.ExpectQuery("SELECT id, email, phone, name, password, birth_date, gender, bio, city, artist, quote, is_verified, last_active, created_at, updated_at FROM \"user\" WHERE email = \\$1").
		WithArgs(email).
		WillReturnError(pgx.ErrNoRows)

	user, err := repo.GetByEmail(context.Background(), email)
	assert.NoError(t, err)
	assert.Nil(t, user)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserRepository_GetByPhone(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	log := mocks.NewMockLogger()
	repo := NewUserRepository(mock, log)

	userID := uuid.New()
	phone := "+1234567890"
	bio := "Test bio"
	city := "Moscow"
	artist := "Test Artist"
	quote := "Test Quote"

	expectedUser := &domain.User{
		ID:         userID,
		Email:      "test@example.com",
		Phone:      &phone,
		Name:       "Test User",
		Password:   "hashed_password",
		BirthDate:  time.Now().Add(-20 * 365 * 24 * time.Hour),
		Gender:     constants.GenderMale,
		Bio:        &bio,
		City:       &city,
		Artist:     &artist,
		Quote:      &quote,
		IsVerified: true,
		LastActive: time.Now(),
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	rows := mock.NewRows([]string{
		"id", "email", "phone", "name", "password", "birth_date", "gender",
		"bio", "city", "artist", "quote", "is_verified", "last_active", "created_at", "updated_at",
	}).AddRow(
		expectedUser.ID, expectedUser.Email, expectedUser.Phone, expectedUser.Name,
		expectedUser.Password, expectedUser.BirthDate, expectedUser.Gender, expectedUser.Bio,
		expectedUser.City, expectedUser.Artist, expectedUser.Quote, expectedUser.IsVerified,
		expectedUser.LastActive, expectedUser.CreatedAt, expectedUser.UpdatedAt,
	)

	mock.ExpectQuery("SELECT \\* FROM \"user\" WHERE phone = \\$1").
		WithArgs(phone).
		WillReturnRows(rows)

	user, err := repo.GetByPhone(context.Background(), phone)
	assert.NoError(t, err)
	assert.Equal(t, expectedUser.ID, user.ID)
	assert.Equal(t, expectedUser.Phone, user.Phone)
	assert.NoError(t, mock.ExpectationsWereMet())
	assert.True(t, log.WasCalled("Trace"))
}

func TestUserRepository_GetByPhone_NotFound(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	log := mocks.NewMockLogger()
	repo := NewUserRepository(mock, log)

	phone := "+1234567890"

	mock.ExpectQuery("SELECT \\* FROM \"user\" WHERE phone = \\$1").
		WithArgs(phone).
		WillReturnError(pgx.ErrNoRows)

	user, err := repo.GetByPhone(context.Background(), phone)
	assert.NoError(t, err)
	assert.Nil(t, user)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserRepository_Update(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	log := mocks.NewMockLogger()
	repo := NewUserRepository(mock, log)

	phone := "+0987654321"
	bio := "Updated bio"
	city := "St. Petersburg"
	artist := "Updated Artist"
	quote := "Updated Quote"

	user := &domain.User{
		ID:         uuid.New(),
		Email:      "updated@example.com",
		Phone:      &phone,
		Name:       "Updated User",
		Password:   "new_hashed_password",
		BirthDate:  time.Now().Add(-25 * 365 * 24 * time.Hour),
		Gender:     constants.GenderFemale,
		Bio:        &bio,
		City:       &city,
		Artist:     &artist,
		Quote:      &quote,
		IsVerified: true,
	}

	rows := mock.NewRows([]string{"updated_at"}).AddRow(time.Now())

	mock.ExpectQuery("UPDATE \"user\"").
		WithArgs(
			user.Email, user.Phone, user.Name, user.Password, user.BirthDate,
			user.Gender, user.Bio, user.City, user.Artist, user.Quote, user.IsVerified, user.ID,
		).
		WillReturnRows(rows)

	err = repo.Update(context.Background(), user)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
	assert.True(t, log.WasCalled("Trace"))
}

func TestUserRepository_Update_Error(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	log := mocks.NewMockLogger()
	repo := NewUserRepository(mock, log)

	phone := "+0987654321"
	bio := "Updated bio"
	city := "St. Petersburg"
	artist := "Updated Artist"
	quote := "Updated Quote"

	user := &domain.User{
		ID:         uuid.New(),
		Email:      "updated@example.com",
		Phone:      &phone,
		Name:       "Updated User",
		Password:   "new_hashed_password",
		BirthDate:  time.Now().Add(-25 * 365 * 24 * time.Hour),
		Gender:     constants.GenderFemale,
		Bio:        &bio,
		City:       &city,
		Artist:     &artist,
		Quote:      &quote,
		IsVerified: true,
	}

	mock.ExpectQuery("UPDATE \"user\"").
		WithArgs(
			user.Email, user.Phone, user.Name, user.Password, user.BirthDate,
			user.Gender, user.Bio, user.City, user.Artist, user.Quote, user.IsVerified, user.ID,
		).
		WillReturnError(pgx.ErrTxClosed)

	err = repo.Update(context.Background(), user)
	assert.Error(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserRepository_UpdateLastActive(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	log := mocks.NewMockLogger()
	repo := NewUserRepository(mock, log)

	userID := uuid.New()

	mock.ExpectExec("UPDATE \"user\" SET last_active = NOW\\(\\) WHERE id = \\$1").
		WithArgs(userID).
		WillReturnResult(pgxmock.NewResult("UPDATE", 1))

	err = repo.UpdateLastActive(context.Background(), userID)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
	assert.True(t, log.WasCalled("Trace"))
}

func TestUserRepository_UpdateLastActive_Error(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	log := mocks.NewMockLogger()
	repo := NewUserRepository(mock, log)

	userID := uuid.New()

	mock.ExpectExec("UPDATE \"user\" SET last_active = NOW\\(\\) WHERE id = \\$1").
		WithArgs(userID).
		WillReturnError(pgx.ErrTxClosed)

	err = repo.UpdateLastActive(context.Background(), userID)
	assert.Error(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserRepository_Delete(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	log := mocks.NewMockLogger()
	repo := NewUserRepository(mock, log)

	userID := uuid.New()

	mock.ExpectExec("DELETE FROM \"user\" WHERE id = \\$1").
		WithArgs(userID).
		WillReturnResult(pgxmock.NewResult("DELETE", 1))

	err = repo.Delete(context.Background(), userID)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
	assert.True(t, log.WasCalled("Trace"))
}

func TestUserRepository_Delete_Error(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	log := mocks.NewMockLogger()
	repo := NewUserRepository(mock, log)

	userID := uuid.New()

	mock.ExpectExec("DELETE FROM \"user\" WHERE id = \\$1").
		WithArgs(userID).
		WillReturnError(pgx.ErrTxClosed)

	err = repo.Delete(context.Background(), userID)
	assert.Error(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserRepository_GetUsersByIDs(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	log := mocks.NewMockLogger()
	repo := NewUserRepository(mock, log)

	userIDs := []uuid.UUID{uuid.New(), uuid.New()}
	phone := "+1234567890"
	bio := "Test bio"
	city := "Moscow"
	artist := "Test Artist"
	quote := "Test Quote"

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
			City:       &city,
			Artist:     &artist,
			Quote:      &quote,
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
			City:       &city,
			Artist:     &artist,
			Quote:      &quote,
			IsVerified: true,
			LastActive: time.Now(),
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		},
	}

	rows := mock.NewRows([]string{
		"id", "email", "phone", "name", "password", "birth_date", "gender",
		"bio", "city", "artist", "quote", "is_verified", "last_active", "created_at", "updated_at",
	})
	for _, user := range expectedUsers {
		rows.AddRow(
			user.ID, user.Email, user.Phone, user.Name, user.Password, user.BirthDate,
			user.Gender, user.Bio, user.City, user.Artist, user.Quote, user.IsVerified,
			user.LastActive, user.CreatedAt, user.UpdatedAt,
		)
	}

	mock.ExpectQuery("SELECT id, email, phone, name, password, birth_date, gender, bio, city, artist, quote, is_verified, last_active, created_at, updated_at FROM \"user\" WHERE id IN").
		WithArgs(userIDs[0], userIDs[1]).
		WillReturnRows(rows)

	users, err := repo.GetUsersByIDs(context.Background(), userIDs)
	assert.NoError(t, err)
	assert.Len(t, users, 2)
	assert.NoError(t, mock.ExpectationsWereMet())
	assert.True(t, log.WasCalled("Trace"))
}

func TestUserRepository_GetUsersByIDs_Empty(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	log := mocks.NewMockLogger()
	repo := NewUserRepository(mock, log)

	users, err := repo.GetUsersByIDs(context.Background(), []uuid.UUID{})
	assert.NoError(t, err)
	assert.Empty(t, users)
}

func TestUserRepository_GetUsersByIDs_ScanError(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	log := mocks.NewMockLogger()
	repo := NewUserRepository(mock, log)

	userIDs := []uuid.UUID{uuid.New()}

	// Создаем данные, которые вызовут ошибку сканирования
	rows := mock.NewRows([]string{
		"id", "email", "phone", "name", "password", "birth_date", "gender",
		"bio", "city", "artist", "quote", "is_verified", "last_active", "created_at", "updated_at",
	}).AddRow(
		"invalid-uuid", "email@test.com", nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, // Invalid UUID to cause scan error
	)

	mock.ExpectQuery("SELECT id, email, phone, name, password, birth_date, gender, bio, city, artist, quote, is_verified, last_active, created_at, updated_at FROM \"user\" WHERE id IN").
		WithArgs(userIDs[0]).
		WillReturnRows(rows)

	users, err := repo.GetUsersByIDs(context.Background(), userIDs)

	// В текущей реализации ошибка сканирования возвращается
	assert.Error(t, err)
	assert.Nil(t, users)
	assert.NoError(t, mock.ExpectationsWereMet())
	assert.True(t, log.WasCalled("Errorf"))
}

func TestUserRepository_GetUsersForFeed(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	log := mocks.NewMockLogger()
	repo := NewUserRepository(mock, log)

	userID := uuid.New()
	otherUserID := uuid.New()
	phone := "+1234567890"
	bio := "Test bio"
	city := "Moscow"
	artist := "Test Artist"
	quote := "Test Quote"
	limit := 10
	offset := 0

	expectedUsers := []domain.User{
		{
			ID:         otherUserID,
			Email:      "other@example.com",
			Phone:      &phone,
			Name:       "Other User",
			Password:   "password",
			BirthDate:  time.Now().Add(-20 * 365 * 24 * time.Hour),
			Gender:     constants.GenderMale,
			Bio:        &bio,
			City:       &city,
			Artist:     &artist,
			Quote:      &quote,
			IsVerified: true,
			LastActive: time.Now(),
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		},
	}

	rows := mock.NewRows([]string{
		"id", "email", "phone", "name", "password", "birth_date", "gender",
		"bio", "city", "artist", "quote", "is_verified", "last_active", "created_at", "updated_at",
	})
	for _, user := range expectedUsers {
		rows.AddRow(
			user.ID, user.Email, user.Phone, user.Name, user.Password, user.BirthDate,
			user.Gender, user.Bio, user.City, user.Artist, user.Quote, user.IsVerified,
			user.LastActive, user.CreatedAt, user.UpdatedAt,
		)
	}

	mock.ExpectQuery("SELECT id, email, phone, name, password, birth_date, gender, bio, city, artist, quote, is_verified, last_active, created_at, updated_at FROM \"user\" u").
		WithArgs(userID, limit, offset).
		WillReturnRows(rows)

	users, err := repo.GetUsersForFeed(context.Background(), userID, limit, offset)
	assert.NoError(t, err)
	assert.Len(t, users, 1)
	assert.Equal(t, otherUserID, users[0].ID)
	assert.NoError(t, mock.ExpectationsWereMet())
	assert.True(t, log.WasCalled("Tracef"))
	assert.True(t, log.WasCalled("Debugf"))
	assert.True(t, log.WasCalled("Trace"))
}

func TestUserRepository_GetUsersForFeed_Error(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	log := mocks.NewMockLogger()
	repo := NewUserRepository(mock, log)

	userID := uuid.New()
	limit := 10
	offset := 0

	mock.ExpectQuery("SELECT id, email, phone, name, password, birth_date, gender, bio, city, artist, quote, is_verified, last_active, created_at, updated_at FROM \"user\" u").
		WithArgs(userID, limit, offset).
		WillReturnError(pgx.ErrTxClosed)

	users, err := repo.GetUsersForFeed(context.Background(), userID, limit, offset)
	assert.Error(t, err)
	assert.Nil(t, users)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserRepository_GetUsersForFeed_ScanError(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	log := mocks.NewMockLogger()
	repo := NewUserRepository(mock, log)

	userID := uuid.New()
	limit := 10
	offset := 0

	// Создаем данные, которые вызовут ошибку сканирования
	rows := mock.NewRows([]string{
		"id", "email", "phone", "name", "password", "birth_date", "gender",
		"bio", "city", "artist", "quote", "is_verified", "last_active", "created_at", "updated_at",
	}).AddRow(
		"invalid-uuid", "email@test.com", nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, // Invalid UUID to cause scan error
	)

	mock.ExpectQuery("SELECT id, email, phone, name, password, birth_date, gender, bio, city, artist, quote, is_verified, last_active, created_at, updated_at FROM \"user\" u").
		WithArgs(userID, limit, offset).
		WillReturnRows(rows)

	users, err := repo.GetUsersForFeed(context.Background(), userID, limit, offset)

	// В текущей реализации ошибка сканирования возвращается
	assert.Error(t, err)
	assert.Nil(t, users)
	assert.NoError(t, mock.ExpectationsWereMet())
	assert.True(t, log.WasCalled("Errorf"))
}
