package middlewares

import (
	"net/http"

	"github.com/clutterpot/clutterpot/auth"
	"github.com/clutterpot/clutterpot/store"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/jwtauth/v5"
)

func New(store *store.Store, auth *auth.Auth) []func(http.Handler) http.Handler {
	return []func(http.Handler) http.Handler{
		middleware.Logger,
		jwtauth.Verifier(auth.JWT()),
		dataloader(store),
	}
}
