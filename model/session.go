package model

import "time"

type Session struct {
	ID        string     `db:"id"`
	UserID    string     `db:"user_id"`
	CreatedAt time.Time  `db:"created_at"`
	ExpiresAt time.Time  `db:"expires_at"`
	DeletedAt *time.Time `db:"deleted_at"`
}

type SessionUser struct {
	User
	Session
}

type SessionInput struct {
	UserID string
}
