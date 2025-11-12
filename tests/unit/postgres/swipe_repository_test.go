package postgres

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

func TestSwipeRepository_Create(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	repo := repository.NewSwipeRepository(mock)

	swiperID := uuid.New()
	targetID := uuid.New()

	swipe := &domain.Swipe{
		SwiperUserID: swiperID,
		TargetUserID: targetID,
		SwipeType:    "like",
	}

	rows := mock.NewRows([]string{"created_at"}).
		AddRow(time.Now())

	mock.ExpectQuery("INSERT INTO swipe").
		WithArgs(swipe.SwiperUserID, swipe.TargetUserID, swipe.SwipeType).
		WillReturnRows(rows)

	err = repo.Create(swipe)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSwipeRepository_GetBySwiperAndTarget(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	repo := repository.NewSwipeRepository(mock)

	swiperID := uuid.New()
	targetID := uuid.New()

	expectedSwipe := &domain.Swipe{
		SwiperUserID: swiperID,
		TargetUserID: targetID,
		SwipeType:    "like",
		CreatedAt:    time.Now(),
	}

	rows := mock.NewRows([]string{
		"swiper_user_id", "target_user_id", "swipe_type", "created_at",
	}).AddRow(
		expectedSwipe.SwiperUserID, expectedSwipe.TargetUserID,
		expectedSwipe.SwipeType, expectedSwipe.CreatedAt,
	)

	mock.ExpectQuery("SELECT \\* FROM swipe WHERE swiper_user_id = \\$1 AND target_user_id = \\$2").
		WithArgs(swiperID, targetID).
		WillReturnRows(rows)

	swipe, err := repo.GetBySwiperAndTarget(swiperID, targetID)
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

	repo := repository.NewSwipeRepository(mock)

	swiperID := uuid.New()
	targetID := uuid.New()

	mock.ExpectQuery("SELECT \\* FROM swipe WHERE swiper_user_id = \\$1 AND target_user_id = \\$2").
		WithArgs(swiperID, targetID).
		WillReturnError(pgx.ErrNoRows)

	swipe, err := repo.GetBySwiperAndTarget(swiperID, targetID)
	assert.NoError(t, err)
	assert.Nil(t, swipe)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSwipeRepository_GetSwipesBySwiper(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	repo := repository.NewSwipeRepository(mock)

	swiperID := uuid.New()
	limit := 10
	offset := 0

	expectedSwipes := []domain.Swipe{
		{
			SwiperUserID: swiperID,
			TargetUserID: uuid.New(),
			SwipeType:    "like",
			CreatedAt:    time.Now(),
		},
		{
			SwiperUserID: swiperID,
			TargetUserID: uuid.New(),
			SwipeType:    "dislike",
			CreatedAt:    time.Now(),
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

	swipes, err := repo.GetSwipesBySwiper(swiperID, limit, offset)
	assert.NoError(t, err)
	assert.Len(t, swipes, 2)
	assert.Equal(t, expectedSwipes[0].SwiperUserID, swipes[0].SwiperUserID)
	assert.Equal(t, expectedSwipes[1].SwiperUserID, swipes[1].SwiperUserID)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSwipeRepository_GetSwipesBySwiper_Empty(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	repo := repository.NewSwipeRepository(mock)

	swiperID := uuid.New()
	limit := 10
	offset := 0

	rows := mock.NewRows([]string{
		"swiper_user_id", "target_user_id", "swipe_type", "created_at",
	})

	mock.ExpectQuery("SELECT \\* FROM swipe WHERE swiper_user_id = \\$1 ORDER BY created_at DESC LIMIT \\$2 OFFSET \\$3").
		WithArgs(swiperID, limit, offset).
		WillReturnRows(rows)

	swipes, err := repo.GetSwipesBySwiper(swiperID, limit, offset)
	assert.NoError(t, err)
	assert.Empty(t, swipes)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSwipeRepository_GetSwipesByTarget(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	repo := repository.NewSwipeRepository(mock)

	targetID := uuid.New()
	limit := 10
	offset := 0

	expectedSwipes := []domain.Swipe{
		{
			SwiperUserID: uuid.New(),
			TargetUserID: targetID,
			SwipeType:    "like",
			CreatedAt:    time.Now(),
		},
		{
			SwiperUserID: uuid.New(),
			TargetUserID: targetID,
			SwipeType:    "super_like",
			CreatedAt:    time.Now(),
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

	swipes, err := repo.GetSwipesByTarget(targetID, limit, offset)
	assert.NoError(t, err)
	assert.Len(t, swipes, 2)
	assert.Equal(t, expectedSwipes[0].TargetUserID, swipes[0].TargetUserID)
	assert.Equal(t, expectedSwipes[1].TargetUserID, swipes[1].TargetUserID)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSwipeRepository_Exists(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	repo := repository.NewSwipeRepository(mock)

	swiperID := uuid.New()
	targetID := uuid.New()
	exists := true

	rows := mock.NewRows([]string{"exists"}).AddRow(exists)

	mock.ExpectQuery("SELECT EXISTS").
		WithArgs(swiperID, targetID).
		WillReturnRows(rows)

	result, err := repo.Exists(swiperID, targetID)
	assert.NoError(t, err)
	assert.True(t, result)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSwipeRepository_Exists_NotFound(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	repo := repository.NewSwipeRepository(mock)

	swiperID := uuid.New()
	targetID := uuid.New()
	exists := false

	rows := mock.NewRows([]string{"exists"}).AddRow(exists)

	mock.ExpectQuery("SELECT EXISTS").
		WithArgs(swiperID, targetID).
		WillReturnRows(rows)

	result, err := repo.Exists(swiperID, targetID)
	assert.NoError(t, err)
	assert.False(t, result)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSwipeRepository_GetSwipesStats(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	repo := repository.NewSwipeRepository(mock)

	userID := uuid.New()
	likesCount := 5
	dislikesCount := 3
	superLikesCount := 2

	rows := mock.NewRows([]string{"likes_count", "dislikes_count", "super_likes_count"}).
		AddRow(likesCount, dislikesCount, superLikesCount)

	mock.ExpectQuery("SELECT").
		WithArgs(userID).
		WillReturnRows(rows)

	likes, dislikes, superLikes, err := repo.GetSwipesStats(userID)
	assert.NoError(t, err)
	assert.Equal(t, likesCount, likes)
	assert.Equal(t, dislikesCount, dislikes)
	assert.Equal(t, superLikesCount, superLikes)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSwipeRepository_GetSwipesStats_Zero(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	repo := repository.NewSwipeRepository(mock)

	userID := uuid.New()
	likesCount := 0
	dislikesCount := 0
	superLikesCount := 0

	rows := mock.NewRows([]string{"likes_count", "dislikes_count", "super_likes_count"}).
		AddRow(likesCount, dislikesCount, superLikesCount)

	mock.ExpectQuery("SELECT").
		WithArgs(userID).
		WillReturnRows(rows)

	likes, dislikes, superLikes, err := repo.GetSwipesStats(userID)
	assert.NoError(t, err)
	assert.Equal(t, likesCount, likes)
	assert.Equal(t, dislikesCount, dislikes)
	assert.Equal(t, superLikesCount, superLikes)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSwipeRepository_GetMutualLikes(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	repo := repository.NewSwipeRepository(mock)

	userID := uuid.New()
	expectedUserIDs := []uuid.UUID{uuid.New(), uuid.New()}

	rows := mock.NewRows([]string{"swiper_user_id"})
	for _, id := range expectedUserIDs {
		rows.AddRow(id)
	}

	mock.ExpectQuery("SELECT s1.swiper_user_id").
		WithArgs(userID).
		WillReturnRows(rows)

	userIDs, err := repo.GetMutualLikes(userID)
	assert.NoError(t, err)
	assert.Len(t, userIDs, 2)
	assert.Equal(t, expectedUserIDs[0], userIDs[0])
	assert.Equal(t, expectedUserIDs[1], userIDs[1])
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSwipeRepository_GetMutualLikes_Empty(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	repo := repository.NewSwipeRepository(mock)

	userID := uuid.New()

	rows := mock.NewRows([]string{"swiper_user_id"})

	mock.ExpectQuery("SELECT s1.swiper_user_id").
		WithArgs(userID).
		WillReturnRows(rows)

	userIDs, err := repo.GetMutualLikes(userID)
	assert.NoError(t, err)
	assert.Empty(t, userIDs)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSwipeRepository_Create_Error(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	repo := repository.NewSwipeRepository(mock)

	swiperID := uuid.New()
	targetID := uuid.New()

	swipe := &domain.Swipe{
		SwiperUserID: swiperID,
		TargetUserID: targetID,
		SwipeType:    "like",
	}

	mock.ExpectQuery("INSERT INTO swipe").
		WithArgs(swipe.SwiperUserID, swipe.TargetUserID, swipe.SwipeType).
		WillReturnError(pgx.ErrNoRows)

	err = repo.Create(swipe)
	assert.Error(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSwipeRepository_GetSwipesStats_Error(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	repo := repository.NewSwipeRepository(mock)

	userID := uuid.New()

	mock.ExpectQuery("SELECT").
		WithArgs(userID).
		WillReturnError(pgx.ErrNoRows)

	likes, dislikes, superLikes, err := repo.GetSwipesStats(userID)
	assert.Error(t, err)
	assert.Equal(t, 0, likes)
	assert.Equal(t, 0, dislikes)
	assert.Equal(t, 0, superLikes)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSwipeRepository_GetMutualLikes_Error(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	repo := repository.NewSwipeRepository(mock)

	userID := uuid.New()

	mock.ExpectQuery("SELECT s1.swiper_user_id").
		WithArgs(userID).
		WillReturnError(pgx.ErrNoRows)

	userIDs, err := repo.GetMutualLikes(userID)
	assert.Error(t, err)
	assert.Nil(t, userIDs)
	assert.NoError(t, mock.ExpectationsWereMet())
}
