// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package model

import (
	"fmt"
	"io"
	"strconv"
	"time"
)

type DeleteFilePayload struct {
	// Unique file ID
	ID string `json:"id"`
	// Time of deletion
	DeletedAt time.Time `json:"deletedAt"`
}

type File struct {
	// Unique file ID
	ID string `json:"id"`
	// Filename
	Filename string `json:"filename"`
	// File [mime type](https://www.iana.org/assignments/media-types/media-types.xhtml)
	MimeType string `json:"mimeType"`
	// File extension
	Extension string `json:"extension"`
	// File size in bytes
	Size int64 `json:"size"`
	// File tags
	Tags []*Tag `json:"tags"`
	// Time of creation
	CreatedAt time.Time `json:"createdAt"`
	// Time of last update
	UpdatedAt time.Time `json:"updatedAt"`
	// Time of deletion
	DeletedAt *time.Time `json:"deletedAt"`
}

type LoginPayload struct {
	// Access token
	AccessToken string `json:"accessToken"`
	// Time of expiration
	ExpiresAt time.Time `json:"expiresAt"`
	// Refresh token
	RefreshToken string `json:"refreshToken"`
}

type PageInfo struct {
	HasNextPage     bool    `json:"hasNextPage"`
	HasPreviousPage bool    `json:"hasPreviousPage"`
	StartCursor     *string `json:"startCursor"`
	EndCursor       *string `json:"endCursor"`
}

type RefreshAccessTokenPayload struct {
	// Access token
	AccessToken string `json:"accessToken"`
	// Time of expiration
	ExpiresAt time.Time `json:"expiresAt"`
}

type RevokeRefreshTokenPayload struct {
	// Refresh token
	RefreshToken string `json:"refreshToken"`
	// Time of deletion
	DeletedAt time.Time `json:"deletedAt"`
}

type Tag struct {
	// Unique tag ID
	ID string `json:"id"`
	// Unique tag name
	Name string `json:"name"`
}

type User struct {
	// Unique user ID
	ID string `json:"id"`
	// Unique username
	Username string `json:"username"`
	// User email
	Email string `json:"email"`
	// User kind
	Kind UserKind `json:"kind"`
	// User display name
	DisplayName *string `json:"displayName"`
	// User bio
	Bio *string `json:"bio"`
	// Time of creation
	CreatedAt time.Time `json:"createdAt"`
	// Time of last update
	UpdatedAt time.Time `json:"updatedAt"`
}

type UserConnection struct {
	// User edges
	Edges []*UserEdge `json:"edges"`
	// User nodes
	Nodes []*User `json:"nodes"`
	// Page info
	PageInfo *PageInfo `json:"pageInfo"`
}

type UserEdge struct {
	// Pagination cursor
	Cursor string `json:"cursor"`
	// User node
	Node *User `json:"node"`
}

type Order string

const (
	OrderAsc  Order = "ASC"
	OrderDesc Order = "DESC"
)

var AllOrder = []Order{
	OrderAsc,
	OrderDesc,
}

func (e Order) IsValid() bool {
	switch e {
	case OrderAsc, OrderDesc:
		return true
	}
	return false
}

func (e Order) String() string {
	return string(e)
}

func (e *Order) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = Order(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid Order", str)
	}
	return nil
}

func (e Order) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}
