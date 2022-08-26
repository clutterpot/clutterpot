package store

import (
	"database/sql"
	"fmt"

	"github.com/clutterpot/clutterpot/model"

	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
)

type fileStore struct {
	db *sqlx.DB
}

func newFileStore(db *sqlx.DB) *fileStore {
	return &fileStore{db: db}
}

func (fs *fileStore) Create(input model.FileInput) (*model.File, error) {
	/*
		WITH f AS (
			INSERT INTO files(id, filename, mime_type, extension, size) VALUES ( <values> )
			RETURNING *
		)[, ft AS (
			INSERT INTO file_tags(file_id, tag_id) VALUES ( (SELECT id FROM f), unnest(array [ <values> ]) )
		)] SELECT * FROM f
	*/
	query := sq.Select("*").From("f").
		Prefix("WITH f AS (?)", sq.Insert("files").SetMap(sq.Eq{
			"id":        model.NewID(),
			"filename":  input.Filename,
			"mime_type": "test",
			"extension": ".test",
			"size":      0,
		}).Suffix("RETURNING *"))

	if input.Tags != nil && len(*input.Tags) > 0 {
		// Convert []string to []any
		tagsIDs := make([]any, len(*input.Tags))
		for i, tagID := range *input.Tags {
			tagsIDs[i] = tagID
		}

		// Append file tags insert query
		query = query.Prefix(", ft AS (?)", sq.Insert("file_tags").SetMap(sq.Eq{
			"file_id": sq.Expr("(?)", sq.Select("id").From("f")),
			"tag_id":  sq.Expr(fmt.Sprintf("unnest(array[%s])", sq.Placeholders(len(*input.Tags))), tagsIDs...),
		}))
	}

	createQuery, args, err := query.PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return nil, fmt.Errorf("an error occurred while creating a file")
	}

	var f model.File
	if err = fs.db.Get(&f, createQuery, args...); err != nil {
		return nil, fmt.Errorf("an error occurred while creating a file")
	}

	return &f, nil
}

func (fs *fileStore) Update(id string, input model.FileUpdateInput) (*model.File, error) {
	/*
		WITH f AS (
			UPDATE files SET ( <values> ) WHERE id = $1 AND deleted_at IS NULL RETURNING *
		)[, ft AS (
			INSERT INTO file_tags(file_id, tag_id) VALUES ( (SELECT id FROM f), unnest(array [ <values> ]) )
		)] SELECT * FROM f
	*/
	query := sq.Select("*").From("f")
	fileUpdateQuery := sq.Update("files").Set("updated_at", "now()")

	if input == (model.FileUpdateInput{}) {
		return nil, fmt.Errorf("file update input cannot be empty")
	}
	if input.Filename != nil {
		fileUpdateQuery = fileUpdateQuery.Set("filename", *input.Filename)
	}

	// Append file update query
	query = query.Prefix("WITH f AS (?)",
		fileUpdateQuery.Where(sq.And{sq.Eq{"id": id}, sq.Eq{"deleted_at": nil}}).Suffix("RETURNING *"),
	)

	if input.Tags != nil && len(*input.Tags) > 0 {
		// Convert []string to []any
		tagsIDs := make([]any, len(*input.Tags))
		for i, tagID := range *input.Tags {
			tagsIDs[i] = tagID
		}

		// Append file tags insert query
		query = query.Prefix(", ft AS (?)", sq.Insert("file_tags").SetMap(sq.Eq{
			"file_id": sq.Expr("(?)", sq.Select("id").From("f")),
			"tag_id":  sq.Expr(fmt.Sprintf("unnest(array[%s])", sq.Placeholders(len(*input.Tags))), tagsIDs...),
		}))
	}

	updateQuery, args, err := query.PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return nil, fmt.Errorf("an error occurred while updating a file")
	}

	var f model.File
	if err = fs.db.Get(&f, updateQuery, args...); err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("file not found")
		}

		return nil, fmt.Errorf("an error occurred while updating a file")
	}

	return &f, nil
}

func (fs *fileStore) SoftDelete(id string) (*model.DeleteFilePayload, error) {
	var d model.DeleteFilePayload

	if err := sq.Update("files").Set("deleted_at", "now()").
		Where(sq.And{sq.Eq{"id": id}, sq.Eq{"deleted_at": nil}}).
		Suffix("RETURNING deleted_at").
		PlaceholderFormat(sq.Dollar).RunWith(fs.db).QueryRow().
		Scan(&d.DeletedAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("file not found")
		}

		return nil, fmt.Errorf("an error occurred while deleting a file by id")
	}

	return &d, nil
}

func (fs *fileStore) Delete(id string) error {
	if _, err := sq.Delete("files").Where("id = ?", id).
		PlaceholderFormat(sq.Dollar).RunWith(fs.db).Query(); err != nil {
		return fmt.Errorf("an error occurred while deleting a file")
	}

	return nil
}

func (fs *fileStore) RemoveTags(id string, tagIDs []string) (*model.RemoveTagsFromFilePayload, error) {
	if _, err := sq.Delete("file_tags").Where(sq.Eq{
		"file_id": id,
		"tag_id":  tagIDs,
	}).PlaceholderFormat(sq.Dollar).RunWith(fs.db).Query(); err != nil {
		return nil, fmt.Errorf("an error occurred while removing file tags")
	}

	return &model.RemoveTagsFromFilePayload{FileID: id, TagIDs: tagIDs}, nil
}

func (fs *fileStore) GetByID(id string) (*model.File, error) {
	query, args, err := sq.Select("*").From("files").
		Where(sq.And{sq.Eq{"id": id}, sq.Eq{"deleted_at": nil}}).PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return nil, fmt.Errorf("an error occurred while getting a file by id")
	}

	var f model.File
	if err := fs.db.Get(&f, query, args...); err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("file not found")
		}

		return nil, fmt.Errorf("an error occurred while getting a file by id")
	}

	return &f, nil
}

func (fs *fileStore) GetAll(after, before *string, first, last *int, filter *model.FileFilter, sort *model.FileSort, order *model.Order) (*model.FileConnection, error) {
	return getAll[model.File](fs.db, sq.Select("*").From("files").Where(filter.GetConj()), "files", after, before, first, last, string(*sort), *order)
}

func (fs *fileStore) GetAllByIDs(ids []string) ([]*model.File, []error) {
	query, args, err := sq.Select("*").From("files").
		Where(sq.And{sq.Eq{"id": ids}, sq.Eq{"deleted_at": nil}}).PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return nil, []error{fmt.Errorf("an error occurred while getting all files by ids")}
	}

	var f []*model.File
	if err := fs.db.Select(&f, query, args...); err != nil || len(f) == 0 {
		if len(f) == 0 {
			return nil, []error{fmt.Errorf("files not found")}
		}

		return nil, []error{fmt.Errorf("an error occurred while getting all files by ids")}
	}

	return f, nil
}
