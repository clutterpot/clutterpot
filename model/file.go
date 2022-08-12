package model

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
