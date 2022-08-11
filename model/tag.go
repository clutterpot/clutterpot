package model

type TagInput struct {
	OwnerID *string `json:"-"`
	Name    string  `json:"name"`
	Private bool    `json:"private"`
}
