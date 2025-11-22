package postgres

import (
	"context"
	"testing"
	"time"

	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/domain/constants"
	domain "github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/domain/entities"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/pashagolub/pgxmock/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSwipeRepository_Create(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	repo := NewSwipeRepository(mock)

	swipe := &domain.Swipe{
		SwiperUserID: uuid.New(),
		TargetUserID: uuid.New(),
		SwipeType:    constants.SwipeTypeLike,
	}

	rows := mock.NewRows([]string{"created_at"}).
		AddRow(time.Now())

	mock.ExpectQuery("INSERT INTO swipe").
		WithArgs(swipe.SwiperUserID, swipe.TargetUserID, swipe.SwipeType).
		WillReturnRows(rows)

	err = repo.Create(context.Background(), swipe)
	assert.NoError(t, err)
	assert.False(t, swipe.CreatedAt.IsZero())
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSwipeRepository_Create_Error(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	repo := NewSwipeRepository(mock)

	swipe := &domain.Swipe{
		SwiperUserID: uuid.New(),
		TargetUserID: uuid.New(),
		SwipeType:    constants.SwipeTypeLike,
	}

	mock.ExpectQuery("INSERT INTO swipe").
		WithArgs(swipe.SwiperUserID, swipe.TargetUserID, swipe.SwipeType).
		WillReturnError(pgx.ErrTxClosed)

	err = repo.Create(context.Background(), swipe)
	assert.Error(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSwipeRepository_GetBySwiperAndTarget(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	repo := NewSwipeRepository(mock)

	swiperID := uuid.New()
	targetID := uuid.New()

	expectedSwipe := &domain.Swipe{
		SwiperUserID: swiperID,
		TargetUserID: targetID,
		SwipeType:    constants.SwipeTypeLike,
		CreatedAt:    time.Now(),
	}

	rows := mock.NewRows([]string{
		"swiper_user_id", "target_user_id", "swipe_type", "created_at",
	}).AddRow(
		expectedSwipe.SwiperUserID, expectedSwipe.TargetUserID, expectedSwipe.SwipeType, expectedSwipe.CreatedAt,
	)

	mock.ExpectQuery("SELECT \\* FROM swipe WHERE swiper_user_id = \\$1 AND target_user_id = \\$2").
		WithArgs(swiperID, targetID).
		WillReturnRows(rows)

	swipe, err := repo.GetBySwiperAndTarget(context.Background(), swiperID, targetID)
	assert.NoError(t, err)
	assert.Equal(t, expectedSwipe.SwiperUserID, swipe.SwiperUserID)
	assert.Equal(t, expectedSwipe.TargetUserID, swipe.TargetUserID)
	assert.Equal(t, expectedSwipe.SwipeType, swipe.SwipeType)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSwipeRepository_GetBySwiperAndTarget_NotFound(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	repo := NewSwipeRepository(mock)

	swiperID := uuid.New()
	targetID := uuid.New()

	mock.ExpectQuery("SELECT \\* FROM swipe WHERE swiper_user_id = \\$1 AND target_user_id = \\$2").
		WithArgs(swiperID, targetID).
		WillReturnError(pgx.ErrNoRows)

	swipe, err := repo.GetBySwiperAndTarget(context.Background(), swiperID, targetID)
	assert.NoError(t, err)
	assert.Nil(t, swipe)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSwipeRepository_GetBySwiperAndTarget_Error(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	repo := NewSwipeRepository(mock)

	swiperID := uuid.New()
	targetID := uuid.New()

	mock.ExpectQuery("SELECT \\* FROM swipe WHERE swiper_user_id = \\$1 AND target_user_id = \\$2").
		WithArgs(swiperID, targetID).
		WillReturnError(pgx.ErrTxClosed)

	swipe, err := repo.GetBySwiperAndTarget(context.Background(), swiperID, targetID)
	assert.Error(t, err)
	assert.Nil(t, swipe)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSwipeRepository_GetSwipesBySwiper(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	repo := NewSwipeRepository(mock)

	swiperID := uuid.New()
	limit := 10
	offset := 0

	expectedSwipes := []domain.Swipe{
		{
			SwiperUserID: swiperID,
			TargetUserID: uuid.New(),
			SwipeType:    constants.SwipeTypeLike,
			CreatedAt:    time.Now(),
		},
		{
			SwiperUserID: swiperID,
			TargetUserID: uuid.New(),
			SwipeType:    constants.SwipeTypeDislike,
			CreatedAt:    time.Now().Add(-time.Hour),
		},
	}

	rows := mock.NewRows([]string{
		"swiper_user_id", "target_user_id", "swipe_type", "created_at",
	})
	for _, swipe := range expectedSwipes {
		rows.AddRow(
			swipe.SwiperUserID, swipe.TargetUserID, swipe.SwipeType, swipe.CreatedAt,
		)
	}

	mock.ExpectQuery("SELECT \\* FROM swipe WHERE swiper_user_id = \\$1 ORDER BY created_at DESC LIMIT \\$2 OFFSET \\$3").
		WithArgs(swiperID, limit, offset).
		WillReturnRows(rows)

	swipes, err := repo.GetSwipesBySwiper(context.Background(), swiperID, limit, offset)
	assert.NoError(t, err)
	assert.Len(t, swipes, 2)
	assert.Equal(t, expectedSwipes[0].TargetUserID, swipes[0].TargetUserID)
	assert.Equal(t, expectedSwipes[1].TargetUserID, swipes[1].TargetUserID)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSwipeRepository_GetSwipesBySwiper_Empty(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	repo := NewSwipeRepository(mock)

	swiperID := uuid.New()
	limit := 10
	offset := 0

	rows := mock.NewRows([]string{
		"swiper_user_id", "target_user_id", "swipe_type", "created_at",
	})

	mock.ExpectQuery("SELECT \\* FROM swipe WHERE swiper_user_id = \\$1 ORDER BY created_at DESC LIMIT \\$2 OFFSET \\$3").
		WithArgs(swiperID, limit, offset).
		WillReturnRows(rows)

	swipes, err := repo.GetSwipesBySwiper(context.Background(), swiperID, limit, offset)
	assert.NoError(t, err)
	assert.Empty(t, swipes)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSwipeRepository_GetSwipesBySwiper_QueryError(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	repo := NewSwipeRepository(mock)

	swiperID := uuid.New()
	limit := 10
	offset := 0

	mock.ExpectQuery("SELECT \\* FROM swipe WHERE swiper_user_id = \\$1 ORDER BY created_at DESC LIMIT \\$2 OFFSET \\$3").
		WithArgs(swiperID, limit, offset).
		WillReturnError(pgx.ErrTxClosed)

	swipes, err := repo.GetSwipesBySwiper(context.Background(), swiperID, limit, offset)
	assert.Error(t, err)
	assert.Nil(t, swipes)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSwipeRepository_GetSwipesByTarget(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	repo := NewSwipeRepository(mock)

	targetID := uuid.New()
	limit := 10
	offset := 0

	expectedSwipes := []domain.Swipe{
		{
			SwiperUserID: uuid.New(),
			TargetUserID: targetID,
			SwipeType:    constants.SwipeTypeLike,
			CreatedAt:    time.Now(),
		},
		{
			SwiperUserID: uuid.New(),
			TargetUserID: targetID,
			SwipeType:    constants.SwipeTypeSuperLike,
			CreatedAt:    time.Now().Add(-time.Hour),
		},
	}

	rows := mock.NewRows([]string{
		"swiper_user_id", "target_user_id", "swipe_type", "created_at",
	})
	for _, swipe := range expectedSwipes {
		rows.AddRow(
			swipe.SwiperUserID, swipe.TargetUserID, swipe.SwipeType, swipe.CreatedAt,
		)
	}

	mock.ExpectQuery("SELECT \\* FROM swipe WHERE target_user_id = \\$1 ORDER BY created_at DESC LIMIT \\$2 OFFSET \\$3").
		WithArgs(targetID, limit, offset).
		WillReturnRows(rows)

	swipes, err := repo.GetSwipesByTarget(context.Background(), targetID, limit, offset)
	assert.NoError(t, err)
	assert.Len(t, swipes, 2)
	assert.Equal(t, expectedSwipes[0].SwiperUserID, swipes[0].SwiperUserID)
	assert.Equal(t, expectedSwipes[1].SwiperUserID, swipes[1].SwiperUserID)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSwipeRepository_GetSwipesByTarget_Empty(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	repo := NewSwipeRepository(mock)

	targetID := uuid.New()
	limit := 10
	offset := 0

	rows := mock.NewRows([]string{
		"swiper_user_id", "target_user_id", "swipe_type", "created_at",
	})

	mock.ExpectQuery("SELECT \\* FROM swipe WHERE target_user_id = \\$1 ORDER BY created_at DESC LIMIT \\$2 OFFSET \\$3").
		WithArgs(targetID, limit, offset).
		WillReturnRows(rows)

	swipes, err := repo.GetSwipesByTarget(context.Background(), targetID, limit, offset)
	assert.NoError(t, err)
	assert.Empty(t, swipes)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSwipeRepository_GetSwipesByTarget_QueryError(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	repo := NewSwipeRepository(mock)

	targetID := uuid.New()
	limit := 10
	offset := 0

	mock.ExpectQuery("SELECT \\* FROM swipe WHERE target_user_id = \\$1 ORDER BY created_at DESC LIMIT \\$2 OFFSET \\$3").
		WithArgs(targetID, limit, offset).
		WillReturnError(pgx.ErrTxClosed)

	swipes, err := repo.GetSwipesByTarget(context.Background(), targetID, limit, offset)
	assert.Error(t, err)
	assert.Nil(t, swipes)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSwipeRepository_Exists(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	repo := NewSwipeRepository(mock)

	swiperID := uuid.New()
	targetID := uuid.New()

	rows := mock.NewRows([]string{"exists"}).AddRow(true)

	mock.ExpectQuery("SELECT EXISTS\\(SELECT 1 FROM swipe WHERE swiper_user_id = \\$1 AND target_user_id = \\$2\\)").
		WithArgs(swiperID, targetID).
		WillReturnRows(rows)

	exists, err := repo.Exists(context.Background(), swiperID, targetID)
	assert.NoError(t, err)
	assert.True(t, exists)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSwipeRepository_Exists_False(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	repo := NewSwipeRepository(mock)

	swiperID := uuid.New()
	targetID := uuid.New()

	rows := mock.NewRows([]string{"exists"}).AddRow(false)

	mock.ExpectQuery("SELECT EXISTS\\(SELECT 1 FROM swipe WHERE swiper_user_id = \\$1 AND target_user_id = \\$2\\)").
		WithArgs(swiperID, targetID).
		WillReturnRows(rows)

	exists, err := repo.Exists(context.Background(), swiperID, targetID)
	assert.NoError(t, err)
	assert.False(t, exists)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSwipeRepository_Exists_Error(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	repo := NewSwipeRepository(mock)

	swiperID := uuid.New()
	targetID := uuid.New()

	mock.ExpectQuery("SELECT EXISTS\\(SELECT 1 FROM swipe WHERE swiper_user_id = \\$1 AND target_user_id = \\$2\\)").
		WithArgs(swiperID, targetID).
		WillReturnError(pgx.ErrTxClosed)

	exists, err := repo.Exists(context.Background(), swiperID, targetID)
	assert.Error(t, err)
	assert.False(t, exists)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSwipeRepository_GetSwipesStats(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	repo := NewSwipeRepository(mock)

	userID := uuid.New()
	expectedLikes := 5
	expectedDislikes := 3
	expectedSuperLikes := 2

	rows := mock.NewRows([]string{"likes_count", "dislikes_count", "super_likes_count"}).
		AddRow(expectedLikes, expectedDislikes, expectedSuperLikes)

	mock.ExpectQuery("SELECT COUNT\\(CASE WHEN swipe_type = 'like' THEN 1 END\\) as likes_count, COUNT\\(CASE WHEN swipe_type = 'dislike' THEN 1 END\\) as dislikes_count, COUNT\\(CASE WHEN swipe_type = 'super_like' THEN 1 END\\) as super_likes_count FROM swipe WHERE swiper_user_id = \\$1").
		WithArgs(userID).
		WillReturnRows(rows)

	likes, dislikes, superLikes, err := repo.GetSwipesStats(context.Background(), userID)
	assert.NoError(t, err)
	assert.Equal(t, expectedLikes, likes)
	assert.Equal(t, expectedDislikes, dislikes)
	assert.Equal(t, expectedSuperLikes, superLikes)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSwipeRepository_GetSwipesStats_Error(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	repo := NewSwipeRepository(mock)

	userID := uuid.New()

	mock.ExpectQuery("SELECT COUNT\\(CASE WHEN swipe_type = 'like' THEN 1 END\\) as likes_count, COUNT\\(CASE WHEN swipe_type = 'dislike' THEN 1 END\\) as dislikes_count, COUNT\\(CASE WHEN swipe_type = 'super_like' THEN 1 END\\) as super_likes_count FROM swipe WHERE swiper_user_id = \\$1").
		WithArgs(userID).
		WillReturnError(pgx.ErrTxClosed)

	likes, dislikes, superLikes, err := repo.GetSwipesStats(context.Background(), userID)
	assert.Error(t, err)
	assert.Equal(t, 0, likes)
	assert.Equal(t, 0, dislikes)
	assert.Equal(t, 0, superLikes)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSwipeRepository_GetMutualLikes(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	repo := NewSwipeRepository(mock)

	userID := uuid.New()
	expectedUserIDs := []uuid.UUID{uuid.New(), uuid.New(), uuid.New()}

	rows := mock.NewRows([]string{"swiper_user_id"})
	for _, id := range expectedUserIDs {
		rows.AddRow(id)
	}

	mock.ExpectQuery("SELECT s1.swiper_user_id FROM swipe s1 INNER JOIN swipe s2 ON s1.swiper_user_id = s2.target_user_id AND s1.target_user_id = s2.swiper_user_id WHERE s1.target_user_id = \\$1 AND s1.swipe_type = 'like' AND s2.swipe_type = 'like' AND NOT EXISTS \\( SELECT 1 FROM match m WHERE \\(m.user1_id = s1.swiper_user_id AND m.user2_id = s1.target_user_id\\) OR \\(m.user1_id = s1.target_user_id AND m.user2_id = s1.swiper_user_id\\) \\)").
		WithArgs(userID).
		WillReturnRows(rows)

	userIDs, err := repo.GetMutualLikes(context.Background(), userID)
	assert.NoError(t, err)
	assert.Len(t, userIDs, 3)
	assert.Equal(t, expectedUserIDs[0], userIDs[0])
	assert.Equal(t, expectedUserIDs[1], userIDs[1])
	assert.Equal(t, expectedUserIDs[2], userIDs[2])
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSwipeRepository_GetMutualLikes_Empty(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	repo := NewSwipeRepository(mock)

	userID := uuid.New()

	rows := mock.NewRows([]string{"swiper_user_id"})

	mock.ExpectQuery("SELECT s1.swiper_user_id FROM swipe s1 INNER JOIN swipe s2 ON s1.swiper_user_id = s2.target_user_id AND s1.target_user_id = s2.swiper_user_id WHERE s1.target_user_id = \\$1 AND s1.swipe_type = 'like' AND s2.swipe_type = 'like' AND NOT EXISTS \\( SELECT 1 FROM match m WHERE \\(m.user1_id = s1.swiper_user_id AND m.user2_id = s1.target_user_id\\) OR \\(m.user1_id = s1.target_user_id AND m.user2_id = s1.swiper_user_id\\) \\)").
		WithArgs(userID).
		WillReturnRows(rows)

	userIDs, err := repo.GetMutualLikes(context.Background(), userID)
	assert.NoError(t, err)
	assert.Empty(t, userIDs)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSwipeRepository_GetMutualLikes_QueryError(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	repo := NewSwipeRepository(mock)

	userID := uuid.New()

	mock.ExpectQuery("SELECT s1.swiper_user_id FROM swipe s1 INNER JOIN swipe s2 ON s1.swiper_user_id = s2.target_user_id AND s1.target_user_id = s2.swiper_user_id WHERE s1.target_user_id = \\$1 AND s1.swipe_type = 'like' AND s2.swipe_type = 'like' AND NOT EXISTS \\( SELECT 1 FROM match m WHERE \\(m.user1_id = s1.swiper_user_id AND m.user2_id = s1.target_user_id\\) OR \\(m.user1_id = s1.target_user_id AND m.user2_id = s1.swiper_user_id\\) \\)").
		WithArgs(userID).
		WillReturnError(pgx.ErrTxClosed)

	userIDs, err := repo.GetMutualLikes(context.Background(), userID)
	assert.Error(t, err)
	assert.Nil(t, userIDs)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSwipeRepository_GetMutualLikes_ScanError(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	repo := NewSwipeRepository(mock)

	userID := uuid.New()

	// Добавляем лишнюю колонку чтобы вызвать ошибку
	rows := mock.NewRows([]string{"swiper_user_id", "extra_column"}).
		AddRow(uuid.New(), "extra_value")

	mock.ExpectQuery("SELECT s1.swiper_user_id FROM swipe s1 INNER JOIN swipe s2 ON s1.swiper_user_id = s2.target_user_id AND s1.target_user_id = s2.swiper_user_id WHERE s1.target_user_id = \\$1 AND s1.swipe_type = 'like' AND s2.swipe_type = 'like' AND NOT EXISTS \\( SELECT 1 FROM match m WHERE \\(m.user1_id = s1.swiper_user_id AND m.user2_id = s1.target_user_id\\) OR \\(m.user1_id = s1.target_user_id AND m.user2_id = s1.swiper_user_id\\) \\)").
		WithArgs(userID).
		WillReturnRows(rows)

	userIDs, err := repo.GetMutualLikes(context.Background(), userID)
	assert.Error(t, err)
	assert.Nil(t, userIDs)
	assert.NoError(t, mock.ExpectationsWereMet())
}
