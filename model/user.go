package model

import (
	"log"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID          string
	Username    string
	Email       string
	Kind        UserKind
	DisplayName *string
	Bio         *string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type UserKind int8

const (
	UserKindBanned UserKind = iota
	UserKindUser
	UserKindAdmin
)

type UserInput struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserUpdateInput struct {
	Username    *string
	Email       *string
	Password    *string
	Kind        *UserKind
	DisplayName *string
	Bio         *string
}

func HashPassword(password string) string {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Panic(err)
	}

	return string(hash)
}
