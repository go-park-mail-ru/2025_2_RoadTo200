package postgres

import (
	"context"
	"testing"
	"time"

	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/domain/constants"
	domain "github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/domain/entities"
	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/tests/mocks"
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

	log := mocks.NewMockLogger()
	repo := NewUserPreferenceRepository(mock, log)

	userID := uuid.New()
	preference := &domain.UserPreference{
		UserID:       userID,
		ShowGender:   constants.GenderPrefBoth,
		AgeMin:       18,
		AgeMax:       35,
		MaxDistance:  50,
		GlobalSearch: true,
	}

	rows := mock.NewRows([]string{"created_at", "updated_at"}).
		AddRow(time.Now(), time.Now())

	mock.ExpectQuery("INSERT INTO user_preference").
		WithArgs(
			preference.UserID, preference.ShowGender, preference.AgeMin, preference.AgeMax,
			preference.MaxDistance, preference.GlobalSearch,
		).
		WillReturnRows(rows)

	err = repo.Create(context.Background(), preference)
	assert.NoError(t, err)
	assert.False(t, preference.CreatedAt.IsZero())
	assert.False(t, preference.UpdatedAt.IsZero())
	assert.NoError(t, mock.ExpectationsWereMet())
	assert.True(t, log.WasCalled("Trace"))
}

