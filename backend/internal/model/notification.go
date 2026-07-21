package model

import (
	"time"
)

type Notification struct {
	ID         string    `json:"id" db:"id"`
	UserID     string    `json:"user_id" db:"user_id"`
	Type       string    `json:"type" db:"type"`
	Title      string    `json:"title" db:"title"`
	Message    string    `json:"message" db:"message"`
	IsRead     bool      `json:"is_read" db:"is_read"`
	ActionLink *string   `json:"action_link" db:"action_link"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
}
