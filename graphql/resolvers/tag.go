package resolvers

import (
	"context"

	"github.com/clutterpot/clutterpot/model"
)

// Mutation resolvers

func (r *mutationResolver) CreateTag(ctx context.Context, input model.TagInput) (*model.Tag, error) {
	if input.Private {
		input.OwnerID = &r.Auth.ForContext(ctx).UserID
	}

	return r.Store.Tag.Create(input)
}
