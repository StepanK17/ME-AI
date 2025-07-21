package models

import (
	"me-ai/pkg/db"
	"time"
)

type Message struct {
	ID             string    `json:"id"`
	ConversationID int       `json:"conversation_id" db:"conversation_id"`
	UserID         int       `json:"user_id" db:"user_id"`
	Content        string    `json:"content"`
	Role           string    `json:"role"`
	Timestamp      time.Time `json:"timestamp"`
	CreatedAt      time.Time `json:"created_at" db:"created_at"`
}

type ChatRequest struct {
	Message string `json:"message"`
}

type ChatResponse struct {
	Message   string `json:"message"`
	Timestamp string `json:"timestamp"`
}

type WebSocketMessage struct {
	Type    string `json:"type"`
	Content string `json:"content"`
	Role    string `json:"role"`
}

type DeleteMessageRequest struct {
	ID int `json:"id"`
}

type MessageRepository struct{}

func (r *MessageRepository) Create(msg *Message) (*Message, error) {
	query := `INSERT INTO messages (conversation_id, user_id, role, content) VALUES ($1, $2, $3, $4) RETURNING id, created_at`
	err := db.DB.QueryRowx(query, msg.ConversationID, msg.UserID, msg.Role, msg.Content).Scan(&msg.ID, &msg.Timestamp)
	if err != nil {
		return nil, err
	}
	return msg, nil
}

func (r *MessageRepository) ListByConversation(convoID int) ([]Message, error) {
	var msgs []Message
	err := db.DB.Select(&msgs, "SELECT * FROM messages WHERE conversation_id=$1 ORDER BY created_at ASC", convoID)
	if err != nil {
		return nil, err
	}
	return msgs, nil
}

func (r *MessageRepository) Delete(id, userID int) error {
	_, err := db.DB.Exec("DELETE FROM messages WHERE id=$1 AND user_id=$2", id, userID)
	return err
}
