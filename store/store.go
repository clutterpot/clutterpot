package store

import (
	"github.com/clutterpot/clutterpot/validator"
	"github.com/jmoiron/sqlx"
)

type Store struct {
	db   *sqlx.DB
	val  *validator.Validator
	User *userStore
	File *fileStore
}

func New(db *sqlx.DB, val *validator.Validator) *Store {
	return &Store{
		db:   db,
		val:  val,
		User: newUserStore(db, val),
		File: newFileStore(db, val),
	}
}
