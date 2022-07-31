package handlers

import (
	"github.com/clutterpot/clutterpot/graphql/resolvers"
	"github.com/clutterpot/clutterpot/graphql/server"
	"github.com/clutterpot/clutterpot/store"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/go-chi/chi/v5"
)

func GQLHandler(store *store.Store) *chi.Mux {
	r := chi.NewRouter()

	gqlServer := handler.NewDefaultServer(server.NewExecutableSchema(server.Config{
		Resolvers: resolvers.New(store),
	}))

	r.Handle("/graphql", gqlServer)
	r.Handle("/playground", playground.Handler("GraphQL playground", "/graphql"))

	return r
}
