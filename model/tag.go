package model

import (
	"time"

	sq "github.com/Masterminds/squirrel"
)

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

func (t *TagFilter) GetConj() sq.And {
	var conj sq.And

	if t == nil {
		return conj
	}

	conj = t.ID.GetConj(conj, "id")
	conj = t.Name.GetConj(conj, "name")

	if t.And != nil {
		conj = append(conj, t.And.GetConj()...)
	}
	if t.Or != nil {
		conj = sq.And{sq.Or{conj, t.Or.GetConj()}}
	}

	return conj
}

type TagSort string

const (
	TagSortCreated TagSort = "created_at"
	TagSortUpdated TagSort = "updated_at"
)

type TagInput struct {
	OwnerID *string `json:"-"`
	Name    string  `json:"name"`
	Global  bool    `json:"global"`
}
