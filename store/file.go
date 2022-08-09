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
		return nil, fmt.Errorf("an error occurred while creating a file")
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

	if err := query.Where(sq.And{sq.Eq{"id": id}, sq.Eq{"deleted_at": nil}}).
		Suffix("RETURNING id, filename, mime_type, extension, size, created_at, updated_at, deleted_at").
		PlaceholderFormat(sq.Dollar).RunWith(fs.db).QueryRow().
		Scan(&f.ID, &f.Filename, &f.MimeType, &f.Extension, &f.Size, &f.CreatedAt, &f.UpdatedAt, &f.DeletedAt); err != nil {
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

func (fs *fileStore) GetByID(id string) (*model.File, error) {
	var f model.File

	if err := sq.Select("id", "filename", "mime_type", "extension", "size", "created_at", "updated_at", "deleted_at").
		From("files").Where(sq.And{sq.Eq{"id": id}, sq.Eq{"deleted_at": nil}}).
		PlaceholderFormat(sq.Dollar).RunWith(fs.db).QueryRow().
		Scan(&f.ID, &f.Filename, &f.MimeType, &f.Extension, &f.Size, &f.CreatedAt, &f.UpdatedAt, &f.DeletedAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("file not found")
		}

		return nil, fmt.Errorf("an error occurred while getting a file by id")
	}

	return &f, nil
}
