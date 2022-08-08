package model

type AuthUser struct {
	User
	Password string `json:"password"`
}
