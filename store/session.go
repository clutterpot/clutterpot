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
	var s model.Session

	if err := sq.Insert("sessions").SetMap(sq.Eq{
		"id":      model.NewID(),
		"user_id": input.UserID,
	}).Suffix("RETURNING id, user_id, created_at, expires_at").
		PlaceholderFormat(sq.Dollar).RunWith(ss.db).QueryRow().
		Scan(&s.ID, &s.UserID, &s.CreatedAt, &s.ExpiresAt); err != nil {
		return nil, fmt.Errorf("an error occurred while creating a session")
	}

	return &s, nil
}

func (ss *sessionStore) Delete(id string) error {
	if _, err := sq.Delete("sessions").Where("id = ?", id).
		PlaceholderFormat(sq.Dollar).RunWith(ss.db).Query(); err != nil {
		return fmt.Errorf("an error occurred while deleting a session")
	}

	return nil
}

func (ss *sessionStore) GetByID(id string) (*model.SessionUser, error) {
	var s model.SessionUser

	if err := sq.Select("s.id", "s.created_at", "s.expires_at", "u.id", "u.username", "u.kind").
		From("sessions as s").Where("s.id = ?", id).Where(sq.Gt{"expires_at": "now()"}).
		LeftJoin("users as u ON u.id = s.user_id").
		PlaceholderFormat(sq.Dollar).RunWith(ss.db).QueryRow().
		Scan(&s.Session.ID, &s.Session.CreatedAt, &s.Session.ExpiresAt, &s.User.ID, &s.User.Username, &s.User.Kind); err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("session not found")
		}

		return nil, fmt.Errorf("an error occurred while getting a session by id")
	}

	return &s, nil
}
