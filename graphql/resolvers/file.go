package resolvers

import (
	"context"

	"github.com/clutterpot/clutterpot/model"
	"github.com/vektah/gqlparser/v2/gqlerror"
)

// Query resolvers

func (r *queryResolver) File(ctx context.Context, id string) (*model.File, error) {
	return r.Store.File.GetByID(id)
}

// Mutation resolvers

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

func (r *mutationResolver) DeleteFile(ctx context.Context, id string) (*model.DeleteFilePayload, error) {
	payload, err := r.Store.File.SoftDelete(id)
	if err != nil {
		return nil, err
	}
	payload.ID = id

	return payload, nil
}

func (r *mutationResolver) RemoveTagsFromFile(ctx context.Context, id string, tagIDs []string) (*model.RemoveTagsFromFilePayload, error) {
	return r.Store.File.RemoveTags(id, tagIDs)
}

// File resolvers

func (r *fileResolver) Tags(ctx context.Context, obj *model.File) ([]*model.Tag, error) {
	return r.Store.Tag.GetByFileID(obj.ID)
}

// RemoveTagsFromFilePayload resolvers

func (r *removeTagsFromFilePayloadResolver) File(ctx context.Context, obj *model.RemoveTagsFromFilePayload) (*model.File, error) {
	return r.Store.File.GetByID(obj.FileID)
}

func (r *removeTagsFromFilePayloadResolver) Tags(ctx context.Context, obj *model.RemoveTagsFromFilePayload) ([]*model.Tag, error) {
	return nil, gqlerror.Errorf("not implemented")
}
