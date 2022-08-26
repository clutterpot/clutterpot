package app

import (
	"github.com/clutterpot/clutterpot/auth"
	"github.com/clutterpot/clutterpot/db"
	"github.com/clutterpot/clutterpot/store"
	"github.com/clutterpot/clutterpot/validator"

	"github.com/go-chi/chi/v5"
	"github.com/jmoiron/sqlx"
)

type App struct {
	db    *sqlx.DB
	http  *chi.Mux
	auth  *auth.Auth
	val   *validator.Validator
	store *store.Store
}

func New() *App {
	return &App{}
}

func (app *App) Init() *App {
	app.db = db.Connect()
	app.http = chi.NewRouter()
	app.auth = auth.New()
	app.val = validator.New()
	app.store = store.New(app.db)

	return app
}

func (app *App) Run() {
	app.migrateDB()
	app.registerMiddlewares()
	app.registerHandlers()
	app.serveHTTP()
}
