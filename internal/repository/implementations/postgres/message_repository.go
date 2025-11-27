package postgres

import (
	"context"
	"strings"

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
	query := `
		INSERT INTO message (match_id, sender_id, receiver_id, content, is_read)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at`

	err := r.pool.QueryRow(ctx, query,
		message.MatchID, message.SenderID, message.ReceiverID, message.Content, message.IsRead).
		Scan(&message.ID, &message.CreatedAt)

	if err != nil {
		return err
	}
	return nil
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
		AND m.is_read = FALSE`

	var count int
	err := r.pool.QueryRow(ctx, query, userID).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

// GetConversations returns all conversations (matches with last message) for user
func (r *MessageRepository) GetConversations(ctx context.Context, userID uuid.UUID, searchQuery string) ([]domain.Conversation, error) {
	// Trim пробелы из поискового запроса
	searchQuery = strings.TrimSpace(searchQuery)

	query := `
		SELECT 
			m.id as match_id,
			CASE 
				WHEN m.user1_id = $1 THEN m.user2_id 
				ELSE m.user1_id 
			END as other_user_id,
			u.name as other_user_name,
			COALESCE(p.photo_url, '') as other_user_photo,
			COALESCE(msg.content, '') as last_message,
			COALESCE(msg.created_at, '1970-01-01'::timestamp) as last_message_time,
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
			WHERE user_id = u.id AND is_approved = TRUE
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
		  AND m.is_active = true`

	args := []interface{}{userID}

	// Добавляем условие поиска если запрос не пустой
	if searchQuery != "" {
		query += " AND LOWER(u.name) LIKE LOWER($2)"
		args = append(args, "%"+searchQuery+"%")
	}

	query += `
		ORDER BY CASE 
			WHEN msg.created_at IS NULL THEN 1 
			ELSE 0 
		END, msg.created_at DESC`

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var conversations []domain.Conversation
	for rows.Next() {
		var conv domain.Conversation

		err := rows.Scan(
			&conv.MatchID,
			&conv.OtherUserID,
			&conv.OtherUserName,
			&conv.OtherUserPhoto,
			&conv.LastMessage,
			&conv.LastMessageTime,
			&conv.UnreadCount,
		)
		if err != nil {
			return nil, err
		}

		conversations = append(conversations, conv)
	}

	return conversations, nil
}
