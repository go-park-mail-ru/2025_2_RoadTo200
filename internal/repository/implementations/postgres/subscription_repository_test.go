package postgres

// import (
// 	"context"
// 	"testing"
// 	"time"

// 	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/domain/constants"
// 	domain "github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/domain/entities"
// 	"github.com/google/uuid"
// 	"github.com/jackc/pgx/v5"
// 	"github.com/pashagolub/pgxmock/v4"
// 	"github.com/stretchr/testify/assert"
// 	"github.com/stretchr/testify/require"
// )

// func TestSubscriptionRepository_Create(t *testing.T) {
// 	mock, err := pgxmock.NewPool()
// 	require.NoError(t, err)
// 	defer mock.Close()

// 	repo := NewSubscriptionRepository(mock)

// 	subscription := &domain.Subscription{
// 		UserID:    uuid.New(),
// 		PlanType:  constants.PlanTypePremium,
// 		StartDate: time.Now(),
// 		EndDate:   time.Now().AddDate(0, 1, 0), // 1 month later
// 		IsActive:  true,
// 	}

// 	mock.ExpectExec("INSERT INTO subscription").
// 		WithArgs(
// 			subscription.UserID, subscription.PlanType, subscription.StartDate,
// 			subscription.EndDate, subscription.IsActive,
// 		).
// 		WillReturnResult(pgxmock.NewResult("INSERT", 1))

// 	err = repo.Create(context.Background(), subscription)
// 	assert.NoError(t, err)
// 	assert.NoError(t, mock.ExpectationsWereMet())
// }

// func TestSubscriptionRepository_Create_Error(t *testing.T) {
// 	mock, err := pgxmock.NewPool()
// 	require.NoError(t, err)
// 	defer mock.Close()

// 	repo := NewSubscriptionRepository(mock)

// 	subscription := &domain.Subscription{
// 		UserID:    uuid.New(),
// 		PlanType:  constants.PlanTypePremium,
// 		StartDate: time.Now(),
// 		EndDate:   time.Now().AddDate(0, 1, 0),
// 		IsActive:  true,
// 	}

// 	mock.ExpectExec("INSERT INTO subscription").
// 		WithArgs(
// 			subscription.UserID, subscription.PlanType, subscription.StartDate,
// 			subscription.EndDate, subscription.IsActive,
// 		).
// 		WillReturnError(pgx.ErrTxClosed)

// 	err = repo.Create(context.Background(), subscription)
// 	assert.Error(t, err)
// 	assert.NoError(t, mock.ExpectationsWereMet())
// }

// func TestSubscriptionRepository_GetByUserID(t *testing.T) {
// 	mock, err := pgxmock.NewPool()
// 	require.NoError(t, err)
// 	defer mock.Close()

// 	repo := NewSubscriptionRepository(mock)

// 	userID := uuid.New()
// 	expectedSubscription := &domain.Subscription{
// 		UserID:    userID,
// 		PlanType:  constants.PlanTypePremium,
// 		StartDate: time.Now().AddDate(0, -1, 0), // Started 1 month ago
// 		EndDate:   time.Now().AddDate(0, 0, 15), // Ends in 15 days
// 		IsActive:  true,
// 		CreatedAt: time.Now().AddDate(0, -1, 0),
// 	}

// 	rows := mock.NewRows([]string{
// 		"user_id", "plan_type", "start_date", "end_date", "is_active", "created_at",
// 	}).AddRow(
// 		expectedSubscription.UserID, expectedSubscription.PlanType, expectedSubscription.StartDate,
// 		expectedSubscription.EndDate, expectedSubscription.IsActive, expectedSubscription.CreatedAt,
// 	)

// 	mock.ExpectQuery("SELECT \\* FROM subscription WHERE user_id = \\$1").
// 		WithArgs(userID).
// 		WillReturnRows(rows)

// 	subscription, err := repo.GetByUserID(context.Background(), userID)
// 	assert.NoError(t, err)
// 	assert.Equal(t, expectedSubscription.UserID, subscription.UserID)
// 	assert.Equal(t, expectedSubscription.PlanType, subscription.PlanType)
// 	assert.Equal(t, expectedSubscription.IsActive, subscription.IsActive)
// 	assert.NoError(t, mock.ExpectationsWereMet())
// }

