package resolvers

import (
	"context"

	"github.com/clutterpot/clutterpot/auth"
	"github.com/clutterpot/clutterpot/dataloaders"
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

func (r *queryResolver) Tags(ctx context.Context, after, before *string, first, last *int, filter *model.TagFilter, sort *model.TagSort, order *model.Order) (*model.TagConnection, error) {
	var ownerID string
	_, claims, err := jwtauth.FromContext(ctx)
	if err == nil {
		ownerID = claims["uid"].(string)
	}

	return r.Store.Tag.GetAll(ownerID, after, before, first, last, filter, sort, order)
}

// Mutation resolvers

func (r *mutationResolver) CreateTag(ctx context.Context, input model.TagInput) (*model.Tag, error) {
	if !input.Global {
		input.OwnerID = &auth.ForContext(ctx).UserID
	}

	return r.Store.Tag.Create(input)
}

// Tag resolvers

func (r *tagResolver) Global(ctx context.Context, obj *model.Tag) (bool, error) {
	return obj.OwnerID == nil, nil
}

func (r *tagResolver) Owner(ctx context.Context, obj *model.Tag) (*model.User, error) {
	if obj.OwnerID == nil {
		return nil, nil
	}

	return dataloaders.ForContext(ctx).User.GetByID().Load(*obj.OwnerID)
}
