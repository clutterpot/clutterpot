package model

import (
	"fmt"
	"log"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID          string    `json:"id"`
	Username    string    `json:"username"`
	Email       string    `json:"email"`
	Kind        UserKind  `json:"kind"`
	DisplayName *string   `json:"displayName"`
	Bio         *string   `json:"bio"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

type UserKind int8

const (
	UserKindBanned UserKind = iota
	UserKindUser
	UserKindAdmin
)

type UserInput struct {
	Username string `json:"username" validate:"required,username,min=4,max=32"`
	Email    string `json:"email" validate:"required,email,max=254"`
	Password string `json:"password" validate:"required,password,min=8,max=255"`
}

type UserUpdateInput struct {
	Username    *string   `json:"username" validate:"omitempty,username,min=4,max=32"`
	Email       *string   `json:"email" validate:"omitempty,email,max=254"`
	Password    *string   `json:"password" validate:"omitempty,password,min=8,max=255"`
	Kind        *UserKind `json:"kind"`
	DisplayName *string   `json:"displayName" validate:"omitempty,displayname,max=32"`
	Bio         *string   `json:"bio" validate:"omitempty,printunicode,max=160"`
}

func HashPassword(password string) string {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Panic(err)
	}

	return string(hash)
}

func CompareHashAndPassword(hash, password string) error {
	if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)); err != nil {
		return fmt.Errorf("invalid email or password")
	}

	return nil
}
