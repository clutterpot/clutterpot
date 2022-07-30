package app

import (
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5/middleware"
)

func (app *App) serveHTTP() {
	app.http.Use(middleware.Logger)

	app.http.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Clutterpot"))
	})

	if err := http.ListenAndServe(os.ExpandEnv(":${PORT}"), app.http); err != nil {
		log.Fatal(err)
	}
}
