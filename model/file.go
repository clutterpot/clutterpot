package model

import "time"

type File struct {
	ID        string     `json:"id"`
	Filename  string     `json:"filename"`
	MimeType  string     `json:"mimeType"`
	Extension string     `json:"extension"`
	Size      int64      `json:"size"`
	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt time.Time  `json:"updatedAt"`
	DeletedAt *time.Time `json:"deletedAt"`
}

type FileInput struct {
	Filename string `json:"filename" validate:"required,filename,printunicode,min=1,max=255"`
}

type FileUpdateInput struct {
	Filename *string `json:"filename" validate:"omitempty,filename,printunicode,min=1,max=255"`
}
