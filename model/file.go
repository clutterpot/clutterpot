package model

import "time"

type File struct {
	ID        string
	Name      string
	MimeType  string
	Extension string
	Size      int64
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}

type FileInput struct {
	Name string
}
