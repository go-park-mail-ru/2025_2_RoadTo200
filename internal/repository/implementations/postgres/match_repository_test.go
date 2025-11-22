package postgres

import (
	"context"
	"testing"
	"time"

	domain "github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/domain/entities"
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

	repo := NewMatchRepository(mock)

	user1ID := uuid.MustParse("00000000-0000-0000-0000-000000000001")
	user2ID := uuid.MustParse("00000000-0000-0000-0000-000000000002")

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

	err = repo.Create(context.Background(), match)
	assert.NoError(t, err)
	assert.False(t, match.MatchedAt.IsZero())
	assert.Equal(t, user1ID, match.User1ID) // user1ID should remain the smaller one
	assert.Equal(t, user2ID, match.User2ID)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestMatchRepository_Create_WithReorderedUsers(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	repo := NewMatchRepository(mock)

	// user2ID is lexicographically smaller than user1ID
	user1ID := uuid.MustParse("00000000-0000-0000-0000-000000000002")
	user2ID := uuid.MustParse("00000000-0000-0000-0000-000000000001")

	match := &domain.Match{
		User1ID:  user1ID,
		User2ID:  user2ID,
		IsActive: true,
	}

	rows := mock.NewRows([]string{"matched_at"}).
		AddRow(time.Now())

	// Expect the users to be reordered (user2ID, user1ID) since user2ID < user1ID
	mock.ExpectQuery("INSERT INTO match").
		WithArgs(user2ID, user1ID, match.IsActive).
		WillReturnRows(rows)

	err = repo.Create(context.Background(), match)
	assert.NoError(t, err)
	assert.False(t, match.MatchedAt.IsZero())
	// After creation, User1ID should be the smaller one and User2ID the larger one
	assert.Equal(t, user2ID, match.User1ID) // user2ID is now User1ID because it's smaller
	assert.Equal(t, user1ID, match.User2ID) // user1ID is now User2ID because it's larger
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestMatchRepository_Create_Error(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	repo := NewMatchRepository(mock)

	user1ID := uuid.New()
	user2ID := uuid.New()

	match := &domain.Match{
		User1ID:  user1ID,
		User2ID:  user2ID,
		IsActive: true,
	}

	mock.ExpectQuery("INSERT INTO match").
		WithArgs(user1ID, user2ID, match.IsActive).
		WillReturnError(pgx.ErrTxClosed)

	err = repo.Create(context.Background(), match)
	assert.Error(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestMatchRepository_GetByUsers(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	repo := NewMatchRepository(mock)

	user1ID := uuid.MustParse("00000000-0000-0000-0000-000000000001")
	user2ID := uuid.MustParse("00000000-0000-0000-0000-000000000002")

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

	match, err := repo.GetByUsers(context.Background(), user1ID, user2ID)
	assert.NoError(t, err)
	assert.Equal(t, expectedMatch.User1ID, match.User1ID)
	assert.Equal(t, expectedMatch.User2ID, match.User2ID)
	assert.Equal(t, expectedMatch.IsActive, match.IsActive)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestMatchRepository_GetByUsers_WithReorderedUsers(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	repo := NewMatchRepository(mock)

	// user2ID is lexicographically smaller than user1ID
	user1ID := uuid.MustParse("00000000-0000-0000-0000-000000000002")
	user2ID := uuid.MustParse("00000000-0000-0000-0000-000000000001")

	expectedMatch := &domain.Match{
		User1ID:   user2ID, // user2ID becomes User1ID because it's smaller
		User2ID:   user1ID, // user1ID becomes User2ID because it's larger
		IsActive:  true,
		MatchedAt: time.Now(),
	}

	rows := mock.NewRows([]string{
		"user1_id", "user2_id", "is_active", "matched_at",
	}).AddRow(
		expectedMatch.User1ID, expectedMatch.User2ID, expectedMatch.IsActive, expectedMatch.MatchedAt,
	)

	// Expect the users to be reordered in the query
	mock.ExpectQuery("SELECT \\* FROM match WHERE user1_id = \\$1 AND user2_id = \\$2").
		WithArgs(user2ID, user1ID).
		WillReturnRows(rows)

	match, err := repo.GetByUsers(context.Background(), user1ID, user2ID)
	assert.NoError(t, err)
	assert.Equal(t, expectedMatch.User1ID, match.User1ID)
	assert.Equal(t, expectedMatch.User2ID, match.User2ID)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestMatchRepository_GetByUsers_NotFound(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	repo := NewMatchRepository(mock)

	user1ID := uuid.New()
	user2ID := uuid.New()

	mock.ExpectQuery("SELECT \\* FROM match WHERE user1_id = \\$1 AND user2_id = \\$2").
		WithArgs(user1ID, user2ID).
		WillReturnError(pgx.ErrNoRows)

	match, err := repo.GetByUsers(context.Background(), user1ID, user2ID)
	assert.NoError(t, err)
	assert.Nil(t, match)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestMatchRepository_GetByUsers_Error(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	repo := NewMatchRepository(mock)

	user1ID := uuid.New()
	user2ID := uuid.New()

	mock.ExpectQuery("SELECT \\* FROM match WHERE user1_id = \\$1 AND user2_id = \\$2").
		WithArgs(user1ID, user2ID).
		WillReturnError(pgx.ErrTxClosed)

	match, err := repo.GetByUsers(context.Background(), user1ID, user2ID)
	assert.Error(t, err)
	assert.Nil(t, match)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestMatchRepository_GetUserMatches(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	repo := NewMatchRepository(mock)

	userID := uuid.New()
	limit := 10
	offset := 0

	expectedMatches := []domain.Match{
		{
			User1ID:   uuid.New(),
			User2ID:   userID,
			IsActive:  true,
			MatchedAt: time.Now(),
		},
		{
			User1ID:   userID,
			User2ID:   uuid.New(),
			IsActive:  true,
			MatchedAt: time.Now().Add(-time.Hour),
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

	mock.ExpectQuery("SELECT user1_id, user2_id, is_active, matched_at FROM match WHERE \\(user1_id = \\$1 OR user2_id = \\$1\\) AND is_active = true ORDER BY matched_at DESC LIMIT \\$2 OFFSET \\$3").
		WithArgs(userID, limit, offset).
		WillReturnRows(rows)

	matches, err := repo.GetUserMatches(context.Background(), userID, limit, offset)
	assert.NoError(t, err)
	assert.Len(t, matches, 2)
	assert.Equal(t, expectedMatches[0].User1ID, matches[0].User1ID)
	assert.Equal(t, expectedMatches[1].User1ID, matches[1].User1ID)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestMatchRepository_GetUserMatches_Empty(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	repo := NewMatchRepository(mock)

	userID := uuid.New()
	limit := 10
	offset := 0

	rows := mock.NewRows([]string{
		"user1_id", "user2_id", "is_active", "matched_at",
	})

	mock.ExpectQuery("SELECT user1_id, user2_id, is_active, matched_at FROM match WHERE \\(user1_id = \\$1 OR user2_id = \\$1\\) AND is_active = true ORDER BY matched_at DESC LIMIT \\$2 OFFSET \\$3").
		WithArgs(userID, limit, offset).
		WillReturnRows(rows)

	matches, err := repo.GetUserMatches(context.Background(), userID, limit, offset)
	assert.NoError(t, err)
	assert.Empty(t, matches)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestMatchRepository_GetUserMatches_QueryError(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	repo := NewMatchRepository(mock)

	userID := uuid.New()
	limit := 10
	offset := 0

	mock.ExpectQuery("SELECT user1_id, user2_id, is_active, matched_at FROM match WHERE \\(user1_id = \\$1 OR user2_id = \\$1\\) AND is_active = true ORDER BY matched_at DESC LIMIT \\$2 OFFSET \\$3").
		WithArgs(userID, limit, offset).
		WillReturnError(pgx.ErrTxClosed)

	matches, err := repo.GetUserMatches(context.Background(), userID, limit, offset)
	assert.Error(t, err)
	assert.Nil(t, matches)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestMatchRepository_UpdateActive(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	repo := NewMatchRepository(mock)

	user1ID := uuid.MustParse("00000000-0000-0000-0000-000000000001")
	user2ID := uuid.MustParse("00000000-0000-0000-0000-000000000002")
	isActive := false

	mock.ExpectExec("UPDATE match SET is_active = \\$1 WHERE user1_id = \\$2 AND user2_id = \\$3").
		WithArgs(isActive, user1ID, user2ID).
		WillReturnResult(pgxmock.NewResult("UPDATE", 1))

	err = repo.UpdateActive(context.Background(), user1ID, user2ID, isActive)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestMatchRepository_UpdateActive_WithReorderedUsers(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	repo := NewMatchRepository(mock)

	// user2ID is lexicographically smaller than user1ID
	user1ID := uuid.MustParse("00000000-0000-0000-0000-000000000002")
	user2ID := uuid.MustParse("00000000-0000-0000-0000-000000000001")
	isActive := false

	// Expect the users to be reordered in the query
	mock.ExpectExec("UPDATE match SET is_active = \\$1 WHERE user1_id = \\$2 AND user2_id = \\$3").
		WithArgs(isActive, user2ID, user1ID).
		WillReturnResult(pgxmock.NewResult("UPDATE", 1))

	err = repo.UpdateActive(context.Background(), user1ID, user2ID, isActive)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestMatchRepository_UpdateActive_Error(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	repo := NewMatchRepository(mock)

	user1ID := uuid.New()
	user2ID := uuid.New()
	isActive := false

	mock.ExpectExec("UPDATE match SET is_active = \\$1 WHERE user1_id = \\$2 AND user2_id = \\$3").
		WithArgs(isActive, user1ID, user2ID).
		WillReturnError(pgx.ErrTxClosed)

	err = repo.UpdateActive(context.Background(), user1ID, user2ID, isActive)
	assert.Error(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestMatchRepository_Delete(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	repo := NewMatchRepository(mock)

	user1ID := uuid.MustParse("00000000-0000-0000-0000-000000000001")
	user2ID := uuid.MustParse("00000000-0000-0000-0000-000000000002")

	mock.ExpectExec("DELETE FROM match WHERE user1_id = \\$1 AND user2_id = \\$2").
		WithArgs(user1ID, user2ID).
		WillReturnResult(pgxmock.NewResult("DELETE", 1))

	err = repo.Delete(context.Background(), user1ID, user2ID)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestMatchRepository_Delete_WithReorderedUsers(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	repo := NewMatchRepository(mock)

	// user2ID is lexicographically smaller than user1ID
	user1ID := uuid.MustParse("00000000-0000-0000-0000-000000000002")
	user2ID := uuid.MustParse("00000000-0000-0000-0000-000000000001")

	// Expect the users to be reordered in the query
	mock.ExpectExec("DELETE FROM match WHERE user1_id = \\$1 AND user2_id = \\$2").
		WithArgs(user2ID, user1ID).
		WillReturnResult(pgxmock.NewResult("DELETE", 1))

	err = repo.Delete(context.Background(), user1ID, user2ID)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestMatchRepository_Delete_Error(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	repo := NewMatchRepository(mock)

	user1ID := uuid.New()
	user2ID := uuid.New()

	mock.ExpectExec("DELETE FROM match WHERE user1_id = \\$1 AND user2_id = \\$2").
		WithArgs(user1ID, user2ID).
		WillReturnError(pgx.ErrTxClosed)

	err = repo.Delete(context.Background(), user1ID, user2ID)
	assert.Error(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestMatchRepository_CheckMutualLike(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	repo := NewMatchRepository(mock)

	user1ID := uuid.New()
	user2ID := uuid.New()

	rows := mock.NewRows([]string{"exists"}).AddRow(true)

	mock.ExpectQuery("SELECT EXISTS\\( SELECT 1 FROM swipe s1 INNER JOIN swipe s2 ON s1.swiper_user_id = s2.target_user_id AND s1.target_user_id = s2.swiper_user_id WHERE s1.swiper_user_id = \\$1 AND s1.target_user_id = \\$2 AND s2.swiper_user_id = \\$2 AND s2.target_user_id = \\$1 AND s1.swipe_type = 'like' AND s2.swipe_type = 'like' \\)").
		WithArgs(user1ID, user2ID).
		WillReturnRows(rows)

	exists, err := repo.CheckMutualLike(context.Background(), user1ID, user2ID)
	assert.NoError(t, err)
	assert.True(t, exists)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestMatchRepository_CheckMutualLike_False(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	repo := NewMatchRepository(mock)

	user1ID := uuid.New()
	user2ID := uuid.New()

	rows := mock.NewRows([]string{"exists"}).AddRow(false)

	mock.ExpectQuery("SELECT EXISTS\\( SELECT 1 FROM swipe s1 INNER JOIN swipe s2 ON s1.swiper_user_id = s2.target_user_id AND s1.target_user_id = s2.swiper_user_id WHERE s1.swiper_user_id = \\$1 AND s1.target_user_id = \\$2 AND s2.swiper_user_id = \\$2 AND s2.target_user_id = \\$1 AND s1.swipe_type = 'like' AND s2.swipe_type = 'like' \\)").
		WithArgs(user1ID, user2ID).
		WillReturnRows(rows)

	exists, err := repo.CheckMutualLike(context.Background(), user1ID, user2ID)
	assert.NoError(t, err)
	assert.False(t, exists)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestMatchRepository_CheckMutualLike_Error(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	repo := NewMatchRepository(mock)

	user1ID := uuid.New()
	user2ID := uuid.New()

	mock.ExpectQuery("SELECT EXISTS\\( SELECT 1 FROM swipe s1 INNER JOIN swipe s2 ON s1.swiper_user_id = s2.target_user_id AND s1.target_user_id = s2.swiper_user_id WHERE s1.swiper_user_id = \\$1 AND s1.target_user_id = \\$2 AND s2.swiper_user_id = \\$2 AND s2.target_user_id = \\$1 AND s1.swipe_type = 'like' AND s2.swipe_type = 'like' \\)").
		WithArgs(user1ID, user2ID).
		WillReturnError(pgx.ErrTxClosed)

	exists, err := repo.CheckMutualLike(context.Background(), user1ID, user2ID)
	assert.Error(t, err)
	assert.False(t, exists)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestMatchRepository_Integration_CRUD(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	repo := NewMatchRepository(mock)

	user1ID := uuid.MustParse("00000000-0000-0000-0000-000000000001")
	user2ID := uuid.MustParse("00000000-0000-0000-0000-000000000002")

	// Test Create
	match := &domain.Match{
		User1ID:  user1ID,
		User2ID:  user2ID,
		IsActive: true,
	}

	rows := mock.NewRows([]string{"matched_at"}).AddRow(time.Now())

	mock.ExpectQuery("INSERT INTO match").
		WithArgs(user1ID, user2ID, match.IsActive).
		WillReturnRows(rows)

	err = repo.Create(context.Background(), match)
	assert.NoError(t, err)

	// Test GetByUsers
	expectedMatch := &domain.Match{
		User1ID:   user1ID,
		User2ID:   user2ID,
		IsActive:  true,
		MatchedAt: time.Now(),
	}

	rows = mock.NewRows([]string{
		"user1_id", "user2_id", "is_active", "matched_at",
	}).AddRow(
		expectedMatch.User1ID, expectedMatch.User2ID, expectedMatch.IsActive, expectedMatch.MatchedAt,
	)

	mock.ExpectQuery("SELECT \\* FROM match WHERE user1_id = \\$1 AND user2_id = \\$2").
		WithArgs(user1ID, user2ID).
		WillReturnRows(rows)

	retrievedMatch, err := repo.GetByUsers(context.Background(), user1ID, user2ID)
	assert.NoError(t, err)
	assert.Equal(t, expectedMatch.User1ID, retrievedMatch.User1ID)

	// Test UpdateActive
	mock.ExpectExec("UPDATE match SET is_active = \\$1 WHERE user1_id = \\$2 AND user2_id = \\$3").
		WithArgs(false, user1ID, user2ID).
		WillReturnResult(pgxmock.NewResult("UPDATE", 1))

	err = repo.UpdateActive(context.Background(), user1ID, user2ID, false)
	assert.NoError(t, err)

	// Test Delete
	mock.ExpectExec("DELETE FROM match WHERE user1_id = \\$1 AND user2_id = \\$2").
		WithArgs(user1ID, user2ID).
		WillReturnResult(pgxmock.NewResult("DELETE", 1))

	err = repo.Delete(context.Background(), user1ID, user2ID)
	assert.NoError(t, err)

	assert.NoError(t, mock.ExpectationsWereMet())
}
