package postgres

import (
	"context"
	"time"

	domain "github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/domain/entities"
	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/repository/interfaces"
	"github.com/google/uuid"
)

var _ interfaces.MessageRepository = (*MessageRepository)(nil)

type MessageRepository struct {
	pool interfaces.PgxIface
}

func NewMessageRepository(pool interfaces.PgxIface) *MessageRepository {
	return &MessageRepository{pool: pool}
}

// Create new message
func (r *MessageRepository) Create(ctx context.Context, message *domain.Message) error {
	// Начинаем транзакцию
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	// Вставляем сообщение
	query := `
		INSERT INTO message (match_id, sender_id, receiver_id, content, is_read)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at`

	err = tx.QueryRow(ctx, query,
		message.MatchID, message.SenderID, message.ReceiverID, message.Content, message.IsRead).
		Scan(&message.ID, &message.CreatedAt)

	if err != nil {
		return err
	}

	// Обновляем match: убираем expires_at (устанавливаем NULL), так как кто-то написал сообщение
	// Матч теперь активен навсегда
	updateMatchQuery := `
		UPDATE match 
		SET expires_at = NULL 
		WHERE id = $1 AND expires_at IS NOT NULL`

	_, err = tx.Exec(ctx, updateMatchQuery, message.MatchID)
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}

// GetByMatchID returns message history for a match with pagination
func (r *MessageRepository) GetByMatchID(ctx context.Context, matchID uuid.UUID, limit, offset int) ([]domain.Message, error) {
	query := `
		SELECT id, match_id, sender_id, receiver_id, content, is_read, created_at
		FROM message 
		WHERE match_id = $1 
		ORDER BY created_at ASC 
		LIMIT $2 OFFSET $3`

	rows, err := r.pool.Query(ctx, query, matchID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []domain.Message
	for rows.Next() {
		var message domain.Message
		err := rows.Scan(
			&message.ID, &message.MatchID, &message.SenderID, &message.ReceiverID,
			&message.Content, &message.IsRead, &message.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		messages = append(messages, message)
	}

	return messages, nil
}

// MarkAsRead marks all messages in a match as read by receiver
func (r *MessageRepository) MarkAsRead(ctx context.Context, matchID, receiverID uuid.UUID) error {
	query := `
		UPDATE message 
		SET is_read = TRUE 
		WHERE match_id = $1 AND receiver_id = $2 AND is_read = FALSE`

	_, err := r.pool.Exec(ctx, query, matchID, receiverID)
	if err != nil {
		return err
	}
	return nil
}

// GetUnreadCount returns total unread message count for user
func (r *MessageRepository) GetUnreadCount(ctx context.Context, userID uuid.UUID) (int, error) {
	query := `
		SELECT COUNT(*) 
		FROM message m
		JOIN match mt ON m.match_id = mt.id
		WHERE (mt.user1_id = $1 OR mt.user2_id = $1)
		AND m.receiver_id = $1
		AND m.is_read = FALSE
		AND mt.is_active = TRUE`

	var count int
	err := r.pool.QueryRow(ctx, query, userID).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

// GetConversations returns all conversations (matches with last message) for user
func (r *MessageRepository) GetConversations(ctx context.Context, userID uuid.UUID, searchQuery string) ([]domain.Conversation, error) {
	query := `
		SELECT 
			m.id as match_id,
			CASE 
				WHEN m.user1_id = $1 THEN m.user2_id 
				ELSE m.user1_id 
			END as other_user_id,
			u.name as other_user_name,
			COALESCE(p.photo_url, '') as other_user_photo,
			msg.content as last_message,
			msg.created_at as last_message_time,
			(SELECT COUNT(*) FROM message m2 
			 WHERE m2.match_id = m.id AND m2.receiver_id = $1 AND m2.is_read = FALSE) as unread_count
		FROM match m
		JOIN "user" u ON u.id = CASE 
			WHEN m.user1_id = $1 THEN m.user2_id 
			ELSE m.user1_id 
		END
		LEFT JOIN LATERAL (
			SELECT photo_url 
			FROM user_photo 
			WHERE user_id = u.id 
			ORDER BY display_order ASC 
			LIMIT 1
		) p ON true
		LEFT JOIN LATERAL (
			SELECT content, created_at 
			FROM message 
			WHERE match_id = m.id 
			ORDER BY created_at DESC 
			LIMIT 1
		) msg ON true
		WHERE (m.user1_id = $1 OR m.user2_id = $1)
		  AND m.is_active = TRUE
	`

	args := []interface{}{userID}
	if searchQuery != "" {
		query += " AND u.name ILIKE $2"
		args = append(args, "%"+searchQuery+"%")
	}

	query += `
		ORDER BY msg.created_at DESC NULLS LAST
	`

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var conversations []domain.Conversation
	for rows.Next() {
		var conv domain.Conversation
		var lastMessage *string
		var lastMessageTime *time.Time

		err := rows.Scan(
			&conv.MatchID, &conv.OtherUserID, &conv.OtherUserName, &conv.OtherUserPhoto,
			&lastMessage, &lastMessageTime, &conv.UnreadCount,
		)
		if err != nil {
			return nil, err
		}

		if lastMessage != nil {
			conv.LastMessage = *lastMessage
		}
		if lastMessageTime != nil {
			conv.LastMessageTime = *lastMessageTime
		}

		conversations = append(conversations, conv)
	}

	return conversations, nil
}
