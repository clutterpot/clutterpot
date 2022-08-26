package middlewares

import (
	"context"
	"net/http"

	"github.com/clutterpot/clutterpot/dataloaders"
	"github.com/clutterpot/clutterpot/store"
)

func dataloader(store *store.Store) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), dataloaders.DataloaderKey, dataloaders.New(store))))
		})
	}
}
