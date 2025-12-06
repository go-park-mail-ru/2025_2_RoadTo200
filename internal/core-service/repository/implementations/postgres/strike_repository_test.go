package postgres

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/core-service/domain/constants"
	domain "github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/core-service/domain/entities"
	errs "github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/domain/errors"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/pashagolub/pgxmock/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStrikeRepository_CreateStrike(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	repo := NewStrikeRepository(mock)

	strike := &domain.Strike{
		ReporterID:   uuid.New(),
		TargetUserID: uuid.New(),
		Reason:       "Test reason",
		Type:         constants.StrikeTypeHarassment,
	}

	mock.ExpectExec("INSERT INTO strike").
		WithArgs(pgxmock.AnyArg(), strike.ReporterID, strike.TargetUserID, strike.Reason, strike.Type).
		WillReturnResult(pgxmock.NewResult("INSERT", 1))

	err = repo.CreateStrike(context.Background(), strike)
	assert.NoError(t, err)
	assert.NotEqual(t, uuid.Nil, strike.ID)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestStrikeRepository_CreateStrike_Error(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	repo := NewStrikeRepository(mock)

	strike := &domain.Strike{
		ReporterID:   uuid.New(),
		TargetUserID: uuid.New(),
		Reason:       "Test reason",
		Type:         constants.StrikeTypeHarassment,
	}

	mock.ExpectExec("INSERT INTO strike").
		WithArgs(pgxmock.AnyArg(), strike.ReporterID, strike.TargetUserID, strike.Reason, strike.Type).
		WillReturnError(errors.New("db error"))

	err = repo.CreateStrike(context.Background(), strike)
	assert.Error(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestStrikeRepository_GetStrikeByID(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	repo := NewStrikeRepository(mock)

	strikeID := uuid.New()
	expectedStrike := &domain.Strike{
		ID:           strikeID,
		ReporterID:   uuid.New(),
		TargetUserID: uuid.New(),
		Type:         constants.StrikeTypeHarassment,
		Reason:       "Test reason",
		Status:       constants.StrikeStatusPending,
		CreatedAt:    time.Now(),
		UpdatedAt:    nil,
	}

	rows := mock.NewRows([]string{"id", "reporter_id", "target_user_id", "type", "reason", "status", "created_at", "updated_at", "moderator_id", "moderator_note"}).
		AddRow(expectedStrike.ID, expectedStrike.ReporterID, expectedStrike.TargetUserID, expectedStrike.Type, expectedStrike.Reason, expectedStrike.Status, expectedStrike.CreatedAt, expectedStrike.UpdatedAt, expectedStrike.ModeratorID, expectedStrike.ModeratorNote)

	mock.ExpectQuery("SELECT id, reporter_id, target_user_id, type, reason, status, created_at, updated_at, moderator_id, moderator_note FROM strike WHERE id = \\$1").
		WithArgs(strikeID.String()).
		WillReturnRows(rows)

	strike, err := repo.GetStrikeByID(context.Background(), strikeID.String())
	assert.NoError(t, err)
	assert.Equal(t, expectedStrike.ID, strike.ID)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestStrikeRepository_GetStrikeByID_NotFound(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	repo := NewStrikeRepository(mock)

	strikeID := uuid.New()

	mock.ExpectQuery("SELECT id, reporter_id, target_user_id, type, reason, status, created_at, updated_at, moderator_id, moderator_note FROM strike WHERE id = \\$1").
		WithArgs(strikeID.String()).
		WillReturnError(pgx.ErrNoRows)

	strike, err := repo.GetStrikeByID(context.Background(), strikeID.String())
	assert.ErrorIs(t, err, errs.ErrStrikeNotFound)
	assert.Nil(t, strike)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestStrikeRepository_GetStrikesByUserID(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	repo := NewStrikeRepository(mock)

	targetUserID := uuid.New()
	strike := domain.Strike{
		ID:           uuid.New(),
		ReporterID:   uuid.New(),
		TargetUserID: targetUserID,
		Type:         constants.StrikeTypeHarassment,
		Reason:       "Test reason",
		Status:       constants.StrikeStatusPending,
		CreatedAt:    time.Now(),
	}

	rows := mock.NewRows([]string{"id", "reporter_id", "target_user_id", "type", "reason", "status", "created_at", "updated_at", "moderator_id", "moderator_note"}).
		AddRow(strike.ID, strike.ReporterID, strike.TargetUserID, strike.Type, strike.Reason, strike.Status, strike.CreatedAt, strike.UpdatedAt, strike.ModeratorID, strike.ModeratorNote)

	mock.ExpectQuery("SELECT id, reporter_id, target_user_id, type, reason, status, created_at, updated_at, moderator_id, moderator_note FROM strike WHERE target_user_id = \\$1 ORDER BY created_at DESC LIMIT \\$2 OFFSET \\$3").
		WithArgs(targetUserID.String(), 10, 0).
		WillReturnRows(rows)

	strikes, err := repo.GetStrikesByUserID(context.Background(), targetUserID.String(), 10, 0)
	assert.NoError(t, err)
	assert.Len(t, strikes, 1)
	assert.Equal(t, strike.ID, strikes[0].ID)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestStrikeRepository_GetStrikesByType(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	repo := NewStrikeRepository(mock)

	strikeType := constants.StrikeTypeHarassment
	strike := domain.Strike{
		ID:           uuid.New(),
		ReporterID:   uuid.New(),
		TargetUserID: uuid.New(),
		Type:         strikeType,
		Reason:       "Test reason",
		Status:       constants.StrikeStatusPending,
		CreatedAt:    time.Now(),
	}

	rows := mock.NewRows([]string{"id", "reporter_id", "target_user_id", "type", "reason", "status", "created_at", "updated_at", "moderator_id", "moderator_note"}).
		AddRow(strike.ID, strike.ReporterID, strike.TargetUserID, strike.Type, strike.Reason, strike.Status, strike.CreatedAt, strike.UpdatedAt, strike.ModeratorID, strike.ModeratorNote)

	mock.ExpectQuery("SELECT id, reporter_id, target_user_id, type, reason, status, created_at, updated_at, moderator_id, moderator_note FROM strike WHERE type = \\$1 ORDER BY created_at DESC LIMIT \\$2 OFFSET \\$3").
		WithArgs(strikeType, 10, 0).
		WillReturnRows(rows)

	strikes, err := repo.GetStrikesByType(context.Background(), strikeType, 10, 0)
	assert.NoError(t, err)
	assert.Len(t, strikes, 1)
	assert.Equal(t, strike.ID, strikes[0].ID)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestStrikeRepository_GetStrikesByDateRange(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	repo := NewStrikeRepository(mock)

	from := time.Now().Add(-24 * time.Hour)
	to := time.Now()
	strike := domain.Strike{
		ID:           uuid.New(),
		ReporterID:   uuid.New(),
		TargetUserID: uuid.New(),
		Type:         constants.StrikeTypeHarassment,
		Reason:       "Test reason",
		Status:       constants.StrikeStatusPending,
		CreatedAt:    time.Now(),
	}

	rows := mock.NewRows([]string{"id", "reporter_id", "target_user_id", "type", "reason", "status", "created_at", "updated_at", "moderator_id", "moderator_note"}).
		AddRow(strike.ID, strike.ReporterID, strike.TargetUserID, strike.Type, strike.Reason, strike.Status, strike.CreatedAt, strike.UpdatedAt, strike.ModeratorID, strike.ModeratorNote)

	mock.ExpectQuery("SELECT id, reporter_id, target_user_id, type, reason, status, created_at, updated_at, moderator_id, moderator_note FROM strike WHERE created_at BETWEEN \\$1 AND \\$2 ORDER BY created_at DESC LIMIT \\$3 OFFSET \\$4").
		WithArgs(from, to, 10, 0).
		WillReturnRows(rows)

	strikes, err := repo.GetStrikesByDateRange(context.Background(), from, to, 10, 0)
	assert.NoError(t, err)
	assert.Len(t, strikes, 1)
	assert.Equal(t, strike.ID, strikes[0].ID)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestStrikeRepository_UpdateStrikeStatus(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	repo := NewStrikeRepository(mock)

	strikeID := uuid.New()
	status := constants.StrikeStatusResolved

	mock.ExpectExec("UPDATE strike SET status = \\$1, updated_at = \\$2 WHERE id = \\$3").
		WithArgs(status, pgxmock.AnyArg(), strikeID.String()).
		WillReturnResult(pgxmock.NewResult("UPDATE", 1))

	err = repo.UpdateStrikeStatus(context.Background(), strikeID.String(), status)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestStrikeRepository_UpdateStrikeStatus_NotFound(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	repo := NewStrikeRepository(mock)

	strikeID := uuid.New()
	status := constants.StrikeStatusResolved

	mock.ExpectExec("UPDATE strike SET status = \\$1, updated_at = \\$2 WHERE id = \\$3").
		WithArgs(status, pgxmock.AnyArg(), strikeID.String()).
		WillReturnResult(pgxmock.NewResult("UPDATE", 0))

	err = repo.UpdateStrikeStatus(context.Background(), strikeID.String(), status)
	assert.ErrorIs(t, err, errs.ErrStrikeNotFound)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestStrikeRepository_DeleteStrike(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	repo := NewStrikeRepository(mock)

	strikeID := uuid.New()

	mock.ExpectExec("DELETE FROM strike WHERE id = \\$1").
		WithArgs(strikeID.String()).
		WillReturnResult(pgxmock.NewResult("DELETE", 1))

	err = repo.DeleteStrike(context.Background(), strikeID.String())
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestStrikeRepository_HasActiveStrikeFromUser(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	repo := NewStrikeRepository(mock)

	reporterID := uuid.New()
	targetUserID := uuid.New()

	rows := mock.NewRows([]string{"count"}).AddRow(1)

	mock.ExpectQuery("SELECT COUNT\\(\\*\\) FROM strike WHERE reporter_id = \\$1 AND target_user_id = \\$2 AND status = \\$3").
		WithArgs(reporterID.String(), targetUserID.String(), constants.StrikeStatusPending).
		WillReturnRows(rows)

	exists, err := repo.HasActiveStrikeFromUser(context.Background(), reporterID.String(), targetUserID.String())
	assert.NoError(t, err)
	assert.True(t, exists)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestStrikeRepository_HasActiveStrikeFromUser_False(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	repo := NewStrikeRepository(mock)

	reporterID := uuid.New()
	targetUserID := uuid.New()

	rows := mock.NewRows([]string{"count"}).AddRow(0)

	mock.ExpectQuery("SELECT COUNT\\(\\*\\) FROM strike WHERE reporter_id = \\$1 AND target_user_id = \\$2 AND status = \\$3").
		WithArgs(reporterID.String(), targetUserID.String(), constants.StrikeStatusPending).
		WillReturnRows(rows)

	exists, err := repo.HasActiveStrikeFromUser(context.Background(), reporterID.String(), targetUserID.String())
	assert.NoError(t, err)
	assert.False(t, exists)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestStrikeRepository_GetStrikesByUserID_Error(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	repo := NewStrikeRepository(mock)

	targetUserID := uuid.New()

	mock.ExpectQuery("SELECT id, reporter_id, target_user_id, type, reason, status, created_at, updated_at, moderator_id, moderator_note FROM strike WHERE target_user_id = \\$1 ORDER BY created_at DESC LIMIT \\$2 OFFSET \\$3").
		WithArgs(targetUserID.String(), 10, 0).
		WillReturnError(errors.New("db error"))

	strikes, err := repo.GetStrikesByUserID(context.Background(), targetUserID.String(), 10, 0)
	assert.Error(t, err)
	assert.Nil(t, strikes)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestStrikeRepository_UpdateStrikeStatus_Error(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	repo := NewStrikeRepository(mock)

	strikeID := uuid.New()
	status := constants.StrikeStatusResolved

	mock.ExpectExec("UPDATE strike SET status = \\$1, updated_at = \\$2 WHERE id = \\$3").
		WithArgs(status, pgxmock.AnyArg(), strikeID.String()).
		WillReturnError(errors.New("db error"))

	err = repo.UpdateStrikeStatus(context.Background(), strikeID.String(), status)
	assert.Error(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}
