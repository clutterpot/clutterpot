package dataloaders

import (
	"context"

	"github.com/clutterpot/clutterpot/store"
)

type contextKey struct {
	name string
}

var DataloaderKey = contextKey{"dataloader"}

type Dataloader struct {
	File *fileLoader
}

func New(store *store.Store) *Dataloader {
	return &Dataloader{
		File: newFileLoader(store),
	}
}

func ForContext(ctx context.Context) *Dataloader {
	return ctx.Value(DataloaderKey).(*Dataloader)
}