// func TestSubscriptionRepository_GetByUserID_NotFound(t *testing.T) {
// 	mock, err := pgxmock.NewPool()
// 	require.NoError(t, err)
// 	defer mock.Close()

// 	repo := NewSubscriptionRepository(mock)

// 	userID := uuid.New()

// 	mock.ExpectQuery("SELECT \\* FROM subscription WHERE user_id = \\$1").
// 		WithArgs(userID).
// 		WillReturnError(pgx.ErrNoRows)

// 	subscription, err := repo.GetByUserID(context.Background(), userID)
// 	assert.NoError(t, err)
// 	assert.Nil(t, subscription)
// 	assert.NoError(t, mock.ExpectationsWereMet())
// }

// func TestSubscriptionRepository_GetByUserID_Error(t *testing.T) {
// 	mock, err := pgxmock.NewPool()
// 	require.NoError(t, err)
// 	defer mock.Close()

// 	repo := NewSubscriptionRepository(mock)

// 	userID := uuid.New()

// 	mock.ExpectQuery("SELECT \\* FROM subscription WHERE user_id = \\$1").
// 		WithArgs(userID).
// 		WillReturnError(pgx.ErrTxClosed)

// 	subscription, err := repo.GetByUserID(context.Background(), userID)
// 	assert.Error(t, err)
// 	assert.Nil(t, subscription)
// 	assert.NoError(t, mock.ExpectationsWereMet())
// }

// func TestSubscriptionRepository_Update(t *testing.T) {
// 	mock, err := pgxmock.NewPool()
// 	require.NoError(t, err)
// 	defer mock.Close()

// 	repo := NewSubscriptionRepository(mock)

// 	subscription := &domain.Subscription{
// 		UserID:    uuid.New(),
// 		PlanType:  constants.PlanTypePremium,
// 		StartDate: time.Now(),
// 		EndDate:   time.Now().AddDate(0, 2, 0), // 2 months later
// 		IsActive:  false,
// 	}

// 	mock.ExpectExec("UPDATE subscription").
// 		WithArgs(
// 			subscription.PlanType, subscription.StartDate, subscription.EndDate,
// 			subscription.IsActive, subscription.UserID,
// 		).
// 		WillReturnResult(pgxmock.NewResult("UPDATE", 1))

// 	err = repo.Update(context.Background(), subscription)
// 	assert.NoError(t, err)
// 	assert.NoError(t, mock.ExpectationsWereMet())
// }

// func TestSubscriptionRepository_Update_Error(t *testing.T) {
// 	mock, err := pgxmock.NewPool()
// 	require.NoError(t, err)
// 	defer mock.Close()

// 	repo := NewSubscriptionRepository(mock)

// 	subscription := &domain.Subscription{
// 		UserID:    uuid.New(),
// 		PlanType:  constants.PlanTypePremium,
// 		StartDate: time.Now(),
// 		EndDate:   time.Now().AddDate(0, 2, 0),
// 		IsActive:  false,
// 	}

// 	mock.ExpectExec("UPDATE subscription").
// 		WithArgs(
// 			subscription.PlanType, subscription.StartDate, subscription.EndDate,
// 			subscription.IsActive, subscription.UserID,
// 		).
// 		WillReturnError(pgx.ErrTxClosed)

// 	err = repo.Update(context.Background(), subscription)
// 	assert.Error(t, err)
// 	assert.NoError(t, mock.ExpectationsWereMet())
// }

// func TestSubscriptionRepository_Delete(t *testing.T) {
// 	mock, err := pgxmock.NewPool()
// 	require.NoError(t, err)
// 	defer mock.Close()

// 	repo := NewSubscriptionRepository(mock)

// 	userID := uuid.New()

// 	mock.ExpectExec("DELETE FROM subscription WHERE user_id = \\$1").
// 		WithArgs(userID).
// 		WillReturnResult(pgxmock.NewResult("DELETE", 1))

// 	err = repo.Delete(context.Background(), userID)
// 	assert.NoError(t, err)
// 	assert.NoError(t, mock.ExpectationsWereMet())
// }

// func TestSubscriptionRepository_Delete_Error(t *testing.T) {
// 	mock, err := pgxmock.NewPool()
// 	require.NoError(t, err)
// 	defer mock.Close()

// 	repo := NewSubscriptionRepository(mock)

// 	userID := uuid.New()

// 	mock.ExpectExec("DELETE FROM subscription WHERE user_id = \\$1").
// 		WithArgs(userID).
// 		WillReturnError(pgx.ErrTxClosed)

