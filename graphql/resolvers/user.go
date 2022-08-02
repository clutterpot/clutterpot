package resolvers

import (
	"context"

	"github.com/clutterpot/clutterpot/model"
)

func (r *queryResolver) User(ctx context.Context, id string) (*model.User, error) {
	return r.Store.User.GetByID(id)
}

func (r *mutationResolver) CreateUser(ctx context.Context, input model.UserInput) (*model.User, error) {
	return r.Store.User.Create(input)
}
func (r *mutationResolver) UpdateUser(ctx context.Context, id string, input model.UserUpdateInput) (*model.User, error) {
	return r.Store.User.Update(id, input)
}
