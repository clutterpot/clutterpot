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
	Username string `validate:"required,username,min=4,max=32"`
	Email    string `validate:"required,email,max=254"`
	Password string `validate:"required,password,min=8,max=255"`
}

type UserUpdateInput struct {
	Username    *string `validate:"omitempty,username,min=4,max=32"`
	Email       *string `validate:"omitempty,email,max=254"`
	Password    *string `validate:"omitempty,password,min=8,max=255"`
	Kind        *UserKind
	DisplayName *string `validate:"omitempty,displayname,max=32"`
	Bio         *string `validate:"omitempty,printunicode,max=160"`
}

func HashPassword(password string) string {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Panic(err)
	}

	return string(hash)
}
