package resolvers

//go:generate go run github.com/99designs/gqlgen generate

import (
	"github.com/clutterpot/clutterpot/auth"
	"github.com/clutterpot/clutterpot/graphql/server"
	"github.com/clutterpot/clutterpot/store"
	"github.com/clutterpot/clutterpot/validator"
)

type Resolver struct {
	Auth      *auth.Auth
	Store     *store.Store
	Validator *validator.Validator
}

func New(auth *auth.Auth, store *store.Store, val *validator.Validator) *Resolver {
	return &Resolver{Auth: auth, Store: store, Validator: val}
}

type queryResolver struct{ *Resolver }

func (r *Resolver) Query() server.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }

func (r *Resolver) Mutation() server.MutationResolver { return &mutationResolver{r} }
