package resolvers

import (
	"context"

	"github.com/clutterpot/clutterpot/model"
)

func (r *queryResolver) User(ctx context.Context, id string) (*model.User, error) {
	return r.Store.User.GetByID(id)
}

func (r *userResolver) Kind(ctx context.Context, obj *model.User) (int, error) {
	return int(obj.Kind), nil
}
