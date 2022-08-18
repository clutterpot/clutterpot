package store

import (
	"database/sql"
	"fmt"

	"github.com/clutterpot/clutterpot/model"

	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
)

type sessionStore struct {
	db *sqlx.DB
}

func newSessionStore(db *sqlx.DB) *sessionStore {
	return &sessionStore{db: db}
}

func (ss *sessionStore) Create(input model.SessionInput) (*model.Session, error) {
	query, args, err := sq.Insert("sessions").SetMap(sq.Eq{
		"id":      model.NewID(),
		"user_id": input.UserID,
	}).Suffix("RETURNING *").PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return nil, fmt.Errorf("an error occurred while creating a session")
	}

	var s model.Session
	if err = ss.db.Get(&s, query, args...); err != nil {
		return nil, err
	}

	return &s, nil
}

func (ss *sessionStore) SoftDelete(id string) (*model.RevokeRefreshTokenPayload, error) {
	var r model.RevokeRefreshTokenPayload

	if err := sq.Update("sessions").Set("deleted_at", "now()").
		Where(sq.And{sq.Eq{"id": id}, sq.Eq{"deleted_at": nil}}).
		Suffix("RETURNING deleted_at").
		PlaceholderFormat(sq.Dollar).RunWith(ss.db).QueryRow().
		Scan(&r.DeletedAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("session not found")
		}

		return nil, fmt.Errorf("an error occurred while deleting a session by id")
	}

	return &r, nil
}

func (ss *sessionStore) Delete(id string) error {
	if _, err := sq.Delete("sessions").Where("id = ?", id).
		PlaceholderFormat(sq.Dollar).RunWith(ss.db).Query(); err != nil {
		return fmt.Errorf("an error occurred while deleting a session")
	}

	return nil
}

func (ss *sessionStore) GetByID(id string) (*model.SessionUser, error) {
	query, args, err := sq.Select("s.*", "u.*").
		From("sessions as s").Where(sq.And{
		sq.Eq{"s.id": id},
		sq.Gt{"expires_at": "now()"},
		sq.Eq{"deleted_at": nil},
	}).LeftJoin("users as u ON u.id = s.user_id").PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return nil, fmt.Errorf("an error occurred while getting a session by id")
	}

	var s model.SessionUser
	if err = ss.db.Get(&s, query, args...); err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("session not found")
		}

		return nil, fmt.Errorf("an error occurred while getting a session by id")
	}

	return &s, nil
}
