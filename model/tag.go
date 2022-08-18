package model

type Tag struct {
	ID      string `json:"id" db:"id" `
	OwnerID string `json:"-" db:"owner_id"`
	Name    string `json:"name" db:"name"`
}

type TagInput struct {
	OwnerID *string `json:"-"`
	Name    string  `json:"name"`
	Private bool    `json:"private"`
}
