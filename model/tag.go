package model

import "time"

type Tag struct {
	ID        string     `json:"id" db:"id" `
	OwnerID   *string    `json:"-" db:"owner_id"`
	Name      string     `json:"name" db:"name"`
	CreatedAt time.Time  `json:"-" db:"created_at"`
	UpdatedAt time.Time  `json:"-" db:"updated_at"`
	DeletedAt *time.Time `json:"-" db:"deleted_at"`
}

func (Tag) IsNode()         {}
func (t Tag) GetID() string { return t.ID }

type TagConnection = Connection[Tag]
type TagEdge = Edge[Tag]

type TagSort string

const (
	TagSortCreated TagSort = "created_at"
	TagSortUpdated TagSort = "updated_at"
)

type TagInput struct {
	OwnerID *string `json:"-"`
	Name    string  `json:"name"`
	Private bool    `json:"private"`
}
