package dataloaders

import (
	"github.com/clutterpot/clutterpot/store"
)

type Dataloader struct {
	File *fileLoader
}

func New(store *store.Store) *Dataloader {
	return &Dataloader{
		File: newFileLoader(store),
	}
}
