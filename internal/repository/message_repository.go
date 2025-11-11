package repository

import (
	"context"

	"github.com/rishabhdeo60-dev/whatsapp-clone/internal/db"
	"github.com/rishabhdeo60-dev/whatsapp-clone/internal/model"
)

type MessageRepository interface {
	// Define methods for message repository here
	SendMessage(ctx context.Context, message *model.Message) (*model.Message, error)
	GetConversation(ctx context.Context, senderID, receiverID int64) ([]*model.Message, error)
	UpdateMessageStatus(ctx context.Context, messageID int64, status string) error
}

type messageRepository struct {
	// Add necessary fields like DB connection here
	db *db.DB
}

// Implement message repository methods here

// SendMessage implements MessageRepository.Save
func (r *messageRepository) SendMessage(ctx context.Context, message *model.Message) (*model.Message, error) {
	query := `INSERT INTO messages (sender_id, receiver_id, content, status) VALUES ($1, $2, $3, $4) RETURNING id, created_at, updated_at`
	return message, r.db.Pool.QueryRow(ctx, query, message.SenderID, message.ReceiverID, message.Content, message.Status).
		Scan(&message.ID, &message.CreatedAt, &message.UpdatedAt)
}

// GetConversation implements MessageRepository.GetConversation
func (r *messageRepository) GetConversation(ctx context.Context, senderID, receiverID int64) ([]*model.Message, error) {
	query := `SELECT id, sender_id, receiver_id, content, status, created_at, updated_at FROM messages WHERE (sender_id = $1 AND receiver_id = $2) OR (sender_id = $2 AND receiver_id = $1) ORDER BY created_at ASC`
	rows, err := r.db.Pool.Query(ctx, query, senderID, receiverID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []*model.Message
	for rows.Next() {
		var msg model.Message
		if err := rows.Scan(&msg.ID, &msg.SenderID, &msg.ReceiverID, &msg.Content,
			&msg.Status, &msg.CreatedAt, &msg.UpdatedAt); err != nil {
			return nil, err
		}
		messages = append(messages, &msg)
	}
	return messages, nil
}

// UpdateMessageStatus implements MessageRepository.UpdateStatus
func (r *messageRepository) UpdateMessageStatus(ctx context.Context, messageID int64, status string) error {
	query := `UPDATE messages SET status = $1, updated_at = NOW() WHERE id = $2`
	_, err := r.db.Pool.Exec(ctx, query, status, messageID)
	return err
}

// NewMessageRepository creates a new instance of MessageRepository
func NewMessageRepository(db *db.DB) MessageRepository {
	return &messageRepository{
		db: db,
	}
}
