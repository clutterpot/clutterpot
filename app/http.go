package app

import (
	"log"
	"net/http"
	"os"

	"github.com/clutterpot/clutterpot/handlers"
	"github.com/clutterpot/clutterpot/middlewares"
)

func (app *App) serveHTTP() {
	app.http.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Clutterpot"))
	})

	log.Println(os.ExpandEnv("Clutterpot server is listening at http://localhost:${PORT}"))

	if err := http.ListenAndServe(os.ExpandEnv(":${PORT}"), app.http); err != nil {
		log.Fatal(err)
	}
}

func (app *App) registerMiddlewares() {
	app.http.Use(middlewares.New(app.store, app.auth)...)
}

func (app *App) registerHandlers() {
	// GraphQL API
	app.http.Mount("/", handlers.GQLHandler(app.auth, app.store, app.val))
}
