package unit

import (
	"testing"
	"time"

	domain "github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/domain/entities"
	repository "github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/repository/implementations/postgres"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
	"github.com/pashagolub/pgxmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUserPreferenceRepository_Create(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	repo := repository.NewUserPreferenceRepository(mock)

	userID := uuid.New()

	preference := &domain.UserPreference{
		UserID:       userID,
		ShowGender:   "female",
		AgeMin:       18,
		AgeMax:       30,
		MaxDistance:  50,
		GlobalSearch: true,
	}

	rows := mock.NewRows([]string{"created_at", "updated_at"}).
		AddRow(time.Now(), time.Now())

	mock.ExpectQuery("INSERT INTO user_preference").
		WithArgs(
			preference.UserID, preference.ShowGender, preference.AgeMin,
			preference.AgeMax, preference.MaxDistance, preference.GlobalSearch,
		).
		WillReturnRows(rows)

	err = repo.Create(preference)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserPreferenceRepository_GetByUserID(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	repo := repository.NewUserPreferenceRepository(mock)

	userID := uuid.New()

	expectedPreference := &domain.UserPreference{
		UserID:       userID,
		ShowGender:   "female",
		AgeMin:       18,
		AgeMax:       30,
		MaxDistance:  50,
		GlobalSearch: true,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	rows := mock.NewRows([]string{
		"user_id", "show_gender", "age_min", "age_max", "max_distance",
		"global_search", "created_at", "updated_at",
	}).AddRow(
		expectedPreference.UserID, expectedPreference.ShowGender, expectedPreference.AgeMin,
		expectedPreference.AgeMax, expectedPreference.MaxDistance, expectedPreference.GlobalSearch,
		expectedPreference.CreatedAt, expectedPreference.UpdatedAt,
	)

	mock.ExpectQuery("SELECT \\* FROM user_preference WHERE user_id = \\$1").
		WithArgs(userID).
		WillReturnRows(rows)

	preference, err := repo.GetByUserID(userID)
	assert.NoError(t, err)
	assert.Equal(t, expectedPreference.UserID, preference.UserID)
	assert.Equal(t, expectedPreference.ShowGender, preference.ShowGender)
	assert.Equal(t, expectedPreference.AgeMin, preference.AgeMin)
	assert.Equal(t, expectedPreference.AgeMax, preference.AgeMax)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserPreferenceRepository_GetByUserID_NotFound(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	repo := repository.NewUserPreferenceRepository(mock)

	userID := uuid.New()

	mock.ExpectQuery("SELECT \\* FROM user_preference WHERE user_id = \\$1").
		WithArgs(userID).
		WillReturnError(pgx.ErrNoRows)

	preference, err := repo.GetByUserID(userID)
	assert.NoError(t, err)
	assert.Nil(t, preference)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserPreferenceRepository_Update(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	repo := repository.NewUserPreferenceRepository(mock)

	userID := uuid.New()

	preference := &domain.UserPreference{
		UserID:       userID,
		ShowGender:   "male",
		AgeMin:       20,
		AgeMax:       35,
		MaxDistance:  100,
		GlobalSearch: false,
	}

	rows := mock.NewRows([]string{"updated_at"}).AddRow(time.Now())

	mock.ExpectQuery("UPDATE user_preference").
		WithArgs(
			preference.ShowGender, preference.AgeMin, preference.AgeMax,
			preference.MaxDistance, preference.GlobalSearch, preference.UserID,
		).
		WillReturnRows(rows)

	err = repo.Update(preference)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserPreferenceRepository_Delete(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	repo := repository.NewUserPreferenceRepository(mock)

	userID := uuid.New()

	mock.ExpectExec("DELETE FROM user_preference WHERE user_id = \\$1").
		WithArgs(userID).
		WillReturnResult(pgxmock.NewResult("DELETE", 1))

	err = repo.Delete(userID)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserPreferenceRepository_Create_Error(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	repo := repository.NewUserPreferenceRepository(mock)

	userID := uuid.New()

	preference := &domain.UserPreference{
		UserID:       userID,
		ShowGender:   "female",
		AgeMin:       18,
		AgeMax:       30,
		MaxDistance:  50,
		GlobalSearch: true,
	}

	mock.ExpectQuery("INSERT INTO user_preference").
		WithArgs(
			preference.UserID, preference.ShowGender, preference.AgeMin,
			preference.AgeMax, preference.MaxDistance, preference.GlobalSearch,
		).
		WillReturnError(pgx.ErrNoRows)

	err = repo.Create(preference)
	assert.Error(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserPreferenceRepository_Update_Error(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	repo := repository.NewUserPreferenceRepository(mock)

	userID := uuid.New()

	preference := &domain.UserPreference{
		UserID:       userID,
		ShowGender:   "male",
		AgeMin:       20,
		AgeMax:       35,
		MaxDistance:  100,
		GlobalSearch: false,
	}

	mock.ExpectQuery("UPDATE user_preference").
		WithArgs(
			preference.ShowGender, preference.AgeMin, preference.AgeMax,
			preference.MaxDistance, preference.GlobalSearch, preference.UserID,
		).
		WillReturnError(pgx.ErrNoRows)

	err = repo.Update(preference)
	assert.Error(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserPreferenceRepository_Delete_Error(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	repo := repository.NewUserPreferenceRepository(mock)

	userID := uuid.New()

	mock.ExpectExec("DELETE FROM user_preference WHERE user_id = \\$1").
		WithArgs(userID).
		WillReturnError(pgx.ErrNoRows)

	err = repo.Delete(userID)
	assert.Error(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserPreferenceRepository_GetByUserID_Error(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	repo := repository.NewUserPreferenceRepository(mock)

	userID := uuid.New()

	mock.ExpectQuery("SELECT \\* FROM user_preference WHERE user_id = \\$1").
		WithArgs(userID).
		WillReturnError(pgx.ErrTxClosed)

	preference, err := repo.GetByUserID(userID)
	assert.Error(t, err)
	assert.Nil(t, preference)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserPreferenceRepository_Create_WithDefaultValues(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	repo := repository.NewUserPreferenceRepository(mock)

	userID := uuid.New()

	preference := &domain.UserPreference{
		UserID:       userID,
		ShowGender:   "both", // default value
		AgeMin:       18,     // default minimum
		AgeMax:       99,     // default maximum
		MaxDistance:  50,     // default distance
		GlobalSearch: false,  // default search scope
	}

	rows := mock.NewRows([]string{"created_at", "updated_at"}).
		AddRow(time.Now(), time.Now())

	mock.ExpectQuery("INSERT INTO user_preference").
		WithArgs(
			preference.UserID, preference.ShowGender, preference.AgeMin,
			preference.AgeMax, preference.MaxDistance, preference.GlobalSearch,
		).
		WillReturnRows(rows)

	err = repo.Create(preference)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserPreferenceRepository_Update_WithExtremeValues(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	repo := repository.NewUserPreferenceRepository(mock)

	userID := uuid.New()

	preference := &domain.UserPreference{
		UserID:       userID,
		ShowGender:   "both",
		AgeMin:       18,
		AgeMax:       99,
		MaxDistance:  500,  // large distance
		GlobalSearch: true, // global search enabled
	}

	rows := mock.NewRows([]string{"updated_at"}).AddRow(time.Now())

	mock.ExpectQuery("UPDATE user_preference").
		WithArgs(
			preference.ShowGender, preference.AgeMin, preference.AgeMax,
			preference.MaxDistance, preference.GlobalSearch, preference.UserID,
		).
		WillReturnRows(rows)

	err = repo.Update(preference)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}
