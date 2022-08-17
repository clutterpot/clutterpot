package dataloaders

import (
	"github.com/clutterpot/clutterpot/model"
	"github.com/clutterpot/clutterpot/store"

	"github.com/vikstrous/dataloadgen"
)

type fileLoader struct {
	store *store.Store
	Tags  *dataloadgen.Loader[string, []*model.Tag]
}

func newFileLoader(store *store.Store) *fileLoader {
	fl := fileLoader{store: store}
	fl.Tags = fl.tags()

	return &fl
}

func (fl *fileLoader) tags() *dataloadgen.Loader[string, []*model.Tag] {
	return dataloadgen.NewLoader(func(keys []string) ([][]*model.Tag, []error) {
		return fl.store.Tag.GetByFileIDs(keys)
	})
}
