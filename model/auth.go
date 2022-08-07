package model

import "time"

type AuthUser struct {
	User
	Password string
}

type AuthPayload struct {
	AccessToken  string
	TokenType    string
	ExpiresIn    time.Time
	RefreshToken *string
	Scope        *string
}
