package model

type FileTag struct {
	FileID string `db:"file_id"`
	Tag
}
