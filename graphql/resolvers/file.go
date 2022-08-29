package resolvers

import (
	"context"

	"github.com/clutterpot/clutterpot/dataloaders"
	"github.com/clutterpot/clutterpot/model"

	"github.com/go-chi/jwtauth/v5"
)

// Query resolvers

func (r *queryResolver) File(ctx context.Context, id string) (*model.File, error) {
	return r.Store.File.GetByID(id)
}

func (r *queryResolver) Files(ctx context.Context, after, before *string, first, last *int, filter *model.FileFilter, sort *model.FileSort, order *model.Order) (*model.FileConnection, error) {
	return r.Store.File.GetAll(after, before, first, last, filter, sort, order)
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

func (r *fileResolver) Tags(ctx context.Context, obj *model.File, after, before *string, first, last *int, filter *model.TagFilter, sort *model.TagSort, order *model.Order) (*model.TagConnection, error) {
	var ownerID string
	_, claims, err := jwtauth.FromContext(ctx)
	if err == nil {
		ownerID = claims["uid"].(string)
	}

	return dataloaders.ForContext(ctx).Tag.GetAllByFileID(&ctx, ownerID, after, before, first, last, filter, sort, order).Load(obj.ID)
}

// RemoveTagsFromFilePayload resolvers

func (r *removeTagsFromFilePayloadResolver) File(ctx context.Context, obj *model.RemoveTagsFromFilePayload) (*model.File, error) {
	return dataloaders.ForContext(ctx).File.GetByID().Load(obj.FileID)
}

func (r *removeTagsFromFilePayloadResolver) Tags(ctx context.Context, obj *model.RemoveTagsFromFilePayload) ([]*model.Tag, error) {
	var ownerID string
	_, claims, err := jwtauth.FromContext(ctx)
	if err == nil {
		ownerID = claims["uid"].(string)
	}

	tags, errs := dataloaders.ForContext(ctx).Tag.GetByID(ownerID).LoadAll(obj.TagIDs)
	if len(errs) > 0 {
		return nil, errs[0]
	}

	return tags, nil
}
