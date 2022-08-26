package dataloaders

import (
	"github.com/clutterpot/clutterpot/model"
	"github.com/clutterpot/clutterpot/store"
	"github.com/vikstrous/dataloadgen"
)

type userLoader struct {
	store *store.Store
	byID  *dataloadgen.Loader[string, *model.User]
}

func newUserLoader(store *store.Store) *userLoader { return &userLoader{store: store} }

func (ul *userLoader) GetByID() *dataloadgen.Loader[string, *model.User] {
	if ul.byID == nil {
		ul.byID = dataloadgen.NewLoader(func(keys []string) ([]*model.User, []error) {
			return ul.store.User.GetAllByIDs(keys)
		}, dataloadgen.WithBatchCapacity(50))
	}

	return ul.byID
}
