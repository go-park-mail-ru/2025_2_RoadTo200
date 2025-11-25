package repository

import (
	"context"
	"fmt"

	domain "github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/domain/entities"
	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/logger"
	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/repository/interfaces"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type MessageRepository struct {
	pool   *pgxpool.Pool
	logger logger.Log
}

func NewMessageRepository(pool *pgxpool.Pool, logger logger.Log) interfaces.MessageRepository {
	return &MessageRepository{
		pool:   pool,
		logger: logger,
	}
}

func (r *MessageRepository) Create(ctx context.Context, message *domain.Message) error {
	query := `
		INSERT INTO message (match_id, sender_id, receiver_id, content, is_read)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at
	`

	err := r.pool.QueryRow(ctx, query,
		message.MatchID,
		message.SenderID,
		message.ReceiverID,
		message.Content,
		message.IsRead,
	).Scan(&message.ID, &message.CreatedAt)

	if err != nil {
		return fmt.Errorf("failed to create message: %w", err)
	}

	return nil
}

func (r *MessageRepository) GetByMatchID(ctx context.Context, matchID uuid.UUID, limit, offset int) ([]domain.Message, error) {
	query := `
		SELECT id, match_id, sender_id, receiver_id, content, is_read, created_at
		FROM message
		WHERE match_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.pool.Query(ctx, query, matchID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get messages: %w", err)
	}
	defer rows.Close()

	var messages []domain.Message
	for rows.Next() {
		var msg domain.Message
		if err := rows.Scan(
			&msg.ID,
			&msg.MatchID,
			&msg.SenderID,
			&msg.ReceiverID,
			&msg.Content,
			&msg.IsRead,
			&msg.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan message: %w", err)
		}
		messages = append(messages, msg)
	}

	// Reverse to return in chronological order if needed, but usually UI handles it.
	// Let's keep DESC as requested by typical chat interfaces (load latest first).

	return messages, nil
}

func (r *MessageRepository) MarkAsRead(ctx context.Context, matchID, receiverID uuid.UUID) error {
	query := `
		UPDATE message
		SET is_read = TRUE
		WHERE match_id = $1 AND receiver_id = $2 AND is_read = FALSE
	`

	_, err := r.pool.Exec(ctx, query, matchID, receiverID)
	if err != nil {
		return fmt.Errorf("failed to mark messages as read: %w", err)
	}

	return nil
}

func (r *MessageRepository) GetUnreadCount(ctx context.Context, userID uuid.UUID) (int, error) {
	query := `
		SELECT COUNT(*)
		FROM message
		WHERE receiver_id = $1 AND is_read = FALSE
	`

	var count int
	err := r.pool.QueryRow(ctx, query, userID).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to get unread count: %w", err)
	}

	return count, nil
}

func (r *MessageRepository) GetConversations(ctx context.Context, userID uuid.UUID, searchQuery string) ([]domain.Conversation, error) {
	query := `
		WITH LastMessages AS (
			SELECT DISTINCT ON (match_id)
				match_id,
				content,
				created_at
			FROM message
			ORDER BY match_id, created_at DESC
		),
		UnreadCounts AS (
			SELECT
				match_id,
				COUNT(*) as unread_count
			FROM message
			WHERE receiver_id = $1 AND is_read = FALSE
			GROUP BY match_id
		)
		SELECT
			m.id as match_id,
			CASE
				WHEN m.user1_id = $1 THEN m.user2_id
				ELSE m.user1_id
			END as other_user_id,
			u.name as other_user_name,
			COALESCE(up.photo_url, '') as other_user_photo,
			COALESCE(lm.content, '') as last_message,
			COALESCE(lm.created_at, '0001-01-01'::timestamp) as last_message_time,
			COALESCE(uc.unread_count, 0) as unread_count
		FROM match m
		JOIN "user" u ON u.id = (CASE WHEN m.user1_id = $1 THEN m.user2_id ELSE m.user1_id END)
		LEFT JOIN LATERAL (
			SELECT photo_url 
			FROM user_photo 
			WHERE user_id = u.id AND is_approved = TRUE
			ORDER BY display_order ASC 
			LIMIT 1
		) up ON true
		LEFT JOIN LastMessages lm ON lm.match_id = m.id
		LEFT JOIN UnreadCounts uc ON uc.match_id = m.id
		WHERE m.user1_id = $1 OR m.user2_id = $1
	`

	args := []interface{}{userID}
	if searchQuery != "" {
		query += " AND u.name ILIKE $2"
		args = append(args, "%"+searchQuery+"%")
	}

	query += `
		ORDER BY last_message_time DESC
	`

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get conversations: %w", err)
	}
	defer rows.Close()

	var conversations []domain.Conversation
	for rows.Next() {
		var conv domain.Conversation
		if err := rows.Scan(
			&conv.MatchID,
			&conv.OtherUserID,
			&conv.OtherUserName,
			&conv.OtherUserPhoto,
			&conv.LastMessage,
			&conv.LastMessageTime,
			&conv.UnreadCount,
		); err != nil {
			return nil, fmt.Errorf("failed to scan conversation: %w", err)
		}
		conversations = append(conversations, conv)
	}

	return conversations, nil
}