// 	err = repo.Delete(context.Background(), userID)
// 	assert.Error(t, err)
// 	assert.NoError(t, mock.ExpectationsWereMet())
// }

// func TestSubscriptionRepository_GetActiveSubscription(t *testing.T) {
// 	mock, err := pgxmock.NewPool()
// 	require.NoError(t, err)
// 	defer mock.Close()

// 	repo := NewSubscriptionRepository(mock)

// 	userID := uuid.New()
// 	currentTime := time.Now()
// 	expectedSubscription := &domain.Subscription{
// 		UserID:    userID,
// 		PlanType:  constants.PlanTypePremium,
// 		StartDate: currentTime.AddDate(0, -1, 0),
// 		EndDate:   currentTime.AddDate(0, 0, 15), // Active for 15 more days
// 		IsActive:  true,
// 		CreatedAt: currentTime.AddDate(0, -1, 0),
// 	}

// 	rows := mock.NewRows([]string{
// 		"user_id", "plan_type", "start_date", "end_date", "is_active", "created_at",
// 	}).AddRow(
// 		expectedSubscription.UserID, expectedSubscription.PlanType, expectedSubscription.StartDate,
// 		expectedSubscription.EndDate, expectedSubscription.IsActive, expectedSubscription.CreatedAt,
// 	)

// 	mock.ExpectQuery("SELECT \\* FROM subscription WHERE user_id = \\$1 AND is_active = true AND end_date > \\$2").
// 		WithArgs(userID, pgxmock.AnyArg()).
// 		WillReturnRows(rows)

// 	subscription, err := repo.GetActiveSubscription(context.Background(), userID)
// 	assert.NoError(t, err)
// 	assert.Equal(t, expectedSubscription.UserID, subscription.UserID)
// 	assert.Equal(t, expectedSubscription.PlanType, subscription.PlanType)
// 	assert.True(t, subscription.IsActive)
// 	assert.True(t, subscription.EndDate.After(currentTime))
// 	assert.NoError(t, mock.ExpectationsWereMet())
// }

// func TestSubscriptionRepository_GetActiveSubscription_NotFound(t *testing.T) {
// 	mock, err := pgxmock.NewPool()
// 	require.NoError(t, err)
// 	defer mock.Close()

// 	repo := NewSubscriptionRepository(mock)

// 	userID := uuid.New()

// 	mock.ExpectQuery("SELECT \\* FROM subscription WHERE user_id = \\$1 AND is_active = true AND end_date > \\$2").
// 		WithArgs(userID, pgxmock.AnyArg()).
// 		WillReturnError(pgx.ErrNoRows)

// 	subscription, err := repo.GetActiveSubscription(context.Background(), userID)
// 	assert.NoError(t, err)
// 	assert.Nil(t, subscription)
// 	assert.NoError(t, mock.ExpectationsWereMet())
// }

// func TestSubscriptionRepository_GetActiveSubscription_Expired(t *testing.T) {
// 	mock, err := pgxmock.NewPool()
// 	require.NoError(t, err)
// 	defer mock.Close()

// 	repo := NewSubscriptionRepository(mock)

// 	userID := uuid.New()

// 	// Mock no rows returned because subscription is expired
// 	mock.ExpectQuery("SELECT \\* FROM subscription WHERE user_id = \\$1 AND is_active = true AND end_date > \\$2").
// 		WithArgs(userID, pgxmock.AnyArg()).
// 		WillReturnError(pgx.ErrNoRows)

// 	subscription, err := repo.GetActiveSubscription(context.Background(), userID)
// 	assert.NoError(t, err)
// 	assert.Nil(t, subscription)
// 	assert.NoError(t, mock.ExpectationsWereMet())
// }

// func TestSubscriptionRepository_GetActiveSubscription_Inactive(t *testing.T) {
// 	mock, err := pgxmock.NewPool()
// 	require.NoError(t, err)
// 	defer mock.Close()

// 	repo := NewSubscriptionRepository(mock)

// 	userID := uuid.New()

// 	// Mock no rows returned because subscription is inactive
// 	mock.ExpectQuery("SELECT \\* FROM subscription WHERE user_id = \\$1 AND is_active = true AND end_date > \\$2").
// 		WithArgs(userID, pgxmock.AnyArg()).
// 		WillReturnError(pgx.ErrNoRows)