func TestUserPreferenceRepository_Create_Error(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	log := mocks.NewMockLogger()
	repo := NewUserPreferenceRepository(mock, log)

	userID := uuid.New()
	preference := &domain.UserPreference{
		UserID:       userID,
		ShowGender:   constants.GenderPrefBoth,
		AgeMin:       18,
		AgeMax:       35,
		MaxDistance:  50,
		GlobalSearch: true,
	}

	mock.ExpectQuery("INSERT INTO user_preference").
		WithArgs(
			preference.UserID, preference.ShowGender, preference.AgeMin, preference.AgeMax,
			preference.MaxDistance, preference.GlobalSearch,
		).
		WillReturnError(pgx.ErrTxClosed)

	err = repo.Create(context.Background(), preference)
	assert.Error(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserPreferenceRepository_GetByUserID(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	log := mocks.NewMockLogger()
	repo := NewUserPreferenceRepository(mock, log)

	userID := uuid.New()

	expectedPreference := &domain.UserPreference{
		UserID:       userID,
		ShowGender:   constants.GenderPrefBoth,
		AgeMin:       20,
		AgeMax:       40,
		MaxDistance:  30,
		GlobalSearch: false,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	rows := mock.NewRows([]string{
		"user_id", "show_gender", "age_min", "age_max", "max_distance", "global_search", "created_at", "updated_at",
	}).AddRow(
		expectedPreference.UserID, expectedPreference.ShowGender, expectedPreference.AgeMin, expectedPreference.AgeMax,
		expectedPreference.MaxDistance, expectedPreference.GlobalSearch, expectedPreference.CreatedAt, expectedPreference.UpdatedAt,
	)

	mock.ExpectQuery("SELECT \\* FROM user_preference WHERE user_id = \\$1").
		WithArgs(userID).
		WillReturnRows(rows)

	preference, err := repo.GetByUserID(context.Background(), userID)
	assert.NoError(t, err)
	assert.Equal(t, expectedPreference.UserID, preference.UserID)
	assert.Equal(t, expectedPreference.ShowGender, preference.ShowGender)
	assert.Equal(t, expectedPreference.AgeMin, preference.AgeMin)
	assert.Equal(t, expectedPreference.AgeMax, preference.AgeMax)
	assert.Equal(t, expectedPreference.MaxDistance, preference.MaxDistance)
	assert.Equal(t, expectedPreference.GlobalSearch, preference.GlobalSearch)
	assert.NoError(t, mock.ExpectationsWereMet())
	assert.True(t, log.WasCalled("Trace"))
}

func TestUserPreferenceRepository_GetByUserID_NotFound(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	log := mocks.NewMockLogger()
	repo := NewUserPreferenceRepository(mock, log)

	userID := uuid.New()

	mock.ExpectQuery("SELECT \\* FROM user_preference WHERE user_id = \\$1").
		WithArgs(userID).
		WillReturnError(pgx.ErrNoRows)

	preference, err := repo.GetByUserID(context.Background(), userID)
	assert.NoError(t, err)
	assert.Nil(t, preference)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserPreferenceRepository_GetByUserID_Error(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	log := mocks.NewMockLogger()
	repo := NewUserPreferenceRepository(mock, log)

	userID := uuid.New()

	mock.ExpectQuery("SELECT \\* FROM user_preference WHERE user_id = \\$1").
		WithArgs(userID).
		WillReturnError(pgx.ErrTxClosed)

	preference, err := repo.GetByUserID(context.Background(), userID)
	assert.Error(t, err)
	assert.Nil(t, preference)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserPreferenceRepository_Update(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	log := mocks.NewMockLogger()
	repo := NewUserPreferenceRepository(mock, log)

	userID := uuid.New()
	preference := &domain.UserPreference{
		UserID:       userID,
		ShowGender:   constants.GenderPrefBoth,
		AgeMin:       25,
		AgeMax:       45,
		MaxDistance:  100,
		GlobalSearch: true,
	}

	rows := mock.NewRows([]string{"updated_at"}).AddRow(time.Now())

	mock.ExpectQuery("UPDATE user_preference").
		WithArgs(
			preference.ShowGender, preference.AgeMin, preference.AgeMax, preference.MaxDistance,
			preference.GlobalSearch, preference.UserID,
		).
		WillReturnRows(rows)

	err = repo.Update(context.Background(), preference)
	assert.NoError(t, err)
	assert.False(t, preference.UpdatedAt.IsZero())
	assert.NoError(t, mock.ExpectationsWereMet())
	assert.True(t, log.WasCalled("Trace"))
}

func TestUserPreferenceRepository_Update_Error(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	log := mocks.NewMockLogger()
	repo := NewUserPreferenceRepository(mock, log)

	userID := uuid.New()
	preference := &domain.UserPreference{
		UserID:       userID,
		ShowGender:   constants.GenderPrefBoth,
		AgeMin:       25,
		AgeMax:       45,
		MaxDistance:  100,
		GlobalSearch: true,
	}

	mock.ExpectQuery("UPDATE user_preference").
		WithArgs(
			preference.ShowGender, preference.AgeMin, preference.AgeMax, preference.MaxDistance,
			preference.GlobalSearch, preference.UserID,
		).
		WillReturnError(pgx.ErrTxClosed)

	err = repo.Update(context.Background(), preference)
	assert.Error(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserPreferenceRepository_Delete(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	log := mocks.NewMockLogger()
	repo := NewUserPreferenceRepository(mock, log)

	userID := uuid.New()

	mock.ExpectExec("DELETE FROM user_preference WHERE user_id = \\$1").
		WithArgs(userID).
		WillReturnResult(pgxmock.NewResult("DELETE", 1))

	err = repo.Delete(context.Background(), userID)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
	assert.True(t, log.WasCalled("Trace"))
}

func TestUserPreferenceRepository_Delete_Error(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	log := mocks.NewMockLogger()
	repo := NewUserPreferenceRepository(mock, log)

	userID := uuid.New()

	mock.ExpectExec("DELETE FROM user_preference WHERE user_id = \\$1").
		WithArgs(userID).
		WillReturnError(pgx.ErrTxClosed)

	err = repo.Delete(context.Background(), userID)
	assert.Error(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserPreferenceRepository_GetInterests(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	log := mocks.NewMockLogger()
	repo := NewUserPreferenceRepository(mock, log)

	userID := uuid.New()

	expectedInterests := []domain.Interest{
		{UserID: userID, Theme: "Music"},
		{UserID: userID, Theme: "Sports"},
		{UserID: userID, Theme: "Travel"},
	}

	rows := mock.NewRows([]string{"user_id", "theme"})
	for _, interest := range expectedInterests {
		rows.AddRow(interest.UserID, interest.Theme)
	}

	mock.ExpectQuery("SELECT i.user_id, i.theme FROM interest i WHERE i.user_id = \\$1").
		WithArgs(userID).
		WillReturnRows(rows)

	interests, err := repo.GetInterests(context.Background(), userID)
	assert.NoError(t, err)
	assert.Len(t, interests, 3)
	assert.Equal(t, expectedInterests[0].Theme, interests[0].Theme)
	assert.Equal(t, expectedInterests[1].Theme, interests[1].Theme)
	assert.Equal(t, expectedInterests[2].Theme, interests[2].Theme)
	assert.NoError(t, mock.ExpectationsWereMet())
	assert.True(t, log.WasCalled("Trace"))
}

func TestUserPreferenceRepository_GetInterests_Empty(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	log := mocks.NewMockLogger()
	repo := NewUserPreferenceRepository(mock, log)

	userID := uuid.New()

	rows := mock.NewRows([]string{"user_id", "theme"})

	mock.ExpectQuery("SELECT i.user_id, i.theme FROM interest i WHERE i.user_id = \\$1").
		WithArgs(userID).
		WillReturnRows(rows)

	interests, err := repo.GetInterests(context.Background(), userID)
	assert.NoError(t, err)
	assert.Empty(t, interests)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserPreferenceRepository_GetInterests_ScanError(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	log := mocks.NewMockLogger()
	repo := NewUserPreferenceRepository(mock, log)

	userID := uuid.New()

	rows := mock.NewRows([]string{"user_id", "theme"}).
		AddRow(nil, nil) // Invalid data to cause scan error

	mock.ExpectQuery("SELECT i.user_id, i.theme FROM interest i WHERE i.user_id = \\$1").
		WithArgs(userID).
		WillReturnRows(rows)

	interests, err := repo.GetInterests(context.Background(), userID)
	assert.NoError(t, err) // Note: In the actual code, scan errors are logged but not returned
	assert.Len(t, interests, 0)
	assert.NoError(t, mock.ExpectationsWereMet())
	assert.True(t, log.WasCalled("Errorf"))
}

func TestUserPreferenceRepository_GetInterests_QueryError(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	log := mocks.NewMockLogger()
	repo := NewUserPreferenceRepository(mock, log)

	userID := uuid.New()

	mock.ExpectQuery("SELECT i.user_id, i.theme FROM interest i WHERE i.user_id = \\$1").
		WithArgs(userID).
		WillReturnError(pgx.ErrTxClosed)

	interests, err := repo.GetInterests(context.Background(), userID)
	assert.Error(t, err)
	assert.Nil(t, interests)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserPreferenceRepository_UpdateInterests(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	log := mocks.NewMockLogger()
	repo := NewUserPreferenceRepository(mock, log)

	userID := uuid.New()
	interests := []domain.Interest{
		{UserID: userID, Theme: "Music"},
		{UserID: userID, Theme: "Sports"},
		{UserID: userID, Theme: "Travel"},
	}

	mock.ExpectBegin()
	mock.ExpectExec("DELETE FROM interest WHERE user_id = \\$1").
		WithArgs(userID).
		WillReturnResult(pgxmock.NewResult("DELETE", 3))

	// For batch operations, we expect multiple Exec calls
	for _, interest := range interests {
		mock.ExpectExec("INSERT INTO interest").
			WithArgs(userID, interest.Theme).
			WillReturnResult(pgxmock.NewResult("INSERT", 1))
	}

	mock.ExpectCommit()

	err = repo.UpdateInterests(context.Background(), userID, interests)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
	assert.True(t, log.WasCalled("Trace"))
}

func TestUserPreferenceRepository_UpdateInterests_Empty(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	log := mocks.NewMockLogger()
	repo := NewUserPreferenceRepository(mock, log)

	userID := uuid.New()

	mock.ExpectBegin()
	mock.ExpectExec("DELETE FROM interest WHERE user_id = \\$1").
		WithArgs(userID).
		WillReturnResult(pgxmock.NewResult("DELETE", 0))
	mock.ExpectCommit()

	err = repo.UpdateInterests(context.Background(), userID, []domain.Interest{})
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserPreferenceRepository_UpdateInterests_BeginError(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	log := mocks.NewMockLogger()
	repo := NewUserPreferenceRepository(mock, log)

	userID := uuid.New()
	interests := []domain.Interest{
		{UserID: userID, Theme: "Music"},
	}

	mock.ExpectBegin().WillReturnError(pgx.ErrTxClosed)

	err = repo.UpdateInterests(context.Background(), userID, interests)
	assert.Error(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserPreferenceRepository_UpdateInterests_DeleteError(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	log := mocks.NewMockLogger()
	repo := NewUserPreferenceRepository(mock, log)

	userID := uuid.New()
	interests := []domain.Interest{
		{UserID: userID, Theme: "Music"},
	}

	mock.ExpectBegin()
	mock.ExpectExec("DELETE FROM interest WHERE user_id = \\$1").
		WithArgs(userID).
		WillReturnError(pgx.ErrTxClosed)
	mock.ExpectRollback()

	err = repo.UpdateInterests(context.Background(), userID, interests)
	assert.Error(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserPreferenceRepository_UpdateInterests_BatchError(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	log := mocks.NewMockLogger()
	repo := NewUserPreferenceRepository(mock, log)

	userID := uuid.New()
	interests := []domain.Interest{
		{UserID: userID, Theme: "Music"},
	}

	mock.ExpectBegin()
	mock.ExpectExec("DELETE FROM interest WHERE user_id = \\$1").
		WithArgs(userID).
		WillReturnResult(pgxmock.NewResult("DELETE", 1))

	// Simulate batch execution error by returning error on first insert
	mock.ExpectExec("INSERT INTO interest").
		WithArgs(userID, interests[0].Theme).
		WillReturnError(pgx.ErrTxClosed)
	mock.ExpectRollback()

	err = repo.UpdateInterests(context.Background(), userID, interests)
	assert.Error(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserPreferenceRepository_UpdateInterests_CommitError(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	log := mocks.NewMockLogger()
	repo := NewUserPreferenceRepository(mock, log)

	userID := uuid.New()
	interests := []domain.Interest{
		{UserID: userID, Theme: "Music"},
	}

	mock.ExpectBegin()
	mock.ExpectExec("DELETE FROM interest WHERE user_id = \\$1").
		WithArgs(userID).
		WillReturnResult(pgxmock.NewResult("DELETE", 1))
	mock.ExpectExec("INSERT INTO interest").
		WithArgs(userID, interests[0].Theme).
		WillReturnResult(pgxmock.NewResult("INSERT", 1))
	mock.ExpectCommit().WillReturnError(pgx.ErrTxClosed)

	err = repo.UpdateInterests(context.Background(), userID, interests)
	assert.Error(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}
