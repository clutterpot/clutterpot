package model

type FileInput struct {
	Filename string `json:"filename" validate:"required,filename,printunicode,min=1,max=255"`
}

type FileUpdateInput struct {
	Filename *string `json:"filename" validate:"omitempty,filename,printunicode,min=1,max=255"`
}