// 	subscription, err := repo.GetActiveSubscription(context.Background(), userID)
// 	assert.NoError(t, err)
// 	assert.Nil(t, subscription)
// 	assert.NoError(t, mock.ExpectationsWereMet())
// }

// func TestSubscriptionRepository_GetActiveSubscription_Error(t *testing.T) {
// 	mock, err := pgxmock.NewPool()
// 	require.NoError(t, err)
// 	defer mock.Close()

// 	repo := NewSubscriptionRepository(mock)

// 	userID := uuid.New()

// 	mock.ExpectQuery("SELECT \\* FROM subscription WHERE user_id = \\$1 AND is_active = true AND end_date > \\$2").
// 		WithArgs(userID, pgxmock.AnyArg()).
// 		WillReturnError(pgx.ErrTxClosed)

// 	subscription, err := repo.GetActiveSubscription(context.Background(), userID)
// 	assert.Error(t, err)
// 	assert.Nil(t, subscription)
// 	assert.NoError(t, mock.ExpectationsWereMet())
// }

// func TestSubscriptionRepository_Integration_CRUD(t *testing.T) {
// 	mock, err := pgxmock.NewPool()
// 	require.NoError(t, err)
// 	defer mock.Close()

// 	repo := NewSubscriptionRepository(mock)

// 	userID := uuid.New()

// 	// Test Create
// 	subscription := &domain.Subscription{
// 		UserID:    userID,
// 		PlanType:  constants.PlanTypePremium,
// 		StartDate: time.Now(),
// 		EndDate:   time.Now().AddDate(0, 1, 0),
// 		IsActive:  true,
// 	}

// 	mock.ExpectExec("INSERT INTO subscription").
// 		WithArgs(
// 			subscription.UserID, subscription.PlanType, subscription.StartDate,
// 			subscription.EndDate, subscription.IsActive,
// 		).
// 		WillReturnResult(pgxmock.NewResult("INSERT", 1))

// 	err = repo.Create(context.Background(), subscription)
// 	assert.NoError(t, err)

// 	// Test GetByUserID
// 	expectedSubscription := &domain.Subscription{
// 		UserID:    userID,
// 		PlanType:  constants.PlanTypePremium,
// 		StartDate: subscription.StartDate,
// 		EndDate:   subscription.EndDate,
// 		IsActive:  true,
// 		CreatedAt: time.Now(),
// 	}

// 	rows := mock.NewRows([]string{
// 		"user_id", "plan_type", "start_date", "end_date", "is_active", "created_at",
// 	}).AddRow(
// 		expectedSubscription.UserID, expectedSubscription.PlanType, expectedSubscription.StartDate,
// 		expectedSubscription.EndDate, expectedSubscription.IsActive, expectedSubscription.CreatedAt,
// 	)

// 	mock.ExpectQuery("SELECT \\* FROM subscription WHERE user_id = \\$1").
// 		WithArgs(userID).
// 		WillReturnRows(rows)

// 	retrievedSubscription, err := repo.GetByUserID(context.Background(), userID)
// 	assert.NoError(t, err)
// 	assert.Equal(t, expectedSubscription.UserID, retrievedSubscription.UserID)

// 	// Test Update
// 	updatedSubscription := &domain.Subscription{
// 		UserID:    userID,
// 		PlanType:  constants.PlanTypePremium,
// 		StartDate: time.Now(),
// 		EndDate:   time.Now().AddDate(0, 2, 0),
// 		IsActive:  false,
// 	}

// 	mock.ExpectExec("UPDATE subscription").
// 		WithArgs(
// 			updatedSubscription.PlanType, updatedSubscription.StartDate, updatedSubscription.EndDate,
// 			updatedSubscription.IsActive, updatedSubscription.UserID,
// 		).
// 		WillReturnResult(pgxmock.NewResult("UPDATE", 1))

// 	err = repo.Update(context.Background(), updatedSubscription)
// 	assert.NoError(t, err)

// 	// Test Delete
// 	mock.ExpectExec("DELETE FROM subscription WHERE user_id = \\$1").
// 		WithArgs(userID).
// 		WillReturnResult(pgxmock.NewResult("DELETE", 1))

// 	err = repo.Delete(context.Background(), userID)
// 	assert.NoError(t, err)

// 	assert.NoError(t, mock.ExpectationsWereMet())
// }
