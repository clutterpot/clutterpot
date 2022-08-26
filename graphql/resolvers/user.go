package resolvers

import (
	"context"

	"github.com/clutterpot/clutterpot/auth"
	"github.com/clutterpot/clutterpot/dataloaders"
	"github.com/clutterpot/clutterpot/model"
)

// Query resolvers

func (r *queryResolver) User(ctx context.Context, id string) (*model.User, error) {
	return dataloaders.ForContext(ctx).User.GetByID().Load(id)
}

func (r *queryResolver) Users(ctx context.Context, after, before *string, first, last *int, filter *model.UserFilter, sort *model.UserSort, order *model.Order) (*model.UserConnection, error) {
	return r.Store.User.GetAll(after, before, first, last, filter, sort, order)
}

// Mutation resolvers

func (r *mutationResolver) CreateUser(ctx context.Context, input model.UserInput) (*model.User, error) {
	if err := r.Validator.Validate(ctx, input); err != nil {
		return nil, nil
	}

	return r.Store.User.Create(input)
}

func (r *mutationResolver) UpdateUser(ctx context.Context, id *string, input model.UserUpdateInput) (*model.User, error) {
	if err := r.Validator.Validate(ctx, input); err != nil {
		return nil, nil
	}
	if id == nil {
		// Get user id from access token claims
		id = &auth.ForContext(ctx).UserID
	}

	return r.Store.User.Update(*id, input)
}
