// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package model

import (
	"time"
)

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
	// Time the access token expires
	ExpiresIn time.Time `json:"expiresIn"`
	// Refresh token
	RefreshToken string `json:"refreshToken"`
}

type RefreshAccessTokenPayload struct {
	// Access token
	AccessToken string `json:"accessToken"`
	// Time the access token expires
	ExpiresIn time.Time `json:"expiresIn"`
}

type RevokeRefreshTokenPayload struct {
	// Refresh token
	RefreshToken string `json:"refreshToken"`
	// Time of deletion
	DeletedAt time.Time `json:"deletedAt"`
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
