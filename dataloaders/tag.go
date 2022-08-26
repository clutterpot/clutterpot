package dataloaders

import (
	"context"

	"github.com/clutterpot/clutterpot/model"
	"github.com/clutterpot/clutterpot/store"

	"github.com/vikstrous/dataloadgen"
)

type tagLoader struct {
	store       *store.Store
	byID        *dataloadgen.Loader[string, *model.Tag]
	allByFileID *dataloadgen.Loader[string, *model.TagConnection]
}

func newTagLoader(store *store.Store) *tagLoader { return &tagLoader{store: store} }

func (tl *tagLoader) GetByID(ownerID string) *dataloadgen.Loader[string, *model.Tag] {
	if tl.byID == nil {
		tl.byID = dataloadgen.NewLoader(func(keys []string) ([]*model.Tag, []error) {
			return tl.store.Tag.GetAllByIDs(keys, ownerID)
		}, dataloadgen.WithBatchCapacity(50))
	}

	return tl.byID
}

func (tl *tagLoader) GetAllByFileID(ctx *context.Context, ownerID string, after, before *string, first, last *int, filter *model.TagFilter, sort *model.TagSort, order *model.Order) *dataloadgen.Loader[string, *model.Connection[model.Tag]] {
	if tl.allByFileID == nil {
		tl.allByFileID = dataloadgen.NewLoader(func(keys []string) ([]*model.TagConnection, []error) {
			return tl.store.Tag.GetAllByFileIDs(keys, ownerID, after, before, first, last, filter, sort, order)
		}, dataloadgen.WithBatchCapacity(50))
	}

	return tl.allByFileID
}
