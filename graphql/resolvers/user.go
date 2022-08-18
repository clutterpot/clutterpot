package resolvers

import (
	"context"

	"github.com/clutterpot/clutterpot/model"
	"github.com/clutterpot/clutterpot/store/pagination"
)

// Query resolvers

func (r *queryResolver) User(ctx context.Context, id string) (*model.User, error) {
	return r.Store.User.GetByID(id)
}

func (r *queryResolver) Users(ctx context.Context, after, before *string, first, last *int, sort *model.UserSort, order *model.Order) (*model.UserConnection, error) {
	users, pageInfo, err := r.Store.User.GetAll(after, before, first, last, sort, order)
	if err != nil {
		return nil, err
	}

	uc := model.UserConnection{
		Nodes:    users,
		PageInfo: pageInfo,
	}

	uc.Edges = make([]*model.UserEdge, len(users))
	for i, u := range users {
		uc.Edges[i] = &model.UserEdge{
			Cursor: pagination.EncodeCursor(u.ID),
			Node:   u,
		}
	}

	return &uc, nil
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
