package store

import "github.com/jmoiron/sqlx"

type Store struct {
	db   *sqlx.DB
	User *userStore
}

func New(db *sqlx.DB) *Store {
	return &Store{
		db:   db,
		User: newUserStore(db),
	}
}
