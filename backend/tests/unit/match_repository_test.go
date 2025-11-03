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

func TestMatchRepository_Create(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	repo := repository.NewMatchRepository(mock)

	user1ID := uuid.New()
	user2ID := uuid.New()

	match := &domain.Match{
		User1ID:  user1ID,
		User2ID:  user2ID,
		IsActive: true,
	}

	rows := mock.NewRows([]string{"matched_at"}).
		AddRow(time.Now())

	mock.ExpectQuery("INSERT INTO match").
		WithArgs(user1ID, user2ID, match.IsActive).
		WillReturnRows(rows)

	err = repo.Create(match)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestMatchRepository_Create_WithOrderNormalization(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	repo := repository.NewMatchRepository(mock)

	// user1ID будет "больше" user2ID в лексикографическом порядке
	user1ID := uuid.MustParse("ffffffff-ffff-ffff-ffff-ffffffffffff")
	user2ID := uuid.MustParse("00000000-0000-0000-0000-000000000000")

	match := &domain.Match{
		User1ID:  user1ID, // больший UUID
		User2ID:  user2ID, // меньший UUID
		IsActive: true,
	}

	rows := mock.NewRows([]string{"matched_at"}).
		AddRow(time.Now())

	// Ожидаем, что в запросе IDs будут в нормализованном порядке (user2ID, user1ID)
	mock.ExpectQuery("INSERT INTO match").
		WithArgs(user2ID, user1ID, match.IsActive).
		WillReturnRows(rows)

	err = repo.Create(match)
	assert.NoError(t, err)
	// Проверяем, что IDs в объекте match были нормализованы
	assert.Equal(t, user2ID, match.User1ID)
	assert.Equal(t, user1ID, match.User2ID)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestMatchRepository_GetByUsers(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	repo := repository.NewMatchRepository(mock)

	user1ID := uuid.New()
	user2ID := uuid.New()

	expectedMatch := &domain.Match{
		User1ID:   user1ID,
		User2ID:   user2ID,
		IsActive:  true,
		MatchedAt: time.Now(),
	}

	rows := mock.NewRows([]string{
		"user1_id", "user2_id", "is_active", "matched_at",
	}).AddRow(
		expectedMatch.User1ID, expectedMatch.User2ID, expectedMatch.IsActive, expectedMatch.MatchedAt,
	)

	mock.ExpectQuery("SELECT \\* FROM match WHERE user1_id = \\$1 AND user2_id = \\$2").
		WithArgs(user1ID, user2ID).
		WillReturnRows(rows)

	match, err := repo.GetByUsers(user1ID, user2ID)
	assert.NoError(t, err)
	assert.Equal(t, expectedMatch.User1ID, match.User1ID)
	assert.Equal(t, expectedMatch.User2ID, match.User2ID)
	assert.Equal(t, expectedMatch.IsActive, match.IsActive)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestMatchRepository_GetByUsers_WithOrderNormalization(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	repo := repository.NewMatchRepository(mock)

	// user1ID будет "больше" user2ID в лексикографическом порядке
	user1ID := uuid.MustParse("ffffffff-ffff-ffff-ffff-ffffffffffff")
	user2ID := uuid.MustParse("00000000-0000-0000-0000-000000000000")

	expectedMatch := &domain.Match{
		User1ID:   user2ID, // нормализованный порядок
		User2ID:   user1ID,
		IsActive:  true,
		MatchedAt: time.Now(),
	}

	rows := mock.NewRows([]string{
		"user1_id", "user2_id", "is_active", "matched_at",
	}).AddRow(
		expectedMatch.User1ID, expectedMatch.User2ID, expectedMatch.IsActive, expectedMatch.MatchedAt,
	)

	// Ожидаем, что в запросе IDs будут в нормализованном порядке
	mock.ExpectQuery("SELECT \\* FROM match WHERE user1_id = \\$1 AND user2_id = \\$2").
		WithArgs(user2ID, user1ID).
		WillReturnRows(rows)

	match, err := repo.GetByUsers(user1ID, user2ID)
	assert.NoError(t, err)
	assert.Equal(t, expectedMatch.User1ID, match.User1ID)
	assert.Equal(t, expectedMatch.User2ID, match.User2ID)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestMatchRepository_GetByUsers_NotFound(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	repo := repository.NewMatchRepository(mock)

	user1ID := uuid.New()
	user2ID := uuid.New()

	mock.ExpectQuery("SELECT \\* FROM match WHERE user1_id = \\$1 AND user2_id = \\$2").
		WithArgs(user1ID, user2ID).
		WillReturnError(pgx.ErrNoRows)

	match, err := repo.GetByUsers(user1ID, user2ID)
	assert.NoError(t, err)
	assert.Nil(t, match)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestMatchRepository_GetUserMatches(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	repo := repository.NewMatchRepository(mock)

	userID := uuid.New()
	limit := 10
	offset := 0

	expectedMatches := []domain.Match{
		{
			User1ID:   userID,
			User2ID:   uuid.New(),
			IsActive:  true,
			MatchedAt: time.Now(),
		},
		{
			User1ID:   uuid.New(),
			User2ID:   userID,
			IsActive:  true,
			MatchedAt: time.Now(),
		},
	}

	rows := mock.NewRows([]string{
		"user1_id", "user2_id", "is_active", "matched_at",
	})
	for _, match := range expectedMatches {
		rows.AddRow(
			match.User1ID, match.User2ID, match.IsActive, match.MatchedAt,
		)
	}

	mock.ExpectQuery("SELECT \\* FROM match").
		WithArgs(userID, limit, offset).
		WillReturnRows(rows)

	matches, err := repo.GetUserMatches(userID, limit, offset)
	assert.NoError(t, err)
	assert.Len(t, matches, 2)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestMatchRepository_GetUserMatches_Empty(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	repo := repository.NewMatchRepository(mock)

	userID := uuid.New()
	limit := 10
	offset := 0

	rows := mock.NewRows([]string{
		"user1_id", "user2_id", "is_active", "matched_at",
	})

	mock.ExpectQuery("SELECT \\* FROM match").
		WithArgs(userID, limit, offset).
		WillReturnRows(rows)

	matches, err := repo.GetUserMatches(userID, limit, offset)
	assert.NoError(t, err)
	assert.Empty(t, matches)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestMatchRepository_UpdateActive(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	repo := repository.NewMatchRepository(mock)

	user1ID := uuid.New()
	user2ID := uuid.New()
	isActive := false

	mock.ExpectExec("UPDATE match SET is_active = \\$1 WHERE user1_id = \\$2 AND user2_id = \\$3").
		WithArgs(isActive, user1ID, user2ID).
		WillReturnResult(pgxmock.NewResult("UPDATE", 1))

	err = repo.UpdateActive(user1ID, user2ID, isActive)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestMatchRepository_UpdateActive_WithOrderNormalization(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	repo := repository.NewMatchRepository(mock)

	// user1ID будет "больше" user2ID в лексикографическом порядке
	user1ID := uuid.MustParse("ffffffff-ffff-ffff-ffff-ffffffffffff")
	user2ID := uuid.MustParse("00000000-0000-0000-0000-000000000000")
	isActive := false

	// Ожидаем, что в запросе IDs будут в нормализованном порядке
	mock.ExpectExec("UPDATE match SET is_active = \\$1 WHERE user1_id = \\$2 AND user2_id = \\$3").
		WithArgs(isActive, user2ID, user1ID).
		WillReturnResult(pgxmock.NewResult("UPDATE", 1))

	err = repo.UpdateActive(user1ID, user2ID, isActive)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestMatchRepository_Delete(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	repo := repository.NewMatchRepository(mock)

	user1ID := uuid.New()
	user2ID := uuid.New()

	mock.ExpectExec("DELETE FROM match WHERE user1_id = \\$1 AND user2_id = \\$2").
		WithArgs(user1ID, user2ID).
		WillReturnResult(pgxmock.NewResult("DELETE", 1))

	err = repo.Delete(user1ID, user2ID)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestMatchRepository_Delete_WithOrderNormalization(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	repo := repository.NewMatchRepository(mock)

	// user1ID будет "больше" user2ID в лексикографическом порядке
	user1ID := uuid.MustParse("ffffffff-ffff-ffff-ffff-ffffffffffff")
	user2ID := uuid.MustParse("00000000-0000-0000-0000-000000000000")

	// Ожидаем, что в запросе IDs будут в нормализованном порядке
	mock.ExpectExec("DELETE FROM match WHERE user1_id = \\$1 AND user2_id = \\$2").
		WithArgs(user2ID, user1ID).
		WillReturnResult(pgxmock.NewResult("DELETE", 1))

	err = repo.Delete(user1ID, user2ID)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestMatchRepository_CheckMutualLike(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	repo := repository.NewMatchRepository(mock)

	user1ID := uuid.New()
	user2ID := uuid.New()
	exists := true

	rows := mock.NewRows([]string{"exists"}).AddRow(exists)

	mock.ExpectQuery("SELECT EXISTS").
		WithArgs(user1ID, user2ID).
		WillReturnRows(rows)

	mutualLike, err := repo.CheckMutualLike(user1ID, user2ID)
	assert.NoError(t, err)
	assert.True(t, mutualLike)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestMatchRepository_CheckMutualLike_NoMutualLike(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	repo := repository.NewMatchRepository(mock)

	user1ID := uuid.New()
	user2ID := uuid.New()
	exists := false

	rows := mock.NewRows([]string{"exists"}).AddRow(exists)

	mock.ExpectQuery("SELECT EXISTS").
		WithArgs(user1ID, user2ID).
		WillReturnRows(rows)

	mutualLike, err := repo.CheckMutualLike(user1ID, user2ID)
	assert.NoError(t, err)
	assert.False(t, mutualLike)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestMatchRepository_Create_Error(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	repo := repository.NewMatchRepository(mock)

	user1ID := uuid.New()
	user2ID := uuid.New()

	match := &domain.Match{
		User1ID:  user1ID,
		User2ID:  user2ID,
		IsActive: true,
	}

	mock.ExpectQuery("INSERT INTO match").
		WithArgs(user1ID, user2ID, match.IsActive).
		WillReturnError(pgx.ErrNoRows)

	err = repo.Create(match)
	assert.Error(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestMatchRepository_UpdateActive_Error(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	repo := repository.NewMatchRepository(mock)

	user1ID := uuid.New()
	user2ID := uuid.New()
	isActive := false

	mock.ExpectExec("UPDATE match SET is_active = \\$1 WHERE user1_id = \\$2 AND user2_id = \\$3").
		WithArgs(isActive, user1ID, user2ID).
		WillReturnError(pgx.ErrNoRows)

	err = repo.UpdateActive(user1ID, user2ID, isActive)
	assert.Error(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestMatchRepository_CheckMutualLike_Error(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	repo := repository.NewMatchRepository(mock)

	user1ID := uuid.New()
	user2ID := uuid.New()

	mock.ExpectQuery("SELECT EXISTS").
		WithArgs(user1ID, user2ID).
		WillReturnError(pgx.ErrNoRows)

	mutualLike, err := repo.CheckMutualLike(user1ID, user2ID)
	assert.Error(t, err)
	assert.False(t, mutualLike)
	assert.NoError(t, mock.ExpectationsWereMet())
}
