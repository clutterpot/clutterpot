package handlers

import (
	"github.com/clutterpot/clutterpot/auth"
	"github.com/clutterpot/clutterpot/graphql/directives"
	"github.com/clutterpot/clutterpot/graphql/resolvers"
	"github.com/clutterpot/clutterpot/graphql/server"
	"github.com/clutterpot/clutterpot/store"
	"github.com/clutterpot/clutterpot/validator"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/go-chi/chi/v5"
)

func GQLHandler(auth *auth.Auth, store *store.Store, val *validator.Validator) *chi.Mux {
	r := chi.NewRouter()

	gqlServer := handler.NewDefaultServer(server.NewExecutableSchema(server.Config{
		Resolvers:  resolvers.New(auth, store, val),
		Directives: directives.New(),
	}))

	r.Handle("/graphql", gqlServer)
	r.Handle("/playground", playground.Handler("GraphQL playground", "/graphql"))

	return r
}
