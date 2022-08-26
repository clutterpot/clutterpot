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
	User *userLoader
	File *fileLoader
	Tag  *tagLoader
}

func New(store *store.Store) *Dataloader {
	return &Dataloader{
		User: newUserLoader(store),
		File: newFileLoader(store),
		Tag:  newTagLoader(store),
	}
}

func ForContext(ctx context.Context) *Dataloader {
	return ctx.Value(DataloaderKey).(*Dataloader)
}
