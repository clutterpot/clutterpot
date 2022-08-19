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

// Mutation resolvers

func (r *mutationResolver) CreateTag(ctx context.Context, input model.TagInput) (*model.Tag, error) {
	if input.Private {
		input.OwnerID = &r.Auth.ForContext(ctx).UserID
	}

	return r.Store.Tag.Create(input)
}
