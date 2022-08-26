package dataloaders

import (
	"github.com/clutterpot/clutterpot/model"
	"github.com/clutterpot/clutterpot/store"

	"github.com/vikstrous/dataloadgen"
)

type tagLoader struct {
	store *store.Store
	owner *dataloadgen.Loader[string, *model.User]
}

func newTagLoader(store *store.Store) *tagLoader { return &tagLoader{store: store} }

func (tl *tagLoader) Owner() *dataloadgen.Loader[string, *model.User] {
	if tl.owner == nil {
		tl.owner = dataloadgen.NewLoader(func(keys []string) ([]*model.User, []error) {
			return tl.store.User.GetAllByIDs(keys)
		}, dataloadgen.WithBatchCapacity(50))
	}

	return tl.owner
}
