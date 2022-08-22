package model

import (
	"fmt"
	"log"
	"time"

	sq "github.com/Masterminds/squirrel"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID          string    `json:"id" db:"id"`
	Username    string    `json:"username" db:"username"`
	Password    string    `json:"-" db:"password"`
	Email       string    `json:"email" db:"email"`
	Kind        UserKind  `json:"kind" db:"kind"`
	DisplayName *string   `json:"displayName" db:"display_name"`
	Bio         *string   `json:"bio" db:"bio"`
	CreatedAt   time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt   time.Time `json:"updatedAt" db:"updated_at"`
}

func (User) IsNode()         {}
func (u User) GetID() string { return u.ID }

type UserConnection = Connection[User]
type UserEdge = Edge[User]
type UserKindFilter = ScalarFilter[UserKind]

func (u *UserFilter) GetConj() sq.And {
	var and sq.And

	and = u.ID.GetConj(and, "id")
	and = u.Username.GetConj(and, "username")
	and = u.Email.GetConj(and, "email")
	and = u.Kind.GetConj(and, "kind")
	and = u.DisplayName.GetConj(and, "display_name")
	and = u.Bio.GetConj(and, "bio")
	and = u.CreatedAt.GetConj(and, "created_at")
	and = u.UpdatedAt.GetConj(and, "updated_at")

	return and
}

type UserKind int8

const (
	UserKindBanned UserKind = iota
	UserKindUser
	UserKindAdmin
)

type UserSort string

const (
	UserSortCreated UserSort = "created_at"
	UserSortUpdated UserSort = "updated_at"
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
