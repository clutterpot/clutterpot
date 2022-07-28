package app

import (
	"github.com/clutterpot/clutterpot/db"
	"github.com/clutterpot/clutterpot/store"

	"github.com/jmoiron/sqlx"
)

type App struct {
	db    *sqlx.DB
	store *store.Store
}

func New() *App {
	return &App{}
}

func (app *App) Init() {
	app.db = db.Connect()
	app.store = store.New(app.db)
}
