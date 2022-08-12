package resolvers

import (
	"context"

	"github.com/clutterpot/clutterpot/model"
)

// Query resolvers

func (r *queryResolver) User(ctx context.Context, id string) (*model.User, error) {
	return r.Store.User.GetByID(id)
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
		id = &r.Auth.ForContext(ctx).UserID
	}

	return r.Store.User.Update(*id, input)
}
