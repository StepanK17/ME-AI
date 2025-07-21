package models

import (
	"me-ai/pkg/db"
)

type User struct {
	ID        int    `json:"id" db:"id"`
	Email     string `json:"email" db:"email"`
	Password  string `json:"password,omitempty" db:"password"`
	Name      string `json:"name" db:"name"`
	CreatedAt string `json:"created_at" db:"created_at"`
}

type UserRepository struct{}

func (r *UserRepository) FindByEmail(email string) (*User, error) {
	var user User
	err := db.DB.Get(&user, "SELECT * FROM users WHERE email=$1", email)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) Create(user *User) (*User, error) {
	query := `INSERT INTO users (email, password, name) VALUES ($1, $2, $3) RETURNING id, created_at`
	err := db.DB.QueryRowx(query, user.Email, user.Password, user.Name).Scan(&user.ID, &user.CreatedAt)
	if err != nil {
		return nil, err
	}
	return user, nil
}
