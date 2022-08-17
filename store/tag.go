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

func (ts *tagStore) GetByFileIDs(fileIDs []string) (tags [][]*model.Tag, errs []error) {
	rows, err := sq.Select("file_id", "id", "name").From("file_tags ft").
		LeftJoin("tags t ON ft.tag_id = t.id").Where(sq.Eq{"file_id": fileIDs}).
		PlaceholderFormat(sq.Dollar).RunWith(ts.db).Query()
	if err != nil {
		return nil, []error{err}
	}
	defer rows.Close()

	t := make(map[string][]*model.Tag)
	for rows.Next() {
		var fileID string
		var tag model.Tag
		if err := rows.Scan(&fileID, &tag.ID, &tag.Name); err != nil {
			errs = append(errs, err)
		}
		t[fileID] = append(t[fileID], &tag)
	}

	tags = make([][]*model.Tag, len(fileIDs))

	for i, id := range fileIDs {
		tags[i] = t[id]
	}

	return tags, nil
}
