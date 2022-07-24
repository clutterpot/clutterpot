package app

import (
	"github.com/clutterpot/clutterpot/db"

	"github.com/jmoiron/sqlx"
)

type App struct {
	db *sqlx.DB
}

func New() *App {
	return &App{
		db: db.Connect(),
	}
}
