package store

import (
	"github.com/jmoiron/sqlx"
)

type Store struct {
	User    *userStore
	File    *fileStore
	Tag     *tagStore
	Session *sessionStore
}

func New(db *sqlx.DB) *Store {
	return &Store{
		User:    newUserStore(db),
		File:    newFileStore(db),
		Tag:     newTagStore(db),
		Session: newSessionStore(db),
	}
}
