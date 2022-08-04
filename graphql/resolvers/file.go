package resolvers

import (
	"context"

	"github.com/clutterpot/clutterpot/model"
)

func (r *queryResolver) File(ctx context.Context, id string) (*model.File, error) {
	return r.Store.File.GetByID(id)
}

func (r *mutationResolver) CreateFile(ctx context.Context, input model.FileInput) (*model.File, error) {
	if err := r.Validator.Validate(ctx, input); err != nil {
		return nil, nil
	}

	return r.Store.File.Create(input)
}

func (r *mutationResolver) UpdateFile(ctx context.Context, id string, input model.FileUpdateInput) (*model.File, error) {
	if err := r.Validator.Validate(ctx, input); err != nil {
		return nil, nil
	}

	return r.Store.File.Update(id, input)
}
