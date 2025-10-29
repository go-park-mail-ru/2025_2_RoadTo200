package postgres

import (
	"context"
	"errors"

	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/domain/dto"
	domain "github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/domain/entities"
	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/repository/interfaces"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type messageRepository struct {
	pool *pgxpool.Pool
}

func NewMessageRepository(pool *pgxpool.Pool) interfaces.MessageRepository {
	return &messageRepository{pool: pool}
}

func (r *messageRepository) Create(ctx context.Context, message *domain.Message) error {
	query := `
		INSERT INTO message (match_id, sender_id, message_text, status)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at`

	err := r.pool.QueryRow(ctx, query,
		message.MatchID, message.SenderID, message.MessageText, message.Status).
		Scan(&message.ID, &message.CreatedAt)

	if err != nil {
		return err
	}
	return nil
}

func (r *messageRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Message, error) {
	var message domain.Message
	query := `SELECT * FROM message WHERE id = $1`

	err := r.pool.QueryRow(ctx, query, id).Scan(
		&message.ID, &message.MatchID, &message.SenderID, &message.MessageText,
		&message.Status, &message.CreatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &message, nil
}

func (r *messageRepository) GetByMatchID(ctx context.Context, matchID uuid.UUID, limit, offset int) ([]domain.Message, error) {
	query := `
		SELECT * FROM message 
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
			&message.ID, &message.MatchID, &message.SenderID, &message.MessageText,
			&message.Status, &message.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		messages = append(messages, message)
	}

	return messages, nil
}

func (r *messageRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status int) error {
	query := `UPDATE message SET status = $1 WHERE id = $2`

	_, err := r.pool.Exec(ctx, query, status, id)
	if err != nil {
		return err
	}
	return nil
}

func (r *messageRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM message WHERE id = $1`

	_, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		return err
	}
	return nil
}

func (r *messageRepository) MarkMessagesAsRead(ctx context.Context, matchID, userID uuid.UUID) error {
	query := `
		UPDATE message 
		SET status = 1 
		WHERE match_id = $1 AND sender_id != $2 AND status = 0`

	_, err := r.pool.Exec(ctx, query, matchID, userID)
	if err != nil {
		return err
	}
	return nil
}

func (r *messageRepository) GetUnreadCount(ctx context.Context, userID uuid.UUID) (int, error) {
	query := `
		SELECT COUNT(*) 
		FROM message m
		JOIN match mt ON m.match_id = mt.id
		WHERE (mt.user1_id = $1 OR mt.user2_id = $1)
		AND m.sender_id != $1
		AND m.status = 0`

	var count int
	err := r.pool.QueryRow(ctx, query, userID).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (r *messageRepository) GetConversations(ctx context.Context, userID uuid.UUID) ([]dto.Conversation, error) {
	query := `
		SELECT 
			m.id as match_id,
			CASE 
				WHEN m.user1_id = $1 THEN m.user2_id 
				ELSE m.user1_id 
			END as other_user_id,
			u.name as other_user_name,
			msg.message_text as last_message,
			msg.created_at as last_message_time,
			(SELECT COUNT(*) FROM message m2 
			 WHERE m2.match_id = m.id AND m2.sender_id != $1 AND m2.status = 0) as unread_count
		FROM match m
		JOIN "user" u ON u.id = CASE 
			WHEN m.user1_id = $1 THEN m.user2_id 
			ELSE m.user1_id 
		END
		LEFT JOIN LATERAL (
			SELECT message_text, created_at 
			FROM message 
			WHERE match_id = m.id 
			ORDER BY created_at DESC 
			LIMIT 1
		) msg ON true
		WHERE (m.user1_id = $1 OR m.user2_id = $1) AND m.is_active = true
		ORDER BY msg.created_at DESC NULLS LAST`

	rows, err := r.pool.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var conversations []dto.Conversation
	for rows.Next() {
		var conv dto.Conversation
		err := rows.Scan(
			&conv.MatchID, &conv.OtherUserID, &conv.OtherUserName,
			&conv.LastMessage, &conv.LastMessageTime, &conv.UnreadCount,
		)
		if err != nil {
			return nil, err
		}
		conversations = append(conversations, conv)
	}

	return conversations, nil
}
