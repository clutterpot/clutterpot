package model

import "time"

type File struct {
	ID        string     `json:"id" db:"id"`
	Filename  string     `json:"filename" db:"filename"`
	MimeType  string     `json:"mimeType" db:"mime_type"`
	Extension string     `json:"extension" db:"extension"`
	Size      int64      `json:"size" db:"size"`
	Tags      []*Tag     `json:"tags" db:"-"`
	CreatedAt time.Time  `json:"createdAt" db:"created_at"`
	UpdatedAt time.Time  `json:"updatedAt" db:"updated_at"`
	DeletedAt *time.Time `json:"deletedAt" db:"deleted_at"`
}

type FileInput struct {
	Filename string    `json:"filename" validate:"required,filename,printunicode,min=1,max=255"`
	Tags     *[]string `json:"tags" validate:"omitempty"`
}

type FileUpdateInput struct {
	Filename *string   `json:"filename" validate:"omitempty,filename,printunicode,min=1,max=255"`
	Tags     *[]string `json:"tags" validate:"omitempty"`
}

type RemoveTagsFromFilePayload struct {
	FileID string
	TagIDs []string
	File   *File  `json:"file"`
	Tags   []*Tag `json:"tags"`
}
