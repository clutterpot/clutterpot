package resolvers

import (
	"context"

	"github.com/clutterpot/clutterpot/model"

	"github.com/go-chi/jwtauth/v5"
)

// Query resolvers

func (r *queryResolver) Tag(ctx context.Context, id string) (*model.Tag, error) {
	var ownerID string
	_, claims, err := jwtauth.FromContext(ctx)
	if err == nil {
		ownerID = claims["uid"].(string)
	}

	return r.Store.Tag.GetByID(id, ownerID)
}

func (r *queryResolver) Tags(ctx context.Context, after, before *string, first, last *int, sort *model.TagSort, order *model.Order) (*model.TagConnection, error) {
	var ownerID string
	_, claims, err := jwtauth.FromContext(ctx)
	if err == nil {
		ownerID = claims["uid"].(string)
	}

	return r.Store.Tag.GetAll(ownerID, after, before, first, last, sort, order)
}

// Mutation resolvers

func (r *mutationResolver) CreateTag(ctx context.Context, input model.TagInput) (*model.Tag, error) {
	if input.Private {
		input.OwnerID = &r.Auth.ForContext(ctx).UserID
	}

	return r.Store.Tag.Create(input)
}
