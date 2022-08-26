package dataloaders

import (
	"context"

	"github.com/clutterpot/clutterpot/model"
	"github.com/clutterpot/clutterpot/store"

	"github.com/vikstrous/dataloadgen"
)

type fileLoader struct {
	store *store.Store
	tags  *dataloadgen.Loader[string, *model.TagConnection]
}

func newFileLoader(store *store.Store) *fileLoader { return &fileLoader{store: store} }

func (fl *fileLoader) Tags(ctx *context.Context, ownerID string, after, before *string, first, last *int, filter *model.TagFilter, sort *model.TagSort, order *model.Order) *dataloadgen.Loader[string, *model.Connection[model.Tag]] {
	if fl.tags == nil {
		fl.tags = dataloadgen.NewLoader(func(keys []string) ([]*model.TagConnection, []error) {
			return fl.store.Tag.GetAllByFileIDs(keys, ownerID, after, before, first, last, filter, sort, order)
		}, dataloadgen.WithBatchCapacity(50))
	}

	return fl.tags
}
