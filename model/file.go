package model

import (
	"time"

	sq "github.com/Masterminds/squirrel"
)

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

func (File) IsNode()         {}
func (f File) GetID() string { return f.ID }

type FileConnection = Connection[File]
type FileEdge = Edge[File]

func (f *FileFilter) GetConj() sq.And {
	var conj sq.And

	conj = f.ID.GetConj(conj, "id")
	conj = f.Filename.GetConj(conj, "filename")
	conj = f.MimeType.GetConj(conj, "mime_type")
	conj = f.Extension.GetConj(conj, "extension")
	conj = f.Size.GetConj(conj, "size")
	conj = f.CreatedAt.GetConj(conj, "created_at")
	conj = f.UpdatedAt.GetConj(conj, "updated_at")
	conj = f.DeletedAt.GetConj(conj, "deleted_at")

	if f.And != nil {
		conj = append(conj, f.And.GetConj()...)
	}
	if f.Or != nil {
		conj = sq.And{sq.Or{conj, f.Or.GetConj()}}
	}

	return conj
}

type FileSort string

const (
	FileSortCreated FileSort = "created_at"
	FileSortUpdated FileSort = "updated_at"
)

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
