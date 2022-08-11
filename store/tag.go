package store

import (
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
