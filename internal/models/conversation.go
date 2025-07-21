package models

import (
	"me-ai/pkg/db"
	"time"
)

type Conversation struct {
	ID        int       `json:"id" db:"id"`
	UserID    int       `json:"user_id" db:"user_id"`
	Title     string    `json:"title" db:"title"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
	Messages  []Message `json:"messages,omitempty"`
}

type CreateConversationRequest struct {
	Title string `json:"title" binding:"required,min=1,max=255"`
}

type ConversationResponse struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type DeleteConversationRequest struct {
	ID int `json:"id"`
}

func (c *Conversation) ToResponse() ConversationResponse {
	return ConversationResponse{
		ID:        c.ID,
		Title:     c.Title,
		CreatedAt: c.CreatedAt,
		UpdatedAt: c.UpdatedAt,
	}
}

type ConversationRepository struct{}

func (r *ConversationRepository) Create(convo *Conversation) (*Conversation, error) {
	query := `INSERT INTO conversations (user_id, title) VALUES ($1, $2) RETURNING id, created_at, updated_at`
	err := db.DB.QueryRowx(query, convo.UserID, convo.Title).Scan(&convo.ID, &convo.CreatedAt, &convo.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return convo, nil
}

func (r *ConversationRepository) ListByUser(userID int) ([]Conversation, error) {
	var convos []Conversation
	err := db.DB.Select(&convos, "SELECT * FROM conversations WHERE user_id=$1 ORDER BY updated_at DESC", userID)
	if err != nil {
		return nil, err
	}
	return convos, nil
}

func (r *ConversationRepository) Delete(id, userID int) error {
	_, err := db.DB.Exec("DELETE FROM conversations WHERE id=$1 AND user_id=$2", id, userID)
	return err
}

func (r *ConversationRepository) UpdateTitle(id, userID int, title string) error {
	_, err := db.DB.Exec("UPDATE conversations SET title=$1, updated_at=NOW() WHERE id=$2 AND user_id=$3", title, id, userID)
	return err
}
