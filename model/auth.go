package model

import "time"

type AuthPayload struct {
	AccessToken  string
	TokenType    string
	ExpiresIn    time.Time
	RefreshToken *string
	Scope        *string
}
