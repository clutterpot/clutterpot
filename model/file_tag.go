package model

type FileTag struct {
	FileID string `db:"file_id"`
	Tag
}

func (ft FileTag) GetKey() string { return ft.FileID }
func (ft FileTag) GetNode() *Tag  { return &ft.Tag }
