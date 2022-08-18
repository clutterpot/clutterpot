package model

type Tag struct {
	ID      string `json:"id" db:"id" `
	OwnerID string `json:"-" db:"owner_id"`
	Name    string `json:"name" db:"name"`
}

func (Tag) IsNode()         {}
func (t Tag) GetID() string { return t.ID }

type TagInput struct {
	OwnerID *string `json:"-"`
	Name    string  `json:"name"`
	Private bool    `json:"private"`
}
