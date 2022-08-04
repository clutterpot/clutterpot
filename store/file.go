package store

import (
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
	var f model.File

	if err := sq.Insert("files").SetMap(sq.Eq{
		"id":        model.NewID(),
		"filename":  input.Filename,
		"mime_type": "test",
		"extension": ".test",
		"size":      0,
	}).Suffix("RETURNING id, filename, mime_type, extension, size, created_at, updated_at, deleted_at").
		PlaceholderFormat(sq.Dollar).RunWith(fs.db).QueryRow().
		Scan(&f.ID, &f.Filename, &f.MimeType, &f.Extension, &f.Size, &f.CreatedAt, &f.UpdatedAt, &f.DeletedAt); err != nil {
		return nil, err
	}

	return &f, nil
}

func (fs *fileStore) Update(id string, input model.FileUpdateInput) (*model.File, error) {
	var f model.File

	query := sq.Update("files").Set("updated_at", "now()")

	if input == (model.FileUpdateInput{}) {
		return nil, fmt.Errorf("file update input cannot be empty")
	}
	if input.Filename != nil {
		query = query.Set("filename", *input.Filename)
	}

	if err := query.Where("id = ?", id).
		Suffix("RETURNING id, filename, mime_type, extension, size, created_at, updated_at, deleted_at").
		PlaceholderFormat(sq.Dollar).RunWith(fs.db).QueryRow().
		Scan(&f.ID, &f.Filename, &f.MimeType, &f.Extension, &f.Size, &f.CreatedAt, &f.UpdatedAt, &f.DeletedAt); err != nil {
		return nil, err
	}

	return &f, nil
}

func (fs *fileStore) Delete(id string) error {
	if _, err := sq.Delete("files").Where("id = ?", id).
		PlaceholderFormat(sq.Dollar).RunWith(fs.db).Query(); err != nil {
		return err
	}

	return nil
}

func (fs *fileStore) GetByID(id string) (*model.File, error) {
	var f model.File

	if err := sq.Select("id", "filename", "mime_type", "extension", "size", "created_at", "updated_at", "deleted_at").
		From("files").Where("id = ?", id).
		PlaceholderFormat(sq.Dollar).RunWith(fs.db).QueryRow().
		Scan(&f.ID, &f.Filename, &f.MimeType, &f.Extension, &f.Size, &f.CreatedAt, &f.UpdatedAt, &f.DeletedAt); err != nil {
		return nil, err
	}

	return &f, nil
}
