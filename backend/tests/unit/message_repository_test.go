package unit

import (
	"testing"
	"time"

	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/domain/dto"
	domain "github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/domain/entities"
	repository "github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/repository/implementations/postgres"
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

	repo := repository.NewMessageRepository(mock)

	matchID := uuid.New()
	senderID := uuid.New()

	message := &domain.Message{
		MatchID:     matchID,
		SenderID:    senderID,
		MessageText: "Hello, world!",
		Status:      0,
	}

	rows := mock.NewRows([]string{"id", "created_at"}).
		AddRow(uuid.New(), time.Now())

	mock.ExpectQuery("INSERT INTO message").
		WithArgs(message.MatchID, message.SenderID, message.MessageText, message.Status).
		WillReturnRows(rows)

	err = repo.Create(message)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestMessageRepository_GetByID(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	repo := repository.NewMessageRepository(mock)

	messageID := uuid.New()
	matchID := uuid.New()
	senderID := uuid.New()

	expectedMessage := &domain.Message{
		ID:          messageID,
		MatchID:     matchID,
		SenderID:    senderID,
		MessageText: "Test message",
		Status:      1,
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

	message, err := repo.GetByID(messageID)
	assert.NoError(t, err)
	assert.Equal(t, expectedMessage.ID, message.ID)
	assert.Equal(t, expectedMessage.MatchID, message.MatchID)
	assert.Equal(t, expectedMessage.MessageText, message.MessageText)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestMessageRepository_GetByID_NotFound(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	repo := repository.NewMessageRepository(mock)

	messageID := uuid.New()

	mock.ExpectQuery("SELECT \\* FROM message WHERE id = \\$1").
		WithArgs(messageID).
		WillReturnError(pgx.ErrNoRows)

	message, err := repo.GetByID(messageID)
	assert.NoError(t, err)
	assert.Nil(t, message)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestMessageRepository_GetByMatchID(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	repo := repository.NewMessageRepository(mock)

	matchID := uuid.New()
	limit := 10
	offset := 0

	expectedMessages := []domain.Message{
		{
			ID:          uuid.New(),
			MatchID:     matchID,
			SenderID:    uuid.New(),
			MessageText: "First message",
			Status:      0,
			CreatedAt:   time.Now(),
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
	for _, message := range expectedMessages {
		rows.AddRow(
			message.ID, message.MatchID, message.SenderID,
			message.MessageText, message.Status, message.CreatedAt,
		)
	}

	mock.ExpectQuery("SELECT \\* FROM message WHERE match_id = \\$1 ORDER BY created_at ASC LIMIT \\$2 OFFSET \\$3").
		WithArgs(matchID, limit, offset).
		WillReturnRows(rows)

	messages, err := repo.GetByMatchID(matchID, limit, offset)
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

	repo := repository.NewMessageRepository(mock)

	matchID := uuid.New()
	limit := 10
	offset := 0

	rows := mock.NewRows([]string{
		"id", "match_id", "sender_id", "message_text", "status", "created_at",
	})

	mock.ExpectQuery("SELECT \\* FROM message WHERE match_id = \\$1 ORDER BY created_at ASC LIMIT \\$2 OFFSET \\$3").
		WithArgs(matchID, limit, offset).
		WillReturnRows(rows)

	messages, err := repo.GetByMatchID(matchID, limit, offset)
	assert.NoError(t, err)
	assert.Empty(t, messages)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestMessageRepository_UpdateStatus(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	repo := repository.NewMessageRepository(mock)

	messageID := uuid.New()
	status := 1

	mock.ExpectExec("UPDATE message SET status = \\$1 WHERE id = \\$2").
		WithArgs(status, messageID).
		WillReturnResult(pgxmock.NewResult("UPDATE", 1))

	err = repo.UpdateStatus(messageID, status)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestMessageRepository_Delete(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	repo := repository.NewMessageRepository(mock)

	messageID := uuid.New()

	mock.ExpectExec("DELETE FROM message WHERE id = \\$1").
		WithArgs(messageID).
		WillReturnResult(pgxmock.NewResult("DELETE", 1))

	err = repo.Delete(messageID)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestMessageRepository_MarkMessagesAsRead(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	repo := repository.NewMessageRepository(mock)

	matchID := uuid.New()
	userID := uuid.New()

	mock.ExpectExec("UPDATE message SET status = 1 WHERE match_id = \\$1 AND sender_id != \\$2 AND status = 0").
		WithArgs(matchID, userID).
		WillReturnResult(pgxmock.NewResult("UPDATE", 3))

	err = repo.MarkMessagesAsRead(matchID, userID)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestMessageRepository_GetUnreadCount(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	repo := repository.NewMessageRepository(mock)

	userID := uuid.New()
	expectedCount := 5

	rows := mock.NewRows([]string{"count"}).AddRow(expectedCount)

	mock.ExpectQuery("SELECT COUNT").
		WithArgs(userID).
		WillReturnRows(rows)

	count, err := repo.GetUnreadCount(userID)
	assert.NoError(t, err)
	assert.Equal(t, expectedCount, count)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestMessageRepository_GetUnreadCount_Zero(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	repo := repository.NewMessageRepository(mock)

	userID := uuid.New()
	expectedCount := 0

	rows := mock.NewRows([]string{"count"}).AddRow(expectedCount)

	mock.ExpectQuery("SELECT COUNT").
		WithArgs(userID).
		WillReturnRows(rows)

	count, err := repo.GetUnreadCount(userID)
	assert.NoError(t, err)
	assert.Equal(t, expectedCount, count)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestMessageRepository_GetConversations(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	repo := repository.NewMessageRepository(mock)

	userID := uuid.New()

	expectedConversations := []dto.Conversation{
		{
			MatchID:         uuid.New(),
			OtherUserID:     uuid.New(),
			OtherUserName:   "John Doe",
			LastMessage:     "Hello there!",
			LastMessageTime: time.Now(),
			UnreadCount:     3,
		},
		{
			MatchID:         uuid.New(),
			OtherUserID:     uuid.New(),
			OtherUserName:   "Jane Smith",
			LastMessage:     "How are you?",
			LastMessageTime: time.Now(),
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

	mock.ExpectQuery("SELECT").
		WithArgs(userID).
		WillReturnRows(rows)

	conversations, err := repo.GetConversations(userID)
	assert.NoError(t, err)
	assert.Len(t, conversations, 2)
	assert.Equal(t, expectedConversations[0].MatchID, conversations[0].MatchID)
	assert.Equal(t, expectedConversations[1].MatchID, conversations[1].MatchID)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestMessageRepository_GetConversations_Empty(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	repo := repository.NewMessageRepository(mock)

	userID := uuid.New()

	rows := mock.NewRows([]string{
		"match_id", "other_user_id", "other_user_name", "last_message", "last_message_time", "unread_count",
	})

	mock.ExpectQuery("SELECT").
		WithArgs(userID).
		WillReturnRows(rows)

	conversations, err := repo.GetConversations(userID)
	assert.NoError(t, err)
	assert.Empty(t, conversations)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestMessageRepository_Create_Error(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	repo := repository.NewMessageRepository(mock)

	matchID := uuid.New()
	senderID := uuid.New()

	message := &domain.Message{
		MatchID:     matchID,
		SenderID:    senderID,
		MessageText: "Hello, world!",
		Status:      0,
	}

	mock.ExpectQuery("INSERT INTO message").
		WithArgs(message.MatchID, message.SenderID, message.MessageText, message.Status).
		WillReturnError(pgx.ErrNoRows)

	err = repo.Create(message)
	assert.Error(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestMessageRepository_UpdateStatus_Error(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	repo := repository.NewMessageRepository(mock)

	messageID := uuid.New()
	status := 1

	mock.ExpectExec("UPDATE message SET status = \\$1 WHERE id = \\$2").
		WithArgs(status, messageID).
		WillReturnError(pgx.ErrNoRows)

	err = repo.UpdateStatus(messageID, status)
	assert.Error(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestMessageRepository_GetUnreadCount_Error(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	repo := repository.NewMessageRepository(mock)

	userID := uuid.New()

	mock.ExpectQuery("SELECT COUNT").
		WithArgs(userID).
		WillReturnError(pgx.ErrNoRows)

	count, err := repo.GetUnreadCount(userID)
	assert.Error(t, err)
	assert.Equal(t, 0, count)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestMessageRepository_GetConversations_Error(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	repo := repository.NewMessageRepository(mock)

	userID := uuid.New()

	mock.ExpectQuery("SELECT").
		WithArgs(userID).
		WillReturnError(pgx.ErrNoRows)

	conversations, err := repo.GetConversations(userID)
	assert.Error(t, err)
	assert.Nil(t, conversations)
	assert.NoError(t, mock.ExpectationsWereMet())
}
