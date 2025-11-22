package postgres

import (
	"context"
	"testing"
	"time"

	domain "github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/domain/entities"
	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/tests/mocks"
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

	log := mocks.NewMockLogger()
	repo := NewUserPhotoRepository(mock, log)

	userID := uuid.New()
	photo := &domain.UserPhoto{
		UserID:       userID,
		PhotoURL:     "https://example.com/photo1.jpg",
		DisplayOrder: 1,
		IsApproved:   true,
	}

	rows := mock.NewRows([]string{"id", "created_at"}).
		AddRow(uuid.New(), time.Now())

	mock.ExpectQuery("INSERT INTO user_photo").
		WithArgs(photo.UserID, photo.PhotoURL, photo.DisplayOrder, photo.IsApproved).
		WillReturnRows(rows)

	err = repo.Create(context.Background(), photo)
	assert.NoError(t, err)
	assert.NotEqual(t, uuid.Nil, photo.ID)
	assert.NoError(t, mock.ExpectationsWereMet())
	assert.True(t, log.WasCalled("Trace"))
}

func TestUserPhotoRepository_Create_Error(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	log := mocks.NewMockLogger()
	repo := NewUserPhotoRepository(mock, log)

	userID := uuid.New()
	photo := &domain.UserPhoto{
		UserID:       userID,
		PhotoURL:     "https://example.com/photo1.jpg",
		DisplayOrder: 1,
		IsApproved:   true,
	}

	mock.ExpectQuery("INSERT INTO user_photo").
		WithArgs(photo.UserID, photo.PhotoURL, photo.DisplayOrder, photo.IsApproved).
		WillReturnError(pgx.ErrTxClosed)

	err = repo.Create(context.Background(), photo)
	assert.Error(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserPhotoRepository_GetByID(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	log := mocks.NewMockLogger()
	repo := NewUserPhotoRepository(mock, log)

	photoID := uuid.New()
	userID := uuid.New()

	expectedPhoto := &domain.UserPhoto{
		ID:           photoID,
		UserID:       userID,
		PhotoURL:     "https://example.com/photo1.jpg",
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

	photo, err := repo.GetByID(context.Background(), photoID)
	assert.NoError(t, err)
	assert.Equal(t, expectedPhoto.ID, photo.ID)
	assert.Equal(t, expectedPhoto.UserID, photo.UserID)
	assert.Equal(t, expectedPhoto.PhotoURL, photo.PhotoURL)
	assert.NoError(t, mock.ExpectationsWereMet())
	assert.True(t, log.WasCalled("Trace"))
}

func TestUserPhotoRepository_GetByID_NotFound(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	log := mocks.NewMockLogger()
	repo := NewUserPhotoRepository(mock, log)

	photoID := uuid.New()

	mock.ExpectQuery("SELECT \\* FROM user_photo WHERE id = \\$1").
		WithArgs(photoID).
		WillReturnError(pgx.ErrNoRows)

	photo, err := repo.GetByID(context.Background(), photoID)
	assert.NoError(t, err)
	assert.Nil(t, photo)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserPhotoRepository_GetByID_Error(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	log := mocks.NewMockLogger()
	repo := NewUserPhotoRepository(mock, log)

	photoID := uuid.New()

	mock.ExpectQuery("SELECT \\* FROM user_photo WHERE id = \\$1").
		WithArgs(photoID).
		WillReturnError(pgx.ErrTxClosed)

	photo, err := repo.GetByID(context.Background(), photoID)
	assert.Error(t, err)
	assert.Nil(t, photo)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserPhotoRepository_GetByUserID(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	log := mocks.NewMockLogger()
	repo := NewUserPhotoRepository(mock, log)

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
			photo.ID, photo.UserID, photo.PhotoURL,
			photo.DisplayOrder, photo.IsApproved, photo.CreatedAt,
		)
	}

	mock.ExpectQuery("SELECT \\* FROM user_photo WHERE user_id = \\$1 ORDER BY display_order").
		WithArgs(userID).
		WillReturnRows(rows)

	photos, err := repo.GetByUserID(context.Background(), userID)
	assert.NoError(t, err)
	assert.Len(t, photos, 2)
	assert.Equal(t, expectedPhotos[0].ID, photos[0].ID)
	assert.Equal(t, expectedPhotos[1].ID, photos[1].ID)
	assert.NoError(t, mock.ExpectationsWereMet())
	assert.True(t, log.WasCalled("Trace"))
}

func TestUserPhotoRepository_GetByUserID_Empty(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	log := mocks.NewMockLogger()
	repo := NewUserPhotoRepository(mock, log)

	userID := uuid.New()

	rows := mock.NewRows([]string{
		"id", "user_id", "photo_url", "display_order", "is_approved", "created_at",
	})

	mock.ExpectQuery("SELECT \\* FROM user_photo WHERE user_id = \\$1 ORDER BY display_order").
		WithArgs(userID).
		WillReturnRows(rows)

	photos, err := repo.GetByUserID(context.Background(), userID)
	assert.NoError(t, err)
	assert.Empty(t, photos)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserPhotoRepository_GetByUserID_ScanError(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	log := mocks.NewMockLogger()
	repo := NewUserPhotoRepository(mock, log)

	userID := uuid.New()

	rows := mock.NewRows([]string{
		"id", "user_id", "photo_url", "display_order", "is_approved", "created_at",
	}).AddRow(
		nil, nil, nil, nil, nil, nil, // Invalid data to cause scan error
	)

	mock.ExpectQuery("SELECT \\* FROM user_photo WHERE user_id = \\$1 ORDER BY display_order").
		WithArgs(userID).
		WillReturnRows(rows)

	photos, err := repo.GetByUserID(context.Background(), userID)
	assert.NoError(t, err) // Note: In the actual code, scan errors are logged but not returned
	assert.Len(t, photos, 0)
	assert.NoError(t, mock.ExpectationsWereMet())
	assert.True(t, log.WasCalled("Errorf"))
}

func TestUserPhotoRepository_GetByUserID_QueryError(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	log := mocks.NewMockLogger()
	repo := NewUserPhotoRepository(mock, log)

	userID := uuid.New()

	mock.ExpectQuery("SELECT \\* FROM user_photo WHERE user_id = \\$1 ORDER BY display_order").
		WithArgs(userID).
		WillReturnError(pgx.ErrTxClosed)

	photos, err := repo.GetByUserID(context.Background(), userID)
	assert.Error(t, err)
	assert.Nil(t, photos)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserPhotoRepository_Update(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	log := mocks.NewMockLogger()
	repo := NewUserPhotoRepository(mock, log)

	photo := &domain.UserPhoto{
		ID:           uuid.New(),
		PhotoURL:     "https://example.com/updated.jpg",
		DisplayOrder: 3,
		IsApproved:   false,
	}

	mock.ExpectExec("UPDATE user_photo").
		WithArgs(photo.PhotoURL, photo.DisplayOrder, photo.IsApproved, photo.ID).
		WillReturnResult(pgxmock.NewResult("UPDATE", 1))

	err = repo.Update(context.Background(), photo)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
	assert.True(t, log.WasCalled("Trace"))
}

func TestUserPhotoRepository_Update_Error(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	log := mocks.NewMockLogger()
	repo := NewUserPhotoRepository(mock, log)

	photo := &domain.UserPhoto{
		ID:           uuid.New(),
		PhotoURL:     "https://example.com/updated.jpg",
		DisplayOrder: 3,
		IsApproved:   false,
	}

	mock.ExpectExec("UPDATE user_photo").
		WithArgs(photo.PhotoURL, photo.DisplayOrder, photo.IsApproved, photo.ID).
		WillReturnError(pgx.ErrTxClosed)

	err = repo.Update(context.Background(), photo)
	assert.Error(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserPhotoRepository_Delete(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	log := mocks.NewMockLogger()
	repo := NewUserPhotoRepository(mock, log)

	photoID := uuid.New()

	mock.ExpectExec("DELETE FROM user_photo WHERE id = \\$1").
		WithArgs(photoID).
		WillReturnResult(pgxmock.NewResult("DELETE", 1))

	err = repo.Delete(context.Background(), photoID)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
	assert.True(t, log.WasCalled("Trace"))
}

func TestUserPhotoRepository_Delete_Error(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	log := mocks.NewMockLogger()
	repo := NewUserPhotoRepository(mock, log)

	photoID := uuid.New()

	mock.ExpectExec("DELETE FROM user_photo WHERE id = \\$1").
		WithArgs(photoID).
		WillReturnError(pgx.ErrTxClosed)

	err = repo.Delete(context.Background(), photoID)
	assert.Error(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserPhotoRepository_UpdateDisplayOrder(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	log := mocks.NewMockLogger()
	repo := NewUserPhotoRepository(mock, log)

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

	// Expect two inserts
	mock.ExpectExec("INSERT INTO user_photo").
		WithArgs(userID, photos[0].PhotoURL, photos[0].DisplayOrder, photos[0].IsApproved).
		WillReturnResult(pgxmock.NewResult("INSERT", 1))
	mock.ExpectExec("INSERT INTO user_photo").
		WithArgs(userID, photos[1].PhotoURL, photos[1].DisplayOrder, photos[1].IsApproved).
		WillReturnResult(pgxmock.NewResult("INSERT", 1))

	mock.ExpectCommit()

	err = repo.UpdateDisplayOrder(context.Background(), userID, photos)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
	assert.True(t, log.WasCalled("Trace"))
}

func TestUserPhotoRepository_UpdateDisplayOrder_BeginError(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	log := mocks.NewMockLogger()
	repo := NewUserPhotoRepository(mock, log)

	userID := uuid.New()
	photos := []domain.UserPhoto{
		{
			PhotoURL:     "https://example.com/photo1.jpg",
			DisplayOrder: 1,
			IsApproved:   true,
		},
	}

	mock.ExpectBegin().WillReturnError(pgx.ErrTxClosed)

	err = repo.UpdateDisplayOrder(context.Background(), userID, photos)
	assert.Error(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserPhotoRepository_UpdateDisplayOrder_DeleteError(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	log := mocks.NewMockLogger()
	repo := NewUserPhotoRepository(mock, log)

	userID := uuid.New()
	photos := []domain.UserPhoto{
		{
			PhotoURL:     "https://example.com/photo1.jpg",
			DisplayOrder: 1,
			IsApproved:   true,
		},
	}

	mock.ExpectBegin()
	mock.ExpectExec("DELETE FROM user_photo WHERE user_id = \\$1").
		WithArgs(userID).
		WillReturnError(pgx.ErrTxClosed)
	mock.ExpectRollback()

	err = repo.UpdateDisplayOrder(context.Background(), userID, photos)
	assert.Error(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserPhotoRepository_UpdateDisplayOrder_InsertError(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	log := mocks.NewMockLogger()
	repo := NewUserPhotoRepository(mock, log)

	userID := uuid.New()
	photos := []domain.UserPhoto{
		{
			PhotoURL:     "https://example.com/photo1.jpg",
			DisplayOrder: 1,
			IsApproved:   true,
		},
	}

	mock.ExpectBegin()
	mock.ExpectExec("DELETE FROM user_photo WHERE user_id = \\$1").
		WithArgs(userID).
		WillReturnResult(pgxmock.NewResult("DELETE", 1))
	mock.ExpectExec("INSERT INTO user_photo").
		WithArgs(userID, photos[0].PhotoURL, photos[0].DisplayOrder, photos[0].IsApproved).
		WillReturnError(pgx.ErrTxClosed)
	mock.ExpectRollback()

	err = repo.UpdateDisplayOrder(context.Background(), userID, photos)
	assert.Error(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
	assert.True(t, log.WasCalled("Errorf"))
}

func TestUserPhotoRepository_UpdateDisplayOrder_CommitError(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	log := mocks.NewMockLogger()
	repo := NewUserPhotoRepository(mock, log)

	userID := uuid.New()
	photos := []domain.UserPhoto{
		{
			PhotoURL:     "https://example.com/photo1.jpg",
			DisplayOrder: 1,
			IsApproved:   true,
		},
	}

	mock.ExpectBegin()
	mock.ExpectExec("DELETE FROM user_photo WHERE user_id = \\$1").
		WithArgs(userID).
		WillReturnResult(pgxmock.NewResult("DELETE", 1))
	mock.ExpectExec("INSERT INTO user_photo").
		WithArgs(userID, photos[0].PhotoURL, photos[0].DisplayOrder, photos[0].IsApproved).
		WillReturnResult(pgxmock.NewResult("INSERT", 1))
	mock.ExpectCommit().WillReturnError(pgx.ErrTxClosed)

	err = repo.UpdateDisplayOrder(context.Background(), userID, photos)
	assert.Error(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}
