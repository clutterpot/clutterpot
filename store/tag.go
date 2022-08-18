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
	query, args, err := sq.Insert("tags").SetMap(sq.Eq{
		"id":       model.NewID(),
		"owner_id": input.OwnerID,
		"name":     input.Name,
	}).Suffix("RETURNING id, name").PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return nil, fmt.Errorf("an error occurred while creating a tag")
	}

	var t model.Tag
	if err = ts.db.Get(&t, query, args...); err != nil {
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
	query, args, err := sq.Select("file_id", "id", "name").From("file_tags ft").
		LeftJoin("tags t ON ft.tag_id = t.id").Where(sq.Eq{"file_id": fileIDs}).
		PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return nil, []error{fmt.Errorf("an error occurred while getting file tags")}
	}

	var fts []*model.FileTag
	if err := ts.db.Select(&fts, query, args...); err != nil {
		return nil, []error{fmt.Errorf("an error occurred while getting file tags")}
	}

	t := make(map[string][]*model.Tag)
	for _, ft := range fts {
		t[ft.FileID] = append(t[ft.FileID], &ft.Tag)
	}

	tags = make([][]*model.Tag, len(fileIDs))
	for i, id := range fileIDs {
		tags[i] = t[id]
	}

	return tags, nil
}
