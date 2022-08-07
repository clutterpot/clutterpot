package model

import "time"

type Session struct {
	ID        string
	UserID    string
	CreatedAt time.Time
	ExpiresAt time.Time
}

type SessionUser struct {
	User
	Session
}

type SessionInput struct {
	UserID string
}
