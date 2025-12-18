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
)

func TestNotificationRepository_Create(t *testing.T) {
	userID := uuid.New()
	fromUserID := uuid.New()
	matchID := uuid.New()
	notification := &domain.Notification{
		UserID:     userID,
		Type:       constants.NotificationTypeLike,
		FromUserID: &fromUserID,
		MatchID:    &matchID,
		IsRead:     false,
	}

	expectedID := uuid.New()
	expectedCreatedAt := time.Now()

	tests := []struct {
		name          string
		notification  *domain.Notification
		setupMock     func(mock pgxmock.PgxPoolIface)
		expectedError bool
	}{
		{
			name:         "successful creation with all fields",
			notification: notification,
			setupMock: func(mock pgxmock.PgxPoolIface) {
				rows := mock.NewRows([]string{"id", "created_at"}).
					AddRow(expectedID, expectedCreatedAt)
				mock.ExpectQuery(`INSERT INTO notification`).
					WithArgs(notification.UserID, notification.Type, notification.FromUserID, notification.MatchID, notification.IsRead).
					WillReturnRows(rows)
			},
			expectedError: false,
		},
		{
			name: "successful creation with nil from_user_id",
			notification: &domain.Notification{
				UserID:     userID,
				Type:       constants.NotificationTypeMatch,
				FromUserID: nil,
				MatchID:    &matchID,
				IsRead:     false,
			},
			setupMock: func(mock pgxmock.PgxPoolIface) {
				rows := mock.NewRows([]string{"id", "created_at"}).
					AddRow(expectedID, expectedCreatedAt)
				mock.ExpectQuery(`INSERT INTO notification`).
					WithArgs(userID, constants.NotificationTypeMatch, (*uuid.UUID)(nil), &matchID, false).
					WillReturnRows(rows)
			},
			expectedError: false,
		},
		{
			name:         "database error",
			notification: notification,
			setupMock: func(mock pgxmock.PgxPoolIface) {
				mock.ExpectQuery(`INSERT INTO notification`).
					WithArgs(notification.UserID, notification.Type, notification.FromUserID, notification.MatchID, notification.IsRead).
					WillReturnError(pgx.ErrNoRows)
			},
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock, err := pgxmock.NewPool()
			assert.NoError(t, err)
			defer mock.Close()

			repo := NewNotificationRepository(mock)
			tt.setupMock(mock)

			err = repo.Create(context.Background(), tt.notification)

			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, expectedID, tt.notification.ID)
				assert.Equal(t, expectedCreatedAt, tt.notification.CreatedAt)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestNotificationRepository_GetByUserID(t *testing.T) {
	userID := uuid.New()
	fromUserID := uuid.New()
	matchID := uuid.New()
	limit := 10
	offset := 5

	expectedNotifications := []domain.Notification{
		{
			ID:         uuid.New(),
			UserID:     userID,
			Type:       constants.NotificationTypeLike,
			FromUserID: &fromUserID,
			MatchID:    &matchID,
			IsRead:     false,
			CreatedAt:  time.Now(),
		},
		{
			ID:         uuid.New(),
			UserID:     userID,
			Type:       constants.NotificationTypeMatch,
			FromUserID: nil,
			MatchID:    nil,
			IsRead:     true,
			CreatedAt:  time.Now().Add(-time.Hour),
		},
	}

	tests := []struct {
		name                  string
		userID                uuid.UUID
		limit                 int
		offset                int
		setupMock             func(mock pgxmock.PgxPoolIface)
		expectedNotifications []domain.Notification
		expectedError         bool
	}{
		{
			name:   "successful get notifications",
			userID: userID,
			limit:  limit,
			offset: offset,
			setupMock: func(mock pgxmock.PgxPoolIface) {
				rows := mock.NewRows([]string{"id", "user_id", "type", "from_user_id", "match_id", "is_read", "created_at"}).
					AddRow(expectedNotifications[0].ID, expectedNotifications[0].UserID, expectedNotifications[0].Type,
						expectedNotifications[0].FromUserID, expectedNotifications[0].MatchID, expectedNotifications[0].IsRead, expectedNotifications[0].CreatedAt).
					AddRow(expectedNotifications[1].ID, expectedNotifications[1].UserID, expectedNotifications[1].Type,
						expectedNotifications[1].FromUserID, expectedNotifications[1].MatchID, expectedNotifications[1].IsRead, expectedNotifications[1].CreatedAt)
				mock.ExpectQuery(`SELECT n\.id, n\.user_id, n\.type, n\.from_user_id, n\.match_id, n\.is_read, n\.created_at.*FROM notification n.*LEFT JOIN match m.*WHERE n\.user_id.*`).
					WithArgs(userID, limit, offset).
					WillReturnRows(rows)
			},
			expectedNotifications: expectedNotifications,
			expectedError:         false,
		},
		{
			name:   "empty result",
			userID: userID,
			limit:  limit,
			offset: offset,
			setupMock: func(mock pgxmock.PgxPoolIface) {
				rows := mock.NewRows([]string{"id", "user_id", "type", "from_user_id", "match_id", "is_read", "created_at"})
				mock.ExpectQuery(`SELECT n\.id, n\.user_id, n\.type, n\.from_user_id, n\.match_id, n\.is_read, n\.created_at.*FROM notification n.*LEFT JOIN match m.*WHERE n\.user_id.*`).
					WithArgs(userID, limit, offset).
					WillReturnRows(rows)
			},
			expectedNotifications: []domain.Notification{},
			expectedError:         false,
		},
		{
			name:   "database error",
			userID: userID,
			limit:  limit,
			offset: offset,
			setupMock: func(mock pgxmock.PgxPoolIface) {
				mock.ExpectQuery(`SELECT n\.id, n\.user_id, n\.type, n\.from_user_id, n\.match_id, n\.is_read, n\.created_at.*FROM notification n.*LEFT JOIN match m.*WHERE n\.user_id.*`).
					WithArgs(userID, limit, offset).
					WillReturnError(assert.AnError)
			},
			expectedNotifications: nil,
			expectedError:         true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock, err := pgxmock.NewPool()
			assert.NoError(t, err)
			defer mock.Close()

			repo := NewNotificationRepository(mock)
			tt.setupMock(mock)

			notifications, err := repo.GetByUserID(context.Background(), tt.userID, tt.limit, tt.offset)

			if tt.expectedError {
				assert.Error(t, err)
				assert.Nil(t, notifications)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedNotifications, notifications)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestNotificationRepository_GetUnreadCount(t *testing.T) {
	userID := uuid.New()
	expectedCount := 5

	tests := []struct {
		name          string
		userID        uuid.UUID
		setupMock     func(mock pgxmock.PgxPoolIface)
		expectedCount int
		expectedError bool
	}{
		{
			name:   "successful get unread count",
			userID: userID,
			setupMock: func(mock pgxmock.PgxPoolIface) {
				rows := mock.NewRows([]string{"count"}).AddRow(expectedCount)
				mock.ExpectQuery(`SELECT COUNT\(\*\) FROM notification`).
					WithArgs(userID).
					WillReturnRows(rows)
			},
			expectedCount: expectedCount,
			expectedError: false,
		},
		{
			name:   "zero unread notifications",
			userID: userID,
			setupMock: func(mock pgxmock.PgxPoolIface) {
				rows := mock.NewRows([]string{"count"}).AddRow(0)
				mock.ExpectQuery(`SELECT COUNT\(\*\) FROM notification`).
					WithArgs(userID).
					WillReturnRows(rows)
			},
			expectedCount: 0,
			expectedError: false,
		},
		{
			name:   "database error",
			userID: userID,
			setupMock: func(mock pgxmock.PgxPoolIface) {
				mock.ExpectQuery(`SELECT COUNT\(\*\) FROM notification`).
					WithArgs(userID).
					WillReturnError(assert.AnError)
			},
			expectedCount: 0,
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock, err := pgxmock.NewPool()
			assert.NoError(t, err)
			defer mock.Close()

			repo := NewNotificationRepository(mock)
			tt.setupMock(mock)

			count, err := repo.GetUnreadCount(context.Background(), tt.userID)

			if tt.expectedError {
				assert.Error(t, err)
				assert.Equal(t, 0, count)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedCount, count)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestNotificationRepository_MarkAsRead(t *testing.T) {
	notificationID := uuid.New()

	tests := []struct {
		name           string
		notificationID uuid.UUID
		setupMock      func(mock pgxmock.PgxPoolIface)
		expectedError  bool
	}{
		{
			name:           "successful mark as read",
			notificationID: notificationID,
			setupMock: func(mock pgxmock.PgxPoolIface) {
				mock.ExpectExec(`UPDATE notification SET is_read = true WHERE id = \$1`).
					WithArgs(notificationID).
					WillReturnResult(pgxmock.NewResult("UPDATE", 1))
			},
			expectedError: false,
		},
		{
			name:           "notification not found (no rows affected)",
			notificationID: notificationID,
			setupMock: func(mock pgxmock.PgxPoolIface) {
				mock.ExpectExec(`UPDATE notification SET is_read = true WHERE id = \$1`).
					WithArgs(notificationID).
					WillReturnResult(pgxmock.NewResult("UPDATE", 0))
			},
			expectedError: false, // UPDATE doesn't return error if no rows affected
		},
		{
			name:           "database error",
			notificationID: notificationID,
			setupMock: func(mock pgxmock.PgxPoolIface) {
				mock.ExpectExec(`UPDATE notification SET is_read = true WHERE id = \$1`).
					WithArgs(notificationID).
					WillReturnError(assert.AnError)
			},
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock, err := pgxmock.NewPool()
			assert.NoError(t, err)
			defer mock.Close()

			repo := NewNotificationRepository(mock)
			tt.setupMock(mock)

			err = repo.MarkAsRead(context.Background(), tt.notificationID)

			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestNotificationRepository_MarkAllAsRead(t *testing.T) {
	userID := uuid.New()

	tests := []struct {
		name          string
		userID        uuid.UUID
		setupMock     func(mock pgxmock.PgxPoolIface)
		expectedError bool
	}{
		{
			name:   "successful mark all as read",
			userID: userID,
			setupMock: func(mock pgxmock.PgxPoolIface) {
				mock.ExpectExec(`UPDATE notification SET is_read = true WHERE user_id = \$1 AND is_read = false`).
					WithArgs(userID).
					WillReturnResult(pgxmock.NewResult("UPDATE", 3))
			},
			expectedError: false,
		},
		{
			name:   "no unread notifications",
			userID: userID,
			setupMock: func(mock pgxmock.PgxPoolIface) {
				mock.ExpectExec(`UPDATE notification SET is_read = true WHERE user_id = \$1 AND is_read = false`).
					WithArgs(userID).
					WillReturnResult(pgxmock.NewResult("UPDATE", 0))
			},
			expectedError: false,
		},
		{
			name:   "database error",
			userID: userID,
			setupMock: func(mock pgxmock.PgxPoolIface) {
				mock.ExpectExec(`UPDATE notification SET is_read = true WHERE user_id = \$1 AND is_read = false`).
					WithArgs(userID).
					WillReturnError(assert.AnError)
			},
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock, err := pgxmock.NewPool()
			assert.NoError(t, err)
			defer mock.Close()

			repo := NewNotificationRepository(mock)
			tt.setupMock(mock)

			err = repo.MarkAllAsRead(context.Background(), tt.userID)

			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
