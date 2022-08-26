package store

import (
	"database/sql"
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

func (ts *tagStore) GetByID(id, ownerID string) (*model.Tag, error) {
	query, args, err := sq.Select("*").From("tags").
		Where(sq.And{
			sq.Eq{"id": id},
			sq.Or{
				sq.Eq{"owner_id": ownerID},
				sq.Eq{"owner_id": nil},
			}}).PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return nil, fmt.Errorf("an error occurred while getting a tag by id")
	}

	var t model.Tag
	if err := ts.db.Get(&t, query, args...); err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("tag not found")
		}

		return nil, fmt.Errorf("an error occurred while getting a tag by id")
	}

	return &t, nil
}

func (ts *tagStore) GetAll(ownerID string, after, before *string, first, last *int, filter *model.TagFilter, sort *model.TagSort, order *model.Order) (*model.TagConnection, error) {
	return getAll[model.Tag](ts.db, sq.Select("*").From("tags").Where(filter.GetConj()).
		Where(sq.Or{sq.Eq{"owner_id": ownerID}, sq.Eq{"owner_id": nil}}), "tags", after, before, first, last, string(*sort), *order)
}

func (ts *tagStore) GetAllByFileIDs(fileIDs []string, ownerID string, after, before *string, first, last *int, filter *model.TagFilter, sort *model.TagSort, order *model.Order) ([]*model.TagConnection, []error) {
	return getAllByKeys[model.Tag, model.FileTag](
		ts.db,
		sq.Select("file_id, t.*").From("file_tags ft").Where(sq.Eq{"ft.file_id": fileIDs}),
		sq.Select("*").From("tags").Where(sq.Or{sq.Eq{"owner_id": ownerID}, sq.Eq{"owner_id": nil}}).Where(filter.GetConj()),
		"JOIN LATERAL (?) t ON ft.tag_id = t.id", "file tags", fileIDs, after, before, first, last, string(*sort), *order)
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
