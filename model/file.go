package model

import "time"

type File struct {
	ID        string
	Filename  string
	MimeType  string
	Extension string
	Size      int64
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}

type FileInput struct {
	Filename string `validate:"required,filename,printunicode,min=1,max=255"`
}

type FileUpdateInput struct {
	Filename *string `validate:"omitempty,filename,printunicode,min=1,max=255"`
}
