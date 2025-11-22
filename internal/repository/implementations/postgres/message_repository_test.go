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

func TestMessageRepository_Create(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	repo := NewMessageRepository(mock)

	message := &domain.Message{
		MatchID:     uuid.New(),
		SenderID:    uuid.New(),
		MessageText: "Hello, world!",
		Status:      0, // unread
	}

	rows := mock.NewRows([]string{"id", "created_at"}).
		AddRow(uuid.New(), time.Now())

	mock.ExpectQuery("INSERT INTO message").
		WithArgs(message.MatchID, message.SenderID, message.MessageText, message.Status).
		WillReturnRows(rows)

	err = repo.Create(context.Background(), message)
	assert.NoError(t, err)
	assert.NotEqual(t, uuid.Nil, message.ID)
	assert.False(t, message.CreatedAt.IsZero())
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestMessageRepository_Create_Error(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	repo := NewMessageRepository(mock)

	message := &domain.Message{
		MatchID:     uuid.New(),
		SenderID:    uuid.New(),
		MessageText: "Hello, world!",
		Status:      0,
	}

	mock.ExpectQuery("INSERT INTO message").
		WithArgs(message.MatchID, message.SenderID, message.MessageText, message.Status).
		WillReturnError(pgx.ErrTxClosed)

	err = repo.Create(context.Background(), message)
	assert.Error(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestMessageRepository_GetByID(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	repo := NewMessageRepository(mock)

	messageID := uuid.New()
	matchID := uuid.New()
	senderID := uuid.New()

	expectedMessage := &domain.Message{
		ID:          messageID,
		MatchID:     matchID,
		SenderID:    senderID,
		MessageText: "Test message",
		Status:      1, // read
		CreatedAt:   time.Now(),
	}

	rows := mock.NewRows([]string{
		"id", "match_id", "sender_id", "message_text", "status", "created_at",
	}).AddRow(
		expectedMessage.ID, expectedMessage.MatchID, expectedMessage.SenderID,
		expectedMessage.MessageText, expectedMessage.Status, expectedMessage.CreatedAt,
	)

	mock.ExpectQuery("SELECT \\* FROM message WHERE id = \\$1").
		WithArgs(messageID).
		WillReturnRows(rows)

	message, err := repo.GetByID(context.Background(), messageID)
	assert.NoError(t, err)
	assert.Equal(t, expectedMessage.ID, message.ID)
	assert.Equal(t, expectedMessage.MatchID, message.MatchID)
	assert.Equal(t, expectedMessage.SenderID, message.SenderID)
	assert.Equal(t, expectedMessage.MessageText, message.MessageText)
	assert.Equal(t, expectedMessage.Status, message.Status)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestMessageRepository_GetByID_NotFound(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	repo := NewMessageRepository(mock)

	messageID := uuid.New()

	mock.ExpectQuery("SELECT \\* FROM message WHERE id = \\$1").
		WithArgs(messageID).
		WillReturnError(pgx.ErrNoRows)

	message, err := repo.GetByID(context.Background(), messageID)
	assert.NoError(t, err)
	assert.Nil(t, message)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestMessageRepository_GetByID_Error(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	repo := NewMessageRepository(mock)

	messageID := uuid.New()

	mock.ExpectQuery("SELECT \\* FROM message WHERE id = \\$1").
		WithArgs(messageID).
		WillReturnError(pgx.ErrTxClosed)

	message, err := repo.GetByID(context.Background(), messageID)
	assert.Error(t, err)
	assert.Nil(t, message)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestMessageRepository_GetByMatchID(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	repo := NewMessageRepository(mock)

	matchID := uuid.New()
	limit := 20
	offset := 0

	expectedMessages := []domain.Message{
		{
			ID:          uuid.New(),
			MatchID:     matchID,
			SenderID:    uuid.New(),
			MessageText: "First message",
			Status:      0,
			CreatedAt:   time.Now().Add(-time.Hour),
		},
		{
			ID:          uuid.New(),
			MatchID:     matchID,
			SenderID:    uuid.New(),
			MessageText: "Second message",
			Status:      1,
			CreatedAt:   time.Now(),
		},
	}

	rows := mock.NewRows([]string{
		"id", "match_id", "sender_id", "message_text", "status", "created_at",
	})
	for _, msg := range expectedMessages {
		rows.AddRow(
			msg.ID, msg.MatchID, msg.SenderID, msg.MessageText, msg.Status, msg.CreatedAt,
		)
	}

	mock.ExpectQuery("SELECT \\* FROM message WHERE match_id = \\$1 ORDER BY created_at ASC LIMIT \\$2 OFFSET \\$3").
		WithArgs(matchID, limit, offset).
		WillReturnRows(rows)

	messages, err := repo.GetByMatchID(context.Background(), matchID, limit, offset)
	assert.NoError(t, err)
	assert.Len(t, messages, 2)
	assert.Equal(t, expectedMessages[0].ID, messages[0].ID)
	assert.Equal(t, expectedMessages[1].ID, messages[1].ID)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestMessageRepository_GetByMatchID_Empty(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	repo := NewMessageRepository(mock)

	matchID := uuid.New()
	limit := 20
	offset := 0

	rows := mock.NewRows([]string{
		"id", "match_id", "sender_id", "message_text", "status", "created_at",
	})

	mock.ExpectQuery("SELECT \\* FROM message WHERE match_id = \\$1 ORDER BY created_at ASC LIMIT \\$2 OFFSET \\$3").
		WithArgs(matchID, limit, offset).
		WillReturnRows(rows)

	messages, err := repo.GetByMatchID(context.Background(), matchID, limit, offset)
	assert.NoError(t, err)
	assert.Empty(t, messages)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestMessageRepository_GetByMatchID_QueryError(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	repo := NewMessageRepository(mock)

	matchID := uuid.New()
	limit := 20
	offset := 0

	mock.ExpectQuery("SELECT \\* FROM message WHERE match_id = \\$1 ORDER BY created_at ASC LIMIT \\$2 OFFSET \\$3").
		WithArgs(matchID, limit, offset).
		WillReturnError(pgx.ErrTxClosed)

	messages, err := repo.GetByMatchID(context.Background(), matchID, limit, offset)
	assert.Error(t, err)
	assert.Nil(t, messages)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestMessageRepository_UpdateStatus(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	repo := NewMessageRepository(mock)

	messageID := uuid.New()
	newStatus := 1 // read

	mock.ExpectExec("UPDATE message SET status = \\$1 WHERE id = \\$2").
		WithArgs(newStatus, messageID).
		WillReturnResult(pgxmock.NewResult("UPDATE", 1))

	err = repo.UpdateStatus(context.Background(), messageID, newStatus)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestMessageRepository_UpdateStatus_Error(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	repo := NewMessageRepository(mock)

	messageID := uuid.New()
	newStatus := 1

	mock.ExpectExec("UPDATE message SET status = \\$1 WHERE id = \\$2").
		WithArgs(newStatus, messageID).
		WillReturnError(pgx.ErrTxClosed)

	err = repo.UpdateStatus(context.Background(), messageID, newStatus)
	assert.Error(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestMessageRepository_Delete(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	repo := NewMessageRepository(mock)

	messageID := uuid.New()

	mock.ExpectExec("DELETE FROM message WHERE id = \\$1").
		WithArgs(messageID).
		WillReturnResult(pgxmock.NewResult("DELETE", 1))

	err = repo.Delete(context.Background(), messageID)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestMessageRepository_Delete_Error(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	repo := NewMessageRepository(mock)

	messageID := uuid.New()

	mock.ExpectExec("DELETE FROM message WHERE id = \\$1").
		WithArgs(messageID).
		WillReturnError(pgx.ErrTxClosed)

	err = repo.Delete(context.Background(), messageID)
	assert.Error(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestMessageRepository_MarkMessagesAsRead(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	repo := NewMessageRepository(mock)

	matchID := uuid.New()
	userID := uuid.New()

	mock.ExpectExec("UPDATE message SET status = 1 WHERE match_id = \\$1 AND sender_id != \\$2 AND status = 0").
		WithArgs(matchID, userID).
		WillReturnResult(pgxmock.NewResult("UPDATE", 3))

	err = repo.MarkMessagesAsRead(context.Background(), matchID, userID)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestMessageRepository_MarkMessagesAsRead_Error(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	repo := NewMessageRepository(mock)

	matchID := uuid.New()
	userID := uuid.New()

	mock.ExpectExec("UPDATE message SET status = 1 WHERE match_id = \\$1 AND sender_id != \\$2 AND status = 0").
		WithArgs(matchID, userID).
		WillReturnError(pgx.ErrTxClosed)

	err = repo.MarkMessagesAsRead(context.Background(), matchID, userID)
	assert.Error(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestMessageRepository_GetUnreadCount(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	repo := NewMessageRepository(mock)

	userID := uuid.New()
	expectedCount := 5

	rows := mock.NewRows([]string{"count"}).AddRow(expectedCount)

	mock.ExpectQuery("SELECT COUNT\\(\\*\\) FROM message m JOIN match mt ON m.match_id = mt.id WHERE \\(mt.user1_id = \\$1 OR mt.user2_id = \\$1\\) AND m.sender_id != \\$1 AND m.status = 0").
		WithArgs(userID).
		WillReturnRows(rows)

	count, err := repo.GetUnreadCount(context.Background(), userID)
	assert.NoError(t, err)
	assert.Equal(t, expectedCount, count)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestMessageRepository_GetUnreadCount_Zero(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	repo := NewMessageRepository(mock)

	userID := uuid.New()

	rows := mock.NewRows([]string{"count"}).AddRow(0)

	mock.ExpectQuery("SELECT COUNT\\(\\*\\) FROM message m JOIN match mt ON m.match_id = mt.id WHERE \\(mt.user1_id = \\$1 OR mt.user2_id = \\$1\\) AND m.sender_id != \\$1 AND m.status = 0").
		WithArgs(userID).
		WillReturnRows(rows)

	count, err := repo.GetUnreadCount(context.Background(), userID)
	assert.NoError(t, err)
	assert.Equal(t, 0, count)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestMessageRepository_GetUnreadCount_Error(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	repo := NewMessageRepository(mock)

	userID := uuid.New()

	mock.ExpectQuery("SELECT COUNT\\(\\*\\) FROM message m JOIN match mt ON m.match_id = mt.id WHERE \\(mt.user1_id = \\$1 OR mt.user2_id = \\$1\\) AND m.sender_id != \\$1 AND m.status = 0").
		WithArgs(userID).
		WillReturnError(pgx.ErrTxClosed)

	count, err := repo.GetUnreadCount(context.Background(), userID)
	assert.Error(t, err)
	assert.Equal(t, 0, count)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestMessageRepository_GetConversations(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	repo := NewMessageRepository(mock)

	userID := uuid.New()

	expectedConversations := []domain.Conversation{
		{
			MatchID:         uuid.New(),
			OtherUserID:     uuid.New(),
			OtherUserName:   "John Doe",
			LastMessage:     "Hello there!",
			LastMessageTime: time.Now(),
			UnreadCount:     2,
		},
		{
			MatchID:         uuid.New(),
			OtherUserID:     uuid.New(),
			OtherUserName:   "Jane Smith",
			LastMessage:     "How are you?",
			LastMessageTime: time.Now().Add(-time.Hour),
			UnreadCount:     0,
		},
	}

	rows := mock.NewRows([]string{
		"match_id", "other_user_id", "other_user_name", "last_message", "last_message_time", "unread_count",
	})
	for _, conv := range expectedConversations {
		rows.AddRow(
			conv.MatchID, conv.OtherUserID, conv.OtherUserName,
			conv.LastMessage, conv.LastMessageTime, conv.UnreadCount,
		)
	}

	mock.ExpectQuery("SELECT m.id as match_id, CASE WHEN m.user1_id = \\$1 THEN m.user2_id ELSE m.user1_id END as other_user_id, u.name as other_user_name, msg.message_text as last_message, msg.created_at as last_message_time, \\(SELECT COUNT\\(\\*\\) FROM message m2 WHERE m2.match_id = m.id AND m2.sender_id != \\$1 AND m2.status = 0\\) as unread_count FROM match m JOIN \"user\" u ON u.id = CASE WHEN m.user1_id = \\$1 THEN m.user2_id ELSE m.user1_id END LEFT JOIN LATERAL \\( SELECT message_text, created_at FROM message WHERE match_id = m.id ORDER BY created_at DESC LIMIT 1 \\) msg ON true WHERE \\(m.user1_id = \\$1 OR m.user2_id = \\$1\\) AND m.is_active = true ORDER BY msg.created_at DESC NULLS LAST").
		WithArgs(userID).
		WillReturnRows(rows)

	conversations, err := repo.GetConversations(context.Background(), userID)
	assert.NoError(t, err)
	assert.Len(t, conversations, 2)
	assert.Equal(t, expectedConversations[0].MatchID, conversations[0].MatchID)
	assert.Equal(t, expectedConversations[0].OtherUserName, conversations[0].OtherUserName)
	assert.Equal(t, expectedConversations[0].UnreadCount, conversations[0].UnreadCount)
	assert.Equal(t, expectedConversations[1].MatchID, conversations[1].MatchID)
	assert.Equal(t, expectedConversations[1].OtherUserName, conversations[1].OtherUserName)
	assert.Equal(t, expectedConversations[1].UnreadCount, conversations[1].UnreadCount)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestMessageRepository_GetConversations_Empty(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	repo := NewMessageRepository(mock)

	userID := uuid.New()

	rows := mock.NewRows([]string{
		"match_id", "other_user_id", "other_user_name", "last_message", "last_message_time", "unread_count",
	})

	mock.ExpectQuery("SELECT m.id as match_id, CASE WHEN m.user1_id = \\$1 THEN m.user2_id ELSE m.user1_id END as other_user_id, u.name as other_user_name, msg.message_text as last_message, msg.created_at as last_message_time, \\(SELECT COUNT\\(\\*\\) FROM message m2 WHERE m2.match_id = m.id AND m2.sender_id != \\$1 AND m2.status = 0\\) as unread_count FROM match m JOIN \"user\" u ON u.id = CASE WHEN m.user1_id = \\$1 THEN m.user2_id ELSE m.user1_id END LEFT JOIN LATERAL \\( SELECT message_text, created_at FROM message WHERE match_id = m.id ORDER BY created_at DESC LIMIT 1 \\) msg ON true WHERE \\(m.user1_id = \\$1 OR m.user2_id = \\$1\\) AND m.is_active = true ORDER BY msg.created_at DESC NULLS LAST").
		WithArgs(userID).
		WillReturnRows(rows)

	conversations, err := repo.GetConversations(context.Background(), userID)
	assert.NoError(t, err)
	assert.Empty(t, conversations)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestMessageRepository_GetConversations_QueryError(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	repo := NewMessageRepository(mock)

	userID := uuid.New()

	mock.ExpectQuery("SELECT m.id as match_id, CASE WHEN m.user1_id = \\$1 THEN m.user2_id ELSE m.user1_id END as other_user_id, u.name as other_user_name, msg.message_text as last_message, msg.created_at as last_message_time, \\(SELECT COUNT\\(\\*\\) FROM message m2 WHERE m2.match_id = m.id AND m2.sender_id != \\$1 AND m2.status = 0\\) as unread_count FROM match m JOIN \"user\" u ON u.id = CASE WHEN m.user1_id = \\$1 THEN m.user2_id ELSE m.user1_id END LEFT JOIN LATERAL \\( SELECT message_text, created_at FROM message WHERE match_id = m.id ORDER BY created_at DESC LIMIT 1 \\) msg ON true WHERE \\(m.user1_id = \\$1 OR m.user2_id = \\$1\\) AND m.is_active = true ORDER BY msg.created_at DESC NULLS LAST").
		WithArgs(userID).
		WillReturnError(pgx.ErrTxClosed)

	conversations, err := repo.GetConversations(context.Background(), userID)
	assert.Error(t, err)
	assert.Nil(t, conversations)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestMessageRepository_GetConversations_ScanError(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	repo := NewMessageRepository(mock)

	userID := uuid.New()

	rows := mock.NewRows([]string{
		"match_id", "other_user_id", "other_user_name", "last_message", "last_message_time", "unread_count",
	}).AddRow(
		nil, nil, nil, nil, nil, nil, // Invalid data to cause scan error
	)

	mock.ExpectQuery("SELECT m.id as match_id, CASE WHEN m.user1_id = \\$1 THEN m.user2_id ELSE m.user1_id END as other_user_id, u.name as other_user_name, msg.message_text as last_message, msg.created_at as last_message_time, \\(SELECT COUNT\\(\\*\\) FROM message m2 WHERE m2.match_id = m.id AND m2.sender_id != \\$1 AND m2.status = 0\\) as unread_count FROM match m JOIN \"user\" u ON u.id = CASE WHEN m.user1_id = \\$1 THEN m.user2_id ELSE m.user1_id END LEFT JOIN LATERAL \\( SELECT message_text, created_at FROM message WHERE match_id = m.id ORDER BY created_at DESC LIMIT 1 \\) msg ON true WHERE \\(m.user1_id = \\$1 OR m.user2_id = \\$1\\) AND m.is_active = true ORDER BY msg.created_at DESC NULLS LAST").
		WithArgs(userID).
		WillReturnRows(rows)

	conversations, err := repo.GetConversations(context.Background(), userID)
	assert.Error(t, err)
	assert.Nil(t, conversations)
	assert.NoError(t, mock.ExpectationsWereMet())
}
