package store

import (
	"fmt"

	"github.com/clutterpot/clutterpot/model"

	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
)

type tagStore struct {
	db *sqlx.DB
}

func newTagStore(db *sqlx.DB) *tagStore {
	return &tagStore{db: db}
}

func (ts *tagStore) Create(input model.TagInput) (*model.Tag, error) {
	var t model.Tag

	if err := sq.Insert("tags").SetMap(sq.Eq{
		"id":       model.NewID(),
		"owner_id": input.OwnerID,
		"name":     input.Name,
	}).Suffix("RETURNING id, name").
		PlaceholderFormat(sq.Dollar).RunWith(ts.db).QueryRow().
		Scan(&t.ID, &t.Name); err != nil {
		return nil, fmt.Errorf("an error occurred while creating a tag")
	}

	return &t, nil
}

func (ts *tagStore) Delete(id string) error {
	if _, err := sq.Delete("tags").Where("id = ?", id).
		PlaceholderFormat(sq.Dollar).RunWith(ts.db).Query(); err != nil {
		return fmt.Errorf("an error occurred while deleting a tag")
	}

	return nil
}

func (ts *tagStore) GetByFileID(fileID string) ([]*model.Tag, error) {
	var t []*model.Tag

	tagsIDs := sq.StatementBuilder.Select("tag_id").From("file_tags").
		Where("file_id = ?", fileID).PlaceholderFormat(sq.Dollar)

	rows, err := sq.Select("name").From("tags").Where(tagsIDs.Prefix("id IN (").Suffix(")")).
		PlaceholderFormat(sq.Dollar).RunWith(ts.db).Query()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var temp model.Tag
		if err := rows.Scan(&temp.Name); err != nil {
			return nil, err
		}
		t = append(t, &temp)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return t, nil
}
