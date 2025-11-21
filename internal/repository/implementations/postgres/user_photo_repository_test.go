package postgres

import (
	"testing"
	"time"

	domain "github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/domain/entities"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
	"github.com/pashagolub/pgxmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUserPhotoRepository_Create(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	repo := NewUserPhotoRepository(mock)

	photo := &domain.UserPhoto{
		UserID:       uuid.New(),
		PhotoURL:     "https://example.com/photo.jpg",
		DisplayOrder: 1,
		IsApproved:   true,
	}

	rows := mock.NewRows([]string{"id", "created_at"}).
		AddRow(uuid.New(), time.Now())

	mock.ExpectQuery("INSERT INTO user_photo").
		WithArgs(photo.UserID, photo.PhotoURL, photo.DisplayOrder, photo.IsApproved).
		WillReturnRows(rows)

	err = repo.Create(photo)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserPhotoRepository_GetByID(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	repo := NewUserPhotoRepository(mock)

	photoID := uuid.New()
	userID := uuid.New()

	expectedPhoto := &domain.UserPhoto{
		ID:           photoID,
		UserID:       userID,
		PhotoURL:     "https://example.com/photo.jpg",
		DisplayOrder: 1,
		IsApproved:   true,
		CreatedAt:    time.Now(),
	}

	rows := mock.NewRows([]string{
		"id", "user_id", "photo_url", "display_order", "is_approved", "created_at",
	}).AddRow(
		expectedPhoto.ID, expectedPhoto.UserID, expectedPhoto.PhotoURL,
		expectedPhoto.DisplayOrder, expectedPhoto.IsApproved, expectedPhoto.CreatedAt,
	)

	mock.ExpectQuery("SELECT \\* FROM user_photo WHERE id = \\$1").
		WithArgs(photoID).
		WillReturnRows(rows)

	photo, err := repo.GetByID(photoID)
	assert.NoError(t, err)
	assert.Equal(t, expectedPhoto.ID, photo.ID)
	assert.Equal(t, expectedPhoto.UserID, photo.UserID)
	assert.Equal(t, expectedPhoto.PhotoURL, photo.PhotoURL)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserPhotoRepository_GetByID_NotFound(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	repo := NewUserPhotoRepository(mock)

	photoID := uuid.New()

	mock.ExpectQuery("SELECT \\* FROM user_photo WHERE id = \\$1").
		WithArgs(photoID).
		WillReturnError(pgx.ErrNoRows)

	photo, err := repo.GetByID(photoID)
	assert.NoError(t, err)
	assert.Nil(t, photo)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserPhotoRepository_GetByUserID(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	repo := NewUserPhotoRepository(mock)

	userID := uuid.New()
	expectedPhotos := []domain.UserPhoto{
		{
			ID:           uuid.New(),
			UserID:       userID,
			PhotoURL:     "https://example.com/photo1.jpg",
			DisplayOrder: 1,
			IsApproved:   true,
			CreatedAt:    time.Now(),
		},
		{
			ID:           uuid.New(),
			UserID:       userID,
			PhotoURL:     "https://example.com/photo2.jpg",
			DisplayOrder: 2,
			IsApproved:   false,
			CreatedAt:    time.Now(),
		},
	}

	rows := mock.NewRows([]string{
		"id", "user_id", "photo_url", "display_order", "is_approved", "created_at",
	})
	for _, photo := range expectedPhotos {
		rows.AddRow(
			photo.ID, photo.UserID, photo.PhotoURL, photo.DisplayOrder,
			photo.IsApproved, photo.CreatedAt,
		)
	}

	mock.ExpectQuery("SELECT \\* FROM user_photo WHERE user_id = \\$1 ORDER BY display_order ASC").
		WithArgs(userID).
		WillReturnRows(rows)

	photos, err := repo.GetByUserID(userID)
	assert.NoError(t, err)
	assert.Len(t, photos, 2)
	assert.Equal(t, expectedPhotos[0].ID, photos[0].ID)
	assert.Equal(t, expectedPhotos[1].ID, photos[1].ID)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserPhotoRepository_GetByUserID_Empty(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	repo := NewUserPhotoRepository(mock)

	userID := uuid.New()

	rows := mock.NewRows([]string{
		"id", "user_id", "photo_url", "display_order", "is_approved", "created_at",
	})

	mock.ExpectQuery("SELECT \\* FROM user_photo WHERE user_id = \\$1 ORDER BY display_order ASC").
		WithArgs(userID).
		WillReturnRows(rows)

	photos, err := repo.GetByUserID(userID)
	assert.NoError(t, err)
	assert.Empty(t, photos)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserPhotoRepository_Update(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	repo := NewUserPhotoRepository(mock)

	photo := &domain.UserPhoto{
		ID:           uuid.New(),
		PhotoURL:     "https://example.com/updated-photo.jpg",
		DisplayOrder: 2,
		IsApproved:   false,
	}

	mock.ExpectExec("UPDATE user_photo").
		WithArgs(photo.PhotoURL, photo.DisplayOrder, photo.IsApproved, photo.ID).
		WillReturnResult(pgxmock.NewResult("UPDATE", 1))

	err = repo.Update(photo)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserPhotoRepository_Delete(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	repo := NewUserPhotoRepository(mock)

	photoID := uuid.New()

	mock.ExpectExec("DELETE FROM user_photo WHERE id = \\$1").
		WithArgs(photoID).
		WillReturnResult(pgxmock.NewResult("DELETE", 1))

	err = repo.Delete(photoID)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserPhotoRepository_UpdateDisplayOrder(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	repo := NewUserPhotoRepository(mock)

	userID := uuid.New()
	photos := []domain.UserPhoto{
		{
			PhotoURL:     "https://example.com/photo1.jpg",
			DisplayOrder: 1,
			IsApproved:   true,
		},
		{
			PhotoURL:     "https://example.com/photo2.jpg",
			DisplayOrder: 2,
			IsApproved:   false,
		},
	}

	mock.ExpectBegin()
	mock.ExpectExec("DELETE FROM user_photo WHERE user_id = \\$1").
		WithArgs(userID).
		WillReturnResult(pgxmock.NewResult("DELETE", 2))

	for _, photo := range photos {
		mock.ExpectExec("INSERT INTO user_photo").
			WithArgs(userID, photo.PhotoURL, photo.DisplayOrder, photo.IsApproved).
			WillReturnResult(pgxmock.NewResult("INSERT", 1))
	}

	mock.ExpectCommit()

	err = repo.UpdateDisplayOrder(userID, photos)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserPhotoRepository_UpdateDisplayOrder_EmptyPhotos(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	repo := NewUserPhotoRepository(mock)

	userID := uuid.New()

	mock.ExpectBegin()
	mock.ExpectExec("DELETE FROM user_photo WHERE user_id = \\$1").
		WithArgs(userID).
		WillReturnResult(pgxmock.NewResult("DELETE", 0))
	mock.ExpectCommit()

	err = repo.UpdateDisplayOrder(userID, []domain.UserPhoto{})
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserPhotoRepository_Create_Error(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	repo := NewUserPhotoRepository(mock)

	photo := &domain.UserPhoto{
		UserID:       uuid.New(),
		PhotoURL:     "https://example.com/photo.jpg",
		DisplayOrder: 1,
		IsApproved:   true,
	}

	mock.ExpectQuery("INSERT INTO user_photo").
		WithArgs(photo.UserID, photo.PhotoURL, photo.DisplayOrder, photo.IsApproved).
		WillReturnError(pgx.ErrNoRows)

	err = repo.Create(photo)
	assert.Error(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserPhotoRepository_Update_Error(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	repo := NewUserPhotoRepository(mock)

	photo := &domain.UserPhoto{
		ID:           uuid.New(),
		PhotoURL:     "https://example.com/updated-photo.jpg",
		DisplayOrder: 2,
		IsApproved:   false,
	}

	mock.ExpectExec("UPDATE user_photo").
		WithArgs(photo.PhotoURL, photo.DisplayOrder, photo.IsApproved, photo.ID).
		WillReturnError(pgx.ErrNoRows)

	err = repo.Update(photo)
	assert.Error(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserPhotoRepository_UpdateDisplayOrder_TransactionError(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	repo := NewUserPhotoRepository(mock)

	userID := uuid.New()
	photos := []domain.UserPhoto{
		{
			PhotoURL:     "https://example.com/photo1.jpg",
			DisplayOrder: 1,
			IsApproved:   true,
		},
	}

	mock.ExpectBegin().WillReturnError(pgx.ErrTxClosed)

	err = repo.UpdateDisplayOrder(userID, photos)
	assert.Error(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}
