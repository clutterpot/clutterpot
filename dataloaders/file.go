package dataloaders

import (
	"github.com/clutterpot/clutterpot/model"
	"github.com/clutterpot/clutterpot/store"

	"github.com/vikstrous/dataloadgen"
)

type fileLoader struct {
	store *store.Store
	byID  *dataloadgen.Loader[string, *model.File]
}

func newFileLoader(store *store.Store) *fileLoader { return &fileLoader{store: store} }

func (fl *fileLoader) GetByID() *dataloadgen.Loader[string, *model.File] {
	if fl.byID == nil {
		fl.byID = dataloadgen.NewLoader(func(keys []string) ([]*model.File, []error) {
			return fl.store.File.GetAllByIDs(keys)
		}, dataloadgen.WithBatchCapacity(50))
	}

	return fl.byID
}
