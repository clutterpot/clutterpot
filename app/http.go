package app

import (
	"log"
	"net/http"
	"os"

	"github.com/clutterpot/clutterpot/handlers"
	"github.com/go-chi/chi/v5/middleware"
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

func (app *App) registerHandlers() {
	app.http.Use(middleware.Logger)

	// GraphQL API
	app.http.Mount("/", handlers.GQLHandler(app.store, app.val))
}
